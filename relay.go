package main

import (
	"log"
	"time"

	"github.com/jpas/saddupe/hid"
	"github.com/spf13/cobra"
)

var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "relays a joyous controller to an alternating device",
	Run:   relayRun,
}

func init() {
	rootCmd.AddCommand(relayCmd)
}

func relayRun(cmd *cobra.Command, args []string) {
	log.Println("dialing gamepad")
	controller, err := hid.BtDial("B8:78:26:64:B8:80")
	if err != nil {
		fatal(err)
	}
	defer controller.Close()
	log.Println("connected to controller")

	time.Sleep(5 * time.Second)

	log.Println("dialing console")
	console, err := hid.BtDial("B8:8A:EC:44:7E:AA")
	if err != nil {
		fatal(err)
	}
	defer console.Close()
	log.Println("connected to console")

	relay := NewRelay(controller, console)
	err = relay.Go()
	if err != nil {
		fatal(err)
	}
}

type Relay struct {
	gamepad *hid.Device
	console *hid.Device
}

func NewRelay(gamepad, console *hid.Device) *Relay {
	return &Relay{gamepad, console}
}

func (r *Relay) Go() error {
	return waitAll(
		r.forwarder("gamepad -> console", r.gamepad, r.console),
		r.forwarder("console -> gamepad", r.console, r.gamepad),
	)
}

func (r *Relay) forwarder(s string, read, write *hid.Device) func(<-chan struct{}) error {
	return func(stop <-chan struct{}) error {
		for {
			select {
			case <-stop:
				return nil
			default:
				r, err := read.Read()
				if err != nil {
					return err
				}

				log.Printf("%s: %02x", s, r)

				if err := write.Write(r); err != nil {
					return err
				}
			}
		}

	}
}
