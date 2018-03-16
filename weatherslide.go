package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type WeatherSlide struct {
    ZipCode int
}

const WEATHER_API_URL = "http://api.wunderground.com/api"
const WEATHER_API_QUERY = "/%s/conditions/forecast/q/%d.json"

const SUNNYVALE_ZIP = 94086

func NewWeatherSlide(zipCode int) *WeatherSlide {
    sl := new(WeatherSlide)
    sl.ZipCode = zipCode
    return sl
}

func (sl WeatherSlide) Preload() {

    // Load live Data from Weather Underground
    resp, httpErr := http.Get(WEATHER_API_URL +
        fmt.Sprintf(WEATHER_API_QUERY, WEATHER_API_KEY, sl.ZipCode))
    if httpErr != nil {
        fmt.Printf("Error loading Weather data: %s\n", httpErr)
        return
        // TODO Display error on screen
    }

    // Parse response to JSON
    respBuf := new(bytes.Buffer)
    respBuf.ReadFrom(resp.Body)
    var respData WeatherApiResponse
    jsonErr := json.Unmarshal(respBuf.Bytes(), &respData)
    if jsonErr != nil {
        fmt.Printf("Error interpreting Weather data: %s\n", jsonErr)
        return
        // TODO Display error on screen
    }
    fmt.Printf("Weather result is %+v", respData)
}

func (sl WeatherSlide) Draw(s *Surface) {

    s.Clear()
    white := Color{255, 255, 255}
    green := Color{0, 255, 0}
    yellow := Color{255, 255, 0}
    s.WriteString("NOW:", white, ALIGN_LEFT, 0, 0)
    s.WriteString("65°", green, ALIGN_RIGHT, s.Width-1, 0)
    s.WriteString("CLEAR", green, ALIGN_RIGHT, s.Width-1, 8)
    s.WriteString("TMRW:", white, ALIGN_LEFT, 0, 16)
    s.WriteString("60°", yellow, ALIGN_RIGHT, s.Width-1, 16)
    s.WriteString("RAIN", yellow, ALIGN_RIGHT, s.Width-1, 24)
}

// Data structures used by the Weather Underground API
type WeatherApiResponse struct {
    Observations WeatherObservations `json:"current_observation"`
    Forecast     WeatherForecast     `json:"forecast"`
}

type WeatherObservations struct {
    TempF float64 `json:"temp_f"`
    Icon  string  `json:"icon"`
}

type WeatherForecast struct {
    SimpleForecast WeatherSimpleForecast `json:"simpleforecast"`
}

type WeatherSimpleForecast struct {
    ForecastDay []WeatherForecastDay `json:"forecastday"`
}

type WeatherForecastDay struct {
    Date WeatherForecastDayDate        `json:"date"`
    High WeatherForecastDayTemperature `json:"high"`
    Low  WeatherForecastDayTemperature `json:"low"`
    Icon string                        `json:"icon"`
}

type WeatherForecastDayDate struct {
    Epoch        string `json:"epoch"`
    WeekdayShort string `json:"weekday_short"`
}

type WeatherForecastDayTemperature struct {
    Fahrenheit string `json:"fahrenheit"`
    Celsius    string `json:"celsius"`
}
