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
	duration := time.Duration(-1 * time.Second)
	end := time.Now().Add(limit)
	paused := false
	nowTime := time.Now()
	update := func() bool {
		if !paused {
			duration += time.Second

		}
		nowTime = nowTime.Add(time.Second)

		clearT()
		putTime(durationToStr(duration))
		putText( /* "Current time: "+ */ nowTime.Format(timeFormat), positionButtom, termbox.ColorDarkGray+termbox.AttrBold)

		if paused {
			putText("Paused", positionButtomP1,
				termbox.AttrBold+termbox.ColorRed)
		} else if limit != 0 {
			putText(
				fmt.Sprintf("Timmer ends: %s", end.Format(timeFormat)),
				positionButtomP1,
				termbox.ColorDarkGray+termbox.AttrBold,
			)
		}

		flush()

		if limit != 0 && duration >= limit {
			if showNotifications {
				go notify("Time out", fmt.Sprintf("%s is over", limit.String()))
			}
			if !silence {
				go playSound(writeNoti())
			}
			time.Sleep(time.Second)
			return true // break
		}
		return false
	}

loop:
	for {
		select {
		case ev := <-queues:
			inputTime := time.Now()
			if inputTime.Sub(lastInput) > inputDelay {
				if isQuit(ev) {
					if confirm(queues, "Stop timmer?", false) {
						break loop
					}
				} else if ev.Ch == 'r' || ev.Ch == 'R' {
					duration = time.Duration(time.Second * -1)
					clearT()
					putText("Reseting...", positionMiddle, termbox.ColorRed+termbox.AttrBold)
					flush()
				} else if ev.Key == termbox.KeySpace {
					paused = !paused
					if !paused {
						end = time.Now().Add(limit - duration)
					}
					update()
				}
			}
			lastInput = inputTime
		case <-ticker.C:
			if update() {
				break loop
			}
		}
	}
}
