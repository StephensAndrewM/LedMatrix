package main

import (
    "image"
    "image/color"
    "time"
)

type TimeSlide struct {
    RedrawTicker *time.Ticker
}

func NewTimeSlide() *TimeSlide {
    sl := new(TimeSlide)
    return sl
}

func (this *TimeSlide) Initialize() {

}

func (this *TimeSlide) Terminate() {

}

func (this *TimeSlide) StartDraw(d Display) {
    this.RedrawTicker = DrawEverySecond(d, this.Draw)
}

func (this *TimeSlide) StopDraw() {
    this.RedrawTicker.Stop()
}

func (this *TimeSlide) IsEnabled() bool {
    return true // Always enabled
}

func (this *TimeSlide) Draw(img *image.RGBA) {
    t := time.Now()
    l1 := t.Format("Monday January 2")
    l2 := t.Format("3:04:05 PM")
    c1 := color.RGBA{255, 255, 255, 255}
    c2 := color.RGBA{255, 255, 0, 255}
    WriteString(img, l1, c1, ALIGN_CENTER, GetLeftOfCenterX(img), 7)
    WriteString(img, l2, c2, ALIGN_CENTER, GetLeftOfCenterX(img), 17)
}
