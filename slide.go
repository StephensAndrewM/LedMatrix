package main

import (
    "bytes"
    "fmt"
    "time"
    "net/http"
)

type Slide interface {
    Preload()
    Draw(s *Surface)
}

// Display a quick error message on screen
func ShowError(s *Surface, space int, code int) {
    yellow := Color{255, 255, 0}
    s.Clear()
    msg := fmt.Sprintf("E #%02d-%02d", space, code)
    s.WriteString(msg, yellow, ALIGN_LEFT, 0, 0)
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
        // TODO Display error on screen
    }

    respBuf := new(bytes.Buffer)
    respBuf.ReadFrom(resp.Body)

    this.LastFetchTime = time.Now()
    this.CachedResponse = respBuf.Bytes()

    return respBuf.Bytes(), true
}
