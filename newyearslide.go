package main

import (
    "fmt"
    "image"
    "image/color"
    "math"
    "time"
)

type NewYearSlide struct {
    Midnight  time.Time
    Fireworks []*Firework

    RedrawTicker *time.Ticker
}

const FPS = 4.0

func NewNewYearSlide() *NewYearSlide {
    sl := new(NewYearSlide)
    return sl
}

func (this *NewYearSlide) Initialize() {
    t := time.Now()
    year := t.Year() + 1
    // If it's January, the new year just passed so we want to count to the
    // current year instead (and show zeros).
    if t.Month() == time.January {
        year = t.Year()
    }
    this.Midnight = time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
}

func (this *NewYearSlide) Terminate() {

}

func (this *NewYearSlide) StartDraw(d Display) {
    this.RedrawTicker = DrawEveryInterval((1000/FPS)*time.Millisecond, d, this.Draw)
    this.Fireworks = []*Firework{
        this.createFirework(10, 8, 255, 0, 0),
        this.createFirework(27, 12, 255, 255, 0),
        this.createFirework(112, 6, 0, 255, 255),
        this.createFirework(103, 10, 255, 0, 255),
    }
}

func (this *NewYearSlide) StopDraw() {
    this.RedrawTicker.Stop()
}

func (this *NewYearSlide) IsEnabled() bool {
    diff := time.Until(this.Midnight)
    return diff > (-1 * time.Hour)
}

func (this *NewYearSlide) Draw(img *image.RGBA) {
    c0 := color.RGBA{0, 255, 255, 255}
    c1 := color.RGBA{255, 255, 255, 255}
    c2 := color.RGBA{0, 255, 0, 255}

    diff := time.Until(this.Midnight)
    if diff < 0 {
        diff = 0
    }

    WriteString(img, this.fmtDuration(diff), c0, ALIGN_CENTER, GetLeftOfCenterX(img), 4)
    WriteString(img, "UNTIL", c1, ALIGN_CENTER, GetLeftOfCenterX(img), 14)
    WriteString(img, fmt.Sprintf("%d", this.Midnight.Year()), c2, ALIGN_CENTER, GetLeftOfCenterX(img)-1, 23)

    for _, f := range this.Fireworks {
        f.Draw(img)
    }
}

func (this *NewYearSlide) fmtDuration(d time.Duration) string {
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    d -= m * time.Minute
    s := d / time.Second
    return fmt.Sprintf("%02d : %02d : %02d", h, m, s)
}

type Firework struct {
    x        int
    y        int
    embers   []*FireworkEmber
    color    color.RGBA
    hasBurst bool
}

func (this *NewYearSlide) createFirework(x, y int, r, g, b uint8) *Firework {
    yellow := color.RGBA{255, 255, 0, 255}

    return &Firework{
        x: x,
        y: y,
        embers: []*FireworkEmber{
            &FireworkEmber{
                x:      float64(x),
                y:      float64(y + SCREEN_HEIGHT),
                xspeed: 0,
                yspeed: -7.5,
                color:  yellow,
            },
        },
        color:    color.RGBA{r, g, b, 255},
        hasBurst: false,
    }
}

func (this *Firework) Draw(img *image.RGBA) {
    // Apply speed first, since we take different action based on that.
    for _, e := range this.embers {
        e.applyPhysics()
    }

    // If it's still flying upwards, check if it's reached expected height.
    if !this.hasBurst {
        if int(this.embers[0].y) <= this.y {
            this.hasBurst = true
            this.initEmbers()
        }
    }

    // Always draw whatever is stored
    for _, e := range this.embers {
        e.draw(img)
    }
}

func (this *Firework) initEmbers() {
    this.embers = nil
    for v := 0.4; v <= 1.6; v += 0.4 {
        embersInRing := 16
        if v < 0.5 {
            embersInRing = 8
        }
        for i := 0; i < embersInRing; i++ {
            angle := (float64(i) / (float64(embersInRing) / 2)) * math.Pi
            yspeed := (math.Cos(angle) * v)
            xspeed := math.Sin(angle) * v

            fmt.Printf("New ember position is %.2f,%.2f, speed is %.2f,%.2f\n", this.x, this.y, xspeed, yspeed)

            this.embers = append(this.embers, &FireworkEmber{
                x:      float64(this.x),
                y:      float64(this.y),
                xspeed: xspeed,
                yspeed: yspeed,
                color:  this.color,
            })
        }
    }
}

type FireworkEmber struct {
    x      float64
    y      float64
    xspeed float64
    yspeed float64
    color  color.RGBA
}

func (this *FireworkEmber) draw(img *image.RGBA) {
    img.SetRGBA(int(math.Round(this.x)), int(math.Round(this.y)), this.color)
}

func (this *FireworkEmber) applyPhysics() {
    this.x += this.xspeed
    this.y += this.yspeed
    this.yspeed += (0.4 / FPS) // gravity
    this.xspeed *= 0.96        // friction
}
