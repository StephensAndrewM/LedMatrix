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

func (f *FormatterDef) WriteGlyph(s *Surface, c Color, g Glyph, x uint64, y uint64) {
	for j,row := range g.Layout {
		for i,val := range row {
			if val != 0 {
				s.SetValue(x+uint64(i),y+uint64(j), c)
			}
		}
	}
}

func (f *FormatterDef) WriteString(s *Surface, c Color, str string, x uint64, y uint64) {
	offset := 0
	for _,char := range str {
		glyph, ok := f.Glpyhs[char]
		if !ok {
			glyph = f.Unknown
		}
		f.WriteGlyph(s, c, glyph, x + uint64(offset), y)
		offset += glyph.Width
	}
}

var Formatter = NewFormatter()