package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "image"
    "image/color"
    "math"
    "time"
)

type ChristmasSlide struct {
}

func NewChristmasSlide() *ChristmasSlide {
    sl := new(ChristmasSlide)
    return sl
}

func (this *ChristmasSlide) Initialize() {

}

func (this *ChristmasSlide) Terminate() {

}

func (this *ChristmasSlide) StartDraw(d Display) {
    DrawOnce(d, this.Draw)
}

func (this *ChristmasSlide) StopDraw() {

}

func (this *ChristmasSlide) Draw(img *image.RGBA) {
    r := color.RGBA{255, 0, 0, 255}
    g := color.RGBA{0, 255, 0, 255}
    tz, err := time.LoadLocation("America/New_York")
    if err != nil {
        // No idea why this would ever happen
        log.Warn("Could not load time zone.")
        return
    }
    ptoDate := time.Date(2019, time.December, 21, 0, 0, 0, 0, tz)
    ptoDiff := time.Until(ptoDate).Hours() / 24.0

    DrawEmptyBox(img, r, 23, 1, 18, 13)
    WriteString(img, fmt.Sprintf("%d", int(math.Ceil(ptoDiff))), r, ALIGN_CENTER, 32, 4)
    WriteString(img, "DAYS UNTIL", g, ALIGN_CENTER, 32, 16)
    WriteString(img, "VACATION", g, ALIGN_CENTER, 32, 24)

    xmasDate := time.Date(2019, time.December, 25, 0, 0, 0, 0, tz)
    xmasDiff := time.Until(xmasDate).Hours() / 24.0

    DrawEmptyBox(img, r, 87, 1, 18, 13)
    WriteString(img, fmt.Sprintf("%d", int(math.Ceil(xmasDiff))), r, ALIGN_CENTER, 96, 4)
    WriteString(img, "DAYS UNTIL", g, ALIGN_CENTER, 96, 16)
    WriteString(img, "CHRISTMAS", g, ALIGN_CENTER, 96, 24)
}
