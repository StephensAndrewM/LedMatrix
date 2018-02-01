package main

type Glyph struct {
	Character rune
	Width int
	Layout [][]uint8
}

type FormatterDef struct {
	Glpyhs map[rune]Glyph
	Unknown Glyph
}

func (f *FormatterDef) InitGlyph(c rune, g Glyph) {
	g.Character = c
	g.Width = len(g.Layout[0])
	f.Glpyhs[c] = g
}

func NewFormatter() *FormatterDef {

	f := new(FormatterDef)
	f.Glpyhs = make(map[rune]Glyph)

	a := Glyph {
		Layout: [][]uint8{
			{0, 1, 1, 1, 0},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 1, 1, 1, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1}}}
	f.InitGlyph('a', a)

	unknown := Glyph {
		Width: 5,
		Layout: [][]uint8{
			{1, 0, 1, 0, 1},
			{0, 1, 0, 1, 0},
			{1, 0, 1, 0, 1},
			{0, 1, 0, 1, 0},
			{1, 0, 1, 0, 1},
			{0, 1, 0, 1, 0},
			{1, 0, 1, 0, 1},
			{0, 1, 0, 1, 0}}}
	f.Unknown = unknown

	return f
	
}

func (f *FormatterDef) WriteGlyph(pg *PixelGrid, p Pixel, g Glyph, x uint64, y uint64) {
	for j,row := range g.Layout {
		for i,val := range row {
			if val != 0 {
				pg.SetValue(x+uint64(i),y+uint64(j), p)
			}
		}
	}
}

func (f *FormatterDef) WriteString(pg *PixelGrid, p Pixel, s string, x uint64, y uint64) {
	offset := 0
	for _,c := range s {
		g, ok := f.Glpyhs[c]
		if !ok {
			g = f.Unknown
		}
		f.WriteGlyph(pg, p, g, x + uint64(offset), y)
		offset += g.Width
	}
}

var Formatter = NewFormatter()