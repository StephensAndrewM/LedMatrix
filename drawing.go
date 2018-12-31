package main

import (
    "encoding/hex"
    "fmt"
    log "github.com/sirupsen/logrus"
    "image"
    "image/color"
)

type Alignment int

const (
    ALIGN_LEFT Alignment = iota
    ALIGN_CENTER
    ALIGN_RIGHT
)

func WriteString(img *image.RGBA, str string, c color.RGBA, align Alignment, x int, y int) {
    WriteStringBoxed(img, str, c, align, x, y, 0)
}

func WriteStringBoxed(img *image.RGBA, str string, c color.RGBA, align Alignment, x int, y int, max int) {
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

func DrawError(img *image.RGBA, space int, code int) {
    yellow := color.RGBA{255, 255, 0, 255}
    msg := fmt.Sprintf("E #%02d-%02d", space, code)
    WriteString(img, msg, yellow, ALIGN_LEFT, 0, 0)
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
