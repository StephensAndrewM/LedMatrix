package main

type Surface struct {
	Width uint64
	Height uint64
	Grid [][]Color
}

func NewSurface(width, height uint64) *Surface {
	g := new(Surface)
	g.Width = width
	g.Height = height
	g.Grid = make([][]Color, height)
	for i := range g.Grid {
	    g.Grid[i] = make([]Color, width)
	}
	return g
}

func (g *Surface) GetValue(x, y uint64) Color {
	return g.Grid[y][x]
}

func (g *Surface) SetValue(x, y uint64, p Color) {
	g.Grid[y][x] = p
}

type Color struct {
	R byte
	G byte
	B byte
}