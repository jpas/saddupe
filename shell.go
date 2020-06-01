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
	in   *Scanner
	out  *Printer
}

var ErrShellExited = errors.New("shell exited")

func NewShell(dupe *Dupe, in io.Reader, out io.Writer) *Shell {
	sh := &Shell{
		dupe: dupe,
		in:   NewScanner(in),
		out:  NewPrinter(out),
	}
	return sh
}

func (sh *Shell) prompt() {
	sh.out.Print("> ")
}

func (sh *Shell) Run() error {
	for {
		sh.prompt()
		select {
		case <-sh.dupe.Done():
			return sh.dupe.Err()
		case ok := <-sh.in.Scan():
			if !ok {
				return nil
			}

			args := strings.Fields(sh.in.Text())
			if len(args) == 0 {
				continue
			}

			err := sh.handleCmd(args[0], args[1:]...)
			if errors.Is(err, ErrShellExited) {
				return nil
			}

			if err != nil {
				sh.out.Println(errors.Wrap(err, "shell failed"))
			}
		}
	}
}

func (sh *Shell) handleCmd(cmd string, args ...string) error {
	var err error
	switch cmd {
	case "e", "exit":
		return ErrShellExited
	case "h", "hold":
		err = sh.handleHold(args)
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
func (sh *Shell) handleHold(args []string) error {
	for _, button := range sh.buttonsByName(args...) {
		button.Hold()
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

	millis := 100
	if len(args) == 2 {
		millis, err = strconv.Atoi(args[1])
		if err != nil {
			return usage
		}
	}

	button.Tap(time.Duration(millis) * time.Millisecond)
	return nil
}

type Scanner struct {
	*bufio.Scanner
	scan chan bool
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		Scanner: bufio.NewScanner(r),
		scan:    make(chan bool),
	}
}

func (s *Scanner) Scan() <-chan bool {
	go func() {
		s.scan <- s.Scanner.Scan()
	}()
	return s.scan
}

type Printer struct {
	w io.Writer
}

func NewPrinter(w io.Writer) *Printer {
	return &Printer{w}
}

func (p Printer) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(p.w, a...)
}

func (p Printer) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(p.w, format, a...)
}

func (p Printer) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(p.w, a...)
}
