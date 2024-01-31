package main

import (
	"github.com/nsf/termbox-go"
)

func main() {
	// initializing the terminal
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	runPomodoro()

}

/*
	queues := make(chan termbox.Event)
	go func() {
		for {
			queues <- termbox.PollEvent()
		}
	}()

	paused := false

loop:
	for {
		select {
		case ev := <-queues:
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC || ev.Ch == 'q' {
				break loop
			} else if ev.Key == termbox.KeySpace {
				paused = !paused
			}
		case t := <-time.Tick(time.Microsecond):
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			if paused {
				putPaused()
			} else {
				putTime(t.Format("03:04:05 PM"))
			}
			termbox.Flush()
		}

	}
*/
