package main

import (
	"fmt"
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

// if limit == 0 then the loop runs for infinitely
func timer(limit time.Duration) {
	queues := make(chan termbox.Event)
	go func() {
		for {
			queues <- termbox.PollEvent()
		}
	}()

	lastInput := time.Now()
	ticker := time.NewTicker(time.Second)
	duration := time.Duration(0)
	end := time.Now().Add(limit)
	paused := false

loop:
	for {
		select {
		case ev := <-queues:
			inputTime := time.Now()
			if inputTime.Sub(lastInput) > inputDelay {
				if isQuit(ev) {
					break loop
				} else if ev.Ch == 'r' || ev.Ch == 'R' {
					duration = time.Duration(0)
					clearT()
					putText("Reseting...", positionMiddle, termbox.ColorRed+termbox.AttrBold)
					flush()
				} else if ev.Key == termbox.KeySpace {
					paused = !paused
					if paused {
						clearT()
						putTime(durationToStr(duration))
						putText("Paused", positionButtom, termbox.AttrBold+termbox.ColorRed)
						flush()
					} else {
						end = time.Now().Add(limit - duration)
					}
				}
			}
			lastInput = inputTime
		case <-ticker.C:
			if paused {
				continue loop
			}
			clearT()
			putTime(durationToStr(duration))
			putText( /* "Current time: "+ */ time.Now().Format(timeFormat), positionButtom, termbox.ColorDarkGray+termbox.AttrBold)
			if limit != 0 {
				putText(
					fmt.Sprintf("Timmer ends: %s", end.Format(timeFormat)),
					positionButtomP1,
					termbox.ColorDarkGray+termbox.AttrBold,
				)
			}
			duration += time.Second
			flush()

			if limit != 0 && duration > limit {
				go notify("Time out", fmt.Sprintf("%s is over", limit.String()))
				go playSound(writeNoti())
				time.Sleep(1 * time.Second)
				break loop
			}
		}
	}
}
