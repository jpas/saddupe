package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jpas/saddupe/internal"
	"github.com/spf13/cobra"
)

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

	if len(args) == 0 {
		fatal(nil)
	}

	host := args[0]
	console, err := internal.Pair(host)
	if err != nil {
		fatal(err)
	}
	log.Printf("paired with %s", console)
}
