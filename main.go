package main

import (
    "flag"
    log "github.com/sirupsen/logrus"
    "time"
)

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
const ADVANCE_INTERVAL = 15 * time.Second
const SCREEN_WIDTH = 128
const SCREEN_HEIGHT = 32

// Control debug settings
const DEBUG_DRAW = false
const DEBUG_HTTP = false

func main() {
    // Init flags for use everywhere
    flag.Parse()

    // Set global settings for logging
    log.SetLevel(log.DebugLevel)
    log.SetFormatter(&log.TextFormatter{
        FullTimestamp: true,
    })

    // Set up the glyph mappings
    InitGlyphs()

    // Set up the display - use hardware as default
    var d Display
    if *useWebDisplayFlag {
        d = NewWebDisplay()
    } else {
        d = NewLedDisplay()
    }
    d.Initialize()

    // Set up the slideshow (controls drawing and advancing)
    s := NewSlideshow(d, GetSlides())
    s.Start()

    // Start the HTTP show controller
    NewController(s)

    // Keep running forever
    select {}
}
