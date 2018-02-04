package main

import (
	"fmt"
	"time"
)

func main() {
	d := NewWebDisplay()
	d.Initialize()
	for {
		time.Sleep(5 * time.Second)
		RedrawDisplay(d)
	}
}

func RedrawDisplay(d Display) {
	s := NewSurface(32,32)
	c := Color{0,255,0}
	s.WriteString("ABCDEF", c, ALIGN_LEFT, 0, 0)
	s.WriteString("12345", c, ALIGN_LEFT, 0, 8)
	s.WriteString("67890", c, ALIGN_LEFT, 0, 16)
	s.WriteString("❤:_/", c, ALIGN_LEFT, 0, 24)
	fmt.Println("RedrawDisplay")
	d.Redraw(s)
}

func UpdateLoop() {
	time.Sleep(5 * time.Second)
}