package main

import (
    "encoding/hex"
    "fmt"
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
    glyphs := make([]Glyph, len(str))
    width := 0
    for i, char := range str {
        g := GetGlyph(char)
        width += g.Width + 1
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
}

func WriteStringBoxed(img *image.RGBA, str string, c color.RGBA, align Alignment, x int, y int, max int) {
    glyphs := make([]Glyph, len(str))
    width := 0
    for i, char := range str {
        g := GetGlyph(char)
        width += g.Width + 1
        // If we exceed how much the box can hold, stop
        if width > max {
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
        DrawEmptyBox(img, aqua, x, y, max, 7)
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

func DrawError(img *image.RGBA, space int, code int) {
    yellow := color.RGBA{255, 255, 0, 255}
    msg := fmt.Sprintf("E #%02d-%02d", space, code)
    WriteString(img, msg, yellow, ALIGN_LEFT, 0, 0)
}

func ColorFromHex(s string) color.RGBA {
    rStr := s[0:2]
    r, rErr := hex.DecodeString(rStr)
    gStr := s[2:4]
    g, gErr := hex.DecodeString(gStr)
    bStr := s[4:6]
    b, bErr := hex.DecodeString(bStr)
    if rErr != nil || gErr != nil || bErr != nil {
        fmt.Printf("Error parsing color %s to RGB.")
    }
    return color.RGBA{r[0], g[0], b[0], 255}
}

func GetLeftOfCenterX(img *image.RGBA) int {
    return img.Bounds().Dx() / 2
}