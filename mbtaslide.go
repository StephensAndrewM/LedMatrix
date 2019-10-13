package main

import (
    "encoding/json"
    "fmt"
    log "github.com/sirupsen/logrus"
    "image"
    "image/color"
    "math"
    "sort"
    "strconv"
    "strings"
    "time"
)

type MbtaSlide struct {
    StationName string
    HttpHelper  *HttpHelper
    Predictions []MbtaPrediction
}

const MBTA_SLIDE_ERROR_SPACE = 3

// Station names - used in constructor
// This is not an exhaustive list, just some easy ones
const MBTA_STATION_ID_DAVIS = "place-davis"
const MBTA_STATION_ID_PARK = "place-pktrm"
const MBTA_STATION_ID_MGH = "place-chmnl"
const MBTA_STATION_ID_GOVCTR = "place-gover"
const MBTA_STATION_ID_HARVARD = "place-harsq"

var MBTA_STATION_NAME_MAP = map[string]string{
    MBTA_STATION_ID_DAVIS:   "DAVIS SQUARE",
    MBTA_STATION_ID_PARK:    "PARK STREET",
    MBTA_STATION_ID_MGH:     "CHARLES/MGH",
    MBTA_STATION_ID_GOVCTR:  "GOVERNMENT CENTER",
    MBTA_STATION_ID_HARVARD: "HARVARD SQUARE",
}

func NewMbtaSlide(stationId string) *MbtaSlide {
    this := new(MbtaSlide)
    name, ok := MBTA_STATION_NAME_MAP[stationId]
    if !ok {
        log.WithFields(log.Fields{
            "stationId": stationId,
        }).Warn("Could not find station name.")
        name = "?????"
    }
    this.StationName = name

    // Set up HTTP fetcher
    url := fmt.Sprintf("https://api-v3.mbta.com/predictions"+
        "?include=route,trip&filter[stop]=%s", stationId)
    refresh := 60 * time.Second
    this.HttpHelper = NewHttpHelper(url, refresh, this.Parse)

    return this
}

func (this *MbtaSlide) Initialize() {
    this.HttpHelper.StartLoop()
}

func (this *MbtaSlide) Terminate() {
    this.HttpHelper.StopLoop()
}

func (this *MbtaSlide) Parse(respBytes []byte) bool {
    // Parse response to JSON
    var resp MbtaApiResponse
    jsonErr := json.Unmarshal(respBytes, &resp)
    if jsonErr != nil {
        log.WithFields(log.Fields{
            "error": jsonErr,
        }).Warn("Error unmarshalling MBTA data.")
        return false
    }

    // We need to resolve a trip ID into structured information
    routeByTripId := this.BuildTripIdToRouteMap(resp.Included)

    // Create a mapping of predicted times by route
    predictionsByRoute := this.BuildRouteToPredictionsMap(resp.Data, routeByTripId)

    // Flatten the predictions into what will be displayed
    this.Predictions = this.FlattenPredictions(routeByTripId, predictionsByRoute)

    return true
}

func (this *MbtaSlide) BuildTripIdToRouteMap(resources []MbtaApiResource) map[string]MbtaRoute {
    // Iterate through all provided "route" resources, building a mapping
    // of route ID (string) to structure route object (with name and color).
    // These route objects do *not* have the destination property set.
    routeDefs := make(map[string]MbtaRoute)
    for _, r := range resources {
        if r.Type == "route" {
            routeDef := MbtaRoute{}
            // We deliberately don't set Dest since we don't know it here
            routeDef.Id = r.Id
            routeDef.Color = r.Attributes.Color
            switch r.Attributes.Type {
            case 0:
                routeDef.Type = MbtaRouteTypeLightRail
            case 1:
                routeDef.Type = MbtaRouteTypeHeavyRail
            case 2:
                routeDef.Type = MbtaRouteTypeCommuterRail
            case 3:
                routeDef.Type = MbtaRouteTypeBus
            default:
                routeDef.Type = MbtaRouteTypeUnknown
            }
            routeDefs[r.Id] = routeDef
        }
    }

    // Iterate again, looking at "trip" resources instead. This resource has
    // the headsign attribute that we use to set Destination on the route,
    // and the ID that we use as the map's key.
    m := make(map[string]MbtaRoute)
    for _, r := range resources {
        if r.Type == "trip" {
            route, ok := routeDefs[r.Relationships.Route.Data.Id]
            if ok {
                m[r.Id] = MbtaRoute{
                    Id:          route.Id,
                    Color:       route.Color,
                    Type:        route.Type,
                    Destination: r.Attributes.Headsign,
                }
            } else {
                log.WithFields(log.Fields{
                    "tripId": r.Relationships.Route.Data.Id,
                }).Warn("Could not find MBTA route data.")
            }
        }
    }
    return m
}

func (this *MbtaSlide) BuildRouteToPredictionsMap(data []MbtaApiResource, routeByTripId map[string]MbtaRoute) map[string][]time.Time {
    predictions := make(map[string][]time.Time)
    for _, r := range data {
        if r.Type == "prediction" {
            // Some vehicles don't give departure estimates - ignore them
            if len(r.Attributes.DepartureTime) == 0 {
                continue
            }
            // Parse the time into a standard format
            t, tErr := time.Parse(time.RFC3339, r.Attributes.DepartureTime)
            if tErr != nil {
                log.WithFields(log.Fields{
                    "error": tErr,
                    "value": r.Attributes.DepartureTime,
                }).Warn("Error interpreting MBTA time.")
                continue
            }
            // Get data about the trip supplied in the prediction
            tr, ok := routeByTripId[r.Relationships.Trip.Data.Id]
            if !ok {
                log.WithFields(log.Fields{
                    "tripId": r.Relationships.Trip.Data.Id,
                }).Warn("Error interpreting MBTA trip ID in prediction.")
                continue
            }
            // Turn the route and destination struct into a string and store
            k := this.RouteToString(tr)
            predictions[k] = append(predictions[k], t)
        }
    }
    return predictions
}

func (this *MbtaSlide) FlattenPredictions(routeByTripId map[string]MbtaRoute, predictionsByRoute map[string][]time.Time) []MbtaPrediction {
    // Create a lookup map of route string -> object
    routeLookup := make(map[string]MbtaRoute)
    for _, v := range routeByTripId {
        routeLookup[this.RouteToString(v)] = v
    }

    // Create a list of objects containing route objects and times
    var predictions []MbtaPrediction
    for k, v := range predictionsByRoute {
        // Don't include routes with no predictions
        if len(v) == 0 {
            continue
        }
        // Translate the route string (key) into an actual object
        r, ok := routeLookup[k]
        if !ok {
            log.WithFields(log.Fields{
                "routeString": k,
            }).Warn("MBTA route object not found for string.")
            continue
        }
        // Finally, add it to the list
        p := MbtaPrediction{
            Route: r,
            Time:  v,
        }
        predictions = append(predictions, p)
    }

    return predictions
}

// Squash a MbtaRoute object to a simple string representation
func (this *MbtaSlide) RouteToString(r MbtaRoute) string {
    return fmt.Sprintf("%s/%s/%s/%s", r.Type, r.Color, r.Id, r.Destination)
}

// Finds the earliest time in an unsorted array
func (this *MbtaSlide) GetMinTime(t []time.Time) time.Time {
    sort.Slice(t, func(i, j int) bool {
        return t[i].Before(t[j])
    })
    return t[0]
}

// For all the times stored in all prediction objects, remove those in the past
func (this *MbtaSlide) FilterTimesInPast(all []MbtaPrediction) (ret []MbtaPrediction) {
    for _, p := range all {
        var times []time.Time
        for _, t := range p.Time {
            if t.Sub(time.Now()) >= 0 {
                times = append(times, t)
            }
        }
        if len(times) > 0 {
            p.Time = times
            ret = append(ret, p)
        }
    }
    return
}

func (this *MbtaSlide) Draw(img *image.RGBA) {
    if !this.HttpHelper.LastFetchSuccess {
        DrawError(img, MBTA_SLIDE_ERROR_SPACE, 1)
        return
    }

    textColor := color.RGBA{255, 255, 255, 255} // white
    titleColor := color.RGBA{255, 255, 0, 255}  // yellow
    busColor := color.RGBA{255, 255, 0, 255}    // yellow
    timeColor := color.RGBA{0, 255, 255, 255}   // aqua

    WriteString(img, this.StationName, titleColor, ALIGN_CENTER, GetLeftOfCenterX(img), 0)

    filteredPredictions := this.FilterTimesInPast(this.Predictions)

    if len(this.Predictions) == 0 {
        // TODO display something here
        return
    }

    // Resort prediction time sets based on current time
    sort.Slice(filteredPredictions, func(i, j int) bool {
        return this.GetMinTime(filteredPredictions[i].Time).Before(this.GetMinTime(filteredPredictions[j].Time))
    })

    // TODO rotate through destinations if there are more than three
    o := 0
    predictionSubset := filteredPredictions[o:min(o+3, len(filteredPredictions))]

    for i, p := range predictionSubset {
        // Calculate vertical position of line
        y := ((i + 1) * 8)

        var estStrs []string
        // Loop through first three predictions, or all, whichever is less
        for j := 0; j < min(len(p.Time), 3); j++ {
            t := p.Time[j]
            est := t.Sub(time.Now())
            estMin := int(math.Floor(est.Minutes()))
            estStrs = append(estStrs, strconv.Itoa(estMin))
        }
        estStr := strings.Join(estStrs, ",_") + "_min"

        // Draw a box for line color, or a bus number when relevant
        if p.Route.Type == MbtaRouteTypeBus {
            WriteString(img, p.Route.Id, busColor, ALIGN_CENTER, 5, y)
        } else {
            lineColor := ColorFromHex(p.Route.Color)
            reducedLineColor := ReduceColor(lineColor)
            DrawBox(img, reducedLineColor, 0, y, 11, 7)
        }

        // Size of box is different based on how many time digits to display
        destWidth := 116 - GetDisplayWidth(estStr)

        // Destination
        dest := strings.ToUpper(p.Route.Destination)
        WriteStringBoxed(img, dest, textColor, ALIGN_LEFT, 12, y, destWidth)

        // Time estimate
        imgWidth := img.Bounds().Dx()
        WriteString(img, estStr, timeColor, ALIGN_RIGHT, imgWidth-1, y)
    }
}

// Data structures used by the MBTA API - used for parsing responses
type MbtaApiResponse struct {
    Included []MbtaApiResource `json:"included"`
    Data     []MbtaApiResource `json:"data"`
}

type MbtaApiResource struct {
    // Other properties are provided but unused
    Type          string                      `json:"type"`
    Relationships MbtaApiResourceRelationship `json:"relationships"`
    Id            string                      `json:"id"`
    Attributes    MbtaApiResourceAttributes   `json:"attributes"`
}

type MbtaApiResourceRelationship struct {
    Trip  MbtaApiResourceRelationshipTrip  `json:"trip"`
    Route MbtaApiResourceRelationshipRoute `json:"route"`
}

type MbtaApiResourceRelationshipTrip struct {
    Data MbtaApiResourceRelationshipData `json:"data"`
}

type MbtaApiResourceRelationshipRoute struct {
    Data MbtaApiResourceRelationshipData `json:"data"`
}

type MbtaApiResourceRelationshipData struct {
    Type string `json:"type"`
    Id   string `json:"id"`
}

type MbtaApiResourceAttributes struct {
    DepartureTime string `json:"departure_time"` // For prediction
    Headsign      string `json:"headsign"`       // For trip
    Color         string `json:"color"`          // For route
    Type          int    `json:"type"`           // For route
}

// A simpler representation for use internally
type MbtaPrediction struct {
    Route MbtaRoute
    Time  []time.Time
}

type MbtaRouteType int

const (
    MbtaRouteTypeUnknown MbtaRouteType = iota
    MbtaRouteTypeLightRail
    MbtaRouteTypeHeavyRail
    MbtaRouteTypeCommuterRail
    MbtaRouteTypeBus
)

type MbtaRoute struct {
    Type        MbtaRouteType
    Id          string
    Color       string
    Destination string
}

func min(x, y int) int {
    if x < y {
        return x
    } else {
        return y
    }
}
