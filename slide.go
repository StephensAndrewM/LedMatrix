package main

type Slide interface {
    // Called when slideshow is being started
    Initialize()
    // Called when slideshow is being stopped
    Terminate()
    // Called when slide is brought into view
    StartDraw(d Display)
    // Called when a different slide is brought into view
    StopDraw()
    // Controls whether slide will be skipped in slideshow
    IsEnabled() bool
}
