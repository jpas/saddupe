package script

import (
	"github.com/jpas/saddupe/state"
)

var State *state.State = nil
var Args []string = nil

func Run(path string) error { return nil }

func Prompt(s ...string) string { return "" }

func Print(a ...interface{}) {}

func Println(a ...interface{}) {}

func Printf(format string, a ...interface{}) {}
