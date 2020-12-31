package main

import (
    "image"
    "image/color"
    "time"
    "fmt"
)

type NewYearSlide struct {
    Midnight     time.Time
    RedrawTicker *time.Ticker
}

func NewNewYearSlide() *NewYearSlide {
    sl := new(NewYearSlide)
    return sl
}

func (this *NewYearSlide) Initialize() {
    t := time.Now()
    year := t.Year() + 1
    // If it's January, the new year just passed so we want to count to the
    // current year instead (and show zeros).
    if t.Month() == time.January {
        year = t.Year()
    }
    this.Midnight = time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
}

func (this *NewYearSlide) Terminate() {

}

func (this *NewYearSlide) StartDraw(d Display) {
    this.RedrawTicker = DrawEverySecond(d, this.Draw)
}

func (this *NewYearSlide) StopDraw() {
    this.RedrawTicker.Stop()
}

func (this *NewYearSlide) IsEnabled() bool {
    diff := time.Until(this.Midnight)
    return diff > (-1 * time.Hour)
}

func (this *NewYearSlide) Draw(img *image.RGBA) {
    c0 := color.RGBA{0, 255, 255, 255}
    c1 := color.RGBA{255, 255, 255, 255}
    c2 := color.RGBA{255, 255, 0, 255}

    diff := time.Until(this.Midnight)
    if diff < 0 {
        diff = 0
    }

    WriteString(img, this.fmtDuration(diff), c0, ALIGN_CENTER, GetLeftOfCenterX(img), 4)
    WriteString(img, "UNTIL", c1, ALIGN_CENTER, GetLeftOfCenterX(img), 14)
    WriteString(img, fmt.Sprintf("%d", this.Midnight.Year()), c2, ALIGN_CENTER, GetLeftOfCenterX(img)-1, 23)
}

func (this *NewYearSlide) fmtDuration(d time.Duration) string {
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    d -= m * time.Minute
    s := d / time.Second
    return fmt.Sprintf("%02d : %02d : %02d", h, m, s)
}
