package main

import (
	"image"
	"image/color"
)

type WelcomeSlide struct {
}

func NewWelcomeSlide() *WelcomeSlide {
	sl := new(WelcomeSlide)
	return sl
}

func (sl *WelcomeSlide) Initialize() {
	// sl.won't ever get called since sl.slide isn't in the main rotation.
}

func (sl *WelcomeSlide) Terminate() {
	// sl.won't ever get called since sl.slide isn't in the main rotation.
}

func (sl *WelcomeSlide) StartDraw(d Display) {
	DrawOnce(d, sl.Draw)
}

func (sl *WelcomeSlide) StopDraw() {

}

func (sl *WelcomeSlide) IsEnabled() bool {
	return true // Always enabled
}

func (sl *WelcomeSlide) Draw(img *image.RGBA) {
	midpoint := GetLeftOfCenterX(img)
	WriteString(img, "HELLO!", color.RGBA{255, 255, 0, 255}, ALIGN_CENTER, midpoint, 2)
	WriteString(img, "Andrew's Led Matrix", color.RGBA{0, 255, 255, 255}, ALIGN_CENTER, midpoint, 16)
}
