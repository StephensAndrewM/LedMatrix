package main

import (
	"image"
	"image/color"
)

type GlyphTestSlide struct {
	Test GlyphTestType
}

type GlyphTestType int

const (
	TEST_LETTERS GlyphTestType = iota
	TEST_NUMSYM
)

func NewGlyphTestSlide(test GlyphTestType) *GlyphTestSlide {
	sl := new(GlyphTestSlide)
	sl.Test = test
	return sl
}

func (sl *GlyphTestSlide) Initialize() {

}

func (sl *GlyphTestSlide) Terminate() {

}

func (sl *GlyphTestSlide) StartDraw(d Display) {
	DrawOnce(d, sl.Draw)
}

func (sl *GlyphTestSlide) StopDraw() {

}

func (sl *GlyphTestSlide) IsEnabled() bool {
	return true // Always enabled
}

func (sl *GlyphTestSlide) Draw(img *image.RGBA) {
	midpoint := GetLeftOfCenterX(img)
	c := color.RGBA{255, 255, 255, 255}
	if sl.Test == TEST_LETTERS {
		WriteString(img, "THE QUICK BROWN FOX", c, ALIGN_CENTER, midpoint, 0)
		WriteString(img, "JUMPS OVER THE LAZY DOG", c, ALIGN_CENTER, midpoint, 8)
		WriteString(img, "the quick brown fox", c, ALIGN_CENTER, midpoint, 16)
		WriteString(img, "jumps over the lazy dog", c, ALIGN_CENTER, midpoint, 24)

	} else if sl.Test == TEST_NUMSYM {
		WriteString(img, "1234567890", c, ALIGN_CENTER, midpoint, 4)
		WriteString(img, "1/2 30° ❤ 6:30", c, ALIGN_CENTER, midpoint, 20)
	}
}
