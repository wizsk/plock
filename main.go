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
	version    string = "1.5"
	defFPS            = 30

	timmerProgDoneCharDef     = '█'
	timmerProgNotDoneCharDef  = '░'
	timmerProgDoneColorDef    = termbox.ColorWhite
	timmerProgNotDoneColorDef = termbox.ColorDarkGray
)

var usages string = `Usage of plock [<session len> <break>] [OPTIONS..]:
A small pomodoro clock from the terminal

OPTIONS:
  -p  pomodoro timer length (default "45m")
  -b  break length (default "10m")
  -c  clock mode
  -t  timer mode or count up form 0 seconds
  -u  timer mode or count up form 0 seconds until specified time. eg. (1m30s, 3:04PM or 15:04)
  -d  Countdown mode
  -e  don't show "Ends at: ` + timeFormat + `"
  -s  silence. play no sounds
  -n  show notifications
  -f  set fps (default: ` + strconv.Itoa(defFPS) + `)
  -sf
        show avg fps
  -prog
        Timmer progress done symbol (default: "` + string(timmerProgNotDoneCharDef) + `")
  -prog-rem
        Timmer progress remaining symbol (default: "` + string(timmerProgDoneCharDef) + `")

`

func usage() {
	fmt.Print(usages)
	os.Exit(1)
}

var (
	silence, showNotifications bool
	showFps                    bool
	fpsFlag                    uint

	timmerProgDoneChar    rune
	timmerProgNotDoneChar rune
)

func main() {
	var timerLen, timerCountUntil, countDown, interm string
	var clcokMode, timerMode, showEndTime, showVersion bool

	flag.BoolVar(&clcokMode, "c", false, "clock mode")

	flag.BoolVar(&timerMode, "t", false, "timer mode or count up form 0 seconds")
	flag.StringVar(&timerCountUntil, "u", "", "timer mode or count up form 0 seconds until specified time. eg. (1m30s, 3:04PM or 15:04)")
	flag.StringVar(&countDown, "d", "", "countdown mode. format: 1m30s, 3:04PM or 15:04")

	flag.BoolVar(&showEndTime, "e", false, `don't show "Ends at: `+timeFormat+`"`)

	flag.BoolVar(&showVersion, "v", false, "show version and exit")

	flag.StringVar(&timerLen, "p", "45m", "pomodoro timer length")
	flag.StringVar(&interm, "b", "10m", "break length")

	flag.BoolVar(&silence, "s", false, "silence. play no sounds")
	flag.BoolVar(&showNotifications, "n", false, "show notifications")

	flag.UintVar(&fpsFlag, "f", defFPS, "Set fps")
	flag.BoolVar(&showFps, "sf", false, "show fps")

	var tmDone, tmNotDone string
	flag.StringVar(&tmDone, "prog", string(timmerProgDoneCharDef), "Timmer progress done symbol")
	flag.StringVar(&tmNotDone, "prog-rem", string(timmerProgNotDoneCharDef), "Timmer progress remaining symbol")

	flag.Usage = usage
	flag.Parse()

	tmDoneRune, tmNotDoneRune := []rune(tmDone), []rune(tmNotDone)
	if len(tmDoneRune) > 1 || len(tmNotDoneRune) > 1 {
		fmt.Printf("done-char or nondone-char can not be more than 1 char long, have %q, %q\n", tmDone, tmNotDone)
		os.Exit(1)
	}

	timmerProgDoneChar = tmDoneRune[0]
	timmerProgNotDoneChar = tmNotDoneRune[0]

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
		timer(time.Duration(0), false)
		return
	}

	cDown := false
	duration := ""
	if timerCountUntil != "" {
		duration = timerCountUntil
	} else if countDown != "" {
		duration = countDown
		cDown = true
	}

	if duration != "" {
		d, err := time.ParseDuration(duration)
		if err != nil {
			d, err = parseTime(duration)
		}
		if err != nil {
			termbox.Close()
			fmt.Println(err)
			fmt.Println()
			flag.Usage()
			os.Exit(1)
		}

		timer(d, cDown)
		return
	}

	args := flag.Args()
	if len(args) >= 2 {
		timerLen = args[0]
		interm = args[1]
	}

	runPomodoro(timerLen, interm, !showEndTime)
}
