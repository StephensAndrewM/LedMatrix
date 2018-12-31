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
    Preload()
    Draw(base *image.RGBA)
}

type HttpHelper struct {
    // Object settings
    BaseUrl         string
    RefreshInternal time.Duration

    // Internal vars
    LastFetchTime  time.Time
    CachedResponse []byte
}

func NewHttpHelper(baseUrl string, refreshInterval time.Duration) *HttpHelper {
    h := new(HttpHelper)
    h.BaseUrl = baseUrl
    h.RefreshInternal = refreshInterval
    log.WithFields(log.Fields{
        "url":      baseUrl,
        "interval": refreshInterval,
    }).Debug("HttpHelper initialized.")
    return h
}

func (this *HttpHelper) Fetch() ([]byte, bool) {
    now := time.Now()
    if now.Before(this.LastFetchTime.Add(this.RefreshInternal)) {
        log.WithFields(log.Fields{
            "cacheTime": this.LastFetchTime,
            "cacheTimeDiff": now.Sub(this.LastFetchTime),
        }).Debug("Returning cached response.")
        return this.CachedResponse, true
    }

    resp, httpErr := http.Get(this.BaseUrl)
    if httpErr != nil {
        log.WithFields(log.Fields{
            "error": httpErr,
        }).Warn("HTTP error in HttpHelper.")
        return nil, false
    }

    respBuf := new(bytes.Buffer)
    respBuf.ReadFrom(resp.Body)

    this.LastFetchTime = time.Now()
    this.CachedResponse = respBuf.Bytes()

    // Output debug file
    if DEBUG_HTTP {
        ioutil.WriteFile(fmt.Sprintf("debug/%d.txt", time.Now().Unix()), respBuf.Bytes(), os.FileMode(770))
    }

    return respBuf.Bytes(), true
}
