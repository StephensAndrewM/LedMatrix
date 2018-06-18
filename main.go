package main

import (
	"fmt"
	"time"
)

const PRELOAD_SEC = 2
const SEC_PER_SLIDE = 15
const SURFACE_WIDTH = 128
const SURFACE_HEIGHT = 32
const RELOAD_INTERVAL = 30

func main() {
	// Create the physical display device
	d := NewWebDisplay()
	d.Initialize()

	// Create the virtual drawing surface
	s := NewSurface(SURFACE_WIDTH, SURFACE_HEIGHT)
	RunMultiSlide(d, s)
}

func RunSingleSlide(d Display, s *Surface) {
	slide := NewMbtaSlide(MBTA_STATION_ID_PARK, MBTA_STATION_NAME_PARK)
	fmt.Printf("Initially loading slide\n")
	go slide.Preload()
	time.Sleep(1 * time.Second)

	elapsedTime := 0
	for {
		slide.Draw(s)
		d.Redraw(s)
		elapsedTime++
		// Call preload() again if it's been long enough
		if elapsedTime >= RELOAD_INTERVAL {
			fmt.Printf("Reloading slide\n")
			go slide.Preload()
			elapsedTime = 0
		}
		time.Sleep(1 * time.Second)
	}

}

func RunMultiSlide(d Display, s *Surface) {
	// Get all of the slides we'll be using
	slides := GetAllSlides()

	// Initial condition (note slide 0 is not preloaded)
	var currentSlideId, elapsedTime int
	currentSlide := slides[currentSlideId]
	elapsedTime = 12

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
		// NewGlyphTestSlide(TEST_LETTERS),
		// NewGlyphTestSlide(TEST_NUMSYM),
		NewMbtaSlide(MBTA_STATION_ID_PARK, MBTA_STATION_NAME_PARK),
		NewMbtaSlide(MBTA_STATION_ID_GOVCTR, MBTA_STATION_NAME_GOVCTR),
		NewMbtaSlide(MBTA_STATION_ID_HARVARD, MBTA_STATION_NAME_HARVARD),
		// NewWeatherSlide(SUNNYVALE_ZIP),
	}
}

func NextSlideId(current int, total int) int {
	current++
	if current >= total {
		current = 0
	}
	return current
}
