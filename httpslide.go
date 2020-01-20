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

// Interface for a slide that makes periodic HTTP requests
type HttpSlide interface {
    // Indicates the frequency that the slide should refresh
    GetRefreshInterval() time.Duration
    // Provides the HTTP request object to send
    BuildRequest() (*http.Request, error)

    // Callback to parse and convert an incoming HTTP response for validity
    Parse(resBytes []byte) bool
}

type HttpHelper struct {
    Slide            HttpSlide
    LastFetchSuccess bool
    Client           *http.Client
    RefreshTicker    *time.Ticker
}

func NewHttpHelper(slide HttpSlide) *HttpHelper {
    h := new(HttpHelper)
    h.Slide = slide
    h.Client = &http.Client{}
    return h
}

func (this *HttpHelper) StartLoop() {
    req, err := this.Slide.BuildRequest()
    log.WithFields(log.Fields{
        "req":      req,
        "err":      err,
        "interval": this.Slide.GetRefreshInterval(),
    }).Debug("HttpHelper refresh loop started.")
    if err != nil {
        return
    }

    // Set up period refresh of the data
    this.RefreshTicker = time.NewTicker(this.Slide.GetRefreshInterval())
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
    req, reqErr := this.Slide.BuildRequest()
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

    this.LastFetchSuccess = this.Slide.Parse(resBytes)

    log.WithFields(log.Fields{
        "req":          req,
        "fetchSuccess": this.LastFetchSuccess,
    }).Debug("Fetch complete.")

    // Output debug file, maybe
    if DEBUG_HTTP {
        outFile := fmt.Sprintf("debug/%d.txt", time.Now().Unix())
        log.WithFields(log.Fields{
            "req":     req,
            "outFile": outFile,
        }).Debug("Logged HTTP response data.")
        ioutil.WriteFile(outFile, resBytes, os.FileMode(770))
    }
}
