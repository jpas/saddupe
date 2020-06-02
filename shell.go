package main

import (
	"io"
	"reflect"
	"time"

	"github.com/containous/yaegi/interp"
	"github.com/jpas/saddupe/state"
)

type Shell struct {
	state  *state.State
	interp *interp.Interpreter
	in     io.Reader
	out    io.Writer
}

func NewShell(st *state.State, in io.Reader, out io.Writer) (*Shell, error) {
	i := interp.New(interp.Options{})

	i.Use(map[string]map[string]reflect.Value{
		"saddupe": {
			"State": reflect.ValueOf(st),
			"Hold": reflect.ValueOf(
				func(name string) error {
					b, err := st.ButtonByName(name)
					if err != nil {
						return err
					}
					b.Hold()
					return nil
				}),
			"Release": reflect.ValueOf(
				func(name string) error {
					b, err := st.ButtonByName(name)
					if err != nil {
						return err
					}
					b.Release()
					return nil
				}),
			"Tap": reflect.ValueOf(
				func(name string, millis ...int) error {
					b, err := st.ButtonByName(name)
					if err != nil {
						return err
					}
					m := 100
					if len(millis) != 0 {
						m = millis[0]
					}
					b.Tap(time.Millisecond * time.Duration(m))
					return nil
				}),
			"Stick": reflect.ValueOf(
				func(name string, direction string) error {
					return nil
				}),
		},
	})

	_, err := i.Eval("import . \"saddupe\"")
	if err != nil {
		return nil, err
	}

	sh := &Shell{
		state:  st,
		interp: i,
		in:     in,
		out:    out,
	}

	return sh, nil
}

func (sh *Shell) REPL() {
	sh.interp.REPL(sh.in, sh.out)
}
