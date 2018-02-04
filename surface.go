package main

import (
    "errors"
)

type Surface struct {
    Width int
    Height int
    Grid [][]Color
    glyphs *GlyphService
}

func NewSurface(width, height int) *Surface {
    s := new(Surface)
    s.Width = width
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
const(
    ALIGN_LEFT Alignment = iota
    ALIGN_CENTER
    ALIGN_RIGHT
)

func (s *Surface) GetValue(x, y int) (Color,error) {
    if x < 0 || x >= s.Width || y < 0 || y >= s.Height  {
        return Color{}, errors.New("Surface.GetValue out of bounds.")
    }
    return s.Grid[y][x], nil
}

func (s *Surface) SetValue(x, y int, p Color) error {
    // fmt.Printf("Attempting to set (%d,%d) to %s", x, y, p)
    if x < 0 || x >= s.Width || y < 0 || y >= s.Height  {
        return errors.New("Surface.SetValue out of bounds.")
    }
    s.Grid[y][x] = p
    return nil
}

func (s *Surface) WriteString(str string, c Color, align Alignment, x int, y int) {
    glyphs := make([]Glyph, len(str))
    width := 0
    for i,char := range str {
        g := s.glyphs.GetGlyph(char)
        width += g.Width
        glyphs[i] = g
    }

    var originX int
    switch(align) {
    case ALIGN_LEFT:
        originX = x
    case ALIGN_RIGHT:
        originX = x - width + 1
    case ALIGN_CENTER:
        originX = x - (width / 2)
    }

    offsetX := 0
    for _,g := range glyphs {
        s.WriteGlyph(g, c, originX + offsetX, y)
        offsetX += g.Width
    }
}

func (s *Surface) WriteGlyph(g Glyph, c Color, x int, y int) {
    for j,row := range g.Layout {
        for i,val := range row {
            if val != 0 {
                s.SetValue(x + int(i), y + int(j), c)
            }
        }
    }
}