package main

import (
	"fmt"
	"time"
)

const PRELOAD_SEC = 5
const SEC_PER_SLIDE = 6
const SURFACE_WIDTH = 128
const SURFACE_HEIGHT = 32

func main() {
	// Create the physical display device
	d := NewWebDisplay()
	d.Initialize()

	// Create the virtual drawing surface
	s := NewSurface(SURFACE_WIDTH, SURFACE_HEIGHT)

	// Get all of the slides we'll be using
	slides := GetAllSlides()

	// Initial condition (note slide 0 is not preloaded)
	var currentSlideId, elapsedTime int
	currentSlide := slides[currentSlideId]

	// Main loop
	for {
		currentSlide.Draw(s)
		d.Redraw(s)
		elapsedTime++

		// Preload what will be the next slide concurrently
		if elapsedTime == (SEC_PER_SLIDE - PRELOAD_SEC) {
			nextSlideId := NextSlideId(currentSlideId, len(slides))
			fmt.Printf("Preloading slide %d\n", nextSlideId)
			go slides[nextSlideId].Preload()
		}

		// Advance the slide when ready
		if elapsedTime >= SEC_PER_SLIDE {
			currentSlideId = NextSlideId(currentSlideId, len(slides))
			elapsedTime = 0
			currentSlide = slides[currentSlideId]
			fmt.Printf("Advancing to slide %d\n", currentSlideId)
		}
		time.Sleep(1 * time.Second)
	}
}

func GetAllSlides() []Slide {
	return []Slide{
		NewTimeSlide(),
		NewMbtaSlide(MBTA_ROUTE_RED, MBTA_STATION_DAVIS),
		NewWeatherSlide(SUNNYVALE_ZIP),
	}
}

func NextSlideId(current int, total int) int {
	current++
	if current >= total {
		current = 0
	}
	return current
}
