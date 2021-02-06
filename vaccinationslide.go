package main

import (
    "bytes"
    "cloud.google.com/go/civil"
    "encoding/csv"
    log "github.com/sirupsen/logrus"
    "image"
    "image/color"
    "strconv"
    "time"
)

type VaccinationSlide struct {
    UsCount map[civil.Date]int
    MaCount map[civil.Date]int
    AzCount map[civil.Date]int

    HttpHelper *HttpHelper
}

func NewVaccinationSlide() *VaccinationSlide {
    this := new(VaccinationSlide)
    this.UsCount = make(map[civil.Date]int)
    this.MaCount = make(map[civil.Date]int)
    this.AzCount = make(map[civil.Date]int)

    this.HttpHelper = NewHttpHelper(HttpConfig{
        SlideId:         "VaccinationSlide",
        RefreshInterval: 6 * time.Hour,
        RequestUrl:      "https://github.com/owid/covid-19-data/raw/master/public/data/vaccinations/us_state_vaccinations.csv",
        ParseCallback:   this.Parse,
    })
    return this
}

func (this *VaccinationSlide) Initialize() {
    this.HttpHelper.StartLoop()
}

func (this *VaccinationSlide) Terminate() {
    this.HttpHelper.StopLoop()
}

func (this *VaccinationSlide) StartDraw(d Display) {
    DrawOnce(d, this.Draw)
}

func (this *VaccinationSlide) StopDraw() {

}

func (this *VaccinationSlide) IsEnabled() bool {
    return true
}

func (this *VaccinationSlide) Parse(respBytes []byte) bool {
    r := csv.NewReader(bytes.NewReader(respBytes))
    rows, err := r.ReadAll()
    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Warn("Parse error in vaccination data.")
        return false
    }

    // We won't draw data before this point
    minDrawDate := civil.DateOf(time.Now().AddDate(0, 0, -HISTORICAL_COVID_DAYS))

    for i, row := range rows {
        // Skip header row
        if i == 0 {
            continue
        }

        // Don't record data before the cutoff for graphing
        d, err := civil.ParseDate(row[0])
        if err != nil {
            log.WithFields(log.Fields{
                "row":   row,
                "value": row[0],
                "err":   err,
            }).Warn("Unparseable date in vaccination CSV")
            continue
        }
        if d.Before(minDrawDate) {
            continue
        }

        // Column 6 is people_vaccinated
        if row[6] == "" {
            continue
        }
        n, err := strconv.ParseFloat(row[6], 64)
        if err != nil {
            log.WithFields(log.Fields{
                "row":   row,
                "value": row[7],
                "error": err,
            }).Warn("Unparseable count in vaccination CSV")
        }
        count := int(n)

        if row[1] == "Massachusetts" {
            this.MaCount[d] = count
        }
        if row[1] == "Arizona" {
            this.AzCount[d] = count
        }
        if row[1] == "United States" {
            this.UsCount[d] = count
        }

    }
    return true
}

func (this *VaccinationSlide) Draw(img *image.RGBA) {
    aqua := color.RGBA{0, 255, 255, 255}
    WriteString(img, "COVID-19 VACCINATIONS", aqua, ALIGN_CENTER, 63, 0)

    green := color.RGBA{0, 255, 0, 255}
    DrawDataRow(img, 8, "US", this.UsCount, green)
    DrawDataRow(img, 16, "Mass", this.MaCount, green)
    DrawDataRow(img, 24, "Ariz", this.AzCount, green)
}
