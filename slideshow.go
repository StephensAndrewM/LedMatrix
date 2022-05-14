package main

import (
	"net/http"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

type Slideshow struct {
	Display         Display
	AdvanceInterval time.Duration
	Slides          []Slide

	Running        bool
	Frozen         bool
	CurrentSlide   Slide
	CurrentSlideId int
	AdvanceTicker  *time.Ticker
}

func NewSlideshow(d Display, config *Config) *Slideshow {
	s := new(Slideshow)
	s.Display = d
	s.AdvanceInterval = config.AdvanceInterval
	s.Slides = config.Slides
	return s
}

func (s *Slideshow) Start() {
	s.Running = true
	s.CurrentSlideId = -1

	// Display the welcome slide while loading
	s.CurrentSlide = NewWelcomeSlide()
	s.CurrentSlide.StartDraw(s.Display)

	// Block until all slides have loaded data
	s.WaitForReadiness()

	log.Info("All slides reported readiness.")

	// Then go to the first slide and run for real
	s.Advance()

	// Increment the slide number periodically and start/stop drawing
	s.AdvanceTicker = time.NewTicker(s.AdvanceInterval)
	go func() {
		for range s.AdvanceTicker.C {
			// Don't advance if the show has been manually frozen
			if !s.Frozen {
				s.Advance()
			}
		}
	}()
}

func (s *Slideshow) Advance() {
	s.CurrentSlide.StopDraw()

	for {
		s.CurrentSlideId = (s.CurrentSlideId + 1) % len(s.Slides)
		s.CurrentSlide = s.Slides[s.CurrentSlideId]
		// If the slide is enabled, stop the loop
		if s.CurrentSlide.IsEnabled() {
			break
		}
		// Otherwise we loop until we find an enabled slide
		// This would probably get stuck if no slides are enabled at all
	}

	s.CurrentSlide.StartDraw(s.Display)
}

func (s *Slideshow) WaitForReadiness() {
	// Don't initialize until internet is available
	WaitForConnection()

	// Attempt to update time before displaying anything calculated
	SyncTime()

	// Initialize all slides (attempt fetching initial content)
	// This call on each slide blocks until request is complete
	for _, sl := range s.Slides {
		sl.Initialize()
	}
}

func (s *Slideshow) Stop() {
	s.Running = false
	s.CurrentSlide.StopDraw()
	s.AdvanceTicker.Stop()

	// Stop any slide-level tickers
	for _, sl := range s.Slides {
		sl.Terminate()
	}

	// Draw a blank image
	s.Display.Redraw(NewBlankImage())
}

func (s *Slideshow) Freeze() {
	s.Frozen = true
}

func (s *Slideshow) Unfreeze() {
	s.Frozen = false
	s.Advance()
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
