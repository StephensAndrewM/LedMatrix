package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Controller struct {
	Slideshow  *Slideshow
	ShutdownCh chan bool
}

func NewController(s *Slideshow) *Controller {
	ctrl := new(Controller)
	ctrl.Slideshow = s
	ctrl.ShutdownCh = make(chan bool)
	return ctrl
}

func (ctrl *Controller) RunUntilShutdown() {
	go func() {
		log.Info("Started HTTP controller endpoint.")
		log.Warn(http.ListenAndServe(":5000", ctrl))
	}()
	// This blocks until shutdown signal is received
	<-ctrl.ShutdownCh
}

func (ctrl *Controller) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
		if !ctrl.Slideshow.Running {
			ctrl.Slideshow.Start()
			ctrl.SendResponse(res, 200, "Starting slideshow")
		} else {
			ctrl.SendResponse(res, 412, "Cannot start, slideshow already running")
		}
	case "/stop":
		if ctrl.Slideshow.Running {
			ctrl.Slideshow.Stop()
			ctrl.SendResponse(res, 200, "Stopping slideshow")
		} else {
			ctrl.SendResponse(res, 412, "Cannot stop, slideshow already stopped")
		}
	case "/freeze":
		if !ctrl.Slideshow.Frozen {
			ctrl.Slideshow.Freeze()
			ctrl.SendResponse(res, 200, "Freezing slideshow")
		} else {
			ctrl.SendResponse(res, 412, "Cannot freeze, slideshow already frozen")
		}
	case "/unfreeze":
		if ctrl.Slideshow.Frozen {
			ctrl.Slideshow.Unfreeze()
			ctrl.SendResponse(res, 200, "Unfreezing slideshow")
		} else {
			ctrl.SendResponse(res, 412, "Cannot unfreeze, slideshow already unfrozen")
		}
	case "/shutdown":
		ctrl.ShutdownCh <- true
		ctrl.SendResponse(res, 200, "Shutting down slideshow controller")
	default:
		log.WithFields(log.Fields{
			"endpoint": req.URL.Path,
		}).Debug("Unknown request type")
		ctrl.SendResponse(res, 400, "Unknown request type")
	}
}

func (ctrl *Controller) SendResponse(res http.ResponseWriter, code int, message string) {
	res.WriteHeader(code)
	res.Write([]byte(message + "\n"))
}
