package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	"math"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
	log "github.com/sirupsen/logrus"
)

// How many days into the past to fetch and graph
// This should be one more than the number of days to graph
// since the graph looks at diffs between days.
var HISTORICAL_COVID_DAYS = 29

type CovidSlide struct {
	UsData DailyData
	MaData DailyData
	AzData DailyData

	FetchTicker           *time.Ticker
	LastFetchSuccessRatio float64
}

func NewCovidSlide() *CovidSlide {
	this := new(CovidSlide)
	this.UsData = NewDailyData("US")
	this.MaData = NewDailyData("Mass")
	this.AzData = NewDailyData("Ariz")
	return this
}

func (this *CovidSlide) Initialize() {
	// Query for new data once immediately
	this.FetchData()

	// Set up a period re-fetch of the data since it's sometimes late
	this.FetchTicker = time.NewTicker(4 * time.Hour)
	go func() {
		for range this.FetchTicker.C {
			this.FetchData()
		}
	}()
}

func (this *CovidSlide) Terminate() {
	this.FetchTicker.Stop()
}

func (this *CovidSlide) StartDraw(d Display) {
	DrawOnce(d, this.Draw)
}

func (this *CovidSlide) StopDraw() {

}

func (this *CovidSlide) IsEnabled() bool {
	// TODO disable this if we all survive
	return true
}

func (this *CovidSlide) Draw(img *image.RGBA) {
	// Stop immediately if we have too many errors
	if this.LastFetchSuccessRatio < 0.5 {
		DrawError(img, "Covid Cases", "Missing data.")
		return
	}

	red := color.RGBA{255, 0, 0, 255}
	WriteString(img, "COVID-19 CASES", red, ALIGN_CENTER, 63, 0)

	yellow := color.RGBA{255, 255, 0, 255}
	DrawDataRow(img, 8, this.UsData, yellow)
	DrawDataRow(img, 16, this.MaData, yellow)
	DrawDataRow(img, 24, this.AzData, yellow)
}

func (this *CovidSlide) FetchData() {
	attempted := 0
	successful := 0

	// Get data up to 15 days in the past
	for i := 1; i <= HISTORICAL_COVID_DAYS; i++ {
		d := civil.DateOf(time.Now().AddDate(0, 0, -i))
		// Check if fetch was successful based on data presence
		_, ok := this.UsData.Totals[d]
		// Refresh if data is 1 or 2 days old, since it might not be stable
		if !ok || i < 3 {
			attempted++
			if this.QueryForDate(d) {
				successful++
			}
		}
	}

	if attempted < successful {
		log.WithFields(log.Fields{
			"attempted":  attempted,
			"successful": successful,
		}).Debug("Some Covid queries failed.")
	}

	this.UsData = CalculateDiffs(this.UsData)
	this.MaData = CalculateDiffs(this.MaData)
	this.AzData = CalculateDiffs(this.AzData)

	this.LastFetchSuccessRatio = float64(successful) / float64(attempted)
}

// Can't use HttpHelper since the data doesn't change frequently
// and we need to do many queries to draw the slide.
func (this *CovidSlide) QueryForDate(d civil.Date) bool {
	url := fmt.Sprintf("https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_daily_reports/%02d-%02d-%04d.csv",
		d.Month, d.Day, d.Year)

	res, err := http.Get(url)
	if err != nil {
		log.WithFields(log.Fields{
			"url":   url,
			"error": err,
		}).Warn("Response error in Covid query.")
		return false
	}

	r := csv.NewReader(res.Body)
	rows, err := r.ReadAll()
	if err != nil {
		log.WithFields(log.Fields{
			"url":   url,
			"error": err,
		}).Warn("Parse error in Covid data.")
		return false
	}

	// Data is broken down by county so we need to aggregate many rows to
	// get the total for an individual state or the entire country.
	usSum := 0
	maSum := 0
	azSum := 0
	for i, row := range rows {
		// Skip header row
		if i == 0 {
			continue
		}

		n, err := strconv.Atoi(row[7])
		if err != nil {
			log.WithFields(log.Fields{
				"row":   row,
				"value": row[7],
				"error": err,
			}).Warn("Error reading cases number from CSV")
		}
		if row[3] == "US" {
			usSum += n
		}
		if row[2] == "Massachusetts" {
			maSum += n
		}
		if row[2] == "Arizona" {
			azSum += n
		}
	}

	if usSum > 0 {
		this.UsData.Totals[d] = usSum
	}
	if maSum > 0 {
		this.MaData.Totals[d] = maSum
	}
	if azSum > 0 {
		this.AzData.Totals[d] = azSum
	}
	return true
}

// Below are various helpers used in both Covid and Vaccination slides

// Limit to only displaying four glyphs using SI notation
func FormatNumber(n int) string {
	// Switch on number of digits in the input
	switch int(math.Log10(float64(n))) + 1 {
	case 4:
		return fmt.Sprintf("%.1fk", float64(n)/float64(1000))
	case 5, 6:
		return fmt.Sprintf("%.0fk", float64(n)/float64(1000))
	case 7:
		return fmt.Sprintf("%.1fM", float64(n)/float64(1000000))
	case 8, 9:
		return fmt.Sprintf("%.0fM", float64(n)/float64(1000000))
	default:
		return fmt.Sprintf("%d", n)
	}
}

// Container for a timeseries of daily-updated data
type DailyData struct {
	Label  string
	Total  int
	Totals map[civil.Date]int
	Diffs  map[civil.Date]int
}

func NewDailyData(label string) DailyData {
	var d DailyData
	d.Label = label
	d.Totals = make(map[civil.Date]int)
	d.Diffs = make(map[civil.Date]int)
	return d
}

func DrawDataRow(img *image.RGBA, y int, data DailyData, highlight color.RGBA) {
	white := color.RGBA{255, 255, 255, 255}
	gray := color.RGBA{128, 128, 128, 255}

	yesterday := civil.DateOf(time.Now().AddDate(0, 0, -1))

	WriteString(img, data.Label, white, ALIGN_LEFT, 1, y)

	if val, ok := data.Totals[yesterday]; ok && val > 0 {
		WriteString(img, FormatNumber(val), highlight, ALIGN_RIGHT, 62, y)
	} else {
		WriteString(img, "?", gray, ALIGN_RIGHT, 62, y)
	}

	if val, ok := data.Diffs[yesterday]; ok && val > 0 {
		WriteString(img, "+"+FormatNumber(val), highlight, ALIGN_RIGHT, 96, y)
	} else {
		WriteString(img, "+?", gray, ALIGN_RIGHT, 96, y)
	}

	DrawSemiAutoNormalizedGraph(img, 128-HISTORICAL_COVID_DAYS, y+6, 7, highlight, ToDiffsForGraph(data.Diffs))
}

func CalculateDiffs(data DailyData) DailyData {
	// Store the last nonzero value to gloss over data gaps.
	var lastVal int

	for i := -HISTORICAL_COVID_DAYS + 1; i < 0; i++ {
		// If there was a value for 1 day prior, use that as the last value.
		dA := civil.DateOf(time.Now().AddDate(0, 0, i-1))
		valA, okA := data.Totals[dA]
		if okA {
			lastVal = valA
		}

		dB := civil.DateOf(time.Now().AddDate(0, 0, i))
		valB, okB := data.Totals[dB]
		// Assuming values continuously increase, we can keep reassigning total.
		if okB {
			data.Total = valB
		}

		data.Diffs[dB] = 0
		if okB && lastVal > 0 {
			data.Diffs[dB] = valB - lastVal
		}
	}
	return data
}

func ToDiffsForGraph(diffsByDate map[civil.Date]int) []float64 {
	var diffs []float64
	for i := -HISTORICAL_COVID_DAYS + 1; i < 0; i++ {
		d := civil.DateOf(time.Now().AddDate(0, 0, i))
		val, ok := diffsByDate[d]
		if !ok {
			val = 0
		}
		diffs = append(diffs, float64(val))
	}
	return diffs
}
