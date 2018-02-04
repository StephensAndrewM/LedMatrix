package main

type TimeSlide struct {

}

func NewTimeSlide() *TimeSlide {
	sl := new(TimeSlide)
	return sl
}

func (sl TimeSlide) Preload() {
	// No preloading needed
}

func (sl TimeSlide) IsEnabled() bool {
	return true
}

func (sl TimeSlide) Draw(s *Surface) {
		
}