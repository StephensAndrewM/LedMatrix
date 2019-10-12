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

func (this *WelcomeSlide) Initialize() {
    // This won't ever get called since this slide isn't in the main rotation.
}

func (this *WelcomeSlide) Draw(img *image.RGBA) {
    midpoint := GetLeftOfCenterX(img)
    WriteString(img, "HELLO!", color.RGBA{255, 255, 255, 255}, ALIGN_CENTER, midpoint, 2)
    WriteString(img, "Andrew's Led Matrix", color.RGBA{0, 255, 255, 255}, ALIGN_CENTER, midpoint, 15)
    WriteString(img, "v_1.2", color.RGBA{0, 255, 0, 255}, ALIGN_CENTER, midpoint, 23)
}
