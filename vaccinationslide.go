package main

import (
	"bytes"
	"encoding/csv"
	"image"
	"image/color"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	log "github.com/sirupsen/logrus"
)

type VaccinationSlide struct {
	UsData DailyData
	MaData DailyData
	AzData DailyData

	HttpHelper *HttpHelper
}

func NewVaccinationSlide() *VaccinationSlide {
	sl := new(VaccinationSlide)
	sl.UsData = NewDailyData("US")
	sl.MaData = NewDailyData("Mass")
	sl.AzData = NewDailyData("Ariz")

	sl.HttpHelper = NewHttpHelper(HttpConfig{
		SlideId:         "VaccinationSlide",
		RefreshInterval: 6 * time.Hour,
		RequestUrl:      "https://github.com/owid/covid-19-data/raw/master/public/data/vaccinations/us_state_vaccinations.csv",
		ParseCallback:   sl.Parse,
	})
	return sl
}

func (sl *VaccinationSlide) Initialize() {
	sl.HttpHelper.StartLoop()
}

func (sl *VaccinationSlide) Terminate() {
	sl.HttpHelper.StopLoop()
}

func (sl *VaccinationSlide) StartDraw(d Display) {
	DrawOnce(d, sl.Draw)
}

func (sl *VaccinationSlide) StopDraw() {

}

func (sl *VaccinationSlide) IsEnabled() bool {
	return true
}

func (sl *VaccinationSlide) Parse(respBytes []byte) bool {
	r := csv.NewReader(bytes.NewReader(respBytes))
	rows, err := r.ReadAll()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Parse error in vaccination data.")
		return false
	}

	// We won't draw data before sl.point
	minDrawDate := civil.DateOf(time.Now().AddDate(0, 0, -HISTORICAL_COVID_DAYS))

	dateCol := -1
	peopleVaccinatedCol := -1
	for i, row := range rows {
		// If header row, find the column IDs for the data we need
		if i == 0 {
			for c, h := range row {
				header := strings.ToUpper(h)
				if header == "DATE" {
					dateCol = c
				}
				if header == "PEOPLE_VACCINATED" {
					peopleVaccinatedCol = c
				}
			}

			if dateCol == -1 || peopleVaccinatedCol == -1 {
				log.WithFields(log.Fields{
					"dateCol":             dateCol,
					"peopleVaccinatedCol": peopleVaccinatedCol,
					"row":                 row,
				}).Warn("Could not find required columns in vaccination CSV.")
				return false
			}
			continue
		}

		d, err := civil.ParseDate(row[dateCol])
		if err != nil {
			log.WithFields(log.Fields{
				"row":   row,
				"value": row[0],
				"err":   err,
			}).Warn("Unparseable date in vaccination CSV.")
			continue
		}
		// Don't save data before the cutoff for graphing
		if d.Before(minDrawDate) {
			continue
		}

		// Treating empty rows as 0 causes problems with diffs so instead we skip them
		if row[peopleVaccinatedCol] == "" {
			continue
		}
		n, err := strconv.ParseFloat(row[peopleVaccinatedCol], 64)
		if err != nil {
			log.WithFields(log.Fields{
				"row":                 row,
				"peopleVaccinatedCol": peopleVaccinatedCol,
				"value":               row[peopleVaccinatedCol],
				"error":               err,
			}).Warn("Unparseable count in vaccination CSV.")
		}
		count := int(n)

		if row[1] == "Massachusetts" {
			sl.MaData.Totals[d] = count
		}
		if row[1] == "Arizona" {
			sl.AzData.Totals[d] = count
		}
		if row[1] == "United States" {
			sl.UsData.Totals[d] = count
		}
	}

	sl.UsData = CalculateDiffs(sl.UsData)
	sl.MaData = CalculateDiffs(sl.MaData)
	sl.AzData = CalculateDiffs(sl.AzData)

	return true
}

func (sl *VaccinationSlide) Draw(img *image.RGBA) {
	// Stop immediately if we have errors
	if !sl.HttpHelper.LastFetchSuccess {
		DrawError(img, "Covid Vaccination", "Missing data.")
		return
	}

	green := color.RGBA{0, 255, 0, 255}
	WriteString(img, "COVID-19 VACCINATIONS", green, ALIGN_CENTER, 63, 0)

	yellow := color.RGBA{255, 255, 0, 255}
	DrawDataRow(img, 8, sl.UsData, yellow)
	DrawDataRow(img, 16, sl.MaData, yellow)
	DrawDataRow(img, 24, sl.AzData, yellow)
}
