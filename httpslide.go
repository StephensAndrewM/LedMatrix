package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type HttpConfig struct {
	SlideId            string
	RefreshInterval    time.Duration
	RequestUrl         string
	RequestUrlCallback func() (*http.Request, error)
	ParseCallback      func([]byte) bool
}

type HttpHelper struct {
	Config           HttpConfig
	LastFetchSuccess bool
	Client           *http.Client
	RefreshTicker    *time.Ticker
}

func NewHttpHelper(config HttpConfig) *HttpHelper {
	h := new(HttpHelper)
	h.Config = config
	h.Client = &http.Client{}
	return h
}

func (h *HttpHelper) StartLoop() {
	if h.RefreshTicker != nil {
		log.WithFields(log.Fields{
			"slide": h.Config.SlideId,
		}).Warn("Attempting to start HTTP loop when already started.")
		return
	}

	// Set up period refresh of the data
	h.RefreshTicker = time.NewTicker(h.Config.RefreshInterval)
	go func() {
		for range h.RefreshTicker.C {
			h.Fetch()
		}
	}()

	// Get the data once now (synchronously)
	h.Fetch()
}

func (h *HttpHelper) StopLoop() {
	if h.RefreshTicker == nil {
		log.WithFields(log.Fields{
			"slide": h.Config.SlideId,
		}).Warn("Attempting to stop HTTP loop when already stopped.")
		return
	}
	h.RefreshTicker.Stop()
	h.RefreshTicker = nil
}

func (h *HttpHelper) BuildRequest() (*http.Request, error) {
	if h.Config.RequestUrlCallback != nil {
		return h.Config.RequestUrlCallback()
	}
	req, err := http.NewRequest("GET", h.Config.RequestUrl, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (h *HttpHelper) Fetch() {
	req, reqErr := h.BuildRequest()
	if reqErr != nil {
		log.WithFields(log.Fields{
			"slide": h.Config.SlideId,
			"req":   req,
			"error": reqErr,
		}).Warn("Request error in HttpHelper.")
		h.LastFetchSuccess = false
		return
	}

	res, resErr := h.Client.Do(req)
	if resErr != nil {
		log.WithFields(log.Fields{
			"slide": h.Config.SlideId,
			"req":   req,
			"res":   res,
			"error": resErr,
		}).Warn("Response error in HttpHelper.")
		h.LastFetchSuccess = false
		return
	}

	if res.StatusCode != 200 {
		log.WithFields(log.Fields{
			"slide": h.Config.SlideId,
			"req":   req,
			"res":   res,
		}).Warn("Got non-200 response code in HttpHelper.")
		h.LastFetchSuccess = false
		return
	}

	resBuf := new(bytes.Buffer)
	resBuf.ReadFrom(res.Body)
	resBytes := resBuf.Bytes()

	h.LastFetchSuccess = h.Config.ParseCallback(resBytes)

	log.WithFields(log.Fields{
		"slide":        h.Config.SlideId,
		"req":          req,
		"fetchSuccess": h.LastFetchSuccess,
	}).Debug("Fetch complete.")

	// Output debug file, maybe
	if *debugHttp {
		outFile := fmt.Sprintf("debug/%d-%s.txt", time.Now().Unix(), h.Config.SlideId)
		log.WithFields(log.Fields{
			"req":     req,
			"outFile": outFile,
		}).Debug("Logged HTTP response data.")
		ioutil.WriteFile(outFile, resBytes, os.FileMode(0770))
	}
}
