package main

import (
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"sync"
)

const R1_PIN = "3"
const G1_PIN = "8"
const B1_PIN = "5"
const R2_PIN = "7"
const G2_PIN = "10"
const B2_PIN = "11"
const A_PIN = "13"
const B_PIN = "12"
const C_PIN = "15"
const D_PIN = "16"
const CLK_PIN = "19"
const OE_PIN = "18"
const LAT_PIN = "21"

type LedDisplay struct {
	Pins         *LedDisplayPins
	Mutex        *sync.Mutex
	SavedSurface *Surface
}

type LedDisplayPins struct {
	R1  *gpio.DirectPinDriver
	G1  *gpio.DirectPinDriver
	B1  *gpio.DirectPinDriver
	R2  *gpio.DirectPinDriver
	G2  *gpio.DirectPinDriver
	B2  *gpio.DirectPinDriver
	A   *gpio.DirectPinDriver
	B   *gpio.DirectPinDriver
	C   *gpio.DirectPinDriver
	D   *gpio.DirectPinDriver
	CLK *gpio.DirectPinDriver
	OE  *gpio.DirectPinDriver
	LAT *gpio.DirectPinDriver
}

func NewLedDisplay() *LedDisplay {
	this := new(LedDisplay)
	this.Pins = new(LedDisplayPins)
	this.Mutex = &sync.Mutex{}
	return this
}

func (this *LedDisplay) Initialize() {
	r := raspi.NewAdaptor()
	this.Pins.R1 = gpio.NewDirectPinDriver(r, "3")
	this.Pins.G1 = gpio.NewDirectPinDriver(r, "8")
	this.Pins.B1 = gpio.NewDirectPinDriver(r, "5")
	this.Pins.R2 = gpio.NewDirectPinDriver(r, "7")
	this.Pins.G2 = gpio.NewDirectPinDriver(r, "10")
	this.Pins.B2 = gpio.NewDirectPinDriver(r, "11")
	this.Pins.A = gpio.NewDirectPinDriver(r, "13")
	this.Pins.B = gpio.NewDirectPinDriver(r, "12")
	this.Pins.C = gpio.NewDirectPinDriver(r, "15")
	this.Pins.D = gpio.NewDirectPinDriver(r, "16")
	this.Pins.CLK = gpio.NewDirectPinDriver(r, "19")
	this.Pins.OE = gpio.NewDirectPinDriver(r, "18")
	this.Pins.LAT = gpio.NewDirectPinDriver(r, "21")

	go this.MainLoop()
}

func (this *LedDisplay) MainLoop() {

	s := this.SavedSurface

	// Loop forever
	for {

		for j := 0; j < s.Height/2; j += 2 {
			for i := 0; i < s.Width; i++ {

				// Get pixel values
				p1, p1Err := s.GetValue(i, j)
				p2, p2Err := s.GetValue(i, j+1)
				if p1Err != nil || p2Err != nil {
					// TODO log error
					return
				}

				// Write the color values to pins
				this.Mutex.Lock()
				this.Pins.CLK.On()
				this.Set(this.Pins.R1, this.ToBinaryColor(p1.R))
				this.Set(this.Pins.G1, this.ToBinaryColor(p1.G))
				this.Set(this.Pins.B1, this.ToBinaryColor(p1.B))
				this.Set(this.Pins.R2, this.ToBinaryColor(p2.R))
				this.Set(this.Pins.G2, this.ToBinaryColor(p2.G))
				this.Set(this.Pins.B2, this.ToBinaryColor(p2.B))
				this.Pins.CLK.Off()
				this.Mutex.Unlock()

				// Need to latch every 32 cols
				if (j+1)%32 == 0 {
					this.Pins.CLK.On()
					this.Pins.LAT.On()
					this.Pins.CLK.Off()
					this.Pins.LAT.Off()
				}
			}
			// Briefly enable output after every row
			this.Pins.CLK.On()
			this.Pins.OE.On()
			this.Pins.CLK.Off()
			this.Pins.OE.Off()
		}
	}
}

func (this *LedDisplay) Redraw(s *Surface) {
	// If this is first time, initialize the struct
	if this.SavedSurface == nil {
		this.SavedSurface = NewSurface(s.Width, s.Height)
	}

	// Copy the pixel values
	this.Mutex.Lock()
	for j := 0; j < s.Height; j++ {
		for i := 0; i < s.Width; i++ {
			c, err := s.GetValue(i, j)
			if err != nil {
				// TODO something with the error
				return
			}
			this.SavedSurface.SetValue(i, j, c)
		}
	}
	this.Mutex.Unlock()

}

// Turn 0 = off, anything else = on
func (this *LedDisplay) ToBinaryColor(c byte) bool {
	return c > 0
}

func (this *LedDisplay) Set(pin *gpio.DirectPinDriver, on bool) {
	if on {
		pin.On()
	} else {
		pin.Off()
	}
}
