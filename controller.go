package main

import(
    log "github.com/sirupsen/logrus"
    "net/http"
)

type Controller struct {
    Slideshow *Slideshow
}

func NewController(s *Slideshow) *Controller {
    this := new(Controller)
    this.Slideshow = s
    go log.Warn(http.ListenAndServe(":5000", this))
    log.Info("Started HTTP controller endpoint.")
    return this
}

func (this *Controller) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    if req.Method != "POST" {
        log.WithFields(log.Fields{
            "endpoint": req.URL.Path,
            "method": req.Method,
        }).Debug("Request with bad method")
        res.WriteHeader(405)
        return
    }

    switch req.URL.Path {
    case "/start":
        if (!this.Slideshow.Running) {
            log.Debug("Restarting slideshow")
            this.Slideshow.Start()
            res.WriteHeader(200)
        } else {
            log.Debug("Cannot start, slideshow already running")
            res.WriteHeader(412)
        }
    case "/stop":
        if (this.Slideshow.Running) {
            log.Debug("Stopping slideshow")
            this.Slideshow.Stop()           
            res.WriteHeader(200)
        } else {
            log.Debug("Cannot stop, slideshow already stopped")
            res.WriteHeader(412)
        }
        default:
            log.WithFields(log.Fields{
                "endpoint": req.URL.Path,
            }).Debug("Unknown request")
            res.WriteHeader(400)
    }
}