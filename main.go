package main

import (
    "flag"
    log "github.com/sirupsen/logrus"
    "image"
    "time"
)

// Global, running sign config
var config = LedSignConfig{
    NightModeStartHour:   23,
    NightModeEndHour:     5,
    SlideAdvanceInterval: 15 * time.Second,
}

// Config for the sign's slides
func GetSlides() []Slide {
    return []Slide{
        NewTimeSlide(),
        NewMbtaSlide(MBTA_STATION_ID_MGH),
        NewWeatherSlide(BOSTON_LATLNG),
    }
}

// Flags that are generally environment-dependent
var useWebDisplayFlag = flag.Bool("use_web_display", false,
    "If true, outputs to simulator instead of hardware.")

// Constants that generally don't need to be configured
const DRAW_INTERVAL = 1 * time.Second
const SCREEN_WIDTH = 128
const SCREEN_HEIGHT = 32

// Control debug settings
const DEBUG_DRAW = false
const DEBUG_HTTP = false

func main() {
    // Init flags for use everywhere
    flag.Parse()

    // Set up the glyph mappings
    InitGlyphs()
    InitLogger()

    // Set up the display - hardware as default
    var d Display
    if *useWebDisplayFlag {
        d = NewWebDisplay()
    } else {
        d = NewLedDisplay()
    }
    d.Initialize()

    RunMultiSlide(d)
}

func InitLogger() {
    log.SetLevel(log.InfoLevel)
    log.SetFormatter(&log.TextFormatter{
        FullTimestamp: true,
    })
}

func RunMultiSlide(d Display) {
    // Hold running state
    slides := GetSlides()
    currentSlideId := -1
    var currentSlide Slide

    // Display *something* while everything starts up
    currentSlide = NewWelcomeSlide()

    // Slide advance ticker - update slide number periodically
    advanceTicker := time.NewTicker(config.SlideAdvanceInterval)
    go func() {
        for range advanceTicker.C {
            currentSlideId = (currentSlideId + 1) % len(slides)
            currentSlide = slides[currentSlideId]
        }
    }()

    // Redraw ticker - update whatever slide is currently displayed
    redrawTicker := time.NewTicker(DRAW_INTERVAL)
    go func() {
        for range redrawTicker.C {
            img := image.NewRGBA(image.Rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT))
            // Leave slide blank if we're currently in night mode
            if !InNightMode(time.Now()) {
                currentSlide.Draw(img)
            }
            d.Redraw(img)
        }
    }()

    // Keep running forever
    select {}
}

func InNightMode(t time.Time) bool {
    return t.Hour() >= config.NightModeStartHour ||
        t.Hour() < config.NightModeEndHour
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
}
