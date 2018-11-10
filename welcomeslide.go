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

func (this *WelcomeSlide) Preload() {
    // DON'T preload anything here.
    // This slide gets displayed as soon as the controller starts.
}

func (this *WelcomeSlide) Draw(img *image.RGBA) {
    midpoint := GetLeftOfCenterX(img)
    WriteString(img, "HELLO!", color.RGBA{255, 255, 255, 255}, ALIGN_CENTER, midpoint, 0)
    WriteString(img, "Andrew's Led Matrix", color.RGBA{0, 255, 255, 255}, ALIGN_CENTER, midpoint, 16)
    WriteString(img, "v_1.0", color.RGBA{0, 255, 0, 255}, ALIGN_CENTER, midpoint, 24)
}
