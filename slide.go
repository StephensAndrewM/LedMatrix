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

type HttpCallback func(resBytes []byte) (result bool)

type HttpHelper struct {
    BaseUrl           string
    RefreshInterval   time.Duration
    Callback          HttpCallback
    BasicAuthUsername string
    BasicAuthPassword string
    LastFetchSuccess  bool
    Client            *http.Client
    RefreshTicker     *time.Ticker
}

func NewHttpHelper(baseUrl string, refreshInterval time.Duration, callback HttpCallback) *HttpHelper {
    h := new(HttpHelper)
    h.BaseUrl = baseUrl
    h.RefreshInterval = refreshInterval
    h.Callback = callback
    h.Client = &http.Client{}
    log.WithFields(log.Fields{
        "url":      baseUrl,
        "interval": refreshInterval,
    }).Debug("HttpHelper created.")
    return h
}

func NewHttpHelperWithAuth(baseUrl string, refreshInterval time.Duration, callback HttpCallback, username, password string) *HttpHelper {
    h := NewHttpHelper(baseUrl, refreshInterval, callback)
    h.BasicAuthUsername = username
    h.BasicAuthPassword = password
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
    req, reqErr := http.NewRequest("GET", this.BaseUrl, nil)
    if reqErr != nil {
        log.WithFields(log.Fields{
            "url":   this.BaseUrl,
            "error": reqErr,
        }).Warn("Request error in HttpHelper")
        this.LastFetchSuccess = false
        return
    }
    // Set up HTTP auth, if needed
    if this.BasicAuthUsername != "" && this.BasicAuthPassword != "" {
        req.SetBasicAuth(this.BasicAuthUsername, this.BasicAuthPassword)
    }

    res, resErr := this.Client.Do(req)
    if resErr != nil {
        log.WithFields(log.Fields{
            "url":   this.BaseUrl,
            "error": resErr,
        }).Warn("Response error in HttpHelper.")
        this.LastFetchSuccess = false
        return
    }

    resBuf := new(bytes.Buffer)
    resBuf.ReadFrom(res.Body)
    resBytes := resBuf.Bytes()

    this.LastFetchSuccess = this.Callback(resBytes)

    log.WithFields(log.Fields{
        "url":     this.BaseUrl,
        "success": this.LastFetchSuccess,
    }).Debug("Fetch complete.")

    // Output debug file, maybe
    if DEBUG_HTTP {
        outFile := fmt.Sprintf("debug/%d.txt", time.Now().Unix())
        log.WithFields(log.Fields{
            "url":     this.BaseUrl,
            "outFile": this.outFile,
        }).Debug("Logged HTTP response data.")
        ioutil.WriteFile(outFile, resBytes, os.FileMode(770))
    }
}
