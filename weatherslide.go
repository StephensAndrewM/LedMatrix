package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	log "github.com/sirupsen/logrus"

	"os"
	"strings"
	"time"
)

type WeatherSlide struct {
	Weather      WeatherData
	WeatherIcons map[string]*image.RGBA

	RealtimeHttpHelper       *HttpHelper
	ForecastHttpHelper       *HttpHelper
	HourlyForecastHttpHelper *HttpHelper
	RedrawTicker             *time.Ticker
}

type WeatherData struct {
	CurrentTemp float64
	CurrentIcon *image.RGBA

	Forecast1Weekday  time.Weekday
	Forecast1Icon     *image.RGBA
	Forecast1HighTemp int
	Forecast1LowTemp  int

	Forecast2Weekday  time.Weekday
	Forecast2Icon     *image.RGBA
	Forecast2HighTemp int
	Forecast2LowTemp  int
}

var weatherIconBaseDirFlag = flag.String("weather_icon_base_dir", "",
	"If specified, base directory to load weather icons from.")

// Latitude/longitude values for API requests
// Obtained using https://api.weather.gov/points/42.3643,-71.0854
const NWS_OFFICE = "BOX/69,76"
const NWS_STATION = "KBOS"

const WEATHER_MAX_HOURLY_DATA = 51

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

func NewWeatherSlide() *WeatherSlide {
	this := new(WeatherSlide)
	this.RealtimeHttpHelper = NewHttpHelper(HttpConfig{
		SlideId:         "WeatherSlide-Realtime",
		RefreshInterval: 5 * time.Minute,
		RequestUrl:      fmt.Sprintf("https://api.weather.gov/stations/%s/observations/latest", NWS_STATION),
		ParseCallback:   this.ParseRealtime,
	})
	this.ForecastHttpHelper = NewHttpHelper(HttpConfig{
		SlideId:         "WeatherSlide-Forecast",
		RefreshInterval: 30 * time.Minute,
		RequestUrl:      fmt.Sprintf("https://api.weather.gov/gridpoints/%s/forecast", NWS_OFFICE),
		ParseCallback:   this.ParseForecast,
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
	this.ForecastHttpHelper.StartLoop()
}

func (this *WeatherSlide) Terminate() {
	this.RealtimeHttpHelper.StopLoop()
	this.ForecastHttpHelper.StopLoop()
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

func (this *WeatherSlide) ParseRealtime(respBytes []byte) bool {
	var respData WeatherGovObservations
	err := json.Unmarshal(respBytes, &respData)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Could not interpret realtime weather JSON.")
		return false
	}

	t, err := time.Parse(time.RFC3339, respData.Properties.Timestamp)
	if err != nil || time.Now().Sub(t) > (6*time.Hour) {
		log.WithFields(log.Fields{
			"Timestamp": respData.Properties.Timestamp,
		}).Warn("Invalid last update time for observations.")
		return false
	}

	tempInCelsius := float64(respData.Properties.Temperature.Value)
	this.Weather.CurrentTemp = (tempInCelsius * (9 / 5)) + 32.0
	this.Weather.CurrentIcon = this.GetIcon(respData.Properties.Icon)
	return true
}

func (this *WeatherSlide) ParseForecast(respBytes []byte) bool {
	var respData WeatherGovForecast
	err := json.Unmarshal(respBytes, &respData)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Could not interpret daily weather JSON.")
		return false
	}

	t, err := time.Parse(time.RFC3339, respData.Properties.UpdateTime)
	if err != nil || time.Now().Sub(t) > (6*time.Hour) {
		log.WithFields(log.Fields{
			"UpdateTime": respData.Properties.UpdateTime,
		}).Warn("Invalid last update time for forecast.")
		return false
	}

	tz, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	// If after 3 PM, show nightly forecast
	if time.Now().Hour() < 15 {
		fTodayEndTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 18, 0, 0, 0, tz)
		fToday := this.GetForecastWithEndTime(fTodayEndTime, respData.Properties.Periods)
		if fToday == nil {
			log.WithFields(log.Fields{
				"fTodayEndTime": fTodayEndTime,
			}).Warn("Could not find forecast with expected end time.")
			return false
		}
		this.Weather.Forecast1HighTemp = fToday.Temperature

	} else {
		this.Weather.Forecast1HighTemp = 0

	}
	fTonightEndTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 6, 0, 0, 0, tz)
	fTonight := this.GetForecastWithEndTime(fTonightEndTime, respData.Properties.Periods)
	if fTonight == nil {
		log.WithFields(log.Fields{
			"fTonightEndTime": fTonightEndTime,
		}).Warn("Could not find forecast with expected end time.")
		return false
	}

	this.Weather.Forecast1Weekday = time.Now().Weekday()
	this.Weather.Forecast1LowTemp = fTonight.Temperature
	this.Weather.Forecast1Icon = this.GetIcon(fTonight.Icon)

	fTomorrowEndTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 18, 0, 0, 0, tz)
	fTomorrow := this.GetForecastWithEndTime(fTomorrowEndTime, respData.Properties.Periods)
	if fTomorrow == nil {
		log.WithFields(log.Fields{
			"fTomorrowEndTime": fTomorrowEndTime,
		}).Warn("Could not find forecast with expected end time.")
		return false
	}
	fTomorrowNightEndTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+2, 6, 0, 0, 0, tz)
	fTomorrowNight := this.GetForecastWithEndTime(fTomorrowNightEndTime, respData.Properties.Periods)
	if fTomorrowNight == nil {
		log.WithFields(log.Fields{
			"fTomorrowNightEndTime": fTomorrowNightEndTime,
		}).Warn("Could not find forecast with expected end time.")
		return false
	}

	this.Weather.Forecast2Weekday = time.Now().Add(time.Hour * 24).Weekday()
	this.Weather.Forecast2HighTemp = fTomorrow.Temperature
	this.Weather.Forecast2LowTemp = fTomorrowNight.Temperature
	this.Weather.Forecast2Icon = this.GetIcon(fTomorrow.Icon)

	return true
}

func (this *WeatherSlide) GetForecastWithEndTime(expectedEndTime time.Time, periods []WeatherGovForecastPeriod) *WeatherGovForecastPeriod {
	for _, period := range periods {
		t, _ := time.Parse(time.RFC3339, period.EndTime)
		if t.Equal(expectedEndTime) {
			return &period
		}
	}
	return nil
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
	if !this.RealtimeHttpHelper.LastFetchSuccess || !this.ForecastHttpHelper.LastFetchSuccess {
		DrawError(img, "Weather", "No data.")
		return
	}

	log.WithFields(log.Fields{
		"this": this.Weather,
	}).Debug("Drawing weather data.")

	this.DrawWeatherBox(img, 21, "NOW", fmt.Sprintf("%.1f°", this.Weather.CurrentTemp), this.Weather.CurrentIcon)

	forecast1Label := strings.ToUpper(this.Weather.Forecast1Weekday.String()[0:3])
	forecast1BottomText := fmt.Sprintf("%d°/%d°", this.Weather.Forecast1HighTemp, this.Weather.Forecast1LowTemp)
	// If high temp is zero, that means it wasn't set and we should only show nightly forecast.
	// Yes technically there's a bug where an actual zero-degree day wouldn't show up correctly.
	if this.Weather.Forecast1HighTemp == 0 {
		forecast1BottomText = fmt.Sprintf("%d°", this.Weather.Forecast1LowTemp)
	}
	this.DrawWeatherBox(img, 63, forecast1Label, forecast1BottomText, this.Weather.Forecast1Icon)

	forecast2Label := strings.ToUpper(this.Weather.Forecast2Weekday.String()[0:3])
	forecast2BottomText := fmt.Sprintf("%d°/%d°", this.Weather.Forecast2HighTemp, this.Weather.Forecast2LowTemp)
	this.DrawWeatherBox(img, 105, forecast2Label, forecast2BottomText, this.Weather.Forecast2Icon)
}

func (this *WeatherSlide) DrawWeatherBox(img *image.RGBA, centerX int, topText, bottomText string, icon *image.RGBA) {
	white := color.RGBA{255, 255, 255, 255}
	aqua := color.RGBA{0, 255, 255, 255}
	WriteString(img, topText, white, ALIGN_CENTER, centerX, 0)
	WriteString(img, bottomText, aqua, ALIGN_CENTER, centerX, 24)
	if icon != nil {
		DrawImageWithColorTransform(img, icon, centerX, 7, aqua)
	}
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

// Data structures used by api.weather.gov JSON feed
type WeatherGovObservations struct {
	Properties WeatherGovObservationsProperties
}

type WeatherGovObservationsProperties struct {
	Timestamp   string
	Icon        string
	Temperature WeatherGovObservationsTemperature
}

type WeatherGovObservationsTemperature struct {
	UnitCode string
	Value    float64
}

type WeatherGovForecast struct {
	Properties WeatherGovForecastProperties
}

type WeatherGovForecastProperties struct {
	UpdateTime string
	Periods    []WeatherGovForecastPeriod
}

type WeatherGovForecastPeriod struct {
	StartTime       string
	EndTime         string
	Temperature     int
	TemperatureUnit string
	Icon            string
}
