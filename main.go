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
	console := "B8:8A:EC:44:7E:AA"
	dupe, err := NewDupe(console)
	if err != nil {
		fatal(err)
	}
	dupe.Run()
}

var pairCmd = &cobra.Command{
	Use:   "pair",
	Short: "Pairs with a alternating device over Bluetooth",
	Run:   pairRun,
}

func init() {
	rootCmd.AddCommand(pairCmd)
}

func pairRun(cmd *cobra.Command, args []string) {
	if os.Geteuid() != 0 {
		fmt.Println("please run as root")
		os.Exit(1)
	}

	host := "80:32:53:37:22:19"
	if err := Pair(host); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Printf("%#v\n", err)
	os.Exit(1)
}
