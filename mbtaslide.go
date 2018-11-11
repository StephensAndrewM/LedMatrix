package main

import (
    "encoding/json"
    "fmt"
    "image"
    "image/color"
    "math"
    "sort"
    "strings"
    "time"
)

type MbtaSlide struct {
    StationName string
    HttpHelper  *HttpHelper
    Predictions []MbtaPrediction

    // Status of loading content
    LastFetchHttpErr bool
    LastFetchJsonErr bool
}

const MBTA_SLIDE_ERROR_SPACE = 3
// Lowest duration of prediction allowed to show
const MBTA_ARRIVAL_THRESHOLD = 5

// Station names - used in constructor
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
    if ok != true {
        fmt.Printf("Could not find station name for %s\n", stationId)
        name = "?????"
    }
    this.StationName = name

    // Set up HTTP fetcher
    url := fmt.Sprintf("https://api-v3.mbta.com/predictions"+
        "?include=route,trip&filter[stop]=%s", stationId)
    refresh := 60 * time.Second
    this.HttpHelper = NewHttpHelper(url, refresh)

    return this
}

func (this *MbtaSlide) Preload() {
    // Reset errors, in case last time wasn't successful
    this.LastFetchHttpErr = false
    this.LastFetchJsonErr = false

    // Load live Data from MBTA
    respBytes, ok := this.HttpHelper.Fetch()
    if !ok {
        fmt.Printf("Error loading MBTA data\n")
        this.LastFetchHttpErr = true
        return
    }

    // Parse response to JSON
    var respData MbtaApiResponse
    jsonErr := json.Unmarshal(respBytes, &respData)
    if jsonErr != nil {
        fmt.Printf("Error interpreting MBTA data: %s\n", jsonErr)
        this.LastFetchJsonErr = true
        return
    }

    this.ParsePredictions(respData)
}

func (this *MbtaSlide) ParsePredictions(resp MbtaApiResponse) {

    trips := this.GetTripDataByTripId(resp.Included)

    var predictions []MbtaPrediction
    for _, r := range resp.Data {
        if r.Type == "prediction" {
            // Some vehicles don't give departure estimates - ignore them
            if len(r.Attributes.DepartureTime) == 0 {
                continue
            }
            // Parse the time into a standard format
            t, tErr := time.Parse(time.RFC3339, r.Attributes.DepartureTime)
            if tErr != nil {
                fmt.Printf(
                    "Error interpreting MBTA time: %s struct: %s\n", tErr, r)
                continue
            }
            // Get data about the trip supplied in the prediction
            tr, ok := trips[r.Relationships.Trip.Data.Id]
            if ok {
                p := MbtaPrediction{
                    Route:       tr.Route,
                    Destination: tr.Headsign,
                    Time:        t,
                }
                predictions = append(predictions, p)
            } else {
                fmt.Printf(
                    "Error interpreting MBTA Trip ID: %s\n",
                    r.Relationships.Trip.Data.Id)
                continue
            }
        }
    }
    sort.Slice(predictions, func(i, j int) bool {
        return predictions[i].Time.Before(predictions[j].Time)
    })
    this.Predictions = predictions
}

func (this *MbtaSlide) GetTripDataByTripId(resources []MbtaApiResource) map[string]MbtaTrip {

    // Build a map of Route data keyed by Route ID
    routeDefs := make(map[string]MbtaRoute)
    for _, r := range resources {
        if r.Type == "route" {
            routeDef := MbtaRoute{}
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

    m := make(map[string]MbtaTrip)
    for _, r := range resources {
        if r.Type == "trip" {
            route, ok := routeDefs[r.Relationships.Route.Data.Id]
            if ok {
                m[r.Id] = MbtaTrip{route, r.Attributes.Headsign}
            } else {
                fmt.Printf("Could not find route data for %s\n",
                    r.Relationships.Route.Data.Id)
            }
        }
    }
    return m
}

func (this *MbtaSlide) Draw(img *image.RGBA) {

    // Stop immediately if we have errors
    if this.LastFetchHttpErr {
        DrawError(img, MBTA_SLIDE_ERROR_SPACE, 1)
        return
    }
    if this.LastFetchJsonErr {
        DrawError(img, MBTA_SLIDE_ERROR_SPACE, 2)
        return
    }

    textColor := color.RGBA{255, 255, 255, 255} // white
    titleColor := color.RGBA{255, 255, 0, 255}  // yellow

    WriteString(img, this.StationName, titleColor, ALIGN_CENTER, GetLeftOfCenterX(img), 1)

    if len(this.Predictions) == 0 {
        return
    }

    n := 0 // Count of valid predictions found - we skip some
    for i := 0; i < len(this.Predictions); i++ {
        p := this.Predictions[i]
        y := ((n + 1) * 8) + 1

        // Get time estimate, and maybe skip
        est := p.Time.Sub(time.Now())
        estMin := int(math.Floor(est.Minutes()))
        // Low predictions aren't useful (uness we run), so don't display any
        // trains less than X minutes away
        if estMin < MBTA_ARRIVAL_THRESHOLD {
            continue
        }

        if p.Route.Type == MbtaRouteTypeBus {
            WriteString(img, p.Route.Id, titleColor, ALIGN_CENTER, 5, y)
        } else {
            lineColor := ColorFromHex(p.Route.Color)
            reducedLineColor := ReduceColor(lineColor)
            DrawBox(img, reducedLineColor, 0, y, 11, 7)
        }

        // Size of box is different based on how many time digits to display
        destWidth := 93
        if estMin > 9 {
            destWidth = 87
        }

        // Destination
        dest := strings.ToUpper(p.Destination)
        WriteStringBoxed(img, dest, textColor, ALIGN_LEFT, 12, y, destWidth)

        // Time estimate
        estStr := fmt.Sprintf("%d_min", estMin)
        imgWidth := img.Bounds().Dx()
        WriteString(img, estStr, textColor, ALIGN_RIGHT, imgWidth-1, y)

        n++
        // We can't display more than 3 predictions on screen so stop
        if n >= 3 {
            break
        }
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
    Route       MbtaRoute
    Destination string
    Time        time.Time
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
    Type  MbtaRouteType
    Id    string
    Color string
}

type MbtaTrip struct {
    Route    MbtaRoute
    Headsign string
}

func min(x, y int) int {
    if x < y {
        return x
    } else {
        return y
    }
}
