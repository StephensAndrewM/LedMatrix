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
	Formatter.WriteString(s, c, "abc", 1, 1)
	fmt.Println("RedrawDisplay")
	d.Redraw(s)
}

func UpdateLoop() {
	time.Sleep(5 * time.Second)
}