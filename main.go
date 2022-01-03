package main

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	AdvanceInterval time.Duration
	Slides          []Slide
}

// Flags that are generally environment-dependent
var generateImagesFlag = flag.Bool("generate_images", false,
	"If true, generates slide images instead of running as slideshow.")
var debugLogFlag = flag.Bool("debug_log", false,
	"If true, prints out debug-level log statements.")
var debugHttp = flag.Bool("debug_http", false,
	"If true, logs all HTTP responses to local directory.")
var debugDraw = flag.Bool("debug_draw", false,
	"If true, displays bounding boxes for drawn elements.")

// Constants that generally don't need to be configured
const SCREEN_WIDTH = 128
const SCREEN_HEIGHT = 32

func main() {
	// Init flags for use everywhere
	flag.Parse()

	// Set global settings for logging
	if *debugLogFlag {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Set up the glyph and icon mappings
	InitGlyphs()
	InitIcons()

	if *generateImagesFlag {
		GenerateImages()
	} else {
		RunAsSlideshow()
	}
}

func RunAsSlideshow() {
	// Grab the global config object to pass elsewhere
	config := GetConfig()

	// Set up the display on hardware
	d := NewLedDisplay()
	d.Initialize()

	// Set up the slideshow (controls drawing and advancing)
	s := NewSlideshow(d, config)
	s.Start()

	// Start the HTTP show controller, which keeps the program running
	c := NewController(s)
	c.RunUntilShutdown()
}

func GenerateImages() {
	config := GetConfig()

	d := NewSaveToFileDisplay()

	// For each slide, initialize then draw once
	for _, s := range config.Slides {
		d.SetSlideId(s)
		s.Initialize()
		s.StartDraw(d)
		s.StopDraw()
	}
}
