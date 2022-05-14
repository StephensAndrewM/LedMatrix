package main

import (
	"image"
	"image/color"
	"strings"
	"time"
)

type TimeSlide struct {
	RedrawTicker *time.Ticker
}

func NewTimeSlide() *TimeSlide {
	sl := new(TimeSlide)
	return sl
}

func (sl *TimeSlide) Initialize() {

}

func (sl *TimeSlide) Terminate() {

}

func (sl *TimeSlide) StartDraw(d Display) {
	sl.RedrawTicker = DrawEverySecond(d, sl.Draw)
}

func (sl *TimeSlide) StopDraw() {
	sl.RedrawTicker.Stop()
}

func (sl *TimeSlide) IsEnabled() bool {
	return true // Always enabled
}

func (sl *TimeSlide) Draw(img *image.RGBA) {
	white := color.RGBA{255, 255, 255, 255}
	yellow := color.RGBA{255, 255, 0, 255}

	t := time.Now()
	d0 := strings.ToUpper(t.Format("Monday"))
	d1 := strings.ToUpper(t.Format("January 2"))
	t0 := t.Format("3:04 PM")

	WriteString(img, d0, white, ALIGN_CENTER, 32, 7)
	WriteString(img, d1, white, ALIGN_CENTER, 32, 17)

	WriteString(img, t0, yellow, ALIGN_CENTER, 96, 12)
}
