package main

import (
    "encoding/json"
    "fmt"
    "math"
    "sort"
    "strings"
    "time"
)

type MbtaSlide struct {
    StationName string
    HttpHelper  *HttpHelper
    Predictions []MbtaPrediction
}

const MBTA_STATION_ID_DAVIS = "place-davis"
const MBTA_STATION_NAME_DAVIS = "DAVIS SQUARE"

func NewMbtaSlide(stationId, stationName string) *MbtaSlide {
    this := new(MbtaSlide)
    this.StationName = stationName

    // Set up HTTP fetcher
    url := fmt.Sprintf("https://api-v3.mbta.com/predictions"+
        "?include=route,trip&filter[stop]=%s", stationId)
    refresh := 60 * time.Second
    this.HttpHelper = NewHttpHelper(url, refresh)

    return this
}

func (this *MbtaSlide) Preload() {

    // Load live Data from MBTA
    respBytes, ok := this.HttpHelper.Fetch()
    if !ok {
        fmt.Printf("Error loading MBTA data\n")
        return
        // TODO Display error on screen
    }

    // Parse response to JSON
    var respData MbtaApiResponse
    jsonErr := json.Unmarshal(respBytes, &respData)
    if jsonErr != nil {
        fmt.Printf("Error interpreting MBTA data: %s\n", jsonErr)
        return
        // TODO Display error on screen
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
            // This logic assumes that a station only serves red line and bus
            // Needs work to support different lines
            if r.Attributes.Type == 1 {
                routeDef.Type = MbtaRouteTypeRedLine
            } else {
                routeDef.Type = MbtaRouteTypeBus
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

func (this *MbtaSlide) Draw(s *Surface) {
    s.Clear()
    white := Color{255, 255, 255}
    yellow := Color{255, 255, 0}
    red := Color{255, 0, 0}
    blank := Color{0, 0, 0}

    s.WriteString(this.StationName, red, ALIGN_CENTER, s.Width/2, 1)

    if len(this.Predictions) == 0 {
        fmt.Println("No predictions")
        return
    }

    n := 0 // Count of valid predictions found - we skip some
    for i := 0; i < len(this.Predictions); i++ {
        p := this.Predictions[i]
        y := ((n + 1) * 8) + 1

        // Get time estimate, and maybe skip
        est := p.Time.Sub(time.Now())
        estMin := int(math.Floor(est.Minutes()))
        // Some predictions go negative - ignore those
        if estMin < 0 {
            continue
        }

        // Route identifier
        if p.Route.Type == MbtaRouteTypeRedLine {
            s.DrawBox(red, 0, y, 11, 7)
            s.WriteString("R", blank, ALIGN_CENTER, 5, y)
        } else {
            s.WriteString(p.Route.Id, yellow, ALIGN_CENTER, 5, y)
        }

        // Destination
        dest := strings.ToUpper(p.Destination)
        s.WriteString(dest, white, ALIGN_LEFT, 13, y)

        // Time estimate
        estStr := fmt.Sprintf("%d min", estMin)
        s.WriteString(estStr, white, ALIGN_RIGHT, s.Width-1, y)

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
    MbtaRouteTypeRedLine MbtaRouteType = iota
    MbtaRouteTypeBus
)

type MbtaRoute struct {
    Type MbtaRouteType
    Id   string
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
