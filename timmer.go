package main

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/nsf/termbox-go"
)

// if limit == 0 then the loop runs for infinitely
func timer(limit time.Duration) {
	queues := make(chan termbox.Event)
	go func() {
		for {
			queues <- termbox.PollEvent()
		}
	}()

	lastInput := time.Now()

	ticker := time.NewTicker(time.Second / time.Duration(fpsFlag)) // 30fps

	duration := atomic.Int64{}
	end := time.Now().Add(limit)
	paused := false

	updateScreenCalled := 0
	var updateScreenFirstCalled time.Time
	var timeElapes int
	updateScreen := func() bool {
		clearT()
		now := time.Now()
		if showFps {
			if updateScreenFirstCalled.IsZero() {
				updateScreenFirstCalled = time.Now()
			}
			updateScreenCalled++
			timeElapes = int(now.Sub(updateScreenFirstCalled).Round(time.Second).Seconds())

			if timeElapes > 0 {
				putText("avg fps: "+strconv.Itoa(updateScreenCalled/timeElapes), positionTop, termbox.ColorRed)
			}
		}
		dur := time.Duration(duration.Load())
		putTime(durationToStr(dur))

		// "Current time: "
		putText(now.Format(timeFormat),
			positionButtom, termbox.ColorDarkGray|termbox.AttrBold)

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

		if limit != 0 && dur >= limit {
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

	timerDone := make(chan struct{})
	timmerTicker := time.NewTicker(time.Second)

	go func(done <-chan struct{}) {
		for {
			select {
			case <-done:
				return
			case <-timmerTicker.C: // when paused it will be stopeed
				duration.Add(int64(time.Second))
			}
		}
	}(timerDone)

loop:
	for {
		select {
		case ev := <-queues:
			inputTime := time.Now()
			if inputTime.Sub(lastInput) > inputDelay {
				if isQuit(ev) {
					timmerTicker.Stop()
					if confirm(queues, "Stop timmer?", true) {
						break loop
					}
					timmerTicker.Reset(time.Second)
				} else if ev.Ch == 'r' || ev.Ch == 'R' {
					duration.Store(0)
					timmerTicker.Reset(time.Second)
					paused = false
					updateScreen()
				} else if ev.Key == termbox.KeySpace {
					paused = !paused
					timmerTicker.Stop()
					if !paused {
						timmerTicker.Reset(time.Second)
						end = time.Now().
							Add(limit - time.Duration(duration.Load()))
					}
					updateScreen()
				}
			}
			lastInput = inputTime

		case <-ticker.C:
			if updateScreen() {
				break loop
			}
		}
	}

	timerDone <- struct{}{}
}
