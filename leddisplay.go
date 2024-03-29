package main

import (
	"image"
	"image/draw"

	rgbmatrix "github.com/mcuadros/go-rpi-rgb-led-matrix"
	log "github.com/sirupsen/logrus"
)

type LedDisplay struct {
	Matrix rgbmatrix.Matrix
	Canvas *rgbmatrix.Canvas
}

func NewLedDisplay() *LedDisplay {
	d := new(LedDisplay)
	config := &rgbmatrix.DefaultConfig
	config.HardwareMapping = "adafruit-hat-pwm"
	config.Rows = 32
	config.Cols = 64
	config.ChainLength = 2
	config.PWMBits = 11
	config.Brightness = 50
	config.ShowRefreshRate = false
	config.PWMLSBNanoseconds = 250
	m, err := rgbmatrix.NewRGBLedMatrix(config)
	if err != nil {
		log.Error("Could not create hardware LED matrix.")
		return nil
	}
	d.Matrix = m
	d.Canvas = rgbmatrix.NewCanvas(m)
	return d
}

func (d *LedDisplay) Initialize() {

}

func (d *LedDisplay) Redraw(img *image.RGBA) {
	draw.Draw(d.Canvas, d.Canvas.Bounds(), img, image.Point{}, draw.Src)
	d.Canvas.Render()
}
