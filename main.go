package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	timeFormat string = "03:04:05 PM"
	version    string = "1.2"
	defFPS            = 30
)

var usages string = `Usage of plock [<session len> <break>] [OPTIONS..]:
A small pomodoro clock from the terminal

OPTIONS:
  -p  pomodoro timer length (default "45m")
  -b  break length (default "10m")
  -c  clock mode
  -t  timer mode or count up form 0 seconds
  -u  timer mode or count up form 0 seconds until specified time. eg. (1m30s, 3:04PM or 15:04)
  -e  don't show "Ends at: ` + timeFormat + `"
  -s  silence. play no sounds
  -n  show notifications
  -f  fps (default: ` + strconv.Itoa(defFPS) + `)
`

func usage() {
	fmt.Print(usages)
	os.Exit(1)
}

var (
	silence, showNotifications bool
	showFps                    bool
	fpsFlag                    uint
)

func main() {
	var timerLen, timerCountUntil, interm string
	var clcokMode, timerMode, showEndTime, showVersion bool

	flag.BoolVar(&clcokMode, "c", false, "clock mode")
	flag.BoolVar(&timerMode, "t", false, "timer mode or count up form 0 seconds")
	flag.StringVar(&timerCountUntil, "u", "", "timer mode or count up form 0 seconds until specified time.")
	flag.BoolVar(&showEndTime, "e", false, "don't show ends at")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.StringVar(&timerLen, "p", "45m", "pomodoro timer length")
	flag.StringVar(&interm, "b", "10m", "break length")
	flag.BoolVar(&silence, "s", false, "silence")
	flag.BoolVar(&showNotifications, "n", false, "silence")
	flag.UintVar(&fpsFlag, "f", defFPS, "silence")
	flag.BoolVar(&showFps, "sf", false, "show fps")
	flag.Usage = usage
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

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
		timer(time.Duration(0))
		return
	}

	if timerCountUntil != "" {
		d, err := time.ParseDuration(timerCountUntil)
		if err != nil {
			d, err = parseTime(timerCountUntil)
		}
		if err != nil {
			termbox.Close()
			fmt.Println(err)
			fmt.Println()
			usage()
		}

		timer(d)
		return
	}

	args := flag.Args()
	if len(args) >= 2 {
		timerLen = args[0]
		interm = args[1]
	}

	runPomodoro(timerLen, interm, !showEndTime)
}
