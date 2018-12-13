package main

import (
    "fmt"
    "image"
    "time"
)

const WELCOME_DURATION = 5 * time.Second
const PRELOAD_DURATION = 2 * time.Second
const DRAW_INTERVAL = 1 * time.Second
const SCREEN_WIDTH = 128
const SCREEN_HEIGHT = 32

// Control debug settings
const DEBUG_DRAW = false
const DEBUG_HTTP = false

func main() {
    // Set up the glyph mappings
    InitGlyphs()

    // Set up the display
    // d := NewWebDisplay()
    d := NewLedDisplay()
    d.Initialize()

    config := LedSignConfig{
        NightModeStartHour:   23,
        NightModeEndHour:     5,
        SlideAdvanceInterval: 15 * time.Second,
        Slides: []Slide{
            NewTimeSlide(),
            NewMbtaSlide(MBTA_STATION_ID_MGH),
            NewChristmasSlide(),
            NewWeatherSlide(BOSTON_LATLNG),
        },
    }
    RunMultiSlide(d, config)
}

func RunMultiSlide(d Display, config LedSignConfig) {

    // Initial condition - set to display welcome slide (not preloaded)
    var currentSlideId int
    var currentSlide Slide
    var timeUntilAdvance time.Duration
    var isNextSlidePreloaded bool

    setInitialCondition := func() {
        currentSlideId = -1
        currentSlide = NewWelcomeSlide()
        timeUntilAdvance = WELCOME_DURATION
        isNextSlidePreloaded = false
    }

    nextSlideId := func() int {
        return (currentSlideId + 1) % len(config.Slides)
    }

    drawNightMode := func() {
        fmt.Printf("Time is %s, sign in night mode.", time.Now().String())
        // Create a blank image and pass it directly to display
        img := image.NewRGBA(image.Rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT))
        d.Redraw(img)
        // Check every hour to see if we're out of night mode
        time.Sleep(1 * time.Hour)
        // Reset the slide to be drawn for when night mode ends
        setInitialCondition()
    }

    // Call the function to set initial condition on startup
    setInitialCondition()

    // Main loop
    for {

        // If we're past the time to switch to night mode, do that instead
        if time.Now().Hour() >= config.NightModeStartHour ||
            time.Now().Hour() < config.NightModeEndHour {
            drawNightMode()
            continue
        }

        // Create a blank image, draw to it, then pass that to the display
        img := image.NewRGBA(image.Rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT))
        currentSlide.Draw(img)
        d.Redraw(img)

        // Preload the next slide (once) in a separate thread
        if timeUntilAdvance <= PRELOAD_DURATION && !isNextSlidePreloaded {
            go config.Slides[nextSlideId()].Preload()
            isNextSlidePreloaded = true
            fmt.Printf("Preloading slide %d\n", nextSlideId())
        }

        // Advance the slide, if enough time has elapsed
        if timeUntilAdvance <= 0 {
            currentSlideId = nextSlideId()
            currentSlide = config.Slides[currentSlideId]
            timeUntilAdvance = config.SlideAdvanceInterval
            isNextSlidePreloaded = false
            fmt.Printf("Advancing to slide %d\n", currentSlideId)
        }

        // Wait until we're ready to redraw
        timeUntilAdvance -= DRAW_INTERVAL
        time.Sleep(DRAW_INTERVAL)
    }
}

func NextSlideId(current int, total int) int {
    current++
    if current >= total {
        current = 0
    }
    return current
}

type LedSignConfig struct {

    // Night Mode: Sign automatically goes dark during the given time interval.
    // Evening hour after which the sign will be enabled (24-hour format).
    NightModeStartHour int
    // Morning hour before which the sign will be enabled (24-hour format).
    NightModeEndHour int

    // If nonzero, how much time before slideshow should advance to next slide.
    // Otherwise, will only stay on one slide
    SlideAdvanceInterval time.Duration

    // List of slides to display in slideshow mode, or single slide to display.
    Slides []Slide
}
