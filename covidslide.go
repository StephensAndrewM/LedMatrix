package main

import (
    "cloud.google.com/go/civil"
    "encoding/csv"
    "fmt"
    log "github.com/sirupsen/logrus"
    "image"
    "image/color"
    "math"
    "net/http"
    "strconv"
    "time"
)

type CovidSlide struct {
    // Store historically retrieved data
    UsCases map[civil.Date]int
    MaCases map[civil.Date]int
}

type DatedCases struct {
    Date  civil.Date
    Cases int32
}

func NewCovidSlide() *CovidSlide {
    this := new(CovidSlide)
    this.UsCases = make(map[civil.Date]int)
    this.MaCases = make(map[civil.Date]int)
    return this
}

// We only re-query on initialization since this data only updates once a day
func (this *CovidSlide) Initialize() {
    // Get data up to 15 days in the past
    for i := -1; i >= -15; i-- {
        d := civil.DateOf(time.Now().AddDate(0, 0, i))
        if _, dOk := this.UsCases[d]; !dOk {
            this.QueryForDate(d)
        }
    }
}

func (this *CovidSlide) Terminate() {

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
    red := color.RGBA{255, 0, 0, 255}

    WriteString(img, "COVID-19 CASES", red, ALIGN_CENTER, 63, 2)

    this.DrawForLocation(img, 13, "United States", this.UsCases)
    this.DrawForLocation(img, 22, "Massachusetts", this.MaCases)
}

func (this *CovidSlide) DrawForLocation(img *image.RGBA, y int, label string, cases map[civil.Date]int) {
    yellow := color.RGBA{255, 255, 0, 255}
    white := color.RGBA{255, 255, 255, 255}

    d1 := civil.DateOf(time.Now().AddDate(0, 0, -1))
    d2 := civil.DateOf(time.Now().AddDate(0, 0, -2))

    WriteString(img, label, white, ALIGN_LEFT, 1, y)
    if n1, ok := cases[d1]; ok {
        WriteString(img, this.Format(n1), yellow, ALIGN_RIGHT, 88, y)
        if n2, ok := cases[d2]; ok {
            WriteString(img, this.Diff(n1, n2), yellow, ALIGN_RIGHT, 110, y)
        }
    }

    var diffValues []float64
    var diffMax float64
    for i := -15; i < -1; i++ {
        dA := civil.DateOf(time.Now().AddDate(0, 0, i))
        dB := civil.DateOf(time.Now().AddDate(0, 0, i+1))
        nA, okA := cases[dA]
        nB, okB := cases[dB]
        diff := float64(nB - nA)
        if okA && okB {
            diffValues = append(diffValues, diff)
        }
        if diffMax < diff {
            diffMax = diff
        }
    }

    DrawNormalizedGraph(img, 113, y+6, 7, 0, diffMax, yellow, diffValues)
}

// Limit to only displaying four glyphs max
func (this *CovidSlide) Format(n int) string {
    // Switch on number of digits in the number
    switch int(math.Log10(float64(n))) + 1 {
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

// Display the difference with a +/- symbol
func (this *CovidSlide) Diff(n1, n2 int) string {
    diff := this.Format(n1 - n2)
    if (n1 - n2) > 0 {
        return fmt.Sprintf("+%s", diff)
    } else {
        return fmt.Sprintf("-%s", diff)
    }
}

// Can't use HttpHelper since the data doesn't change frequently
// and we need to do queries to draw the slide.
func (this *CovidSlide) QueryForDate(d civil.Date) {
    url := fmt.Sprintf("https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_daily_reports/%02d-%02d-%04d.csv",
        d.Month, d.Day, d.Year)

    res, err := http.Get(url)
    if err != nil {
        log.WithFields(log.Fields{
            "url":   url,
            "error": err,
        }).Warn("Response error in Covid query.")
        return
    }

    r := csv.NewReader(res.Body)
    rows, err := r.ReadAll()
    if err != nil {
        log.WithFields(log.Fields{
            "url":   url,
            "error": err,
        }).Warn("Parse error in Covid data.")
        return
    }

    maSum := 0
    usSum := 0
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
    }

    this.UsCases[d] = usSum
    this.MaCases[d] = maSum
}
