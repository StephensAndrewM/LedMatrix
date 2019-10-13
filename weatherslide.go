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
)

type WeatherSlide struct {
    HttpHelper   *HttpHelper
    Weather      WeatherApiResponse
    WeatherIcons map[string]*image.RGBA

    RedrawTicker *time.Ticker
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

    // Set up HTTP fetcher
    url := fmt.Sprintf("https://api.darksky.net/forecast/%s/%s",
        WEATHER_API_KEY, latLng)
    refresh := 2 * time.Minute
    this.HttpHelper = NewHttpHelper(url, refresh, this.Parse)

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

    this.Weather = respData
    return true
}

func (this *WeatherSlide) Draw(img *image.RGBA) {
    // Stop immediately if we have errors
    if !this.HttpHelper.LastFetchSuccess {
        DrawError(img, "Weather", "No data.")
        return
    }

    white := color.RGBA{255, 255, 255, 255}

    WriteString(img, "NOW", white, ALIGN_CENTER, WEATHER_COL_CENTER, 0)
    currentIcon, currentIconExists := this.WeatherIcons[this.Weather.Current.Icon]
    if currentIconExists {
        iconLeftX := WEATHER_COL_CENTER - (WEATHER_ICON_WIDTH / 2)
        DrawImageWithColorTransform(img, currentIcon, iconLeftX, 7, white)
    } else {
        log.WithFields(log.Fields{
            "icon": this.Weather.Current.Icon,
        }).Warn("Missing icon for weather condition.")
    }
    WriteString(img, fmt.Sprintf("%.1f°", this.Weather.Current.Temperature), white, ALIGN_CENTER, WEATHER_COL_CENTER, 24)

    // If afternoon, offset forecasts by one day (don't show today)
    forecastOffset := 0
    if time.Now().Hour() > 12 {
        forecastOffset = forecastOffset + 1
    }

    this.DrawForecast(img, WEATHER_COL_WIDTH, this.Weather.Daily.Data[forecastOffset+0])

    this.DrawTemperatureGraph(img)
    this.DrawPrecipitationGraph(img)
}

func (this *WeatherSlide) DrawTemperatureGraph(img *image.RGBA) {
    white := color.RGBA{255, 255, 255, 255}
    originX := 80
    originY := 12
    height := 10
    width := 48

    // Thermometer symbol
    WriteString(img, "🌡", white, ALIGN_LEFT, originX-8, originY-8)

    var timeValues []int64
    var dataPoints []float64
    for _, val := range this.Weather.Hourly.Data[:48] {
        timeValues = append(timeValues, val.Time)
        dataPoints = append(dataPoints, val.Temperature)
    }

    DrawAutoNormalizedGraph(img, originX, originY-1, height, white, dataPoints)
    this.DrawTimeAxes(img, originX, originY, width, height, timeValues)
}

func (this *WeatherSlide) DrawPrecipitationGraph(img *image.RGBA) {
    white := color.RGBA{255, 255, 255, 255}
    originX := 80
    originY := 28
    height := 10
    width := 48

    // Raindrop symbol
    WriteString(img, "💧", white, ALIGN_LEFT, originX-8, originY-8)

    var timeValues []int64
    var dataPoints []float64
    for _, val := range this.Weather.Hourly.Data[:48] {
        timeValues = append(timeValues, val.Time)
        dataPoints = append(dataPoints, val.PrecipProbability)
    }

    DrawNormalizedGraph(img, originX, originY-1, height, 0.0, 1.0, white, dataPoints)
    this.DrawTimeAxes(img, originX, originY, width, height, timeValues)
}

func (this *WeatherSlide) DrawTimeAxes(img *image.RGBA, originX, originY, width, height int, timeValues []int64) {
    yellow := color.RGBA{255, 255, 0, 255}

    DrawVertLine(img, yellow, originY-height, originY, originX-1)
    DrawHorizLine(img, yellow, originX, originX+width, originY)

    // Draw emphasis on noon/midnight
    for i, val := range timeValues {
        t := time.Unix(val, 0)
        if t.Hour() == 0 {
            DrawVertLine(img, yellow, originY, originY+2, originX+i)
        }
        if t.Hour() == 12 {
            DrawVertLine(img, yellow, originY, originY+1, originX+i)
        }
    }
}

func (this *WeatherSlide) DrawForecast(img *image.RGBA, offsetX int, forecast WeatherApiDailyForecastData) {
    aqua := color.RGBA{0, 255, 255, 255}

    label := strings.ToUpper(time.Unix(forecast.Time, 0).Weekday().String()[0:3])
    WriteString(img, label, aqua, ALIGN_CENTER, offsetX+WEATHER_COL_CENTER, 0)

    icon, iconExists := this.WeatherIcons[forecast.Icon]
    if iconExists {
        iconLeftX := offsetX + (WEATHER_COL_CENTER - (WEATHER_ICON_WIDTH / 2))
        DrawImageWithColorTransform(img, icon, iconLeftX, 7, aqua)
    } else {
        log.WithFields(log.Fields{
            "icon": this.Weather.Current.Icon,
        }).Warn("Missing icon for weather condition.")
    }

    temp := fmt.Sprintf("%d°/%d°", int(forecast.High), int(forecast.Low))
    WriteString(img, temp, aqua, ALIGN_CENTER, offsetX+WEATHER_COL_CENTER, 24)
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
