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
			putTime(t.Format("03:04:05 PM"))
			flush()
		}

	}
}

func timer() {
	queues := make(chan termbox.Event)
	go func() {
		for {
			queues <- termbox.PollEvent()
		}
	}()

	lastInput := time.Now()
	ticker := time.NewTicker(time.Second)
	duration := time.Duration(0)

loop:
	for {
		select {
		case ev := <-queues:
			inputTime := time.Now()

			if inputTime.Sub(lastInput) > inputDelay && isQuit(ev) {
				break loop
			}
			lastInput = inputTime
		case <-ticker.C:
			clearT()
			putTime(durationToStr(duration))
			putText(time.Now().Format("03:04:05 PM"), positionButtom, termbox.ColorBlue+termbox.AttrBold)
			duration += time.Second
			flush()
		}
	}
}
