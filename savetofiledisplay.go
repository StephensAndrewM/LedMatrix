package main

import (
    "fmt"
    "github.com/fogleman/gg"
    log "github.com/sirupsen/logrus"
    "image"
    "image/color"
    "reflect"
)

var DRAW_GRIDLINES = true
var RENDER_SCALE = 8
var DOT_PADDING = 0.75
var MIN_BRIGHTNESS = uint8(40)

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
    // Define the height of the drawing canvas, in real pixels
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
            dc.SetColor(this.FloorColor(img.RGBAAt(i, j)))
            dc.Fill()
        }
    }

    if DRAW_GRIDLINES {
        dc.SetRGB(0, 1.0, 1.0)

        // Draw major center lines
        dc.SetLineWidth(2.0)
        dc.DrawLine(0, float64(dcHeight)/2.0, float64(dcWidth), float64(dcHeight)/2.0)
        dc.DrawLine(float64(dcWidth)/2.0, 0, float64(dcWidth)/2.0, float64(dcHeight))

        // Draw minor 8-dot grid lines
        dc.SetLineWidth(0.5)
        for j := 8; j < SCREEN_HEIGHT; j += 8 {
            dc.DrawLine(0, float64(j*RENDER_SCALE), float64(dcWidth), float64(j*RENDER_SCALE))
        }
        for i := 8; i < SCREEN_WIDTH; i += 8 {
            dc.DrawLine(float64(i*RENDER_SCALE), 0, float64(i*RENDER_SCALE), float64(dcHeight))
        }

        // Put the lines on the canvas
        dc.Stroke()

        // Draw dot counts
        if err := dc.LoadFontFace("/usr/share/fonts/truetype/ubuntu/UbuntuMono-Regular.ttf", 12); err != nil {
           panic(err)
        }
        for i := 0; i < SCREEN_WIDTH; i += 8 {
            dc.DrawString(fmt.Sprintf("%d", i), float64(i*RENDER_SCALE), float64(8))
        }

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

// Set a minimum (gray) RGB value if none is provided
func (this *SaveToFileDisplay) FloorColor(c color.RGBA) color.RGBA {
    r := this.Max(c.R, MIN_BRIGHTNESS)
    g := this.Max(c.G, MIN_BRIGHTNESS)
    b := this.Max(c.B, MIN_BRIGHTNESS)
    return color.RGBA{r, g, b, c.A}
}

func (this *SaveToFileDisplay) Max(a, b uint8) uint8 {
    if a > b {
        return a
    }
    return b
}
