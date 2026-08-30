package main

import (
	"fmt"
	"math"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/nsf/termbox-go"
)

// if limit == 0 then the loop runs for infinitely
func timer(limit time.Duration, down bool) {
	queues := make(chan termbox.Event)
	go func() {
		for {
			queues <- termbox.PollEvent()
		}
	}()

	lastInput := time.Now()

	ticker := time.NewTicker(time.Second / time.Duration(fpsFlag)) // 30fps

	extralines := 4
	if limit == 0 {
		extralines = 2
	}

	duration := atomic.Int64{}
	if down {
		duration.Store(int64(limit))
	}
	end := time.Now().Add(limit)
	paused := false

	updateScreenCalled := 0
	var updateScreenFirstCalled time.Time
	var timeElapes int

	timmerDone := false
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

		nextY := putTime(durationToStr(dur), extralines) + 1

		// "Current time: "
		// putText(now.Format(timeFormat),
		// 	positionButtom, termbox.ColorDarkGray|termbox.AttrBold)

		if limit == 0 {
			// nextY++
		} else {
			tsx, _ := termbox.Size()

			s := 100

			if s > tsx {
				s = tsx - 10
			}

			y := nextY
			nextY += 2

			x := tsx/2 - s/2
			du := float64(dur)
			lim := float64(limit)

			if down {
				du = lim - du
			}

			doneUntil := int(math.Round(du / lim * float64(s)))
			// logFa(time.Duration(duration.Load()), limit)

			const doneChar = '█'
			const notDoneChar = '░'
			const doneColor = termbox.ColorWhite
			const notDoneColor = termbox.ColorDarkGray
			for range doneUntil {
				termbox.SetCell(x, y, doneChar, doneColor, termbox.ColorDefault)
				x++
			}
			for range s - doneUntil {
				termbox.SetCell(x, y, notDoneChar, notDoneColor, termbox.ColorDefault)
				x++
			}
		}

		if paused {
			putTextY("Paused", nextY,
				termbox.AttrBold+termbox.ColorRed)
		} else if limit != 0 {
			putTextY(
				fmt.Sprintf("%s | Ends: %s", now.Format(timeFormat), end.Format(timeFormat)),
				nextY,
				termbox.ColorDarkGray+termbox.AttrBold,
			)
		} else {
			putTextY(
				now.Format(timeFormat),
				nextY,
				termbox.ColorDarkGray+termbox.AttrBold,
			)
		}

		flush()

		if down && dur <= 0 || !down && limit != 0 && dur >= limit {
			timmerDone = true
			if showNotifications {
				go notify("Time out", fmt.Sprintf("%s is over", limit.String()))
			}
			if !silence {
				go playSound(writeNoti())
			}
			time.Sleep(time.Second)
			return true // break
		}
		return timmerDone
	}

	timerDone := make(chan struct{})
	timmerTicker := time.NewTicker(time.Second)

	go func(done <-chan struct{}) {
		for {
			select {
			case <-done:
				return
			case <-timmerTicker.C: // when paused it will be stopeed
				if down {
					duration.Add(-int64(time.Second))
				} else {
					duration.Add(int64(time.Second))
				}
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
				} else if ev.Ch == 'r' || ev.Ch == 'ق' {
					if confirm(queues, "Reset timmer?", true) {
						timmerTicker.Stop()

						if down {
							duration.Store(int64(limit))
						} else {
							duration.Store(0)
						}

						timmerTicker.Reset(time.Second)
						paused = false
						updateScreen()
					}
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
