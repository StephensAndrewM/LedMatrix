package main

import (
    "image"
    "image/color"
    "time"
    "strings"
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
    c0 := color.RGBA{0, 255, 255, 255}
    c1 := color.RGBA{255, 255, 255, 255}
    c2 := color.RGBA{255, 255, 0, 255}
    
    t := time.Now()
    l0 := "WEEKDAY"
    if t.Weekday() == 0 || t.Weekday() == 6 {
    	l0 = "WEEKEND"
    	c0 = color.RGBA{0, 255, 0, 255}
    }
    l1 := strings.ToUpper(t.Format("Monday January 2"))
    l2 := t.Format("3:04 PM")
    
    WriteString(img, l0, c0, ALIGN_CENTER, GetLeftOfCenterX(img), 2)
    WriteString(img, l1, c1, ALIGN_CENTER, GetLeftOfCenterX(img), 14)
    WriteString(img, l2, c2, ALIGN_CENTER, GetLeftOfCenterX(img), 23)
}
