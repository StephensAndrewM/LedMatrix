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

// How many days into the past to fetch and graph
// This should be one more than the number of days to graph
// since the graph looks at diffs between days.
var HISTORICAL_COVID_DAYS = 29

type CovidSlide struct {
    FetchTicker *time.Ticker
    // Store historically retrieved data
    UsCases map[civil.Date]int
    MaCases map[civil.Date]int
    AzCases map[civil.Date]int
}

type DatedCases struct {
    Date  civil.Date
    Cases int32
}

func NewCovidSlide() *CovidSlide {
    this := new(CovidSlide)
    this.UsCases = make(map[civil.Date]int)
    this.MaCases = make(map[civil.Date]int)
    this.AzCases = make(map[civil.Date]int)
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
    red := color.RGBA{255, 0, 0, 255}

    WriteString(img, "COVID-19 CASES", red, ALIGN_CENTER, 63, 0)

    this.DrawForLocation(img, 8, "US", this.UsCases)
    this.DrawForLocation(img, 16, "Mass", this.MaCases)
    this.DrawForLocation(img, 24, "Ariz", this.AzCases)
}

func (this *CovidSlide) DrawForLocation(img *image.RGBA, y int, label string, cases map[civil.Date]int) {
    yellow := color.RGBA{255, 255, 0, 255}
    white := color.RGBA{255, 255, 255, 255}
    gray := color.RGBA{128, 128, 128, 255}

    d1 := civil.DateOf(time.Now().AddDate(0, 0, -1))
    d2 := civil.DateOf(time.Now().AddDate(0, 0, -2))

    WriteString(img, label, white, ALIGN_LEFT, 1, y)
    if n1, ok := cases[d1]; ok && n1 > 0 {
        // First display total number of cases
        WriteString(img, this.Format(n1), yellow, ALIGN_RIGHT, 63, y)

        // Then calculate and display the diff
        if n2, ok := cases[d2]; ok && (n1-n2) > 0 {
            WriteString(img, "+"+this.Format(n1-n2), yellow, ALIGN_RIGHT, 94, y)
        } else {
            WriteString(img, "?", gray, ALIGN_RIGHT, 94, y)
        }
    } else {
        WriteString(img, "?", gray, ALIGN_RIGHT, 63, y)
    }

    var diffValues []float64
    var diffMax float64
    for i := -HISTORICAL_COVID_DAYS; i < -1; i++ {
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

    DrawNormalizedGraph(img, 128-HISTORICAL_COVID_DAYS, y+6, 7, 0, diffMax, yellow, diffValues)
}

// Limit to only displaying four glyphs max
func (this *CovidSlide) Format(n int) string {
    // Switch on number of digits in the number
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

func (this *CovidSlide) FetchData() {
    attempted := 0
    successful := 0

    // Get data up to 15 days in the past
    for i := 1; i <= HISTORICAL_COVID_DAYS; i++ {
        d := civil.DateOf(time.Now().AddDate(0, 0, -i))
        // Check if fetch was successful based on data presence
        _, dOk := this.UsCases[d]
        // Refresh if data is 1 or 2 days old, since it might not be stable
        if !dOk || i < 3 {
            attempted++
            if this.QueryForDate(d) {
                successful++
            }
        }
    }

    log.WithFields(log.Fields{
        "attempted":  attempted,
        "successful": successful,
    }).Info("Fetched latest Covid data.")
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

    // We know that there are some cases so a 0 value indicates an error
    if usSum > 0 {
        this.UsCases[d] = usSum
    }
    if maSum > 0 {
        this.MaCases[d] = maSum
    }
    if azSum > 0 {
        this.AzCases[d] = azSum
    }
    return true
}
