package main

import (
    "fmt"
    "github.com/fogleman/gg"
    log "github.com/sirupsen/logrus"
    "image"
    "reflect"
)

var DRAW_GRIDLINES = false
var RENDER_SCALE = 8
var DOT_PADDING = 0.75

type SaveToFileDisplay struct {
    SlideId string
}

func NewSaveToFileDisplay() *SaveToFileDisplay {
    d := new(SaveToFileDisplay)
    return d
}

func (this *SaveToFileDisplay) Initialize() {

}

func (this *SaveToFileDisplay) Redraw(img *image.RGBA) {
    dcWidth := SCREEN_WIDTH * RENDER_SCALE
    dcHeight := SCREEN_HEIGHT * RENDER_SCALE

    dc := gg.NewContext(dcWidth, dcHeight)
    dc.DrawRectangle(0, 0, float64(dcWidth), float64(dcHeight))
    dc.SetRGB(0, 0, 0)
    dc.Fill()

    // Draw main LED circles
    for j := 0; j < SCREEN_HEIGHT; j++ {
        for i := 0; i < SCREEN_WIDTH; i++ {
            dc.DrawCircle(
                (float64(i)+0.5)*float64(RENDER_SCALE),
                (float64(j)+0.5)*float64(RENDER_SCALE),
                (float64(RENDER_SCALE)*DOT_PADDING)/2)
            dc.SetColor(img.RGBAAt(i, j))
            dc.Fill()
        }
    }

    // Draw gridlines
    if DRAW_GRIDLINES {
        dc.SetRGB(0, 1.0, 1.0)
        dc.SetLineWidth(2.0)
        dc.DrawLine(0, float64(dcHeight)/2.0, float64(dcWidth), float64(dcHeight)/2.0)
        dc.Stroke()
        dc.DrawLine(float64(dcWidth)/2.0, 0, float64(dcWidth)/2.0, float64(dcHeight))
        dc.Stroke()
    }

    filename := fmt.Sprintf("render/%s.png", this.SlideId)
    err := dc.SavePNG(filename)
    if err != nil {
        log.Fatal(err)
    }

    log.WithFields(log.Fields{
        "file": filename,
    }).Info("Saved rendering of slide.")
}

func (this *SaveToFileDisplay) SetSlideId(s Slide) {
    this.SlideId = reflect.TypeOf(s).Elem().Name()
}
