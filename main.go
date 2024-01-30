package main

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

type app struct {
	remaing time.Time
	queues  chan termbox.Event
}

func durToString(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

func main() {
	dur, _ := time.ParseDuration("1h10m10s")
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()
	queues := make(chan termbox.Event)
	go func() {
		for {
			queues <- termbox.PollEvent()
		}
	}()

	paused := false
	ticker := time.NewTicker(time.Second)

loop:
	for {
		select {
		case ev := <-queues:
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC || ev.Ch == 'q' {
				break loop
			} else if ev.Key == termbox.KeySpace {
				paused = !paused
				if paused {
					termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
					putTime(durToString(dur))
					putPaused()
					termbox.Flush()
					ticker.Stop()
				} else {
					ticker.Reset(time.Second)
				}
			}
		case <-ticker.C:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			putTime(durToString(dur))
			dur -= time.Second
			if dur == 0 {
				break loop
			}

			termbox.Flush()
		}

	}

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
