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

    // Status of loading content
    LastFetchHttpErr bool
    LastFetchJsonErr bool
    LastFetchDataErr bool
}

var weatherIconBaseDirFlag = flag.String("weather_icon_base_dir", "",
    "If specified, base directory to load weather icons from.")

// Values to use as parameter for initializing slide
const BOSTON_LATLNG = "42.2129,-71.0349"

// Define a few constants for drawing
const WEATHER_COL_WIDTH = 42
const WEATHER_COL_CENTER = 21
const WEATHER_ICON_WIDTH = 16
const WEATHER_SLIDE_ERROR_SPACE = 4

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
    sl := new(WeatherSlide)

    // Set up HTTP fetcher
    url := fmt.Sprintf("https://api.darksky.net/forecast/%s/%s",
        WEATHER_API_KEY, latLng)
    refresh := 5 * time.Minute
    sl.HttpHelper = NewHttpHelper(url, refresh)
    // Block drawing until we get a response
    sl.LastFetchHttpErr = true

    // Preload all the weather icons
    sl.WeatherIcons = make(map[string]*image.RGBA)
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
        sl.WeatherIcons[k] = imgRgba
    }

    return sl
}

func (this *WeatherSlide) Preload() {
    // Load live Data from API
    respBytes, ok := this.HttpHelper.Fetch()
    if !ok {
        log.Warn("Error loading weather data")
        this.LastFetchHttpErr = true
        return
    }
    this.LastFetchHttpErr = false

    // Parse response to JSON
    var respData WeatherApiResponse
    jsonErr := json.Unmarshal(respBytes, &respData)
    if jsonErr != nil {
        log.WithFields(log.Fields{
            "error": jsonErr,
        }).Warn("Could not interpret weather JSON.")
        this.LastFetchJsonErr = true
        return
    }
    this.LastFetchJsonErr = false

    // Assert that the response contains what we expect
    if respData.Current.Icon == "" ||
        len(respData.Daily.Data) == 0 {
        log.Warn("Weather response data has no data.")
        this.LastFetchDataErr = true
    }
    this.LastFetchDataErr = false
    this.Weather = respData
}

func (this *WeatherSlide) Draw(img *image.RGBA) {

    // Stop immediately if we have errors
    if this.LastFetchHttpErr {
        DrawError(img, WEATHER_SLIDE_ERROR_SPACE, 1)
        return
    }
    if this.LastFetchJsonErr {
        DrawError(img, WEATHER_SLIDE_ERROR_SPACE, 2)
        return
    }
    if this.LastFetchDataErr {
        DrawError(img, WEATHER_SLIDE_ERROR_SPACE, 3)
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
    this.DrawForecast(img, WEATHER_COL_WIDTH*2, this.Weather.Daily.Data[forecastOffset+1])
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
}

type WeatherApiCurrentConditions struct {
    Icon              string  `json:"icon"`
    Temperature       float64 `json:"temperature"`
    PrecipProbability float64 `json:"precipProbability"`
}

type WeatherApiDailyForecast struct {
    Data []WeatherApiDailyForecastData `json:"data"`
}

type WeatherApiDailyForecastData struct {
    Time              int64   `json:"time"`
    Icon              string  `json:"icon"`
    High              float64 `json:"temperatureHigh"`
    Low               float64 `json:"temperatureLow"`
    PrecipProbability float64 `json:"precipProbability"`
}
