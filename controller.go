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

    log.WithFields(log.Fields{
        "endpoint": req.URL.Path,
    }).Debug("Controller received request")

    switch req.URL.Path {
    case "/start":
        if !this.Slideshow.Running {
            this.Slideshow.Start()
            this.SendResponse(res, 200, "Starting slideshow")
        } else {
            this.SendResponse(res, 412, "Cannot start, slideshow already running")
        }
    case "/stop":
        if this.Slideshow.Running {
            this.Slideshow.Stop()
            this.SendResponse(res, 200, "Stopping slideshow")
        } else {
            this.SendResponse(res, 412, "Cannot stop, slideshow already stopped")
        }
    case "/freeze":
        if !this.Slideshow.Frozen {
            this.Slideshow.Freeze()
            this.SendResponse(res, 200, "Freezing slideshow")
        } else {
            this.SendResponse(res, 412, "Cannot freeze, slideshow already frozen")
        }
    case "/unfreeze":
        if this.Slideshow.Frozen {
            this.Slideshow.Unfreeze()
            this.SendResponse(res, 200, "Unfreezing slideshow")
        } else {
            this.SendResponse(res, 412, "Cannot unfreeze, slideshow already unfrozen")
        }
    case "/shutdown":
        this.ShutdownCh <- true
        this.SendResponse(res, 200, "Shutting down slideshow controller")
    default:
        log.WithFields(log.Fields{
            "endpoint": req.URL.Path,
        }).Debug("Unknown request type")
        this.SendResponse(res, 400, "Unknown request type")
    }
}

func (this *Controller) SendResponse(res http.ResponseWriter, code int, message string) {
    res.WriteHeader(code)
    res.Write([]byte(message + "\n"))
}
