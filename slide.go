package main

import (
    "bytes"
    "fmt"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "net/http"
    "os"
    "time"
)

type Slide interface {
    // Called when slideshow is being started
    Initialize()
    // Called when slideshow is being stopped
    Terminate()
    // Called when slide is brought into view
    StartDraw(d Display)
    // Called when a different slide is brought into view
    StopDraw()
    // Controls whether slide will be skipped in slideshow
    IsEnabled() bool
}

type HttpCallback func(respBytes []byte) (result bool)

type HttpHelper struct {
    BaseUrl          string
    RefreshInterval  time.Duration
    Callback         HttpCallback
    LastFetchSuccess bool

    RefreshTicker *time.Ticker
}

func NewHttpHelper(baseUrl string, refreshInterval time.Duration, callback HttpCallback) *HttpHelper {
    h := new(HttpHelper)
    h.BaseUrl = baseUrl
    h.RefreshInterval = refreshInterval
    h.Callback = callback
    log.WithFields(log.Fields{
        "url":      baseUrl,
        "interval": refreshInterval,
    }).Debug("HttpHelper created.")
    return h
}

func (this *HttpHelper) StartLoop() {
    log.WithFields(log.Fields{
        "url":      this.BaseUrl,
        "interval": this.RefreshInterval,
    }).Debug("HttpHelper refresh loop started.")

    // Set up period refresh of the data
    this.RefreshTicker = time.NewTicker(this.RefreshInterval)
    go func() {
        for range this.RefreshTicker.C {
            this.Fetch()
        }
    }()

    // Get the data once now (synchronously)
    this.Fetch()
}

func (this *HttpHelper) StopLoop() {
    this.RefreshTicker.Stop()
}

func (this *HttpHelper) Fetch() {
    resp, httpErr := http.Get(this.BaseUrl)
    if httpErr != nil {
        log.WithFields(log.Fields{
            "url":   this.BaseUrl,
            "error": httpErr,
        }).Warn("HTTP error in HttpHelper.")
        this.LastFetchSuccess = false
        return
    }

    respBuf := new(bytes.Buffer)
    respBuf.ReadFrom(resp.Body)
    respBytes := respBuf.Bytes()

    this.LastFetchSuccess = this.Callback(respBytes)

    log.WithFields(log.Fields{
        "url":     this.BaseUrl,
        "success": this.LastFetchSuccess,
    }).Debug("Fetch complete.")

    // Output debug file, maybe
    if DEBUG_HTTP {
        ioutil.WriteFile(fmt.Sprintf("debug/%d.txt", time.Now().Unix()), respBytes, os.FileMode(770))
    }
}
