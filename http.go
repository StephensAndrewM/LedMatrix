package main

import (
    log "github.com/sirupsen/logrus"
    "net/http"
    "time"
)

const GSTATIC_URL = "http://clients3.google.com/generate_204"

// Checks for internet periodically, not returning until connected.
func WaitForConnection() {
    c := 1
    for {
        if ConnectionPresent() {
            log.WithFields(log.Fields{
                "checks": c,
            }).Info("Internet connection present.")
            return
        }
        time.Sleep(1 * time.Second)
        c++
    }
}

// Sanity check for internet access. Not bulletproof but works.
func ConnectionPresent() bool {
    _, err := http.Get(GSTATIC_URL)
    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Debug("Connection failed.")
    }
    return err == nil
}
