package main

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "saddupe",
	Short: "A sad duplicate of a joyous controller",
	Run: rootRun,
}

func main() {
	rootCmd.Execute()
}

func rootRun(cmd *cobra.Command, args []string) {
	console, err := NewBtAddr("B8:8A:EC:44:7E:AA")
	if err != nil {
		log.Fatal(err)
	}

	dupe, err := NewDupe(console)
	if err != nil {
		log.Fatal(err)
	}
	dupe.Run()
}
