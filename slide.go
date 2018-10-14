package main

import (
    "bytes"
    "fmt"
    "time"
    "net/http"
    "os"
    "io/ioutil"
    "image"
)

type Slide interface {
    Preload()
    Draw(base *image.RGBA)
}

type HttpHelper struct {
    // Object settings
    BaseUrl             string
    RefreshInternal     time.Duration

    // Internal vars
    LastFetchTime       time.Time
    CachedResponse      []byte
}

func NewHttpHelper(baseUrl string, refreshInterval time.Duration) *HttpHelper {
    h := new(HttpHelper)
    h.BaseUrl = baseUrl
    h.RefreshInternal = refreshInterval
    return h
}

func (this *HttpHelper) Fetch() ([]byte, bool) {
    now := time.Now()
    if now.Before(this.LastFetchTime.Add(this.RefreshInternal)) {
        return this.CachedResponse, true
    }

    resp, httpErr := http.Get(this.BaseUrl)
    if httpErr != nil {
        fmt.Printf("Error loading data: %s\n", httpErr)
        return nil, false
    }

    respBuf := new(bytes.Buffer)
    respBuf.ReadFrom(resp.Body)

    this.LastFetchTime = time.Now()
    this.CachedResponse = respBuf.Bytes()

    // Output debug file
    ioutil.WriteFile("debug.txt", respBuf.Bytes(), os.FileMode(770))

    return respBuf.Bytes(), true
}
