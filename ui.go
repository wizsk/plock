package main

import (
	"github.com/nsf/termbox-go"
)

// put time in the middle
func putTime(s string) {
	t := toText(s)
	x, y := termbox.Size()
	x, y = x/2-t.width()/2, y/2-t.height()/2 // middle

	lx := x // last x
	for h := 0; h < t.height(); h++ {
		for _, sym := range t {
			for _, r := range sym[h] {
				termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
				x++
			}
		}
		y++
		x = lx
	}
}

func putPaused() {
	x, y := termbox.Size()
	x, y = x/2-pausedText.width()/2, y*3/4
	lx := x
	for _, l := range pausedText {
		for _, r := range l {
			termbox.SetCell(x, y, r, termbox.ColorRed, termbox.ColorDefault)
			x++
		}
		y++
		x = lx
	}
}
