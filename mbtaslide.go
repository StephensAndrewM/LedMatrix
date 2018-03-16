package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "math"
    "net/http"
    "sort"
    "time"
)

type MbtaSlide struct {
    Station     string
    Line        string
    Predictions []MbtaPrediction
}

const MBTA_API_URL = "https://api-v3.mbta.com/predictions"
const STATION_PREDICTION_QUERY = "?filter[stop]=%s&filter[route]=%s"

const MBTA_STATION_DAVIS = "place-davis"
const MBTA_ROUTE_RED = "Red"

var StationLabels = map[string]string{
    MBTA_STATION_DAVIS: "DAVIS SQ",
}

// TODO find a way to generalize this
var DavisDirectionLabels = map[int]string{
    0: "IN",
    1: "OUT",
}

func NewMbtaSlide(line, station string) *MbtaSlide {
    sl := new(MbtaSlide)
    sl.Station = station
    sl.Line = line
    return sl
}

func (sl *MbtaSlide) Preload() {

    // Load live Data from MBTA
    resp, httpErr := http.Get(MBTA_API_URL +
        fmt.Sprintf(STATION_PREDICTION_QUERY, sl.Station, sl.Line))
    if httpErr != nil {
        fmt.Printf("Error loading MBTA data: %s\n", httpErr)
        return
        // TODO Display error on screen
    }

    // Parse response to JSON
    respBuf := new(bytes.Buffer)
    respBuf.ReadFrom(resp.Body)
    var respData MbtaApiPredictionResponse
    jsonErr := json.Unmarshal(respBuf.Bytes(), &respData)
    if jsonErr != nil {
        fmt.Printf("Error interpreting MBTA data: %s\n", jsonErr)
        return
        // TODO Display error on screen
    }

    // Convert MBTA data structures to a more workable format and sort them
    var predictions []MbtaPrediction
    for _, data := range respData.Data {
        t, tErr := time.Parse(time.RFC3339, data.Attributes.DepartureTime)
        if tErr != nil {
            fmt.Printf(
                "Error interpreting MBTA time: %s struct: %s\n", tErr, data)
            return
        }
        p := MbtaPrediction{
            Direction: data.Attributes.DirectionId,
            Time:      t,
        }
        predictions = append(predictions, p)
    }
    sort.Slice(predictions, func(i, j int) bool {
        return predictions[i].Time.Before(predictions[j].Time)
    })
    sl.Predictions = predictions
}

func (sl *MbtaSlide) Draw(s *Surface) {
    s.Clear()
    if len(sl.Predictions) == 0 {
        fmt.Println("No predictions")
        return
    }
    white := Color{255, 255, 255}
    blank := Color{0, 0, 0}
    red := Color{255, 0, 0}
    s.DrawBox(red, 0, 0, s.Width, 9)
    s.WriteString(StationLabels[sl.Station], blank, ALIGN_CENTER, s.Width/2, 1)
    for i := 0; i < min(3, len(sl.Predictions)); i++ {
        y := ((i + 1) * 8) + 1
        s.WriteString(
            DavisDirectionLabels[sl.Predictions[i].Direction],
            white,
            ALIGN_LEFT,
            0,
            y)
        est := sl.Predictions[i].Time.Sub(time.Now())
        s.WriteString(
            fmt.Sprintf("%d MIN", int(math.Floor(est.Minutes()))),
            white,
            ALIGN_RIGHT,
            s.Width-1,
            y)
    }
}

// Data structures used by the MBTA API - used for parsing responses
type MbtaApiPredictionResponse struct {
    Data []MbtaApiPredictionResource `json:"data"`
}

type MbtaApiPredictionResource struct {
    // Other properties are provided but unused
    Attributes MbtaApiPredictionAttributes `json:"attributes"`
}

type MbtaApiPredictionAttributes struct {
    // Other properties are provided but unused
    DirectionId   int    `json:"direction_id"`
    DepartureTime string `json:"departure_time"`
}

// A simpler representation for use internally
type MbtaPrediction struct {
    Direction int
    Time      time.Time
}

func min(x, y int) int {
    if x < y {
        return x
    } else {
        return y
    }
}
