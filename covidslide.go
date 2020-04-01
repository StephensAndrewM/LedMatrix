package main

import (
    "cloud.google.com/go/civil"
    "fmt"
    log "github.com/sirupsen/logrus"
    "image"
    "image/color"
    "time"
    "encoding/csv"
    "net/http"
    "strconv"
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

func (this *CovidSlide) Initialize() {
    // We only re-query on initialization since this data only updates once a day
    d1 := civil.DateOf(time.Now().AddDate(0, 0, -1))
    d2 := civil.DateOf(time.Now().AddDate(0, 0, -2))

    if _, d1Ok := this.UsCases[d1]; !d1Ok {
        this.QueryForDate(d1)
    }
    if _, d2Ok := this.UsCases[d2]; !d2Ok {
        this.QueryForDate(d2)
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
    r := color.RGBA{255, 0, 0, 255}
    y := color.RGBA{255, 255, 0, 255}
    w := color.RGBA{255, 255, 255, 255}

    WriteString(img, "COVID-19 CASES", r, ALIGN_LEFT, 0, 0)
    WriteString(img, "#", r, ALIGN_RIGHT, 100, 0)
    WriteString(img, "%Δ", r, ALIGN_RIGHT, 127, 0)

    d1 := civil.DateOf(time.Now().AddDate(0, 0, -1))
    d2 := civil.DateOf(time.Now().AddDate(0, 0, -2))

    WriteString(img, "United States", w, ALIGN_LEFT, 0, 12)
    if n1, ok := this.UsCases[d1]; ok {
        WriteString(img, this.Format(n1), y, ALIGN_RIGHT, 100, 12)
        if n2, ok := this.UsCases[d2]; ok {
            WriteString(img, this.Diff(n1, n2), y, ALIGN_RIGHT, 127, 12)
        }
    }

    WriteString(img, "Massachusetts", w, ALIGN_LEFT, 0, 20)
    if n1, ok := this.MaCases[d1]; ok {
        WriteString(img, this.Format(n1), y, ALIGN_RIGHT, 100, 20)
        if n2, ok := this.MaCases[d2]; ok {
            WriteString(img, this.Diff(n1, n2), y, ALIGN_RIGHT, 127, 20)
        }
    }
}

// Limit to only displaying four glyphs max
func (this *CovidSlide) Format(n int) string {
    if n >= 10000000 {
        return fmt.Sprintf("%.0fM", float64(n)/float64(1000000))
    } else if n >= 1000000 {
        return fmt.Sprintf("%.1fM", float64(n)/float64(1000000))
    } else if n > 100000 {
        return fmt.Sprintf("%.0fk", float64(n)/float64(1000))
    } else if n > 10000 {
        return fmt.Sprintf("%.1fk", float64(n)/float64(1000))
    }
    return fmt.Sprintf("%d", n)
}

func (this *CovidSlide) Diff(n1, n2 int) string {
    n1f := float64(n1)
    n2f := float64(n2)
    diff := ((n1f - n2f) / n2f) * 100
    if diff > 0 {
        return fmt.Sprintf("+%0.1f", diff)
    } else {
        return fmt.Sprintf("-%0.1f", diff)
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
    for i,row := range rows {
        // Skip header row
        if i == 0 { continue; }

        n, err := strconv.Atoi(row[7])
        if err != nil {
            log.WithFields(log. Fields{
                "row": row,
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
