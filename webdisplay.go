package main

import (
    "encoding/json"
    "github.com/gorilla/websocket"
    log "github.com/sirupsen/logrus"
    "image"
    "net/http"
)

type WebDisplay struct {
    Conn *websocket.Conn
}

func NewWebDisplay() *WebDisplay {
    d := new(WebDisplay)
    return d
}

func (d *WebDisplay) Initialize() {
    upgrader := websocket.Upgrader{}

    http.Handle("/", http.FileServer(http.Dir("public_html/")))
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        c, err := upgrader.Upgrade(w, r, nil)
        log.Debug("Socket connected.")
        if err != nil {
            log.WithFields(log.Fields{
                "error": err,
            }).Warn("Socket upgrade error.")
            return
        }
        d.Conn = c
    })
    go http.ListenAndServe(":8000", nil)
    log.Info("Started web display.")
}

func (d *WebDisplay) Redraw(img *image.RGBA) {
    width := img.Bounds().Dx()
    height := img.Bounds().Dy()

    // Convert to 2D array of 3-tuples
    data := make([][][]int, height)
    for j := 0; j < height; j++ {
        data[j] = make([][]int, width)
        for i := 0; i < width; i++ {
            data[j][i] = make([]int, 3)
            rgba := img.RGBAAt(i, j)
            data[j][i][0] = int(rgba.R)
            data[j][i][1] = int(rgba.G)
            data[j][i][2] = int(rgba.B)
        }
    }

    if d.Conn != nil {
        json, err := json.Marshal(data)
        if err != nil {
            log.WithFields(log.Fields{
                "error": err,
            }).Warn("Could not serialize JSON for WebDisplay.")
            return
        }
        d.Conn.WriteMessage(websocket.TextMessage, json)
    }
}
