package main

import "github.com/nsf/termbox-go"

func confirm(q <-chan termbox.Event, msg string, ySelected bool) bool {
	if len(msg) == 0 {
		return false
	}

	const yt = "Yes"
	const nt = "No"
	const ntOffset = 10
	const yDistence = 2 // from the text to yes no

	// for the msg
	var of, oft int
	xOffset := len(msg) / 2
	if len(yt)+len(nt)+ntOffset > len(msg) {
		of = ((len(yt) + len(nt) + ntOffset) / 2) - (len(msg) / 2) - 1 // here -1 is for the padding used in x
		xOffset = (len(yt) + len(nt) + ntOffset) / 2
	} else {
		oft = len(msg)/2 - ((len(yt) + len(nt) + ntOffset) / 2)
	}

loop:
	for {
		clearT()
		tw, th := termbox.Size()
		x, y := (tw/2)-(xOffset), (th/2)-yDistence

		for i, r := range msg {
			termbox.SetChar(x+i+of, y+0, r)
		}

		// add the offset for the yes or no text
		x += oft
		y += yDistence
		if ySelected {
			for _, v := range [2]int{-1, 0} {
				termbox.SetBg(x+v, y, termbox.ColorLightRed)
			}
		}
		for i, r := range yt {
			if ySelected {
				termbox.SetBg(x+i+1, y, termbox.ColorLightRed)
			}
			termbox.SetChar(x+i, y, r)
		}

		if !ySelected {
			for _, v := range [3]int{-1, 0, len(nt)} {
				termbox.SetBg(x+ntOffset+v, y, termbox.ColorLightGreen)
			}
		}

		for i, r := range nt {
			if !ySelected {
				termbox.SetCell(x+ntOffset+i, y, r, termbox.ColorBlack, termbox.ColorLightGreen)
			} else {
				termbox.SetChar(x+ntOffset+i, y, r)
			}
		}

		flush()

		select {
		case ev := <-q:
			if /* ev.Ch == 'q' || */ ev.Key == termbox.KeyEnter || ev.Key == termbox.KeyCtrlC {
				clearT()
				break loop
			} else if ev.Type == termbox.EventResize {
				continue loop
			} else {
				ySelected = !ySelected
			}
		}
	}

	return ySelected
}
