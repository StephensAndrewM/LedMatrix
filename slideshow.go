package main

import (
    "image"
    "time"
)

type Slideshow struct {
    Display Display
    Slides  []Slide

    Running        bool
    CurrentSlide   Slide
    CurrentSlideId int
    AdvanceTicker  *time.Ticker
    RedrawTicker   *time.Ticker
}

func NewSlideshow(d Display, slides []Slide) *Slideshow {
    this := new(Slideshow)
    this.Display = d
    this.Slides = slides
    return this
}

func (this *Slideshow) Start() {
    this.Running = true
    this.CurrentSlide = NewWelcomeSlide()
    this.CurrentSlideId = -1

    // Redraw ticker - update whatever slide is currently displayed
    this.RedrawTicker = time.NewTicker(DRAW_INTERVAL)
    go func() {
        for range this.RedrawTicker.C {
            this.DrawCurrent()
        }
    }()

    // Block until all slides have loaded data
    this.WaitForReadiness()

    // Then advance to the first slide and run for real
    this.Advance()

    // Advance ticker - increment the slide number periodically
    this.AdvanceTicker = time.NewTicker(ADVANCE_INTERVAL)
    go func() {
        for range this.AdvanceTicker.C {
            this.Advance()
        }
    }()
}

func (this *Slideshow) Advance() {
    this.CurrentSlideId = (this.CurrentSlideId + 1) % len(this.Slides)
    this.CurrentSlide = this.Slides[this.CurrentSlideId]
    // Draw now just in case this is out of sync with draw timer
    this.DrawCurrent()
}

func (this *Slideshow) DrawCurrent() {
    img := image.NewRGBA(image.Rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT))
    this.CurrentSlide.Draw(img)
    this.Display.Redraw(img)
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
    this.RedrawTicker.Stop()
    this.AdvanceTicker.Stop()

    // Stop any slide-level tickers
    for _, s := range this.Slides {
        s.Terminate()
    }

    // Draw a blank image
    img := image.NewRGBA(image.Rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT))
    this.Display.Redraw(img)
}
