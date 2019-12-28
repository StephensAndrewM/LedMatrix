package main

import (
    log "github.com/sirupsen/logrus"
    "net/http"
)

type Controller struct {
    Slideshow  *Slideshow
    ShutdownCh chan bool
}

func NewController(s *Slideshow) *Controller {
    this := new(Controller)
    this.Slideshow = s
    this.ShutdownCh = make(chan bool)
    return this
}

func (this *Controller) RunUntilShutdown() {
    go func() {
        log.Info("Started HTTP controller endpoint.")
        log.Warn(http.ListenAndServe(":5000", this))
    }()
    // This blocks until shutdown signal is received
    <-this.ShutdownCh
}

func (this *Controller) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    if req.Method != "POST" {
        log.WithFields(log.Fields{
            "endpoint": req.URL.Path,
            "method":   req.Method,
        }).Debug("Request with bad method")
        res.WriteHeader(405)
        return
    }

    switch req.URL.Path {
    case "/start":
        if !this.Slideshow.Running {
            log.Debug("Restarting slideshow")
            this.Slideshow.Start()
            res.WriteHeader(200)
        } else {
            log.Debug("Cannot start, slideshow already running")
            res.WriteHeader(412)
        }
    case "/stop":
        if this.Slideshow.Running {
            log.Debug("Stopping slideshow")
            this.Slideshow.Stop()
            res.WriteHeader(200)
        } else {
            log.Debug("Cannot stop, slideshow already stopped")
            res.WriteHeader(412)
        }
    case "/shutdown":
        log.Debug("Shutting down slideshow controller")
        this.ShutdownCh <- true
        res.WriteHeader(200)
    default:
        log.WithFields(log.Fields{
            "endpoint": req.URL.Path,
        }).Debug("Unknown request")
        res.WriteHeader(400)
    }
}
