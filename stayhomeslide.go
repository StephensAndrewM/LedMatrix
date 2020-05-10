package main

import (
    "fmt"
    "image"
    "image/color"
    "math"
    "time"
)

type StayHomeSlide struct {
}

func NewStayHomeSlide() *StayHomeSlide {
    sl := new(StayHomeSlide)
    return sl
}

func (this *StayHomeSlide) Initialize() {

}

func (this *StayHomeSlide) Terminate() {

}

func (this *StayHomeSlide) StartDraw(d Display) {
    DrawOnce(d, this.Draw)
}

func (this *StayHomeSlide) StopDraw() {

}

func (this *StayHomeSlide) IsEnabled() bool {
    // TODO disable this if we all survive
    return true
}

func (this *StayHomeSlide) Draw(img *image.RGBA) {
    y := color.RGBA{255, 255, 0, 255}
    r := color.RGBA{255, 0, 0, 255}

    start := time.Date(2020, time.March, 10, 0, 0, 0, 0, time.Local)
    diff := int(math.Ceil(time.Since(start).Hours()/24.0)) - 1

    DrawIcon(img, "house-16", r, 8, 2)
    DrawIcon(img, "house-16", r, (128-8-16), 2)

    DrawEmptyBox(img, y, 54, 1, 20, 13)
    WriteString(img, fmt.Sprintf("%d", diff), y, ALIGN_CENTER, 64, 4)

    WriteString(img, "DAYS SINCE", r, ALIGN_CENTER, 64, 16)
    WriteString(img, "OFFICES CLOSED", r, ALIGN_CENTER, 64, 24)
}
