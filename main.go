package main

import (
	"fmt"
	"time"
)

const SEC_PER_SLIDE = 15
const SURFACE_WIDTH = 64
const SURFACE_HEIGHT = 32

func main() {
	// Create the physical display device
	d := NewWebDisplay()
	d.Initialize()

	// Create the virtual drawing surface
	s := NewSurface(SURFACE_WIDTH, SURFACE_HEIGHT)

	// Get all of the slides we'll be using
	slides := GetAllSlides()

	var slideOffset, elapsedTime int
	for {
		fmt.Printf("Slide %d Time %d\n", slideOffset, elapsedTime)
		sl := slides[slideOffset]
		sl.Draw(s)
		d.Redraw(s)
		elapsedTime++
		if (elapsedTime >= SEC_PER_SLIDE) {
			slideOffset = IncrementSlideOffset(slideOffset, len(slides))
			elapsedTime = 0
		}
		time.Sleep(1 * time.Second)
	}
}

func RedrawDisplay(d Display) {
	s := NewSurface(32,32)
	d.Redraw(s)
}

func GetAllSlides() []Slide {
	return []Slide{
		NewTimeSlide(),
		NewWeatherSlide(),
	}
}

func IncrementSlideOffset(current int, total int) int {
	current++
	if current >= total {
		current = 0
	}
	return current
}