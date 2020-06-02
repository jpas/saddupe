package shell

import (
	"io"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
	"github.com/jpas/saddupe/state"
)

type Shell struct {
	state  *state.State
	interp *interp.Interpreter
}

func New(st *state.State) (*Shell, error) {
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	i.Use(map[string]map[string]reflect.Value{
		"github.com/jpas/saddupe/shell/env": {
			"State": reflect.ValueOf(st),
			"Run": reflect.ValueOf(func(path string) error {
				return Run(path, st)
			}),
		},
	})

	// Ensure our shell environment is imported
	_, err := i.Eval("import . \"github.com/jpas/saddupe/shell/env\"")
	if err != nil {
		return nil, err
	}

	return &Shell{state: st, interp: i}, nil
}

func Run(path string, st *state.State) error {
	sh, err := New(st)
	if err != nil {
		return err
	}
	return sh.Run(path)
}

func (sh *Shell) Run(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	_, err = sh.interp.Eval(string(b))
	return err
}

func (sh *Shell) REPL(in io.Reader, out io.Writer) {
	sh.interp.REPL(in, out)
}
