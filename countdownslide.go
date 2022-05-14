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
	sl := new(CountdownSlide)
	sl.events = events
	return sl
}

func (sl *CountdownSlide) Initialize() {

}

func (sl *CountdownSlide) Terminate() {

}

func (sl *CountdownSlide) StartDraw(d Display) {
	DrawOnce(d, sl.Draw)
}

func (sl *CountdownSlide) StopDraw() {

}

func (sl *CountdownSlide) IsEnabled() bool {
	return true
}

func (sl *CountdownSlide) Draw(img *image.RGBA) {
	today := civil.DateOf(time.Now())
	var filteredEvents []CountdownEvent
	for _, event := range sl.events {
		if event.date.Before(today) {
			continue
		}
		filteredEvents = append(filteredEvents, event)
	}

	// Hardcoded y values for each line that provide reasonable spacing
	var yVals []int
	switch len(filteredEvents) {
	case 1:
		yVals = []int{13}
	case 2:
		yVals = []int{6, 19}
	case 3:
		yVals = []int{3, 13, 23}
	default:
		yVals = []int{0, 8, 16, 24}
	}

	numberColor := color.RGBA{255, 255, 255, 255}
	for i, event := range filteredEvents {
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
	date  civil.Date
	label string
	color color.RGBA
}
