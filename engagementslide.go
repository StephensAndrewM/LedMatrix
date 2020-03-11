package main

import (
    "image"
    "image/color"
    "time"
)

type EngagementSlide struct {
    xPos         int
    yPos         int
    xDir         int
    yDir         int
    cIndex       int
    RedrawTicker *time.Ticker
}

func NewEngagementSlide() *EngagementSlide {
    sl := new(EngagementSlide)
    sl.xDir = 1
    sl.yDir = 1
    return sl
}

func (this *EngagementSlide) Initialize() {

}

func (this *EngagementSlide) Terminate() {

}

func (this *EngagementSlide) StartDraw(d Display) {
    this.RedrawTicker = DrawEveryInterval(100*time.Millisecond, d, this.Draw)
}

func (this *EngagementSlide) StopDraw() {
    this.RedrawTicker.Stop()
}

func (this *EngagementSlide) IsEnabled() bool {
    return true
}

func (this *EngagementSlide) Draw(img *image.RGBA) {
    colors := []color.RGBA{
        color.RGBA{255, 255, 255, 255},
        color.RGBA{255, 255, 0, 255},
        color.RGBA{255, 0, 255, 255},
        color.RGBA{0, 255, 255, 255},
        color.RGBA{0, 255, 0, 255},
        color.RGBA{255, 0, 0, 255},
    }

    str := "SHE SAID YES!!!"
    xMax := 128 - GetDisplayWidth(str)
    yMax := 32 - 7

    this.xPos += this.xDir
    this.yPos += this.yDir

    if this.xPos <= 0 || this.xPos >= xMax {
        this.xDir *= -1
    }
    if this.yPos <= 0 || this.yPos >= yMax {
        this.yDir *= -1
    }

    this.cIndex = (this.cIndex + 1) % len(colors)

    WriteString(img, str, colors[this.cIndex], ALIGN_LEFT, this.xPos, this.yPos)
}
