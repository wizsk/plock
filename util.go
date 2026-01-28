// utility funcs
package main

import (
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

const notificationFileName string = "notification_sound.m4a"

//go:embed notification_sound.m4a
var files embed.FS

func writeNoti() string {
	// if there is no mpv just skip writing sounds
	if _, err := exec.LookPath("mpv"); err != nil {
		return ""
	}

	path := filepath.Join(os.TempDir(), notificationFileName)
	if _, err := os.Stat(path); err == nil {
		return path
	}

	f, err := os.Create(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	noti, err := files.Open(notificationFileName)
	if err != nil {
		panic(err)
	}
	defer noti.Close()

	if _, err := io.Copy(f, noti); err != nil {
		return ""
	}
	return f.Name()
}

func warnAboutDependencies() {
	const cmdNotFound = "command %q not found, please install it to get %s\n"
	if !silence {
		if _, err := exec.LookPath("mpv"); err != nil {
			fmt.Fprintf(os.Stderr, cmdNotFound, "mpv", "allart sounds")
		}
	}
	if showNotifications {
		if runtime.GOOS == "linux" {
			if _, err := exec.LookPath("notify-send"); err != nil {
				fmt.Fprintf(os.Stderr, cmdNotFound, "notify-send", "notifications")
			}
		}
	}
}

func isQuit(ev termbox.Event) bool {
	return ev.Ch == 'q' || ev.Ch == 'Q' || ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC
}

func notify(heading, description string) {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("notify-send", "--wait", "--urgency=critical", heading, description)
		cmd.Run()
	}
}

func printErr(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func playSound(path string) {
	if path == "" {
		return
	}
	cmd := exec.Command("mpv", path)
	cmd.Run()
}

// durationToStr formats the app.current to "MM:SS" or "HH:MM:SS"
func durationToStr(d time.Duration) string {
	d = d.Round(time.Second)
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

// form: https://github.com/antonmedv/countdown main.go:205
func parseTime(date string) (time.Duration, error) {
	targetTime, err := time.Parse(time.Kitchen, strings.ToUpper(date))
	if err != nil {
		targetTime, err = time.Parse("15:04", date)
		if err != nil {
			return time.Duration(0), err
		}
	}

	now := time.Now()
	originTime := time.Date(0, time.January, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

	// The time of day has already passed, so target tomorrow.
	if targetTime.Before(originTime) {
		targetTime = targetTime.AddDate(0, 0, 1)
	}

	duration := targetTime.Sub(originTime)

	return duration, err
}
