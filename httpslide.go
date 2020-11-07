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

func (this *HttpHelper) StartLoop() {
    if this.RefreshTicker != nil {
        log.WithFields(log.Fields{
            "slide": this.Config.SlideId,
        }).Warn("Attempting to start HTTP loop when already started.")
        return
    }

    // Set up period refresh of the data
    this.RefreshTicker = time.NewTicker(this.Config.RefreshInterval)
    go func() {
        for range this.RefreshTicker.C {
            this.Fetch()
        }
    }()

    // Get the data once now (synchronously)
    this.Fetch()
}

func (this *HttpHelper) StopLoop() {
    if this.RefreshTicker == nil {
        log.WithFields(log.Fields{
            "slide": this.Config.SlideId,
        }).Warn("Attempting to stop HTTP loop when already stopped.")
        return
    }
    this.RefreshTicker.Stop()
    this.RefreshTicker = nil
}

func (this *HttpHelper) BuildRequest() (*http.Request, error) {
    if this.Config.RequestUrlCallback != nil {
        return this.Config.RequestUrlCallback()
    }
    return http.NewRequest("GET", this.Config.RequestUrl, nil)
}

func (this *HttpHelper) Fetch() {
    req, reqErr := this.BuildRequest()
    if reqErr != nil {
        log.WithFields(log.Fields{
            "req":   req,
            "error": reqErr,
        }).Warn("Request error in HttpHelper")
        this.LastFetchSuccess = false
        return
    }

    res, resErr := this.Client.Do(req)
    if resErr != nil {
        log.WithFields(log.Fields{
            "req":   req,
            "res":   res,
            "error": resErr,
        }).Warn("Response error in HttpHelper.")
        this.LastFetchSuccess = false
        return
    }

    resBuf := new(bytes.Buffer)
    resBuf.ReadFrom(res.Body)
    resBytes := resBuf.Bytes()

    this.LastFetchSuccess = this.Config.ParseCallback(resBytes)

    log.WithFields(log.Fields{
        "req":          req,
        "fetchSuccess": this.LastFetchSuccess,
    }).Debug("Fetch complete.")

    // Output debug file, maybe
    if DEBUG_HTTP {
        outFile := fmt.Sprintf("debug/%d-%s.txt", time.Now().Unix(), this.Config.SlideId)
        log.WithFields(log.Fields{
            "req":     req,
            "outFile": outFile,
        }).Debug("Logged HTTP response data.")
        ioutil.WriteFile(outFile, resBytes, os.FileMode(770))
    }
}
