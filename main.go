package main

import (
	"fmt"
	"time"
    "image"
)

const PRELOAD_SEC = 2
const SEC_PER_SLIDE = 10
// How frequently to re-call Preload() in single-slide mode
const RELOAD_INTERVAL = 30

// Screen dimensions
const SCREEN_WIDTH = 128
const SCREEN_HEIGHT = 32

// Control debug settings
const DEBUG_DRAW = false
const DEBUG_HTTP = true

func main() {
    // Set up the glyph mappings
    InitGlyphs()

    // Set up a disply and run the slides
	d := NewWebDisplay()
	d.Initialize()
	// RunMultiSlide(d)
    RunSingleSlide(d)
}

func RunSingleSlide(d Display) {
	// slide := NewMbtaSlide(MBTA_STATION_ID_PARK, MBTA_STATION_NAME_PARK)
    slide := NewWeatherSlide(BOSTON_LATLNG)
	fmt.Printf("Initially loading slide\n")
	go slide.Preload()
	time.Sleep(1 * time.Second)

	elapsedTime := 0
	for {
        img := image.NewRGBA(image.Rect(0,0,SCREEN_WIDTH,SCREEN_HEIGHT))
		slide.Draw(img)
		d.Redraw(img)
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

func RunMultiSlide(d Display) {
	// Get all of the slides we'll be using
	slides := GetAllSlides()

	// Initial condition (note slide 0 is not preloaded)
	var currentSlideId, elapsedTime int
	currentSlide := slides[currentSlideId]

	// Main loop
	for {
        img := image.NewRGBA(image.Rect(0,0,SCREEN_WIDTH,SCREEN_HEIGHT))
		currentSlide.Draw(img)
		d.Redraw(img)
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
		NewWeatherSlide(BOSTON_LATLNG),
	}
}

func NextSlideId(current int, total int) int {
	current++
	if current >= total {
		current = 0
	}
	return current
}
