package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"
)

var DISPLAY_TALLIES = false

type StayHomeSlide struct {
}

func NewStayHomeSlide() *StayHomeSlide {
	sl := new(StayHomeSlide)
	return sl
}

func (this *StayHomeSlide) Initialize() {

}

func (this *StayHomeSlide) Terminate() {

}

func (this *StayHomeSlide) StartDraw(d Display) {
	DrawOnce(d, this.Draw)
}

func (this *StayHomeSlide) StopDraw() {

}

func (this *StayHomeSlide) IsEnabled() bool {
	// Only display this slide for interesting days.
	return this.GetDayCount()%10 == 0 ||
		this.GetDayCount()%25 == 0 ||
		this.GetDayCount()%365 == 0
}

func (this *StayHomeSlide) Draw(img *image.RGBA) {
	yellow := color.RGBA{255, 255, 0, 255}
	red := color.RGBA{255, 0, 0, 255}

	diff := this.GetDayCount()

	if DISPLAY_TALLIES {

		for j := 0; j <= diff/50; j++ {
			lineDiff := 50
			if diff < ((j + 1) * 50) {
				lineDiff = diff % 50
			}
			for i := 0; i <= lineDiff/5; i++ {
				block := 5
				if lineDiff < ((i + 1) * 5) {
					block = lineDiff % 5
				}
				x := (i * 13) + 1
				y := j * 8
				for t := 0; t < Min(block, 4); t++ {
					vertLineX := x + ((t * 2) + 1)
					DrawVertLine(img, yellow, y, y+6, vertLineX)
				}
				if block == 5 {
					DrawHorizLine(img, yellow, x, x+2, y+2)
					DrawHorizLine(img, yellow, x+2, x+6, y+3)
					DrawHorizLine(img, yellow, x+6, x+8, y+4)
				}
			}
		}

	} else {

		DrawIcon(img, "house-16", red, 8, 1)
		DrawIcon(img, "house-16", red, (128 - 8 - 16), 1)

		// Draw the number and box centered on the slide
		width := GetDisplayWidth(fmt.Sprintf("%d", diff)) + 7
		DrawEmptyBox(img, yellow, 64-(width/2), 1, width, 13)
		WriteString(img, fmt.Sprintf("%d", diff), yellow, ALIGN_CENTER, 64, 4)

	}

	WriteString(img, "DAYS SINCE", red, ALIGN_CENTER, 64, 16)
	WriteString(img, "OFFICES CLOSED", red, ALIGN_CENTER, 64, 24)
}

func (this *StayHomeSlide) GetDayCount() int {
	start := time.Date(2020, time.March, 10, 0, 0, 0, 0, time.Local)
	return int(math.Ceil(time.Since(start).Hours()/24.0)) - 1
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
