package main

import (
	"time"
)

// Provides the customizable options for the slideshow
func GetConfig() *Config {
	return &Config{
		AdvanceInterval: 15 * time.Second,
		Slides: []Slide{
			NewTimeSlide(),
			NewMbtaSlide(MBTA_STATION_ID_MGH),
			NewWeatherSlide(BOSTON_LATLNG),
			NewChristmasSlide(),
		},
	}
}
