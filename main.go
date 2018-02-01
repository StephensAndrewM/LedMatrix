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
	g := NewPixelGrid(32,32)
	p := Pixel{0,255,0}
	Formatter.WriteString(g, p, "abc", 1, 1)
	fmt.Println("RedrawDisplay")
	// fmt.Println(g)
	d.Redraw(g)
}

func UpdateLoop() {
	time.Sleep(5 * time.Second)
}