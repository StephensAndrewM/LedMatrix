package main

import (
    "encoding/json"
    "fmt"
    "time"
    "image"
    "image/color"
)

type WeatherSlide struct {
    HttpHelper *HttpHelper
    Weather WeatherApiResponse
}

const SUNNYVALE_ZIP = 94086
const BOSTON_ZIP = 02114

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
    sl := new(WeatherSlide)

    // Set up HTTP fetcher
    url := fmt.Sprintf("http://api.wunderground.com/api"+
        "/%s/conditions/forecast/q/%d.json", WEATHER_API_KEY, zipCode)
    refresh := 60 * time.Second
    sl.HttpHelper = NewHttpHelper(url, refresh)
    
    return sl
}

func (this *WeatherSlide) Preload() {

    // Load live Data from API
    respBytes, ok := this.HttpHelper.Fetch()
    if !ok {
        fmt.Printf("Error loading weather data\n")
        return
    }

    // Parse response to JSON
    var respData WeatherApiResponse
    jsonErr := json.Unmarshal(respBytes, &respData)
    if jsonErr != nil {
        fmt.Printf("Error interpreting Weather data: %s\n", jsonErr)
        return
    }
    fmt.Printf("Weather result is %+v", respData)

    this.Weather = respData
}

func (this *WeatherSlide) Draw(img *image.RGBA) {
    white := color.RGBA{255, 255, 255, 255}
    yellow := color.RGBA{255, 255, 0, 255}
    
    const dayLabelXOffset = 48
    const tempXOffset = 82
    WriteString(img, "Now:", yellow, ALIGN_RIGHT, dayLabelXOffset, 2)
    this.WriteWeatherString(img, this.Weather.Observations.Icon, 2)
    WriteString(img, 
        fmt.Sprintf("%.1f°", this.Weather.Observations.TempF), white, ALIGN_LEFT, tempXOffset, 2)

    WriteString(img, "Tomorrow:", yellow, ALIGN_RIGHT, dayLabelXOffset, 12)
    this.WriteWeatherString(img, this.Weather.Forecast.SimpleForecast.ForecastDay[0].Icon, 12)
    WriteString(img, 
        fmt.Sprintf("%s°/%s°", 
            this.Weather.Forecast.SimpleForecast.ForecastDay[0].High.Fahrenheit,
            this.Weather.Forecast.SimpleForecast.ForecastDay[0].Low.Fahrenheit), 
        white, ALIGN_LEFT, tempXOffset, 12)

    WriteString(img, this.Weather.Forecast.SimpleForecast.ForecastDay[1].Date.Weekday + ":", yellow, ALIGN_RIGHT, dayLabelXOffset, 22)
    this.WriteWeatherString(img, this.Weather.Forecast.SimpleForecast.ForecastDay[1].Icon, 22)
        WriteString(img, 
        fmt.Sprintf("%s°/%s°", 
            this.Weather.Forecast.SimpleForecast.ForecastDay[1].High.Fahrenheit,
            this.Weather.Forecast.SimpleForecast.ForecastDay[1].Low.Fahrenheit), 
        white, ALIGN_LEFT, tempXOffset, 22)
}

func (this *WeatherSlide) WriteWeatherString(img *image.RGBA, condition string, yOffset int) {
    white := color.RGBA{255, 255, 255, 255}
    const conditionXOffset = 54
    icon, ok := WEATHER_API_ICON_MAP[condition]
    if ok {
        WriteString(img, icon, white, ALIGN_LEFT, conditionXOffset, yOffset)
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
