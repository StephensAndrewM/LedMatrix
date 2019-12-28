package main

import (
    "fmt"
    "image"
    "image/color"
    "math"
    "time"
)

type ChristmasSlide struct {
    PtoDate  time.Time
    XmasDate time.Time
}

func NewChristmasSlide() *ChristmasSlide {
    sl := new(ChristmasSlide)
    return sl
}

func (this *ChristmasSlide) Initialize() {
    t := time.Now()
    this.PtoDate = time.Date(t.Year(), time.December, 21, 0, 0, 0, 0, time.Local)
    this.XmasDate = time.Date(t.Year(), time.December, 25, 0, 0, 0, 0, time.Local)
}

func (this *ChristmasSlide) Terminate() {

}

func (this *ChristmasSlide) StartDraw(d Display) {
    DrawOnce(d, this.Draw)
}

func (this *ChristmasSlide) StopDraw() {

}

func (this *ChristmasSlide) IsEnabled() bool {
    return this.DaysUntil(this.XmasDate) >= 0 && this.DaysUntil(this.XmasDate) > 30
}

func (this *ChristmasSlide) Draw(img *image.RGBA) {
    r := color.RGBA{255, 0, 0, 255}
    g := color.RGBA{0, 255, 0, 255}

    DrawEmptyBox(img, r, 23, 1, 18, 13)
    WriteString(img, fmt.Sprintf("%d", this.DaysUntil(this.PtoDate)), r, ALIGN_CENTER, 32, 4)
    WriteString(img, "DAYS UNTIL", g, ALIGN_CENTER, 32, 16)
    WriteString(img, "VACATION", g, ALIGN_CENTER, 32, 24)

    DrawEmptyBox(img, r, 87, 1, 18, 13)
    WriteString(img, fmt.Sprintf("%d", this.DaysUntil(this.XmasDate)), r, ALIGN_CENTER, 96, 4)
    WriteString(img, "DAYS UNTIL", g, ALIGN_CENTER, 96, 16)
    WriteString(img, "CHRISTMAS", g, ALIGN_CENTER, 96, 24)
}

func (this *ChristmasSlide) DaysUntil(d time.Time) int {
    diff := time.Until(d).Hours() / 24.0
    if diff < 0 {
        return 0
    }
    return int(math.Ceil(diff))
}
