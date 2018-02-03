package main

type Display interface {
	Initialize()
	Redraw(s *Surface)
}