package main

import (
    "cloud.google.com/go/civil"
    "encoding/json"
    "flag"
    "fmt"
    log "github.com/sirupsen/logrus"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "net/http"
    "os"
    "strings"
    "time"
)

type WeatherSlide struct {
    Lat string
    Lng string

    Weather      WeatherData
    WeatherIcons map[string]*image.RGBA

    RealtimeHttpHelper       *HttpHelper
    DailyForecastHttpHelper  *HttpHelper
    HourlyForecastHttpHelper *HttpHelper
    RedrawTicker             *time.Ticker
}

type WeatherData struct {
    CurrentTemp float64
    CurrentIcon *image.RGBA

    ForecastWeekday  time.Weekday
    ForecastIcon     *image.RGBA
    ForecastHighTemp int
    ForecastLowTemp  int

    TimeGraphValues   []time.Time
    TempGraphValues   []float64
    PrecipGraphValues []float64
}

var weatherIconBaseDirFlag = flag.String("weather_icon_base_dir", "",
    "If specified, base directory to load weather icons from.")

// Latitude/longitude values for API requests
const BOSTON_LAT = "42.2129"
const BOSTON_LNG = "-71.0349"

const WEATHER_MAX_HOURLY_DATA = 51

// ClimaCell API provides these possible weather_code values
var WEATHER_API_ICON_MAP = map[string]string{
    "freezing_rain_heavy": "rain_snow",
    "freezing_rain":       "rain_snow",
    "freezing_rain_light": "rain_snow",
    "freezing_drizzle":    "rain_snow",
    "ice_pellets_heavy":   "rain_snow",
    "ice_pellets":         "rain_snow",
    "ice_pellets_light":   "rain_snow",
    "snow_heavy":          "snou",
    "snow":                "snou",
    "snow_light":          "snou",
    "flurries":            "snow_sun",
    "tstorm":              "rain1",
    "rain_heavy":          "rain1",
    "rain":                "rain1",
    "rain_light":          "rain0",
    "drizzle":             "rain0_sun",
    "fog_light":           "cloud_wind",
    "fog":                 "cloud_wind",
    "cloudy":              "clouds",
    "mostly_cloudy":       "clouds",
    "partly_cloudy":       "cloud_sun",
    "mostly_clear":        "sun",
    "clear":               "sun",
}

func NewWeatherSlide(lat, lng string) *WeatherSlide {
    this := new(WeatherSlide)
    this.Lat = lat
    this.Lng = lng
    this.RealtimeHttpHelper = NewHttpHelper(HttpConfig{
        SlideId:            "WeatherSlide-Realtime",
        RefreshInterval:    5 * time.Minute,
        RequestUrlCallback: this.GetRealtimeUrl,
        ParseCallback:      this.ParseRealtime,
    })
    this.HourlyForecastHttpHelper = NewHttpHelper(HttpConfig{
        SlideId:            "WeatherSlide-HourlyForecast",
        RefreshInterval:    30 * time.Minute,
        RequestUrlCallback: this.GetHourlyUrl,
        ParseCallback:      this.ParseHourly,
    })
    this.DailyForecastHttpHelper = NewHttpHelper(HttpConfig{
        SlideId:            "WeatherSlide-DailyForecast",
        RefreshInterval:    30 * time.Minute,
        RequestUrlCallback: this.GetDailyUrl,
        ParseCallback:      this.ParseDaily,
    })

    // Preload all the weather icons
    this.WeatherIcons = make(map[string]*image.RGBA)
    for k := range WEATHER_API_ICON_MAP {
        f := *weatherIconBaseDirFlag +
            "icons/weather/" + WEATHER_API_ICON_MAP[k] + ".xbm.png"
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
    this.RealtimeHttpHelper.StartLoop()
    this.HourlyForecastHttpHelper.StartLoop()
    this.DailyForecastHttpHelper.StartLoop()
}

func (this *WeatherSlide) Terminate() {
    this.RealtimeHttpHelper.StopLoop()
    this.HourlyForecastHttpHelper.StopLoop()
    this.DailyForecastHttpHelper.StopLoop()
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

func (this *WeatherSlide) BuildUrl(endpoint, responseFields string, includeDates bool) (*http.Request, error) {
    extraParams := ""
    if includeDates {
        end_time := time.Now().Add((WEATHER_MAX_HOURLY_DATA + 2) * time.Hour).Format(time.RFC3339)
        extraParams = fmt.Sprintf("&start_time=now&end_time=%s", end_time)
    }
    url := fmt.Sprintf("https://api.climacell.co/v3/weather/%s?lat=%s&lon=%s&unit_system=us&apikey=%s&fields=%s%s", endpoint, this.Lat, this.Lng, WEATHER_API_KEY, responseFields, extraParams)
    return http.NewRequest("GET", url, nil)
}

func (this *WeatherSlide) GetRealtimeUrl() (*http.Request, error) {
    return this.BuildUrl("realtime", "temp,weather_code", false)
}

func (this *WeatherSlide) GetHourlyUrl() (*http.Request, error) {
    return this.BuildUrl("forecast/hourly", "temp,precipitation_probability", true)
}

func (this *WeatherSlide) GetDailyUrl() (*http.Request, error) {
    return this.BuildUrl("forecast/daily", "temp,weather_code", true)
}

func (this *WeatherSlide) ParseRealtime(respBytes []byte) bool {
    var respData WeatherApiRealtimeResponse
    err := json.Unmarshal(respBytes, &respData)
    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Warn("Could not interpret realtime weather JSON.")
        return false
    }

    this.Weather.CurrentTemp = respData.Temp.Value
    this.Weather.CurrentIcon = this.GetIcon(respData.Code.Value)
    return true
}

func (this *WeatherSlide) ParseDaily(respBytes []byte) bool {
    var respData []WeatherApiDailyResponse
    err := json.Unmarshal(respBytes, &respData)
    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Warn("Could not interpret daily weather JSON.")
        return false
    }

    if len(respData) < 2 {
        log.WithFields(log.Fields{
            "length": len(respData),
        }).Warn("Fewer days than expected returned for daily weather request.")
        return false
    }

    // Use data for tomorrow if current time is after noon
    forecastFromApi := respData[0]
    if time.Now().Hour() > 12 {
        forecastFromApi = respData[1]
    }

    log.WithFields(log.Fields{
        "data": forecastFromApi,
    }).Debug("Forecast data")

    forecastDate, err := civil.ParseDate(forecastFromApi.ObservationTime.Value)
    if err != nil {
        log.WithFields(log.Fields{
            "data": forecastFromApi.ObservationTime,
        }).Warn("Could not parse day in daily forecast.")
        return false
    }
    this.Weather.ForecastWeekday = forecastDate.In(time.UTC).Weekday()

    this.Weather.ForecastIcon = this.GetIcon(forecastFromApi.Code.Value)

    low := 0.0
    high := 0.0
    for _, singleForecast := range forecastFromApi.Temp {
        if singleForecast.Min.Value != 0 {
            low = singleForecast.Min.Value
        }
        if singleForecast.Max.Value != 0 {
            high = singleForecast.Max.Value
        }
    }
    if low == 0 || high == 0 {
        log.WithFields(log.Fields{
            "data": forecastFromApi.Temp,
        }).Warn("Could not find high/low temperature in daily forecast.")
        return false
    }

    this.Weather.ForecastHighTemp = int(high)
    this.Weather.ForecastLowTemp = int(low)

    return true
}

func (this *WeatherSlide) ParseHourly(respBytes []byte) bool {
    var respData []WeatherApiHourlyResponse
    err := json.Unmarshal(respBytes, &respData)
    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Warn("Could not interpret hourly weather JSON.")
        return false
    }

    if len(respData) < WEATHER_MAX_HOURLY_DATA {
        log.WithFields(log.Fields{
            "length": len(respData),
        }).Warn("Fewer hours than expected returned for hourly weather request.")
        return false
    }

    for _, val := range respData {
        t, err := time.Parse(time.RFC3339, val.ObservationTime.Value)
        if err != nil {
            log.WithFields(log.Fields{
                "value": val.ObservationTime,
            }).Warn("Could not parse time in hourly forecast.")
            return false
        }
        this.Weather.TimeGraphValues = append(this.Weather.TimeGraphValues, t)
        this.Weather.TempGraphValues = append(this.Weather.TempGraphValues, val.Temp.Value)
        this.Weather.PrecipGraphValues = append(this.Weather.PrecipGraphValues, val.PrecipProbability.Value)
    }

    return true
}

func (this *WeatherSlide) GetIcon(condition string) *image.RGBA {
    icon, ok := this.WeatherIcons[condition]
    if !ok {
        log.WithFields(log.Fields{
            "condition": condition,
        }).Warn("Missing icon for weather condition.")
        return nil
    }
    return icon
}

func (this *WeatherSlide) Draw(img *image.RGBA) {
    // Stop immediately if we have errors
    if !this.RealtimeHttpHelper.LastFetchSuccess || !this.DailyForecastHttpHelper.LastFetchSuccess || !this.HourlyForecastHttpHelper.LastFetchSuccess {
        DrawError(img, "Weather", "No data.")
        return
    }

    white := color.RGBA{255, 255, 255, 255}

    WriteString(img, "NOW", white, ALIGN_CENTER, 16, 0)
    if this.Weather.CurrentIcon != nil {
        DrawImageWithColorTransform(img, this.Weather.CurrentIcon, 8, 7, white)
    }
    WriteString(img, fmt.Sprintf("%.1f°", this.Weather.CurrentTemp), white, ALIGN_CENTER, 16, 24)

    forecastOriginX := 32

    aqua := color.RGBA{0, 255, 255, 255}
    label := strings.ToUpper(this.Weather.ForecastWeekday.String()[0:3])
    WriteString(img, label, aqua, ALIGN_CENTER, forecastOriginX+16, 0)
    if this.Weather.ForecastIcon != nil {
        DrawImageWithColorTransform(img, this.Weather.ForecastIcon, forecastOriginX+8, 7, aqua)
    }
    forecastTemp := fmt.Sprintf("%d°/%d°", this.Weather.ForecastHighTemp, this.Weather.ForecastLowTemp)
    WriteString(img, forecastTemp, aqua, ALIGN_CENTER, forecastOriginX+16, 24)

    graphWidth := WEATHER_MAX_HOURLY_DATA
    graphHeight := 10
    graphOriginX := 128-graphWidth

    // Thermometer symbol
    tempGraphOriginY := 12
    WriteString(img, "🌡", white, ALIGN_LEFT, graphOriginX-6, tempGraphOriginY-8)
    DrawAutoNormalizedGraph(img, graphOriginX, tempGraphOriginY-1, graphHeight, white, this.Weather.TempGraphValues)
    this.DrawTimeAxes(img, graphOriginX, tempGraphOriginY, graphWidth, graphHeight, this.Weather.TimeGraphValues)

    // Raindrop symbol
    rainGraphOriginY := 28
    WriteString(img, "💧", white, ALIGN_LEFT, graphOriginX-6, rainGraphOriginY-8)
    DrawNormalizedGraph(img, graphOriginX, rainGraphOriginY-1, graphHeight, 0.0, 100.0, white, this.Weather.PrecipGraphValues)
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

// Data structures used by the ClimaCell JSON API
type WeatherApiRealtimeResponse struct {
    Temp WeatherApiFloatValue  `json:"temp"`
    Code WeatherApiStringValue `json:"weather_code"`
}

type WeatherApiHourlyResponse struct {
    Temp              WeatherApiFloatValue  `json:"temp"`
    PrecipProbability WeatherApiFloatValue  `json:"precipitation_probability"`
    ObservationTime   WeatherApiStringValue `json:"observation_time"`
}

type WeatherApiDailyResponse struct {
    Temp            []WeatherApiDailySingleForecastTemp `json:"temp"`
    Code            WeatherApiStringValue               `json:"weather_code"`
    ObservationTime WeatherApiStringValue               `json:"observation_time"`
}

type WeatherApiDailySingleForecastTemp struct {
    Min WeatherApiFloatValue `json:"min"`
    Max WeatherApiFloatValue `json:"max"`
}

type WeatherApiFloatValue struct {
    Value float64 `json:"value"`
}

type WeatherApiStringValue struct {
    Value string `json:"value"`
}
