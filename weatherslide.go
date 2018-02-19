package main

type WeatherSlide struct {
}

func NewWeatherSlide() *WeatherSlide {
    sl := new(WeatherSlide)
    return sl
}

func (sl WeatherSlide) Preload() {

}

func (sl WeatherSlide) Draw(s *Surface) {
    s.Clear()
    white := Color{255, 255, 255}
    green := Color{0, 255, 0}
    yellow := Color{255, 255, 0}
    s.WriteString("NOW:", white, ALIGN_LEFT, 0, 0)
    s.WriteString("65°", green, ALIGN_RIGHT, s.Width-1, 0)
    s.WriteString("CLEAR", green, ALIGN_RIGHT, s.Width-1, 8)
    s.WriteString("TMRW:", white, ALIGN_LEFT, 0, 16)
    s.WriteString("60°", yellow, ALIGN_RIGHT, s.Width-1, 16)
    s.WriteString("RAIN", yellow, ALIGN_RIGHT, s.Width-1, 24)
}
