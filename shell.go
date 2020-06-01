package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/jpas/saddupe/state"
	"github.com/pkg/errors"
)

type Shell struct {
	dupe *Dupe
	in   io.Reader
	out  io.Writer
}

var ErrShellExited = errors.New("shell exited")

func NewShell(dupe *Dupe, in io.Reader, out io.Writer) *Shell {
	sh := &Shell{
		dupe: dupe,
		in:   in,
		out:  out,
	}
	return sh
}

func (sh *Shell) prompt() {
	fmt.Fprint(sh.out, "> ")
}

func (sh *Shell) Run() error {
	lines := make(chan string)
	done := make(chan bool)
	go func() {
		defer close(lines)
		in := bufio.NewScanner(sh.in)
		sh.prompt()
		for in.Scan() {
			lines <- in.Text()
			if d, ok := <-done; d || !ok {
				return
			}
			sh.prompt()
		}
		fmt.Fprintln(sh.out)
	}()

	for {
		select {
		case <-sh.dupe.Done():
			close(done)
			return sh.dupe.Err()
		case line, ok := <-lines:
			if !ok {
				return nil
			}

			args := strings.Fields(line)
			if len(args) == 0 {
				done <- false
				continue
			}

			err := sh.handleCmd(args[0], args[1:]...)
			if errors.Is(err, ErrShellExited) {
				done <- true
				return nil
			}

			if err != nil {
				fmt.Fprintln(sh.out, errors.Wrap(err, "shell failed"))
			}

			done <- false
		}
	}
}

func (sh *Shell) handleCmd(cmd string, args ...string) error {
	var err error
	switch cmd {
	case "e", "exit":
		return ErrShellExited
	case "p", "press":
		err = sh.handlePress(args)
	case "r", "release":
		err = sh.handleRelease(args)
	case "s", "stick":
		err = sh.handleStick(args)
	case "t", "tap":
		err = sh.handleTap(args)
	default:
		err = errors.Errorf("unknown command: %s", cmd)
	}
	return err
}

func (sh *Shell) buttonsByName(names ...string) []*state.Button {
	var buttons []*state.Button
	state := sh.dupe.State()
	for _, name := range names {
		b, err := state.ButtonByName(name)
		if err != nil {
			continue
		}
		buttons = append(buttons, b)
	}
	return buttons

}
func (sh *Shell) handlePress(args []string) error {
	for _, button := range sh.buttonsByName(args...) {
		button.Press()
	}
	return nil
}

func (sh *Shell) handleRelease(args []string) error {
	for _, button := range sh.buttonsByName(args...) {
		button.Release()
	}
	return nil
}

func (sh *Shell) handleStick(args []string) error {
	usage := errors.New("usage: stick <side> <direction>")
	if len(args) != 2 {
		return usage
	}

	stick, err := sh.dupe.State().StickByName(args[0])
	if err != nil {
		return err
	}

	// TODO(jpas) cartesian coordinates with x,y
	// TODO(jpas) polar coordinates with r:a

	var x, y float64
	switch args[1] {
	case "u", "up":
		x, y = 0, 1
	case "d", "down":
		x, y = 0, -1
	case "l", "left":
		x, y = -1, 0
	case "r", "right":
		x, y = 1, 0
	case "c", "center":
		x, y = 0, 0
	default:
		return errors.Errorf("bad direction %s", args[1])
	}
	stick.Set(x, y)
	return nil
}

func (sh *Shell) handleTap(args []string) error {
	usage := errors.New("usage: tap <button> [millis]")

	if len(args) == 0 {
		return usage
	}
	button, err := sh.dupe.State().ButtonByName(args[0])
	if err != nil {
		return err
	}

	var millis int
	switch len(args) {
	case 1:
		millis = 100
	case 2:
		millis, err = strconv.Atoi(args[1])
		if err != nil {
			return usage
		}
	default:
		return usage
	}

	button.Press()
	go func() {
		time.Sleep(time.Duration(millis) * time.Millisecond)
		button.Release()
	}()

	return nil
}
