package main

import (
	"flag"

	"github.com/nsf/termbox-go"
)

func main() {
	var timer, interm string
	var clcokMode, stopWatchMode bool
	flag.BoolVar(&clcokMode, "c", false, "clock mode")
	flag.BoolVar(&stopWatchMode, "s", false, "stopwatch or count up")
	flag.StringVar(&timer, "p", "45m", "pomodoro timer length")
	flag.StringVar(&interm, "b", "10m", "break length")
	flag.Parse()

	// initializing the terminal
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	termbox.HideCursor()
	defer termbox.Close()

	if clcokMode {
		clock()
		return
	}

	if stopWatchMode {
		stopWatch()
		return
	}

	args := flag.Args()
	if len(args) >= 2 {
		timer = args[0]
		interm = args[1]
	}

	runPomodoro(timer, interm)
}
