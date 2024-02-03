package main

import (
	"flag"
	"fmt"

	"github.com/nsf/termbox-go"
)

const usages string = `Usage of clock [pomodoro_time break_time]:
  -p
	pomodoro timer length (default "45m")
  -b
	break length (default "10m")
  -c
	clock mode
  -t
	timer mode or count up form 0 seconds
`

func main() {
	var timerLen, interm string
	var clcokMode, timerMode bool
	flag.BoolVar(&clcokMode, "c", false, "clock mode")
	flag.BoolVar(&timerMode, "t", false, "timer mode or count up form 0 seconds")
	flag.StringVar(&timerLen, "p", "45m", "pomodoro timer length")
	flag.StringVar(&interm, "b", "10m", "break length")
	flag.Usage = func() { fmt.Print(usages) }
	flag.Parse()

	if !(clcokMode || timerMode) {
		warnAboutDependencies()
	}

	// initializing the terminal
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	termbox.HideCursor()
	defer termbox.Close()
	defer termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	if clcokMode {
		clock()
		return
	}

	if timerMode {
		timer()
		return
	}

	args := flag.Args()
	if len(args) >= 2 {
		timerLen = args[0]
		interm = args[1]
	}

	runPomodoro(timerLen, interm)
}
