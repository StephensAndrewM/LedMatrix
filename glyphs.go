package main

import (
	log "github.com/sirupsen/logrus"
)

type Glyph struct {
	Character rune
	Width     int
	Layout    [][]uint8
}

var glyphSet map[rune]Glyph

func RegisterGlyph(c rune, layout [][]uint8) {
	g := Glyph{}
	g.Character = c
	g.Width = len(layout[0])
	g.Layout = layout
	glyphSet[c] = g
}

func InitGlyphs() {

	// Initialize the map
	glyphSet = make(map[rune]Glyph)

	// Uppercase Letters
	RegisterGlyph('A', [][]uint8{
		{0, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 1, 1, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1}})
	RegisterGlyph('B', [][]uint8{
		{1, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 1, 1, 0}})
	RegisterGlyph('C', [][]uint8{
		{0, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 0}})
	RegisterGlyph('D', [][]uint8{
		{1, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 1, 1, 0}})
	RegisterGlyph('E', [][]uint8{
		{1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{1, 1, 1, 1, 0},
		{1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{1, 1, 1, 1, 1}})
	RegisterGlyph('F', [][]uint8{
		{1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{1, 1, 1, 1, 0},
		{1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0}})
	RegisterGlyph('G', [][]uint8{
		{0, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 0},
		{1, 0, 1, 1, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 0}})
	RegisterGlyph('H', [][]uint8{
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 1, 1, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1}})
	RegisterGlyph('I', [][]uint8{
		{1, 1, 1},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{1, 1, 1}})
	RegisterGlyph('J', [][]uint8{
		{0, 1, 1, 1},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{1, 0, 1, 0},
		{0, 1, 0, 0}})
	RegisterGlyph('K', [][]uint8{
		{1, 0, 0, 0, 1},
		{1, 0, 0, 1, 0},
		{1, 0, 1, 0, 0},
		{1, 1, 0, 0, 0},
		{1, 0, 1, 0, 0},
		{1, 0, 0, 1, 0},
		{1, 0, 0, 0, 1}})
	RegisterGlyph('L', [][]uint8{
		{1, 0, 0, 0},
		{1, 0, 0, 0},
		{1, 0, 0, 0},
		{1, 0, 0, 0},
		{1, 0, 0, 0},
		{1, 0, 0, 0},
		{1, 1, 1, 1}})
	RegisterGlyph('M', [][]uint8{
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 0, 1, 1},
		{1, 0, 1, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1}})
	RegisterGlyph('N', [][]uint8{
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 0, 0, 1},
		{1, 0, 1, 0, 1},
		{1, 0, 0, 1, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1}})
	RegisterGlyph('O', [][]uint8{
		{0, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 0}})
	RegisterGlyph('P', [][]uint8{
		{1, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 1, 1, 0},
		{1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0}})
	RegisterGlyph('Q', [][]uint8{
		{0, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 1, 0, 1},
		{1, 0, 0, 1, 0},
		{0, 1, 1, 0, 1}})
	RegisterGlyph('R', [][]uint8{
		{1, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 1, 1, 0},
		{1, 0, 1, 0, 0},
		{1, 0, 0, 1, 0},
		{1, 0, 0, 0, 1}})
	RegisterGlyph('S', [][]uint8{
		{0, 1, 1, 1, 1},
		{1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{0, 1, 1, 1, 0},
		{0, 0, 0, 0, 1},
		{0, 0, 0, 0, 1},
		{1, 1, 1, 1, 0}})
	RegisterGlyph('T', [][]uint8{
		{1, 1, 1, 1, 1},
		{0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0}})
	RegisterGlyph('U', [][]uint8{
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 0}})
	RegisterGlyph('V', [][]uint8{
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 0, 1, 0},
		{0, 0, 1, 0, 0}})
	RegisterGlyph('W', [][]uint8{
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1},
		{0, 1, 0, 1, 0}})
	RegisterGlyph('X', [][]uint8{
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 0, 1, 0},
		{0, 0, 1, 0, 0},
		{0, 1, 0, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1}})
	RegisterGlyph('Y', [][]uint8{
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 0, 1, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0}})
	RegisterGlyph('Z', [][]uint8{
		{1, 1, 1, 1, 1},
		{0, 0, 0, 0, 1},
		{0, 0, 0, 1, 0},
		{0, 0, 1, 0, 0},
		{0, 1, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{1, 1, 1, 1, 1}})

	// Lowercase Letters
	RegisterGlyph('a', [][]uint8{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 1},
		{0, 1, 1, 1},
		{1, 0, 0, 1},
		{0, 1, 1, 1}})
	RegisterGlyph('b', [][]uint8{
		{1, 0, 0, 0},
		{1, 0, 0, 0},
		{1, 1, 1, 0},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{1, 1, 1, 0}})
	RegisterGlyph('c', [][]uint8{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 1, 1, 0},
		{1, 0, 0, 0},
		{1, 0, 0, 0},
		{1, 0, 0, 1},
		{0, 1, 1, 0}})
	RegisterGlyph('d', [][]uint8{
		{0, 0, 0, 1},
		{0, 0, 0, 1},
		{0, 1, 1, 1},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{0, 1, 1, 1}})
	RegisterGlyph('e', [][]uint8{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 1, 1, 0},
		{1, 0, 0, 1},
		{1, 1, 1, 1},
		{1, 0, 0, 0},
		{0, 1, 1, 0}})
	RegisterGlyph('f', [][]uint8{
		{0, 0, 1},
		{0, 1, 0},
		{0, 1, 0},
		{1, 1, 1},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0}})
	RegisterGlyph('g', [][]uint8{
		{0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 1},
		{0, 0, 0, 0, 1},
		{0, 1, 1, 1, 0}})
	RegisterGlyph('h', [][]uint8{
		{1, 0, 0, 0},
		{1, 0, 0, 0},
		{1, 0, 1, 0},
		{1, 1, 0, 1},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{1, 0, 0, 1}})
	RegisterGlyph('i', [][]uint8{
		{0, 1, 0},
		{0, 0, 0},
		{1, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{1, 1, 1}})
	RegisterGlyph('j', [][]uint8{
		{0, 0, 0, 1},
		{0, 0, 0, 0},
		{0, 0, 1, 1},
		{0, 0, 0, 1},
		{0, 0, 0, 1},
		{1, 0, 0, 1},
		{0, 1, 1, 0}})
	RegisterGlyph('k', [][]uint8{
		{1, 0, 0, 0},
		{1, 0, 0, 0},
		{1, 0, 0, 1},
		{1, 0, 1, 0},
		{1, 1, 0, 0},
		{1, 0, 1, 0},
		{1, 0, 0, 1}})
	RegisterGlyph('l', [][]uint8{
		{1, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{1, 1, 1}})
	RegisterGlyph('m', [][]uint8{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{1, 1, 0, 1, 0},
		{1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1}})
	RegisterGlyph('n', [][]uint8{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{1, 1, 1, 0},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{1, 0, 0, 1}})
	RegisterGlyph('o', [][]uint8{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 1, 1, 0},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{0, 1, 1, 0}})
	RegisterGlyph('p', [][]uint8{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{1, 1, 1, 0},
		{1, 0, 0, 1},
		{1, 1, 1, 0},
		{1, 0, 0, 0},
		{1, 0, 0, 0}})
	RegisterGlyph('q', [][]uint8{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 1, 1, 1},
		{1, 0, 0, 1},
		{0, 1, 1, 1},
		{0, 0, 0, 1},
		{0, 0, 0, 1}})
	RegisterGlyph('r', [][]uint8{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{1, 1, 1, 0},
		{1, 0, 0, 1},
		{1, 0, 0, 0},
		{1, 0, 0, 0},
		{1, 0, 0, 0}})
	RegisterGlyph('s', [][]uint8{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 1, 1, 1},
		{1, 0, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 1},
		{1, 1, 1, 0}})
	RegisterGlyph('t', [][]uint8{
		{0, 1, 0},
		{0, 1, 0},
		{1, 1, 1},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 1}})
	RegisterGlyph('u', [][]uint8{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{0, 1, 1, 1}})
	RegisterGlyph('v', [][]uint8{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 0, 1, 0},
		{0, 0, 1, 0, 0}})
	RegisterGlyph('w', [][]uint8{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1},
		{0, 1, 0, 1, 0}})
	RegisterGlyph('x', [][]uint8{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{1, 0, 0, 0, 1},
		{0, 1, 0, 1, 0},
		{0, 0, 1, 0, 0},
		{0, 1, 0, 1, 0},
		{1, 0, 0, 0, 1}})
	RegisterGlyph('y', [][]uint8{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 1},
		{0, 0, 0, 0, 1},
		{0, 1, 1, 1, 0}})
	RegisterGlyph('z', [][]uint8{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1},
		{0, 0, 0, 1, 0},
		{0, 0, 1, 0, 0},
		{0, 1, 0, 0, 0},
		{1, 1, 1, 1, 1}})

	// Numbers
	RegisterGlyph('0', [][]uint8{
		{0, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 0}})
	RegisterGlyph('1', [][]uint8{
		{0, 1, 0},
		{1, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{1, 1, 1}})
	RegisterGlyph('2', [][]uint8{
		{0, 2, 2, 2, 0},
		{2, 0, 0, 0, 2},
		{0, 0, 0, 0, 2},
		{0, 0, 0, 2, 0},
		{0, 0, 2, 0, 0},
		{0, 2, 0, 0, 0},
		{2, 2, 2, 2, 2}})
	RegisterGlyph('3', [][]uint8{
		{0, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{0, 0, 0, 0, 1},
		{0, 0, 1, 1, 0},
		{0, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 0}})
	RegisterGlyph('4', [][]uint8{
		{0, 0, 0, 1, 0},
		{0, 0, 1, 1, 0},
		{0, 1, 0, 1, 0},
		{1, 0, 0, 1, 0},
		{1, 1, 1, 1, 1},
		{0, 0, 0, 1, 0},
		{0, 0, 0, 1, 0}})
	RegisterGlyph('5', [][]uint8{
		{1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{1, 1, 1, 1, 0},
		{0, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 0}})
	RegisterGlyph('6', [][]uint8{
		{0, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 0},
		{1, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 0}})
	RegisterGlyph('7', [][]uint8{
		{1, 1, 1, 1, 1},
		{0, 0, 0, 0, 1},
		{0, 0, 0, 1, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0}})
	RegisterGlyph('8', [][]uint8{
		{0, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 0}})
	RegisterGlyph('9', [][]uint8{
		{0, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 1},
		{0, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{0, 1, 1, 1, 0}})

	// Misc Symbols
	RegisterGlyph('/', [][]uint8{
		{0, 0, 1},
		{0, 0, 1},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{1, 0, 0},
		{1, 0, 0}})
	RegisterGlyph(':', [][]uint8{
		{0},
		{0},
		{1},
		{0},
		{1},
		{0},
		{0}})
	RegisterGlyph('°', [][]uint8{
		{1, 1},
		{1, 1},
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 0}})
	RegisterGlyph('❤', [][]uint8{
		{0, 1, 1, 0, 1, 1, 0},
		{1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1},
		{0, 1, 1, 1, 1, 1, 0},
		{0, 0, 1, 1, 1, 0, 0},
		{0, 0, 0, 1, 0, 0, 0}})
	RegisterGlyph('.', [][]uint8{
		{0},
		{0},
		{0},
		{0},
		{0},
		{0},
		{1}})
	RegisterGlyph('#', [][]uint8{
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{1, 1, 1, 1, 1},
		{0, 1, 0, 1, 0},
		{1, 1, 1, 1, 1},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0}})
	RegisterGlyph('%', [][]uint8{
		{1, 1, 0, 0, 1, 0, 0},
		{1, 1, 0, 0, 1, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 1, 0, 0, 1, 1},
		{0, 0, 1, 0, 0, 1, 1}})
	RegisterGlyph('Δ', [][]uint8{
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 1, 0, 1, 0, 0},
		{0, 0, 1, 0, 1, 0, 0},
		{0, 1, 0, 0, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0},
		{1, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1}})
	RegisterGlyph('-', [][]uint8{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{1, 1, 1},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0}})
	RegisterGlyph('\'', [][]uint8{
		{1},
		{1},
		{0},
		{0},
		{0},
		{0},
		{0}})
	RegisterGlyph('!', [][]uint8{
		{1},
		{1},
		{1},
		{1},
		{1},
		{0},
		{1}})
	RegisterGlyph(',', [][]uint8{
		{0},
		{0},
		{0},
		{0},
		{0},
		{0},
		{1},
		{1}})
	RegisterGlyph('+', [][]uint8{
		{0, 0, 0},
		{0, 0, 0},
		{0, 1, 0},
		{1, 1, 1},
		{0, 1, 0},
		{0, 0, 0},
		{0, 0, 0}})
	RegisterGlyph('?', [][]uint8{
		{0, 1, 1, 1, 0},
		{1, 0, 0, 0, 1},
		{0, 0, 0, 0, 1},
		{0, 0, 0, 1, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0}})
	// Thermometer
	RegisterGlyph('🌡', [][]uint8{
		{0, 1, 1, 0},
		{0, 1, 1, 0},
		{0, 1, 1, 0},
		{0, 1, 1, 0},
		{0, 1, 1, 0},
		{1, 0, 0, 1},
		{0, 1, 1, 0}})
	// Raindrops
	RegisterGlyph('💧', [][]uint8{
		{0, 0, 0, 0},
		{0, 1, 0, 1},
		{0, 1, 0, 1},
		{1, 0, 0, 1},
		{1, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 0}})
	// Underscore represents a short space
	RegisterGlyph('_', [][]uint8{
		{},
		{},
		{},
		{},
		{},
		{},
		{}})
	RegisterGlyph(' ', [][]uint8{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0}})
	// U+FFFD is used when the glyph doesn't exist
	RegisterGlyph('�', [][]uint8{
		{1, 0, 1, 0, 1},
		{0, 1, 0, 1, 0},
		{1, 0, 1, 0, 1},
		{0, 1, 0, 1, 0},
		{1, 0, 1, 0, 1},
		{0, 1, 0, 1, 0},
		{1, 0, 1, 0, 1},
		{0, 1, 0, 1, 0}})
}

func GetGlyph(char rune) Glyph {
	glyph, ok := glyphSet[char]
	if !ok {
		glyph, ok = glyphSet['�']
		if !ok {
			log.Error("Could not load fallback character.")
		}
	}
	return glyph
}
