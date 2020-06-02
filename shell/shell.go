package shell

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
	"github.com/jpas/saddupe/state"
	"github.com/pkg/errors"
)

type Shell struct {
	state  *state.State
	interp *interp.Interpreter
	in     io.Reader
	out    io.Writer
}

const envPath = "github.com/jpas/saddupe/shell/env"

func New(st *state.State, in io.Reader, out io.Writer) (*Shell, error) {
	sh := &Shell{
		state:  st,
		interp: interp.New(interp.Options{}),
		in:     in,
		out:    out,
	}
	return sh, nil
}

func (sh *Shell) loadEnv(args ...string) error {
	if len(args) == 0 {
		args = []string{}
	}

	sh.interp.Use(stdlib.Symbols)
	sh.interp.Use(map[string]map[string]reflect.Value{
		envPath: {
			"State": reflect.ValueOf(sh.state),
			"Args":  reflect.ValueOf(args),
			"Run": reflect.ValueOf(func(path string, args ...string) error {
				sub, err := New(sh.state, sh.in, sh.out)
				if err != nil {
					return err
				}
				return sub.Run(path, args...)
			}),
			"Prompt": reflect.ValueOf(func(s ...string) string {
				fmt.Fprint(sh.out, strings.Join(s, " "))
				scanner := bufio.NewScanner(sh.in)
				if !scanner.Scan() {
					return ""
				}
				return scanner.Text()
			}),
			"Print": reflect.ValueOf(func(a ...interface{}) {
				fmt.Fprint(sh.out, a...)
			}),
			"Println": reflect.ValueOf(func(a ...interface{}) {
				fmt.Fprintln(sh.out, a...)
			}),
			"Printf": reflect.ValueOf(func(format string, a ...interface{}) {
				fmt.Fprintf(sh.out, format, a...)
			}),
		},
	})

	// Ensure our shell environment is imported
	_, err := sh.interp.Eval("import . \"github.com/jpas/saddupe/shell/env\"")
	return err
}

func (sh *Shell) Run(path string, args ...string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "script open failed")
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "script read failed")
	}

	err = sh.loadEnv(args...)
	if err != nil {
		return errors.Wrap(err, "load environment failed")
	}

	_, err = sh.interp.Eval(string(b))
	return err
}

func (sh *Shell) REPL() error {
	err := sh.loadEnv()
	if err != nil {
		return errors.Wrap(err, "load environment failed")
	}
	sh.interp.REPL(sh.in, sh.out)
	return nil
}
