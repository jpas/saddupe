package main

import (
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

	dupe, err := NewBtDupe(console)
	if err != nil {
		fatal(err)
	}

	if err := dupe.Run(); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Printf("%#v\n", err)
	os.Exit(1)
}
