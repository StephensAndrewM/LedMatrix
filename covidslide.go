package main

import (
    "cloud.google.com/go/civil"
    "encoding/csv"
    "fmt"
    log "github.com/sirupsen/logrus"
    "image"
    "image/color"
    "net/http"
    "strconv"
    "time"
    "math"
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

    WriteString(img, "COVID-19 CASES", r, ALIGN_CENTER, 63, 2)

    d1 := civil.DateOf(time.Now().AddDate(0, 0, -1))
    d2 := civil.DateOf(time.Now().AddDate(0, 0, -2))

    WriteString(img, "United States", w, ALIGN_LEFT, 1, 13)
    if n1, ok := this.UsCases[d1]; ok {
        WriteString(img, this.Format(n1), y, ALIGN_RIGHT, 96, 13)
        if n2, ok := this.UsCases[d2]; ok {
            WriteString(img, this.Diff(n1, n2), y, ALIGN_RIGHT, 126, 13)
        }
    }

    WriteString(img, "Massachusetts", w, ALIGN_LEFT, 1, 22)
    if n1, ok := this.MaCases[d1]; ok {
        WriteString(img, this.Format(n1), y, ALIGN_RIGHT, 96, 22)
        if n2, ok := this.MaCases[d2]; ok {
            WriteString(img, this.Diff(n1, n2), y, ALIGN_RIGHT, 126, 22)
        }
    }
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
