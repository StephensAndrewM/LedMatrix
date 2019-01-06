package main

import (
    "strings"
    "time"
    "image"
    "image/color"
)

type TimeSlide struct {
}

func NewTimeSlide() *TimeSlide {
    sl := new(TimeSlide)
    return sl
}

func (this *TimeSlide) Draw(img *image.RGBA) {
    t := time.Now()
    l1 := strings.ToUpper(t.Format("Jan 2"))
    l2 := t.Format("3:04:05 PM")
    c1 := color.RGBA{255, 255, 255, 255}
    c2 := color.RGBA{255, 255, 0, 255}
    WriteString(img, l1, c1, ALIGN_CENTER, GetLeftOfCenterX(img), 8)
    WriteString(img, l2, c2, ALIGN_CENTER, GetLeftOfCenterX(img), 16)
}
