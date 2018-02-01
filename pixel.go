package main

import "fmt"

type PixelGrid struct {
	Width uint64
	Height uint64
	Grid [][]Pixel
}

func NewPixelGrid(width, height uint64) *PixelGrid {
	g := new(PixelGrid)
	g.Width = width
	g.Height = height
	g.Grid = make([][]Pixel, height)
	for i := range g.Grid {
	    g.Grid[i] = make([]Pixel, width)
	}
	return g
}

func (g *PixelGrid) GetValue(x, y uint64) Pixel {
	return g.Grid[y][x]
}

func (g *PixelGrid) SetValue(x, y uint64, p Pixel) {
	g.Grid[y][x] = p
	fmt.Println("SetValue")
	// fmt.Println(g)
}

type Pixel struct {
	R byte
	G byte
	B byte
}