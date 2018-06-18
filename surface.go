package main

import (
    "encoding/hex"
    "errors"
    "fmt"
)

type Surface struct {
    Width    int
    Height   int
    Midpoint int
    Grid     [][]Color
    glyphs   *GlyphService
}

func NewSurface(width, height int) *Surface {
    s := new(Surface)
    s.Width = width
    s.Midpoint = width / 2
    s.Height = height
    s.Grid = make([][]Color, height)
    for i := range s.Grid {
        s.Grid[i] = make([]Color, width)
    }
    s.glyphs = NewGlyphService()
    return s
}

type Color struct {
    R byte
    G byte
    B byte
}

type Alignment int

const (
    ALIGN_LEFT Alignment = iota
    ALIGN_CENTER
    ALIGN_RIGHT
)

func (s *Surface) GetValue(x, y int) (Color, error) {
    if x < 0 || x >= s.Width || y < 0 || y >= s.Height {
        return Color{}, errors.New("Surface.GetValue out of bounds.")
    }
    return s.Grid[y][x], nil
}

func (s *Surface) SetValue(x, y int, p Color) error {
    // fmt.Printf("Attempting to set (%d,%d) to %s", x, y, p)
    if x < 0 || x >= s.Width || y < 0 || y >= s.Height {
        return errors.New("Surface.SetValue out of bounds.")
    }
    s.Grid[y][x] = p
    return nil
}

func (s *Surface) WriteString(str string, c Color, align Alignment, x int, y int) {
    glyphs := make([]Glyph, len(str))
    width := 0
    for i, char := range str {
        g := s.glyphs.GetGlyph(char)
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
        s.WriteGlyph(g, c, originX+offsetX, y)
        offsetX += g.Width + 1
    }
}

func (s *Surface) WriteStringBoxed(str string, c Color, align Alignment, x int, y int, max int) {
    glyphs := make([]Glyph, len(str))
    width := 0
    for i, char := range str {
        g := s.glyphs.GetGlyph(char)
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
        s.WriteGlyph(g, c, originX+offsetX, y)
        offsetX += g.Width + 1
    }

    // Draw the debug bounding box over the characters
    // aqua := Color{0, 255, 255}
    // s.DrawEmptyBox(aqua, x, y, max, 7)
}

func (s *Surface) WriteGlyph(g Glyph, c Color, x int, y int) {
    for j, row := range g.Layout {
        for i, val := range row {
            if val != 0 {
                s.SetValue(x+int(i), y+int(j), c)
            }
        }
    }
}

func (s *Surface) Clear() {
    blank := Color{0, 0, 0}
    for j := 0; j < s.Height; j++ {
        for i := 0; i < s.Width; i++ {
            s.SetValue(i, j, blank)
        }
    }
}

func (s *Surface) DrawBox(c Color, x int, y int, width int, height int) {
    for j := y; j < y+height; j++ {
        for i := x; i < x+width; i++ {
            s.SetValue(i, j, c)
        }
    }
}

func (s *Surface) DrawEmptyBox(c Color, x int, y int, width int, height int) {
    for j := y; j < y+height; j++ {
        if j == y || j == y+height-1 {
            for i := x; i < x+width; i++ {
                s.SetValue(i, j, c)
            }
        }
        s.SetValue(x, j, c)
        s.SetValue(x+width, j, c)
    }
}

func ColorFromHex(s string) Color {
    rStr := s[0:2]
    r, rErr := hex.DecodeString(rStr)
    gStr := s[2:4]
    g, gErr := hex.DecodeString(gStr)
    bStr := s[4:6]
    b, bErr := hex.DecodeString(bStr)
    if rErr != nil || gErr != nil || bErr != nil {
        fmt.Printf("Error parsing color %s to RGB.")
    }
    return Color{r[0], g[0], b[0]}
}
