package main

import (
    log "github.com/sirupsen/logrus"
    "net/http"
    "time"
    "os/exec"
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

    log.Info("All slides reported readiness.")

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

    // Attempt to update time before displaying anything calculated
    SyncTime()

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

// Checks for internet periodically, not returning until connected.
func WaitForConnection() {
    c := 1
    for {
        if ConnectionPresent() {
            log.WithFields(log.Fields{
                "checks": c,
            }).Info("Internet connection present.")
            return
        }
        time.Sleep(1 * time.Second)
        c++
    }
}

// Synchronizes with a NTP server, in case Pi lost power for a while
func SyncTime() {
    cmd := exec.Command("/usr/sbin/ntpdate", "-s", "time.google.com")
    err := cmd.Run()
    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Warning("Failed NTP time synchronization.")
    }
}

// Sanity check for internet access. Not bulletproof but works.
func ConnectionPresent() bool {
    _, err := http.Get("http://clients3.google.com/generate_204")
    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Debug("Connection failed.")
    }
    return err == nil
}
