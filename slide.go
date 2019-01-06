package main

import (
    "bytes"
    "fmt"
    log "github.com/sirupsen/logrus"
    "image"
    "io/ioutil"
    "net/http"
    "os"
    "time"
)

type Slide interface {
    Draw(base *image.RGBA)
}

type HttpCallback func(respBytes []byte) (result bool)

type HttpHelper struct {
    BaseUrl          string
    RefreshInterval  time.Duration
    Callback         HttpCallback
    LastFetchSuccess bool
}

func NewHttpHelper(baseUrl string, refreshInterval time.Duration, callback HttpCallback) *HttpHelper {
    h := new(HttpHelper)
    h.BaseUrl = baseUrl
    h.RefreshInterval = refreshInterval
    h.Callback = callback
    log.WithFields(log.Fields{
        "url":      baseUrl,
        "interval": refreshInterval,
    }).Debug("HttpHelper initialized.")

    // Set up period refresh of the data
    ticker := time.NewTicker(refreshInterval)
    go func() {
        for range ticker.C {
            h.Fetch()
        }
    }()

    // Get the data once now so we don't have to wait
    h.Fetch()

    return h
}

func (this *HttpHelper) Fetch() {
    // Don't fetch data while in night mode, unless we're about to wake back
    // up (now + refresh interval), in which case continue.
    if InNightMode(time.Now()) &&
        InNightMode(time.Now().Add(this.RefreshInterval)) {
        return
    }

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
