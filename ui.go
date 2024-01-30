package main

import (
	"unicode/utf8"

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

func putMsg(m symbol, color termbox.Attribute) {
	x, y := termbox.Size()
	x, y = x/2-m.width()/2, y*3/4
	lx := x
	for _, l := range m {
		for _, r := range l {
			termbox.SetCell(x, y, r, color, termbox.ColorDefault)
			x++
		}
		y++
		x = lx
	}
}

type position uint8

const (
	positionTop = iota
	positionMiddle
	positionButtom
)

func putText(t string, p position, color termbox.Attribute) {
	x, y := termbox.Size()
	x = x/2 - utf8.RuneCountInString(t)/2
	switch p {
	case positionTop:
		y /= 7
	case positionMiddle:
		y /= 2
	case positionButtom:
		y = y * 3 / 4
	default:
		return
	}

	for _, r := range t {
		termbox.SetCell(x, y, r, color, termbox.ColorDefault)
		x++
	}
}
