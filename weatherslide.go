package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type WeatherSlide struct {
	Weather      WeatherData
	WeatherIcons map[string]*image.RGBA

	ObservationsHttpHelper   *HttpHelper
	ForecastHttpHelper       *HttpHelper
	HourlyForecastHttpHelper *HttpHelper
	RedrawTicker             *time.Ticker
}

type WeatherData struct {
	CurrentTemp int
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

// Set of possible values provided by https://api.weather.gov/icons
var WEATHER_API_ICON_MAP = map[string]string{
	"day/skc":         "sun",             // Fair/clear
	"night/skc":       "moon",            // Fair/clear
	"day/few":         "cloud_sun",       // A few clouds
	"night/few":       "cloud_moon",      // A few clouds
	"day/sct":         "cloud_sun",       // Partly cloudy
	"night/sct":       "cloud_moon",      // Partly cloudy
	"bkn":             "clouds",          // Mostly cloudy
	"ovc":             "clouds",          // Overcast
	"day/wind_skc":    "sun",             // Fair/clear and windy
	"night/wind_skc":  "moon",            // Fair/clear and windy
	"day/wind_few":    "cloud_wind_sun",  // A few clouds and windy
	"night/wind_few":  "cloud_wind_moon", // A few clouds and windy
	"day/wind_sct":    "cloud_wind_sun",  // Partly cloudy and windy
	"night/wind_sct":  "cloud_wind_moon", // Partly cloudy and windy
	"wind_bkn":        "cloud_wind",      // Mostly cloudy and windy
	"wind_ovc":        "cloud_wind",      // Overcast and windy
	"snow":            "snow",            // Snow
	"rain_snow":       "rain_snow",       // Rain/snow
	"rain_sleet":      "rain_snow",       // Rain/sleet
	"snow_sleet":      "rain_snow",       // Snow/sleet
	"fzra":            "rain1",           // Freezing rain
	"rain_fzra":       "rain1",           // Rain/freezing rain
	"snow_fzra":       "rain_snow",       // Freezing rain/snow
	"sleet":           "rain1",           // Sleet
	"rain":            "rain1",           // Rain
	"rain_showers":    "rain0",           // Rain showers (high cloud cover)
	"rain_showers_hi": "rain0",           // Rain showers (low cloud cover)
	"tsra":            "lightning",       // Thunderstorm (high cloud cover)
	"tsra_sct":        "lightning",       // Thunderstorm (medium cloud cover)
	"tsra_hi":         "lightning",       // Thunderstorm (low cloud cover)
	"blizzard":        "snow",            // Blizzard
	"fog":             "cloud",           // Fog/mist
	// "tornado":         "",                // Tornado
	// "hurricane":       "",                // Hurricane conditions
	// "tropical_storm":  "",                // Tropical storm conditions
	// "dust":            "",                // Dust
	// "smoke":           "",                // Smoke
	// "haze":            "",                // Haze
	// "hot":             "",                // Hot
	// "cold":            "",                // Cold
}

func NewWeatherSlide() *WeatherSlide {
	this := new(WeatherSlide)
	this.ObservationsHttpHelper = NewHttpHelper(HttpConfig{
		SlideId:            "WeatherSlide-Observations",
		RefreshInterval:    5 * time.Minute,
		RequestUrlCallback: this.BuildObservationsUrl,
		ParseCallback:      this.ParseObservations,
	})
	this.ForecastHttpHelper = NewHttpHelper(HttpConfig{
		SlideId:            "WeatherSlide-Forecast",
		RefreshInterval:    30 * time.Minute,
		RequestUrlCallback: this.BuildForecastUrl,
		ParseCallback:      this.ParseForecast,
	})
	return this
}

func (this *WeatherSlide) Initialize() {
	this.ObservationsHttpHelper.StartLoop()
	this.ForecastHttpHelper.StartLoop()
}

func (this *WeatherSlide) Terminate() {
	this.ObservationsHttpHelper.StopLoop()
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

func (this *WeatherSlide) BuildObservationsUrl() (*http.Request, error) {
	return this.BuildUrl(fmt.Sprintf("https://api.weather.gov/stations/%s/observations/latest", NWS_STATION))
}

func (this *WeatherSlide) BuildForecastUrl() (*http.Request, error) {
	return this.BuildUrl(fmt.Sprintf("https://api.weather.gov/gridpoints/%s/forecast", NWS_OFFICE))
}

func (this *WeatherSlide) BuildUrl(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Add required headers to outgoing requests.
	req.Header.Set("User-Agent", "https://github.com/stephensandrewm/LedMatrix")
	req.Header.Set("Accept", "application/ld+json")
	return req, nil
}

func (this *WeatherSlide) ParseObservations(respBytes []byte) bool {
	var respData WeatherGovObservations
	err := json.Unmarshal(respBytes, &respData)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Could not interpret observations weather JSON.")
		return false
	}

	t, err := time.Parse(time.RFC3339, respData.Timestamp)
	if err != nil || time.Now().Sub(t) > (6*time.Hour) {
		log.WithFields(log.Fields{
			"Timestamp": respData.Timestamp,
		}).Warn("Invalid last update time for observations.")
		return false
	}

	tempInCelsius := float64(respData.Temperature.Value)
	this.Weather.CurrentTemp = int((tempInCelsius * (9 / 5)) + 32.0)
	this.Weather.CurrentIcon = this.GetIcon(respData.Icon)
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

	t, err := time.Parse(time.RFC3339, respData.UpdateTime)
	if err != nil || time.Now().Sub(t) > (6*time.Hour) {
		log.WithFields(log.Fields{
			"UpdateTime": respData.UpdateTime,
		}).Warn("Invalid last update time for forecast.")
		return false
	}

	tz, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	// If after 6 PM, show nightly forecast
	if time.Now().Hour() < 18 {
		fTodayEndTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 18, 0, 0, 0, tz)
		fToday := this.GetForecastWithEndTime(fTodayEndTime, respData.Periods)
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
	fTonight := this.GetForecastWithEndTime(fTonightEndTime, respData.Periods)
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
	fTomorrow := this.GetForecastWithEndTime(fTomorrowEndTime, respData.Periods)
	if fTomorrow == nil {
		log.WithFields(log.Fields{
			"fTomorrowEndTime": fTomorrowEndTime,
		}).Warn("Could not find forecast with expected end time.")
		return false
	}
	fTomorrowNightEndTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+2, 6, 0, 0, 0, tz)
	fTomorrowNight := this.GetForecastWithEndTime(fTomorrowNightEndTime, respData.Periods)
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

func (this *WeatherSlide) GetIcon(url string) *image.RGBA {
	r := regexp.MustCompile("\\/icons\\/land\\/([^\\/]+\\/([a-z_]+))")
	m := r.FindStringSubmatch(url)
	if len(m) < 3 || m[1] == "" || m[2] == "" {
		log.WithFields(log.Fields{
			"url": url,
		}).Warn("Could not extract condition from icon URL.")
		return nil
	}

	// Icon could be defined using one of two patterns. Find which one.
	conditionWithTimeOfDay := m[1]
	condition := m[2]
	icon, ok := WEATHER_API_ICON_MAP[conditionWithTimeOfDay]
	if !ok {
		icon, ok = WEATHER_API_ICON_MAP[condition]
		if !ok {
			log.WithFields(log.Fields{
				"url":                    url,
				"condition":              condition,
				"conditionWithTimeOfDay": conditionWithTimeOfDay,
			}).Warn("Conditions did not map to a known weather icon.")
			return nil
		}
	}

	// Once the icon key is found, base64-decode it and store it.
	base64Icon, ok := weatherIcons[icon]
	if !ok {
		log.WithFields(log.Fields{
			"icon": icon,
		}).Warn("Could not find weather icon.")
	}
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Icon))
	img, err := png.Decode(decoder)
	if err != nil {
		log.WithFields(log.Fields{
			"icon": icon,
			"err":  err,
		}).Warn("Could not decode weather icon.")
		return nil
	}

	// Then convert from image.Image to image.RGBA
	b := img.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)
	return rgba

}

func (this *WeatherSlide) Draw(img *image.RGBA) {
	// Stop immediately if we have errors
	if !this.ObservationsHttpHelper.LastFetchSuccess || !this.ForecastHttpHelper.LastFetchSuccess {
		DrawError(img, "Weather", "No data.")
		return
	}

	yellow := color.RGBA{255, 255, 0, 255}
	aqua := color.RGBA{0, 255, 255, 255}

	this.DrawWeatherBox(img, 21, "NOW", fmt.Sprintf("%d°", this.Weather.CurrentTemp), yellow, this.Weather.CurrentIcon)

	forecast1Label := strings.ToUpper(this.Weather.Forecast1Weekday.String()[0:3])
	forecast1BottomText := fmt.Sprintf("%d°/%d°", this.Weather.Forecast1HighTemp, this.Weather.Forecast1LowTemp)
	// If high temp is zero, that means it wasn't set and we should only show nightly forecast.
	// Yes technically there's a bug where an actual zero-degree day wouldn't show up correctly.
	if this.Weather.Forecast1HighTemp == 0 {
		forecast1BottomText = fmt.Sprintf("%d°", this.Weather.Forecast1LowTemp)
	}
	this.DrawWeatherBox(img, 63, forecast1Label, forecast1BottomText, aqua, this.Weather.Forecast1Icon)

	forecast2Label := strings.ToUpper(this.Weather.Forecast2Weekday.String()[0:3])
	forecast2BottomText := fmt.Sprintf("%d°/%d°", this.Weather.Forecast2HighTemp, this.Weather.Forecast2LowTemp)
	this.DrawWeatherBox(img, 105, forecast2Label, forecast2BottomText, aqua, this.Weather.Forecast2Icon)
}

func (this *WeatherSlide) DrawWeatherBox(img *image.RGBA, centerX int, dateText, temperatureText string, dateColor color.RGBA, icon *image.RGBA) {
	white := color.RGBA{255, 255, 255, 255}
	WriteString(img, temperatureText, white, ALIGN_CENTER, centerX, 0)
	if icon != nil {
		DrawImageWithColorTransform(img, icon, centerX-8, 7, white)
	}
	WriteString(img, dateText, dateColor, ALIGN_CENTER, centerX, 24)
}

// Data structures used by api.weather.gov JSON feed
type WeatherGovObservations struct {
	Timestamp   string
	Icon        string
	Temperature WeatherGovObservationsTemperature
}

type WeatherGovObservationsTemperature struct {
	UnitCode string
	Value    float64
}

type WeatherGovForecast struct {
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
