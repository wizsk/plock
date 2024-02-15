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
				}
			}
			lastInput = inputTime
		case <-ticker.C:
			clearT()
			putTime(durationToStr(duration))
			putText(time.Now().Format(timeFormat), positionButtom, termbox.ColorBlue+termbox.AttrBold)
			duration += time.Second
			flush()

			if limit != 0 && duration > limit {
				go notify("Time out", fmt.Sprintf("The timmer set for %s is finished", limit.String()))
				go playSound(writeNoti())
				time.Sleep(1 * time.Second)
				break loop
			}
		}
	}
}
