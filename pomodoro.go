package main

import (
	"context"
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

// formatDuration formats the app.current to "MM:SS" or "HH:MM:SS"
func (a *app) formatDuration() string {
	d := a.current.Round(time.Second)
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

// getPollEvents gets the pulls form termbox.PollEvent and puths them to app.queues
//
// call like `go app.getPollEvents()`
func (a *app) getPollEvents() {
	for {
		a.queues <- termbox.PollEvent()
	}
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
	}
}

func runPomodoro() {
	a := newTimmer("1s", "1s")
	go a.getPollEvents()

	clearT()
	putText("Starting....", positionMiddle, termbox.ColorDefault+termbox.AttrBold)
	flush()

loop:
	for {
		select {
		case ev := <-a.queues:
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

		case <-a.ticker.C:
			clearT()
			putTime(a.formatDuration())
			a.current -= time.Second
			if !a.nextInterm {
				putText(a.printNextSession(), positionButtom, termbox.ColorLightBlue+termbox.AttrBold)
			} else {
				putText(a.printCurrSession(), positionButtom, termbox.ColorYellow+termbox.AttrBold)
			}
			flush()
			if a.current == 0 {
				if a.nextInterm {
					ctx, stopSound := context.WithCancel(context.Background())
					go playSound(ctx)
					defer func() {
						stopSound()
					}()
					go notify(fmt.Sprintf("Session %d done", a.session), "Go do some exercise")
					a.current = a.intermission

				} else {
					a.current = a.timmer
					a.session++
					ctx, cancl := context.WithCancel(context.Background())
					go playSound(ctx)
					go notify("Break finised", "Comeback and contine working")
					clearT()
					putText("Press any key to continue or q to quit...", positionMiddle, termbox.ColorLightRed+termbox.AttrBold)
					flush()
					if ev := <-a.queues; isQuit(ev) {
						cancl()
						break loop
					}
					cancl()
				}
				a.nextInterm = !a.nextInterm
			}
		}
	}
}
