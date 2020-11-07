package main

import (
    "time"
    "image/color"
    "cloud.google.com/go/civil"
)

// Provides the customizable options for the slideshow
func GetConfig() *Config {
    return &Config{
        AdvanceInterval: 15 * time.Second,
        Slides: []Slide{
            NewTimeSlide(),
            NewWeatherSlide(BOSTON_LAT, BOSTON_LNG),
            NewCountdownSlide([]CountdownEvent{
                CountdownEvent{
                    civil.Date{2020, time.September, 7},
                    "LABOR DAY",
                    color.RGBA{0, 0, 255, 255},
                },
                CountdownEvent{
                    civil.Date{2020, time.October, 31},
                    "HALLOWEEN",
                    color.RGBA{255, 0, 255, 255},
                },
                CountdownEvent{
                    civil.Date{2020, time.November, 26},
                    "THANKSGIVING",
                    color.RGBA{0, 255, 0, 255},
                },
            }),
            NewCovidSlide(),
        },
    }
}
