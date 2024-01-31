package main

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/nsf/termbox-go"
)

func isQuit(ev termbox.Event) bool {
	return ev.Ch == 'q' || ev.Ch == 'Q' || ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC
}

func notify(heading, description string) {
	if runtime.GOOS != "linux" {
		return
	}
	cmd := exec.Command("notify-send", heading, description)
	cmd.Run()
}

func playSound(ctx context.Context) {
	d := make(chan struct{})
	cmd := exec.Command("mpv", "--loop=2", "noti.m4a")

	go func(cmd *exec.Cmd, done chan<- struct{}) {
		if err := cmd.Start(); err != nil {
			panic(err)
		}
		d <- struct{}{}
		cmd.Wait()
		d <- struct{}{}
	}(cmd, d)
	<-d

	select {
	case <-ctx.Done():
		fmt.Println("wtf", cmd.Process.Pid)
		cmd.Process.Kill()
	case <-d:
		return
	}
}
