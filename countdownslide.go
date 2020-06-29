package main

import (
    "fmt"
    "image"
    "image/color"
    "math"
    "time"
)

type CountdownSlide struct {
}

func NewCountdownSlide() *CountdownSlide {
    sl := new(CountdownSlide)
    return sl
}

func (this *CountdownSlide) Initialize() {

}

func (this *CountdownSlide) Terminate() {

}

func (this *CountdownSlide) StartDraw(d Display) {
    DrawOnce(d, this.Draw)
}

func (this *CountdownSlide) StopDraw() {

}

func (this *CountdownSlide) IsEnabled() bool {
    return true
}

func (this *CountdownSlide) Draw(img *image.RGBA) {
    jobDate := time.Date(2020, time.July, 17, 0, 0, 0, 0, time.Local)
    this.DrawCountdownLine(img, jobDate, "LYUBA'S LAST DAY", color.RGBA{255, 0, 255, 255}, 6)

    moveDate := time.Date(2020, time.August, 15, 0, 0, 0, 0, time.Local)
    this.DrawCountdownLine(img, moveDate, "NEW APARTMENT", color.RGBA{0, 255, 0, 255}, 19)
}

func (this *CountdownSlide) DrawCountdownLine(img *image.RGBA, d time.Time, event string, c color.RGBA, y int) {
    numberColor := color.RGBA{255, 255, 255, 255}

    WriteString(img, fmt.Sprintf("%d", this.DaysUntil(d)), numberColor, ALIGN_RIGHT, 20, y)
    WriteString(img, event, c, ALIGN_LEFT, 26, y)
}

func (this *CountdownSlide) DaysUntil(d time.Time) int {
    diff := time.Until(d).Hours() / 24.0
    if diff < 0 {
        return 0
    }
    return int(math.Ceil(diff))
}
