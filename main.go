package main

import (
	"fmt"
	"os"

	"github.com/jpas/saddupe/shell"
	"github.com/jpas/saddupe/state"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "saddupe",
	Short: "A sad duplicate of a joyous controller",
	Run:   rootRun,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func rootRun(cmd *cobra.Command, args []string) {
	var console string

	if len(args) == 0 {
		console = "B8:8A:EC:44:7E:AA" // switch
	} else {
		console = args[0]
	}

	st := state.NewState(state.Pro)
	dupe, err := NewBtDupe(st, console)
	if err != nil {
		fatal(err)
	}

	sh, err := shell.NewShell(st)
	if err != nil {
		fatal(err)
	}
	go func() {
		dupe.Started()
		sh.REPL(os.Stdin, os.Stdout)
		dupe.Close()
	}()

	if err := dupe.Run(); err != nil {
		fatal(err)
	}

}

func fatal(err error) {
	fmt.Printf("%#v\n", err)
	os.Exit(1)
}
