package main

type Slide interface {
	Preload()
	Draw(s *Surface)
}