package main

import (
    "fmt"
    "image"
    "image/color"
    "time"
    "cloud.google.com/go/civil"
)

type CountdownSlide struct {
    events []CountdownEvent
}

func NewCountdownSlide(events []CountdownEvent) *CountdownSlide {
    this := new(CountdownSlide)
    this.events = events
    return this
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
    today := civil.DateOf(time.Now())
    var filteredEvents []CountdownEvent
    for _,event := range(this.events) {
        if event.date.Before(today) {
            continue;
        }
        filteredEvents = append(filteredEvents, event)
    }

    // Hardcoded y values for each line that provide reasonable spacing
    var yVals []int
    switch len(filteredEvents) {
    case 1:
        yVals = []int{13}
    case 2:
        yVals = []int{6,19}
    case 3:
        yVals = []int{3, 13, 23}
    default:
        yVals = []int{0,8,16,24}
    }

    numberColor := color.RGBA{255, 255, 255, 255}
    for i,event := range(filteredEvents) {
        y := yVals[i]
        WriteString(img, fmt.Sprintf("%d", event.date.DaysSince(today)), numberColor, ALIGN_RIGHT, 20, y)
        WriteString(img, event.label, event.color, ALIGN_LEFT, 26, y)

        // We don't have room to display more than 4 events, so stop
        if i >= 3 {
            break
        }
    }
}

type CountdownEvent struct {
    date civil.Date
    label string
    color color.RGBA
}