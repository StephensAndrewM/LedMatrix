package main

import (
    "encoding/hex"
    log "github.com/sirupsen/logrus"
    "image"
    "image/color"
    "math"
    "time"
)

type Alignment int

const (
    ALIGN_LEFT Alignment = iota
    ALIGN_CENTER
    ALIGN_RIGHT
)

func NewBlankImage() *image.RGBA {
    img := image.NewRGBA(image.Rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT))
    DrawBox(img, color.RGBA{0, 0, 0, 255}, 0, 0, SCREEN_WIDTH, SCREEN_HEIGHT)
    return img
}

func DrawEverySecond(d Display, drawFn func(*image.RGBA)) *time.Ticker {
    return DrawEveryInterval(1*time.Second, d, drawFn)
}

func DrawEveryInterval(interval time.Duration, d Display, drawFn func(*image.RGBA)) *time.Ticker {
    // First draw right now to avoid lag time
    DrawOnce(d, drawFn)
    // Then set up the period redraw
    t := time.NewTicker(interval)
    go func() {
        for range t.C {
            DrawOnce(d, drawFn)
        }
    }()
    return t
}

func DrawOnce(d Display, drawFn func(*image.RGBA)) {
    img := NewBlankImage()
    drawFn(img)
    d.Redraw(img)
}

func WriteString(img *image.RGBA, str string, c color.RGBA, align Alignment, x int, y int) {
    WriteStringBoxed(img, str, c, align, x, y, 0)
}

func WriteStringBoxed(img *image.RGBA, str string, c color.RGBA, align Alignment, x int, y int, max int) {
    // This shouldn't happen, but is an indicator to just not draw anything
    if max < 0 {
        return
    }

    glyphs := make([]Glyph, len(str))
    width := 0
    for i, char := range str {
        g := GetGlyph(char)
        width += g.Width + 1
        // If we exceed how much the box can hold, stop
        if max > 0 && width > max {
            break
        }
        glyphs[i] = g
    }
    // Remove the kerning on the last letter
    width--

    var originX int
    switch align {
    case ALIGN_LEFT:
        originX = x
    case ALIGN_RIGHT:
        originX = x - width + 1
    case ALIGN_CENTER:
        originX = x - (width / 2)
    }

    offsetX := 0
    for _, g := range glyphs {
        WriteGlyph(img, g, c, originX+offsetX, y)
        offsetX += g.Width + 1
    }

    // Draw the debug bounding box over the characters
    if DEBUG_DRAW {
        aqua := color.RGBA{0, 255, 255, 255}
        if max > 0 {
            // Display the cutoff point of the text
            DrawEmptyBox(img, aqua, originX, y, max-1, 7)
        } else {
            // Otherwise, just display the width of the text
            DrawEmptyBox(img, aqua, originX, y, width-1, 7)
        }
    }
}

func WriteGlyph(img *image.RGBA, g Glyph, c color.RGBA, x int, y int) {
    for j, row := range g.Layout {
        for i, val := range row {
            if val != 0 {
                img.SetRGBA(x+int(i), y+int(j), c)
            }
        }
    }
}

func DrawIcon(img *image.RGBA, iconName string, c color.RGBA, x int, y int) {
    icon := GetIcon(iconName)
    for j, row := range icon.Layout {
        for i, val := range row {
            if val != 0 {
                img.SetRGBA(x+int(i), y+int(j), c)
            }
        }
    }
}

func DrawBox(img *image.RGBA, c color.RGBA, x int, y int, width int, height int) {
    for j := y; j < y+height; j++ {
        for i := x; i < x+width; i++ {
            img.SetRGBA(i, j, c)
        }
    }
}

func DrawEmptyBox(img *image.RGBA, c color.RGBA, x int, y int, width int, height int) {
    for j := y; j < y+height; j++ {
        if j == y || j == y+height-1 {
            for i := x; i < x+width; i++ {
                img.SetRGBA(i, j, c)
            }
        }
        img.SetRGBA(x, j, c)
        img.SetRGBA(x+width, j, c)
    }
}

func DrawHorizLine(img *image.RGBA, c color.RGBA, x1 int, x2 int, y int) {
    for i := x1; i <= x2; i++ {
        img.SetRGBA(i, y, c)
    }
}

func DrawVertLine(img *image.RGBA, c color.RGBA, y1 int, y2 int, x int) {
    for i := y1; i <= y2; i++ {
        img.SetRGBA(x, i, c)
    }
}

func DrawError(img *image.RGBA, slideName, error string) {
    white := color.RGBA{255, 255, 255, 255}
    yellow := color.RGBA{255, 255, 0, 255}
    WriteString(img, slideName, white, ALIGN_CENTER, SCREEN_WIDTH/2, 8)
    WriteString(img, error, yellow, ALIGN_CENTER, SCREEN_WIDTH/2, 16)
}

// Map black pixels to given color, all other colors to transparent
func DrawImageWithColorTransform(canvas *image.RGBA, source *image.RGBA, xOffset int, yOffset int, c color.RGBA) {
    black := color.RGBA{0, 0, 0, 255}
    for j := 0; j < source.Bounds().Dy(); j++ {
        for i := 0; i < source.Bounds().Dx(); i++ {
            if source.At(i, j) == black {
                canvas.SetRGBA(i+xOffset, j+yOffset, c)
            }
        }
    }
}

func GetDisplayWidth(str string) int {
    width := 0
    for _, c := range str {
        g := GetGlyph(c)
        width += g.Width + 1
    }
    // Remove the kerning on the last letter
    width--
    return width
}

func ColorFromHex(s string) color.RGBA {
    rStr := s[0:2]
    r, rErr := hex.DecodeString(rStr)
    gStr := s[2:4]
    g, gErr := hex.DecodeString(gStr)
    bStr := s[4:6]
    b, bErr := hex.DecodeString(bStr)
    if rErr != nil || gErr != nil || bErr != nil {
        log.Warn("Error parsing color %s to RGB.", s)
        return color.RGBA{0, 0, 0, 255}
    }
    return color.RGBA{r[0], g[0], b[0], 255}
}

func ReduceColor(c color.RGBA) color.RGBA {
    round := func(x uint8) uint8 {
        if x > 128 {
            return 255
        } else {
            return 0
        }
    }
    return color.RGBA{round(c.R), round(c.G), round(c.B), 0}
}

func GetLeftOfCenterX(img *image.RGBA) int {
    return img.Bounds().Dx() / 2
}

func DrawAutoNormalizedGraph(img *image.RGBA, x, y, h int, c color.RGBA, data []float64) {
    min := data[0]
    max := data[0]
    for _, val := range data {
        if val < min {
            min = val
        }
        if val > max {
            max = val
        }
    }
    DrawNormalizedGraph(img, x, y, h, min, max, c, data)
}

func DrawNormalizedGraph(img *image.RGBA, x, y, h int, min, max float64, c color.RGBA, data []float64) {
    dataRange := max - min
    var normalized []int
    for _, val := range data {
        normVal := int(Round(((val - min) / dataRange) * float64(h)))
        normalized = append(normalized, normVal)
        // fmt.Printf("%.2f --> %d [%.2f, %.2f]\n", val, normVal, min, max)
    }
    for i, val := range normalized {
        // Zero-check and +1 needed so we don't draw a point for a zero value
        if val > 0 {
            DrawVertLine(img, c, y-val+1, y, x+i)
        }
    }
}

// This is needed for compatibility with Go 1.9
// What kind of language doesn't implement this???
func Round(x float64) float64 {
    t := math.Trunc(x)
    if math.Abs(x-t) >= 0.5 {
        return t + math.Copysign(1, x)
    }
    return t
}
