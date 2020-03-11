package main

import (
    "fmt"
    "image"
    "image/color"
    "math"
    "time"
)

type CoronavirusSlide struct {
}

func NewCoronavirusSlide() *CoronavirusSlide {
    sl := new(CoronavirusSlide)
    return sl
}

func (this *CoronavirusSlide) Initialize() {

}

func (this *CoronavirusSlide) Terminate() {

}

func (this *CoronavirusSlide) StartDraw(d Display) {
    DrawOnce(d, this.Draw)
}

func (this *CoronavirusSlide) StopDraw() {

}

func (this *CoronavirusSlide) IsEnabled() bool {
    // TODO disable this if we all survive
    return true
}

func (this *CoronavirusSlide) Draw(img *image.RGBA) {
    y := color.RGBA{255, 255, 0, 255}
    r := color.RGBA{255, 0, 0, 255}

    start := time.Date(2020, time.March, 10, 0, 0, 0, 0, time.Local)
    diff := int(math.Ceil(time.Since(start).Hours()/24.0)) - 1

    DrawEmptyBox(img, y, 54, 1, 20, 13)
    WriteString(img, fmt.Sprintf("%d", diff), y, ALIGN_CENTER, 64, 4)

    WriteString(img, "DAYS SINCE", r, ALIGN_CENTER, 64, 16)
    WriteString(img, "START OF QUARANTINE", r, ALIGN_CENTER, 64, 24)
}
