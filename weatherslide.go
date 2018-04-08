package main

import (
    "encoding/json"
    "fmt"
    "time"
)

type WeatherSlide struct {
    HttpHelper *HttpHelper
    Weather WeatherApiResponse
}

const SUNNYVALE_ZIP = 94086

var WEATHER_API_ICON_MAP = map[string]string{
    // TODO fill this out with all possible icons
    // TODO make this use actual images
    "rain": "RN",
    "cloudy": "CLO",
    "partlycloudy": "PCL",
    "mostlycloudy": "MCL",
    "clear": "CLR",
}

func NewWeatherSlide(zipCode int) *WeatherSlide {
    this := new(WeatherSlide)

    // Set up HTTP fetcher
    url := fmt.Sprintf("http://api.wunderground.com/api"+
        "/%s/conditions/forecast/q/%d.json", WEATHER_API_KEY, zipCode)
    refresh := 60 * time.Second
    this.HttpHelper = NewHttpHelper(url, refresh)
    
    return this
}

func (this *WeatherSlide) Preload() {

    // Load live Data from MBTA
    respBytes, ok := this.HttpHelper.Fetch()
    if !ok {
        fmt.Printf("Error loading Weather data\n")
        return
        // TODO Display error on screen
    }

    // Parse response to JSON
    var respData WeatherApiResponse
    jsonErr := json.Unmarshal(respBytes, &respData)
    if jsonErr != nil {
        fmt.Printf("Error interpreting Weather data: %s\n", jsonErr)
        return
        // TODO Display error on screen
    }
    fmt.Printf("Weather result is %+v", respData)

    this.Weather = respData
}

func (this *WeatherSlide) Draw(s *Surface) {
    s.Clear()
    white := Color{255, 255, 255}
    yellow := Color{255, 255, 0}
    
    const dayLabelXOffset = 48
    const tempXOffset = 82
    s.WriteString("Now:", yellow, ALIGN_RIGHT, dayLabelXOffset, 2)
    this.WriteWeatherString(s, this.Weather.Observations.Icon, 2)
    s.WriteString(
        fmt.Sprintf("%.1f°", this.Weather.Observations.TempF), white, ALIGN_LEFT, tempXOffset, 2)

    s.WriteString("Tomorrow:", yellow, ALIGN_RIGHT, dayLabelXOffset, 12)
    this.WriteWeatherString(s, this.Weather.Forecast.SimpleForecast.ForecastDay[0].Icon, 12)
    s.WriteString(
        fmt.Sprintf("%s°/%s°", 
            this.Weather.Forecast.SimpleForecast.ForecastDay[0].High.Fahrenheit,
            this.Weather.Forecast.SimpleForecast.ForecastDay[0].Low.Fahrenheit), 
        white, ALIGN_LEFT, tempXOffset, 12)

    s.WriteString(this.Weather.Forecast.SimpleForecast.ForecastDay[1].Date.Weekday + ":", yellow, ALIGN_RIGHT, dayLabelXOffset, 22)
    this.WriteWeatherString(s, this.Weather.Forecast.SimpleForecast.ForecastDay[1].Icon, 22)
        s.WriteString(
        fmt.Sprintf("%s°/%s°", 
            this.Weather.Forecast.SimpleForecast.ForecastDay[1].High.Fahrenheit,
            this.Weather.Forecast.SimpleForecast.ForecastDay[1].Low.Fahrenheit), 
        white, ALIGN_LEFT, tempXOffset, 22)
}

func (this *WeatherSlide) WriteWeatherString(s *Surface, condition string, yOffset int) {
    white := Color{255, 255, 255}
    const conditionXOffset = 54
    icon, ok := WEATHER_API_ICON_MAP[condition]
    if ok {
        s.WriteString(icon, white, ALIGN_LEFT, conditionXOffset, yOffset)
    } else {
        fmt.Printf("Unknown condition %s\n", condition)
    }
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
    Weekday string `json:"weekday"`
}

type WeatherForecastDayTemperature struct {
    Fahrenheit string `json:"fahrenheit"`
    Celsius    string `json:"celsius"`
}
