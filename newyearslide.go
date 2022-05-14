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

func (sl *NewYearSlide) Initialize() {
	t := time.Now()
	year := t.Year() + 1
	// If it's January, the new year just passed so we want to count to the
	// current year instead (and show zeros).
	if t.Month() == time.January {
		year = t.Year()
	}
	sl.Midnight = time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
}

func (sl *NewYearSlide) Terminate() {

}

func (sl *NewYearSlide) StartDraw(d Display) {
	sl.RedrawTicker = DrawEveryInterval((1000/FPS)*time.Millisecond, d, sl.Draw)
	sl.Fireworks = []*Firework{
		sl.createFirework(10, 8, 255, 0, 0),
		sl.createFirework(27, 12, 255, 255, 0),
		sl.createFirework(112, 6, 0, 255, 255),
		sl.createFirework(103, 10, 255, 0, 255),
	}
}

func (sl *NewYearSlide) StopDraw() {
	sl.RedrawTicker.Stop()
}

func (sl *NewYearSlide) IsEnabled() bool {
	diff := time.Until(sl.Midnight)
	return diff > (-1 * time.Hour)
}

func (sl *NewYearSlide) Draw(img *image.RGBA) {
	c0 := color.RGBA{0, 255, 255, 255}
	c1 := color.RGBA{255, 255, 255, 255}
	c2 := color.RGBA{0, 255, 0, 255}

	diff := time.Until(sl.Midnight)
	if diff < 0 {
		diff = 0
	}

	WriteString(img, sl.fmtDuration(diff), c0, ALIGN_CENTER, GetLeftOfCenterX(img), 4)
	WriteString(img, "UNTIL", c1, ALIGN_CENTER, GetLeftOfCenterX(img), 14)
	WriteString(img, fmt.Sprintf("%d", sl.Midnight.Year()), c2, ALIGN_CENTER, GetLeftOfCenterX(img)-1, 23)

	for _, f := range sl.Fireworks {
		f.Draw(img)
	}
}

func (sl *NewYearSlide) fmtDuration(d time.Duration) string {
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

func (sl *NewYearSlide) createFirework(x, y int, r, g, b uint8) *Firework {
	yellow := color.RGBA{255, 255, 0, 255}

	return &Firework{
		x: x,
		y: y,
		embers: []*FireworkEmber{
			{
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

func (sl *Firework) Draw(img *image.RGBA) {
	// Apply speed first, since we take different action based on that.
	for _, e := range sl.embers {
		e.applyPhysics()
	}

	// If it's still flying upwards, check if it's reached expected height.
	if !sl.hasBurst {
		if int(sl.embers[0].y) <= sl.y {
			sl.hasBurst = true
			sl.initEmbers()
		}
	}

	// Always draw whatever is stored
	for _, e := range sl.embers {
		e.draw(img)
	}
}

func (sl *Firework) initEmbers() {
	sl.embers = nil
	for v := 0.4; v <= 1.6; v += 0.4 {
		embersInRing := 16
		if v < 0.5 {
			embersInRing = 8
		}
		for i := 0; i < embersInRing; i++ {
			angle := (float64(i) / (float64(embersInRing) / 2)) * math.Pi
			yspeed := (math.Cos(angle) * v)
			xspeed := math.Sin(angle) * v

			sl.embers = append(sl.embers, &FireworkEmber{
				x:      float64(sl.x),
				y:      float64(sl.y),
				xspeed: xspeed,
				yspeed: yspeed,
				color:  sl.color,
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

func (sl *FireworkEmber) draw(img *image.RGBA) {
	img.SetRGBA(int(math.Round(sl.x)), int(math.Round(sl.y)), sl.color)
}

func (sl *FireworkEmber) applyPhysics() {
	sl.x += sl.xspeed
	sl.y += sl.yspeed
	sl.yspeed += (0.4 / FPS) // gravity
	sl.xspeed *= 0.96        // friction
}
