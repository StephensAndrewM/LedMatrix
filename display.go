package main

import (
    "image"
)

type Display interface {
	Initialize()
	Redraw(img *image.RGBA)
}