package main

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

type app struct {
	session uint
	// remaining time
	current time.Duration
	timmer  time.Duration
	// break lenght
	intermission     time.Duration
	queues           chan termbox.Event
	paused           bool
	waitForUserInput bool
	// next is the intermission
	nextInterm bool
	ticker     *time.Ticker

	lastInput time.Time

	notiPath string // if it is "" then no path provided
}

// returns: Session: X
func (a *app) printCurrSession() string {
	return fmt.Sprintf("Session: %d", a.session)
}

// returns: Session X will start soon..
func (a *app) printNextSession() string {
	return fmt.Sprintf("Session %d will start soon..", a.session+1)
}

// func (a *app) sessionInc() { a.session++ }

// clearT clears he terminal
func clearT() {
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}
}

// flush flushes the terminal
func flush() {
	if err := termbox.Flush(); err != nil {
		panic(err)
	}
}

// formatDuration calls durationToStr func
func (a *app) formatDuration() string {
	return durationToStr(a.current)
}

// getPollEvents gets the pulls form termbox.PollEvent and puths them to app.queues
//
// call like `go app.getPollEvents()`
func (a *app) getPollEvents() {
	for {
		a.queues <- termbox.PollEvent()
	}
}

const inputDelay time.Duration = 500 * time.Microsecond

func (a *app) inputDelayOK() bool {
	return time.Now().Sub(a.lastInput) < inputDelay
}

func newTimmer(duration, interm string) app {
	if duration == "" {
		duration = "45m"
	}
	if interm == "" {
		interm = "10m"
	}

	var dur, in time.Duration
	var err error
	if dur, err = time.ParseDuration(duration); err != nil {
		panic(err)
	}
	if in, err = time.ParseDuration(interm); err != nil {
		panic(err)
	}

	return app{
		session: 1,

		current:      dur,
		timmer:       dur,
		intermission: in,

		queues:     make(chan termbox.Event),
		paused:     false,
		nextInterm: true,
		ticker:     time.NewTicker(time.Second),

		notiPath: writeNoti(),
	}
}

func runPomodoro(timer, interm string) {
	a := newTimmer(timer, interm)
	go a.getPollEvents()

	clearT()
	putText("Starting....", positionMiddle, termbox.ColorDefault+termbox.AttrBold)
	flush()

loop:
	for {
		select {
		case ev := <-a.queues:
			if a.inputDelayOK() {
				continue loop
			}

			if isQuit(ev) {
				break loop
			} else if ev.Ch == 'n' || ev.Ch == 's' {
				clearT()
				if !a.nextInterm {
					putText(fmt.Sprintf("Session %d skipped...", a.session), positionMiddle, termbox.ColorGreen+termbox.AttrBold)
					a.current = a.timmer
					a.session++
				} else {
					putText("Break skipped...", positionMiddle, termbox.ColorGreen+termbox.AttrBold)
					a.current = a.intermission
				}
				flush()
				a.nextInterm = !a.nextInterm
				time.Sleep(800 * time.Millisecond)

			} else if ev.Key == termbox.KeySpace {
				clearT()
				a.paused = !a.paused
				if a.paused {
					putTime(a.formatDuration())
					putText("Paused", positionButtom, termbox.ColorLightRed+termbox.AttrBold)
					a.ticker.Stop()
				} else {
					putText("Resuming...", positionMiddle, termbox.ColorLightGreen+termbox.AttrBold)
					a.ticker.Reset(time.Second)
				}
				flush()
			}
			a.lastInput = time.Now()

		case <-a.ticker.C:
			if a.current == 0 {
				if a.nextInterm {
					go playSound(a.notiPath)
					go notify(fmt.Sprintf("Session %d done", a.session), "Go do some exercise")
					a.current = a.intermission
					a.nextInterm = false

				} else {
					a.current = a.timmer
					a.session++
					go notify("Break finised", "Comeback and contine working")
					clearT()
					putText("Press any key to continue or q to quit...", positionMiddle, termbox.ColorLightRed+termbox.AttrBold)
					flush()
					go playSound(a.notiPath)

					if ev := <-a.queues; isQuit(ev) {
						break loop
					}

					a.nextInterm = true
					continue loop
				}
			}

			clearT()
			putTime(a.formatDuration())
			a.current -= time.Second
			if !a.nextInterm {
				putText(a.printNextSession(), positionButtom, termbox.ColorLightBlue+termbox.AttrBold)
			} else {
				putText(a.printCurrSession(), positionButtom, termbox.ColorYellow+termbox.AttrBold)
			}
			flush()
		}
	}
}
