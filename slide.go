package main

type Slide interface {
	Preload()
	IsEnabled() (bool)
	Draw(s *Surface)
}