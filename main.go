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
	intermission time.Duration

	queues chan termbox.Event
	paused bool

	waitForUserInput bool
	// next is the intermission
	nextInterm bool
	ticker     *time.Ticker
}

func (a *app) currSession() string {
	return fmt.Sprintf("Session: %d", a.session)
}
func (a *app) nextSession() string {
	return fmt.Sprintf("Session %d will start soon..", a.session+1)
}

func (a *app) sessionInc() {
	a.session++
}
func clearT() {
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}
}
func flush() {
	if err := termbox.Flush(); err != nil {
		panic(err)
	}
}

// func (a *app) formatTime() string
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

func (a *app) pollEvent() {
	for {
		a.queues <- termbox.PollEvent()
	}
}

func newTimmer() app {
	dur, _ := time.ParseDuration("9s")
	return app{
		session: 1,

		current: dur,
		timmer:  dur,
		// intermission: 5 * time.Minute,
		intermission: 5 * time.Second,

		queues:     make(chan termbox.Event),
		paused:     false,
		nextInterm: true,
		ticker:     time.NewTicker(time.Second),
	}
}

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()
	a := newTimmer()
	go a.pollEvent()
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
				putText(a.nextSession(), positionButtom, termbox.ColorLightBlue+termbox.AttrBold)
			} else {
				putText(a.currSession(), positionButtom, termbox.ColorYellow+termbox.AttrBold)
			}
			flush()
			if a.current == 0 {
				if a.nextInterm {
					a.current = a.intermission
				} else {
					a.current = a.timmer
					a.session++
					clearT()
					putText("Press any key to continue or q to quit...", positionMiddle, termbox.ColorLightRed+termbox.AttrBold)
					flush()
					if ev := <-a.queues; isQuit(ev) {
						break loop
					}

				}
				a.nextInterm = !a.nextInterm
			}
		}
	}
}

func isQuit(ev termbox.Event) bool {
	return ev.Ch == 'q' || ev.Ch == 'Q' || ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC
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
