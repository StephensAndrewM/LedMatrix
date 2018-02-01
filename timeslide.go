package main

type TimeSlide struct {

}

func NewTimeSlide() *TimeSlide {
	s := new(TimeSlide)
	return s
}

func (s TimeSlide) Preload() {
	// No preloading needed
}

func (s TimeSlide) IsEnabled() bool {
	return true
}

func (s TimeSlide) Draw() {
	
}

func (s TimeSlide) Redraw() {
	
}