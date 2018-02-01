package main

type Slide interface {
	Preload()
	IsEnabled() (bool)
	Draw(d Display)
	Redraw(d Display)
}