package main

import (
    "fmt"
    "image"
    "image/color"
    "math"
    "math/rand"
    "time"
)

var TreeDef = [][]int{
    {10, 10},
    {9, 11},
    {9, 11},
    {8, 12},
    {8, 12},
    {8, 12},
    {7, 13},
    {7, 13},
    {6, 14},
    {6, 14},
    {6, 14},
    {5, 15},
    {5, 15},
    {4, 16},
    {4, 16},
    {4, 16},
    {3, 17},
    {3, 17},
    {2, 18},
    {2, 18},
    {2, 18},
    {1, 19},
    {1, 19},
    {0, 20},
    {0, 20},
}

type ChristmasSlide struct {
    XmasDate     time.Time
    RedrawTicker *time.Ticker
}

func NewChristmasSlide() *ChristmasSlide {
    sl := new(ChristmasSlide)
    return sl
}

func (this *ChristmasSlide) Initialize() {
    t := time.Now()
    this.XmasDate = time.Date(t.Year(), time.December, 25, 0, 0, 0, 0, time.Local)
}

func (this *ChristmasSlide) Terminate() {

}

func (this *ChristmasSlide) StartDraw(d Display) {
    this.RedrawTicker = DrawEverySecond(d, this.Draw)
}

func (this *ChristmasSlide) StopDraw() {
    this.RedrawTicker.Stop()
}

func (this *ChristmasSlide) IsEnabled() bool {
    return this.DaysUntil(this.XmasDate) >= 0 && this.DaysUntil(this.XmasDate) <= 30
}

func (this *ChristmasSlide) Draw(img *image.RGBA) {
    red := color.RGBA{255, 0, 0, 255}
    green := color.RGBA{0, 255, 0, 255}
    darkgreen := color.RGBA{0, 128, 0, 255}
    yellow := color.RGBA{255, 255, 0, 255}
    brown := color.RGBA{255, 128, 0, 255}
    aqua := color.RGBA{0, 255, 255, 255}
    blue := color.RGBA{0, 0, 255, 255}
    orange := color.RGBA{255, 220, 0, 255}

    lights := []color.RGBA{
        aqua, aqua,
        red, red,
        green, green,
        blue, blue,
        orange, orange,
    }

    treeOffsetX := 18
    treeOffsetY := 2

    // Draw the star
    img.SetRGBA(treeOffsetX+10, treeOffsetY-1, yellow)
    // Draw the tree body
    for j, line := range TreeDef {
        DrawHorizLine(img, darkgreen, treeOffsetX+line[0], treeOffsetX+line[1], treeOffsetY+j)
    }
    // Draw the stump
    DrawBox(img, brown, treeOffsetX+9, treeOffsetY+25, 3, 5)

    // Draw some sparkles
    for _, c := range lights {
        x,y := this.GetRandomWithinTree()
        img.SetRGBA(treeOffsetX+x, treeOffsetY+y, c)
    }

    countdownOffsetX := 82
    DrawEmptyBox(img, red, countdownOffsetX-9, 1, 18, 13)
    days := this.DaysUntil(this.XmasDate)
    if days < 0 {
        days = 0
    }
    WriteString(img, fmt.Sprintf("%d", days), red, ALIGN_CENTER, countdownOffsetX, 4)
    WriteString(img, "DAYS UNTIL", green, ALIGN_CENTER, countdownOffsetX, 16)
    WriteString(img, "CHRISTMAS", green, ALIGN_CENTER, countdownOffsetX, 24)
}

func (this *ChristmasSlide) GetRandomWithinTree() (int, int) {
    for {
        x := rand.Intn(21)
        y := rand.Intn(24) + 1 // Don't select top line
        if x > TreeDef[y][0] && x < TreeDef[y][1] {
            return x, y
        }
    }
}

func (this *ChristmasSlide) DaysUntil(d time.Time) int {
    diff := time.Until(d).Hours() / 24.0
    return int(math.Ceil(diff))
}
