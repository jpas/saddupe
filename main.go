package main

import (
	"errors"
	"fmt"
	"os"

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

	var sh *Shell

	dupe, err := NewBtDupe(console)
	if err != nil {
		fatal(err)
	}

	sh = NewShell(dupe, os.Stdin, os.Stdout)
	if err := sh.Run(); err != nil && !errors.Is(err, ErrShellExited) {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Printf("%#v\n", err)
	os.Exit(1)
}
