package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

func clock() {
	queues := make(chan termbox.Event)
	go func() {
		for {
			queues <- termbox.PollEvent()
		}
	}()

	lastInput := time.Now()
	ticker := time.NewTicker(time.Second)

loop:
	for {
		select {
		case ev := <-queues:
			inputTime := time.Now()

			if inputTime.Sub(lastInput) > inputDelay && isQuit(ev) {
				break loop
			}
			lastInput = inputTime
		case t := <-ticker.C:
			clearT()
			putTime(t.Format(timeFormat))
			flush()
		}

	}
}
