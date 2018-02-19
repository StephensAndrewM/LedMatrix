package main

import (
    "strings"
    "time"
)

type TimeSlide struct {
}

func NewTimeSlide() *TimeSlide {
    sl := new(TimeSlide)
    return sl
}

func (sl TimeSlide) Preload() {

}

func (sl TimeSlide) Draw(s *Surface) {
    s.Clear()
    t := time.Now()
    l1 := strings.ToUpper(t.Format("Jan 2"))
    l2 := t.Format("3:04:05 PM")
    c1 := Color{255, 255, 255}
    c2 := Color{255, 255, 0}
    s.WriteString(l1, c1, ALIGN_CENTER, s.Midpoint, 8)
    s.WriteString(l2, c2, ALIGN_CENTER, s.Midpoint, 16)
}
