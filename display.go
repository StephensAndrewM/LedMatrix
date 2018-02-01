package main

type Display interface {
	Initialize()
	Redraw(g *PixelGrid)
}