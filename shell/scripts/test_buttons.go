package main

import (
	"time"

	. "github.com/jpas/saddupe/shell/env"
)

func main() {
	buttons := []string{
		"y",
		"x",
		"b",
		"a",
		"r",
		"sr",
		"zr",
		"l",
		"sl",
		"zl",
		"minus",
		"plus",
		"down",
		"up",
		"right",
		"left",
		"leftstick",
		"rightstick",
	}

	stop := make(chan struct{})
	go func() {
		Prompt("press return to stop")
		close(stop)
	}()

	t := time.Second / 15

	for b := 0; ; b++ {
		button, err := State.ButtonByName(buttons[b%len(buttons)])
		if err != nil {
			continue
		}
		button.Tap(t)

		select {
		case <-time.After(time.Second / 2):
			continue
		case <-stop:
			return
		}
	}
}
