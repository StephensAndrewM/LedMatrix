package main

import (
    "time"
)

type Slideshow struct {
    Display         Display
    AdvanceInterval time.Duration
    Slides          []Slide

    Running        bool
    CurrentSlide   Slide
    CurrentSlideId int
    AdvanceTicker  *time.Ticker
}

func NewSlideshow(d Display, config *Config) *Slideshow {
    this := new(Slideshow)
    this.Display = d
    this.AdvanceInterval = config.AdvanceInterval
    this.Slides = config.Slides
    return this
}

func (this *Slideshow) Start() {
    this.Running = true
    this.CurrentSlideId = -1

    // Display the welcome slide while loading
    this.CurrentSlide = NewWelcomeSlide()
    this.CurrentSlide.StartDraw(this.Display)

    // Block until all slides have loaded data
    this.WaitForReadiness()

    // Then go to the first slide and run for real
    this.Advance()

    // Increment the slide number periodically and start/stop drawing
    this.AdvanceTicker = time.NewTicker(this.AdvanceInterval)
    go func() {
        for range this.AdvanceTicker.C {
            this.Advance()
        }
    }()
}

func (this *Slideshow) Advance() {
    this.CurrentSlide.StopDraw()

    for {
        this.CurrentSlideId = (this.CurrentSlideId + 1) % len(this.Slides)
        this.CurrentSlide = this.Slides[this.CurrentSlideId]
        // If the slide is enabled, stop the loop
        if this.CurrentSlide.IsEnabled() {
            break
        }
        // Otherwise we loop until we find an enabled slide
        // TODO make sure this doesn't get stuck if no slide is enabled
    }

    this.CurrentSlide.StartDraw(this.Display)
}

func (this *Slideshow) WaitForReadiness() {
    // Don't initialize until internet is available
    WaitForConnection()

    // Initialize all slides (attempt fetching initial content)
    // This call on each slide blocks until request is complete
    for _, s := range this.Slides {
        s.Initialize()
    }
}

func (this *Slideshow) Stop() {
    this.Running = false
    this.CurrentSlide.StopDraw()
    this.AdvanceTicker.Stop()

    // Stop any slide-level tickers
    for _, s := range this.Slides {
        s.Terminate()
    }

    // Draw a blank image
    this.Display.Redraw(NewBlankImage())
}
