package main

import (
    "encoding/json"
    "flag"
    "fmt"
    log "github.com/sirupsen/logrus"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "os"
    "strings"
    "time"
    "net/http"
)

type WeatherSlide struct {
    LatLng     string

    // Weather      WeatherApiResponse
    Weather WeatherData
    WeatherIcons map[string]*image.RGBA

    HttpHelper *HttpHelper
    RedrawTicker *time.Ticker
}

type WeatherData struct {
    CurrentTemp float64
    CurrentIcon *image.RGBA

    ForecastWeekday time.Weekday
    ForecastIcon *image.RGBA
    ForecastHighTemp int
    ForecastLowTemp int

    TimeGraphValues []time.Time
    TempGraphValues []float64
    PrecipGraphValues []float64
}

var weatherIconBaseDirFlag = flag.String("weather_icon_base_dir", "",
    "If specified, base directory to load weather icons from.")

// Values to use as parameter for initializing slide
const BOSTON_LATLNG = "42.2129,-71.0349"

// Define a few constants for drawing
const WEATHER_COL_WIDTH = 32
const WEATHER_COL_CENTER = 16
const WEATHER_ICON_WIDTH = 16

var WEATHER_API_ICON_MAP = map[string]string{
    "rain":                "rain1.xbm.png",
    "snow":                "snou.xbm.png",
    "sleet":               "rain2.xbm.png",
    "wind":                "wind.xbm.png",
    "fog":                 "wind.xbm.png",
    "cloudy":              "clouds.xbm.png",
    "partly-cloudy-day":   "cloud_sun.xbm.png",
    "partly-cloudy-night": "cloud_moon.xbm.png",
    "clear-day":           "sun.xbm.png",
    "clear-night":         "moon.xbm.png",
}

func NewWeatherSlide(latLng string) *WeatherSlide {
    this := new(WeatherSlide)
    this.LatLng = latLng
    this.HttpHelper = NewHttpHelper(this)

    // Preload all the weather icons
    this.WeatherIcons = make(map[string]*image.RGBA)
    for k := range WEATHER_API_ICON_MAP {
        f := *weatherIconBaseDirFlag +
            "icons/weather/" + WEATHER_API_ICON_MAP[k]
        // Open the file as binary stream
        reader, err1 := os.Open(f)
        if err1 != nil {
            log.WithFields(log.Fields{
                "file":  f,
                "error": err1,
            }).Warn("Could not open image.")
            continue
        }
        defer reader.Close()
        // Attempt to convert the image to image.Image
        img, err2 := png.Decode(reader)
        if err2 != nil {
            log.WithFields(log.Fields{
                "file":  f,
                "error": err2,
            }).Warn("Could not decode image.")
            continue
        }
        // Then convert that to image.RGBA
        b := img.Bounds()
        imgRgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
        draw.Draw(imgRgba, imgRgba.Bounds(), img, b.Min, draw.Src)
        this.WeatherIcons[k] = imgRgba
    }

    return this
}

func (this *WeatherSlide) Initialize() {
    this.HttpHelper.StartLoop()
}

func (this *WeatherSlide) Terminate() {
    this.HttpHelper.StopLoop()
}

func (this *WeatherSlide) StartDraw(d Display) {
    this.RedrawTicker = DrawEverySecond(d, this.Draw)
}

func (this *WeatherSlide) StopDraw() {
    this.RedrawTicker.Stop()
}

func (this *WeatherSlide) IsEnabled() bool {
    return true // Always enabled
}

func (this *WeatherSlide) GetRefreshInterval() time.Duration {
    return 2 * time.Minute
}

func (this *WeatherSlide) BuildRequest() (*http.Request, error) {
    url := fmt.Sprintf("https://api.darksky.net/forecast/%s/%s",
        WEATHER_API_KEY, this.LatLng)

    return http.NewRequest("GET", url, nil)
}

func (this *WeatherSlide) Parse(respBytes []byte) bool {
    // Parse response to JSON
    var respData WeatherApiResponse
    jsonErr := json.Unmarshal(respBytes, &respData)
    if jsonErr != nil {
        log.WithFields(log.Fields{
            "error": jsonErr,
        }).Warn("Could not interpret weather JSON.")
        return false
    }

    // Assert that the response contains what we expect
    if respData.Current.Icon == "" ||
        len(respData.Daily.Data) == 0 {
        log.WithFields(log.Fields{
            "error": respData.Error,
        }).Warn("Weather response data has no data.")
        return false
    }

    var weather WeatherData

    // Convert data on current conditions
    weather.CurrentTemp = respData.Current.Temperature;
    currentIcon, currentIconExists := this.WeatherIcons[respData.Current.Icon]
    if currentIconExists {
        weather.CurrentIcon = currentIcon
    } else {
        log.WithFields(log.Fields{
            "icon": respData.Current.Icon,
        }).Warn("Missing icon for current.")
    }

    // Convert data on today's/tomorrow's forecast
    forecastFromApi := respData.Daily.Data[0]
    if (time.Now().Hour() > 12) {
        forecastFromApi = respData.Daily.Data[1]
    }
    weather.ForecastWeekday = time.Unix(forecastFromApi.Time, 0).Weekday()
    forecastIcon, forecastIconExists := this.WeatherIcons[forecastFromApi.Icon]
    if forecastIconExists {
        weather.ForecastIcon = forecastIcon
    } else {
        log.WithFields(log.Fields{
            "icon": forecastFromApi.Icon,
        }).Warn("Missing icon for forecast.")
    }
    weather.ForecastHighTemp = int(forecastFromApi.High)
    weather.ForecastLowTemp = int(forecastFromApi.Low)

    // Convert data on hourly temperature/precipitation forecast
    for _, val := range respData.Hourly.Data[:48] {
        weather.TimeGraphValues = append(weather.TimeGraphValues, time.Unix(val.Time, 0))
        weather.TempGraphValues = append(weather.TempGraphValues, val.Temperature)
        weather.PrecipGraphValues = append(weather.PrecipGraphValues, val.PrecipProbability)
    }

    this.Weather = weather;

    return true
}

func (this *WeatherSlide) Draw(img *image.RGBA) {
    // Stop immediately if we have errors
    if !this.HttpHelper.LastFetchSuccess {
        DrawError(img, "Weather", "No data.")
        return
    }

    white := color.RGBA{255, 255, 255, 255}

    WriteString(img, "NOW", white, ALIGN_CENTER, 16, 0)
    if this.Weather.CurrentIcon != nil {
        DrawImageWithColorTransform(img, this.Weather.CurrentIcon, 8, 7, white)
    }
    WriteString(img, fmt.Sprintf("%.1f°", this.Weather.CurrentTemp), white, ALIGN_CENTER, 16, 24)

    forecastOriginX := 33

    aqua := color.RGBA{0, 255, 255, 255}
    label := strings.ToUpper(this.Weather.ForecastWeekday.String()[0:3])
    WriteString(img, label, aqua, ALIGN_CENTER, forecastOriginX+16, 0)
    if this.Weather.ForecastIcon != nil {
        DrawImageWithColorTransform(img, this.Weather.ForecastIcon, forecastOriginX+8, 7, aqua)
    }
    forecastTemp := fmt.Sprintf("%d°/%d°", this.Weather.ForecastHighTemp, this.Weather.ForecastLowTemp)
    WriteString(img, forecastTemp, aqua, ALIGN_CENTER, forecastOriginX+16, 24)

    graphOriginX := 80
    graphHeight := 10
    graphWidth := 48

    // Thermometer symbol
    tempGraphOriginY := 12
    WriteString(img, "🌡", white, ALIGN_LEFT, graphOriginX-6, tempGraphOriginY-8)
    DrawAutoNormalizedGraph(img, graphOriginX, tempGraphOriginY-1, graphHeight, white, this.Weather.TempGraphValues)
    this.DrawTimeAxes(img, graphOriginX, tempGraphOriginY, graphWidth, graphHeight, this.Weather.TimeGraphValues)

    // Raindrop symbol
    rainGraphOriginY := 28
    WriteString(img, "💧", white, ALIGN_LEFT, graphOriginX-6, rainGraphOriginY-8)
    DrawNormalizedGraph(img, graphOriginX, rainGraphOriginY-1, graphHeight, 0.0, 1.0, white, this.Weather.PrecipGraphValues)
    this.DrawTimeAxes(img, graphOriginX, rainGraphOriginY, graphWidth, graphHeight, this.Weather.TimeGraphValues)
}

func (this *WeatherSlide) DrawTimeAxes(img *image.RGBA, originX, originY, width, height int, timeValues []time.Time) {
    yellow := color.RGBA{255, 255, 0, 255}

    DrawVertLine(img, yellow, originY-height, originY, originX-1)
    DrawHorizLine(img, yellow, originX, originX+width-1, originY)

    // Draw emphasis on noon/midnight
    for i, t := range timeValues {
        if t.Hour() == 0 {
            DrawVertLine(img, yellow, originY, originY+2, originX+i)
        }
        if t.Hour() == 12 {
            DrawVertLine(img, yellow, originY, originY+1, originX+i)
        }
    }
}

// Data structures used by the Weather Underground API
type WeatherApiResponse struct {
    Current WeatherApiCurrentConditions `json:"currently"`
    Daily   WeatherApiDailyForecast     `json:"daily"`
    Hourly  WeatherApiHourlyForecast    `json:"hourly"`
    Code    int                         `json:"code"`
    Error   string                      `json:"error"`
}

type WeatherApiCurrentConditions struct {
    Icon              string  `json:"icon"`
    Temperature       float64 `json:"temperature"`
    PrecipProbability float64 `json:"precipProbability"`
}

type WeatherApiDailyForecast struct {
    Data []WeatherApiDailyForecastData `json:"data"`
}

type WeatherApiHourlyForecast struct {
    Data []WeatherApiHourlyForecastData `json:"data"`
}

type WeatherApiDailyForecastData struct {
    Time              int64   `json:"time"`
    Icon              string  `json:"icon"`
    High              float64 `json:"temperatureHigh"`
    Low               float64 `json:"temperatureLow"`
    PrecipProbability float64 `json:"precipProbability"`
}

type WeatherApiHourlyForecastData struct {
    Time              int64   `json:"time"`
    Icon              string  `json:"icon"`
    Temperature       float64 `json:"temperature"`
    PrecipProbability float64 `json:"precipProbability"`
}
