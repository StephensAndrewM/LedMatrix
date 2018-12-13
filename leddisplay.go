package main

import (
    "fmt"
    "github.com/mcuadros/go-rpi-rgb-led-matrix"
    "image"
    "image/draw"
)

type LedDisplay struct {
    Matrix rgbmatrix.Matrix
    Canvas *rgbmatrix.Canvas
}

func NewLedDisplay() *LedDisplay {
    d := new(LedDisplay)
    config := &rgbmatrix.DefaultConfig
    config.HardwareMapping = "adafruit-hat"
    config.Rows = 32
    config.Cols = 64
    config.ChainLength = 2
    config.PWMBits = 8
    config.Brightness = 50
    config.ShowRefreshRate = false
    m, err := rgbmatrix.NewRGBLedMatrix(config)
    if err != nil {
        fmt.Println("Could not create hardware LED matrix.")
        return nil
    }
    d.Matrix = m
    d.Canvas = rgbmatrix.NewCanvas(m)
    return d
}

func (d *LedDisplay) Initialize() {

}

func (d *LedDisplay) Redraw(img *image.RGBA) {
    draw.Draw(d.Canvas, d.Canvas.Bounds(), img, image.ZP, draw.Src)
    d.Canvas.Render()
}
