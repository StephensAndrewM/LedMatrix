package main

// import "fmt"

type Glyph struct {
	Character rune
	Width int
	Layout [][]uint8
}

type GlyphService struct {
	Glpyhs map[rune]Glyph
	Unknown Glyph
}

func (s *GlyphService) Register(c rune, layout [][]uint8) {
	g := Glyph{}
	g.Character = c
	g.Width = len(layout[0])
	g.Layout = layout
	s.Glpyhs[c] = g
}

func NewGlyphService() *GlyphService {

	s := new(GlyphService)
	s.Glpyhs = make(map[rune]Glyph)

	s.Register('a', [][]uint8{
			{0,1,1,1,0},
			{1,0,0,0,1},
			{1,0,0,0,1},
			{1,1,1,1,1},
			{1,0,0,0,1},
			{1,0,0,0,1},
			{1,0,0,0,1}})

	// Manually set up unknown glyph (checkerboard)
	unknown := Glyph {
		Width: 5,
		Layout: [][]uint8{
			{1,0,1,0,1},
			{0,1,0,1,0},
			{1,0,1,0,1},
			{0,1,0,1,0},
			{1,0,1,0,1},
			{0,1,0,1,0},
			{1,0,1,0,1},
			{0,1,0,1,0}}}
	s.Unknown = unknown

	return s
	
}

func (s *GlyphService) GetGlyph(char rune) Glyph {
	glyph, ok := s.Glpyhs[char]
	// fmt.Println("GlyphService.GetGlyph attempting to find " + string(char))
	if !ok {
		glyph = s.Unknown
	}
	return glyph
}