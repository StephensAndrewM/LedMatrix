package main

type GlyphTestSlide struct {
    Test GlyphTestType
}

type GlyphTestType int

const (
    TEST_LETTERS GlyphTestType = iota
    TEST_NUMSYM
)

func NewGlyphTestSlide(test GlyphTestType) *GlyphTestSlide {
    sl := new(GlyphTestSlide)
    sl.Test = test
    return sl
}

func (sl *GlyphTestSlide) Preload() {

}

func (sl *GlyphTestSlide) Draw(s *Surface) {
    s.Clear()
    c := Color{255, 255, 255}
    if sl.Test == TEST_LETTERS {
        s.WriteString("THE QUICK BROWN FOX", c, ALIGN_CENTER, s.Midpoint, 0)
        s.WriteString("JUMPS OVER THE LAZY DOG", c, ALIGN_CENTER, s.Midpoint, 8)
        s.WriteString("the quick brown fox", c, ALIGN_CENTER, s.Midpoint, 16)
        s.WriteString("jumps over the lazy dog", c, ALIGN_CENTER, s.Midpoint, 24)

    } else if sl.Test == TEST_NUMSYM {
        s.WriteString("1234567890", c, ALIGN_CENTER, s.Midpoint, 4)
        s.WriteString("1/2 30° ❤ 6:30", c, ALIGN_CENTER, s.Midpoint, 20)
    }
}
