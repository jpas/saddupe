package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jpas/saddupe/hid"
	"github.com/jpas/saddupe/packet"
	"github.com/jpas/saddupe/state"
	"github.com/pkg/errors"
)

type Dupe struct {
	dev   *hid.Device
	state *state.State
	stop  chan struct{}
}

func NewDupe(dev *hid.Device) (*Dupe, error) {
	d := &Dupe{
		dev:   dev,
		state: state.NewState(),
		stop:  make(chan struct{}),
	}
	return d, nil
}

func NewBtDupe(console string) (*Dupe, error) {
	dev, err := hid.Dial("bt", console)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to console")
	}
	return NewDupe(dev)
}

func (d *Dupe) State() *state.State {
	return d.state
}

func (d *Dupe) Run() error {
	errc := make(chan error)

	go func() {
		for {
			time.Sleep(time.Second / 240)
			d.state.Tick += 1
		}
	}()

	recv := make(chan packet.Packet)
	go func() {
		for {
			p, err := d.recv()
			if err != nil {
				errc <- err
				return
			}
			if p.Header() != hid.OutputReportHeader {
				continue
			}
			recv <- p
		}
	}()

	shell := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				shell <- "exit"
			}
			shell <- scanner.Text()
		}
	}()

	push := time.After(time.Second / 120)
	for {
		select {
		case <-d.stop:
			return nil
		case err := <-errc:
			close(d.stop)
			return err
		case p := <-recv:
			p, err := d.handlePacket(p)
			if err != nil {
				log.Println(errors.Wrap(err, "packet handler failed"))
				continue
			}
			if p == nil {
				continue
			}
			if err := d.send(p); err != nil {
				log.Println(errors.Wrap(err, "send failed"))
				continue
			}
		case <-push:
			push = time.After(time.Second / 120)
			p := &packet.StatePacket{State: *d.state}
			if err := d.send(p); err != nil {
				log.Println(errors.Wrap(err, "send failed"))
				continue
			}
		case line := <-shell:
			args := strings.Fields(line)
			if len(args) == 0 {
				continue
			}
			if err := d.handleShellCmd(args[0], args[1:]...); err != nil {
				log.Println(errors.Wrap(err, "shell failed"))
				continue
			}
		}
	}
}

func (d *Dupe) send(p packet.Packet) error {
	if p.Header() != hid.InputReportHeader {
		return errors.New("not an input report")
	}
	r, err := packet.EncodeReport(p)
	if err != nil {
		return errors.Wrap(err, "packet encode failed")
	}
	return d.dev.Write(r)
}

func (d *Dupe) recv() (packet.Packet, error) {
	r, err := d.dev.Read()
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}
	p, err := packet.DecodeReport(r)
	if err != nil {
		return nil, errors.Wrap(err, "decode failed")
	}
	return p, nil
}

func (d *Dupe) handlePacket(p packet.Packet) (packet.Packet, error) {
	switch p := p.(type) {
	case *packet.CmdPacket:
		return d.handleCmd(p.Cmd)
	case *packet.RumblePacket:
		return nil, nil
	default:
		return nil, errors.Errorf("no packet handler for %s", reflect.TypeOf(p))
	}
}

func (d *Dupe) handleCmd(c packet.Cmd) (packet.Packet, error) {
	var r packet.Ret
	var err error

	switch c := c.(type) {
	case *packet.CmdAuxSetConfig:
		r, err = d.handleCmdAuxSetConfig(c)
	case *packet.CmdButtonTime:
		r, err = d.handleCmdButtonTime(c)
	case *packet.CmdDeviceGetInfo:
		r, err = d.handleCmdDeviceGetInfo(c)
	case *packet.CmdFlashRead:
		r, err = d.handleCmdFlashRead(c)
	case *packet.CmdModeSet:
		r, err = d.handleCmdModeSet(c)
	default:
		r = packet.NewRetAck(c.Op(), true)
	}

	if err != nil {
		r = packet.NewRetAck(c.Op(), false)
	}
	return &packet.RetPacket{State: *d.State(), Ret: r}, nil
}

func (d *Dupe) handleCmdAuxSetConfig(c packet.Cmd) (packet.Ret, error) {
	return &packet.RetAuxSetConfig{}, nil
}

func (d *Dupe) handleCmdButtonTime(c packet.Cmd) (packet.Ret, error) {
	r := &packet.RetButtonTime{
		L:    uint16(d.state.L.Milliseconds() / 10),
		R:    uint16(d.state.R.Milliseconds() / 10),
		ZL:   uint16(d.state.ZL.Milliseconds() / 10),
		ZR:   uint16(d.state.ZR.Milliseconds() / 10),
		Home: uint16(d.state.Home.Milliseconds() / 10),
	}
	return r, nil
}

func (d *Dupe) handleCmdDeviceGetInfo(c packet.Cmd) (packet.Ret, error) {
	ret := &packet.RetDeviceGetInfo{
		Kind:     0x03, // hard coded pro controller
		MAC:      d.dev.MAC,
		HasColor: false,
	}
	return ret, nil
}

func (d *Dupe) handleCmdFlashRead(c packet.Cmd) (packet.Ret, error) {
	cmd := c.(*packet.CmdFlashRead)
	data := make([]byte, cmd.Len)
	if err := d.state.Flash.Read(data, cmd.Addr, cmd.Len); err != nil {
		return nil, errors.Wrap(err, "unable to read flash")
	}
	return &packet.RetFlashRead{Addr: cmd.Addr, Data: data}, nil
}

func (d *Dupe) handleCmdModeSet(c packet.Cmd) (packet.Ret, error) {
	cmd := c.(*packet.CmdModeSet)
	d.state.Mode = cmd.Mode
	return packet.NewRetAck(cmd.Op(), true), nil
}

func (d *Dupe) handleShellCmd(cmd string, args ...string) error {
	var err error
	switch cmd {
	case "e", "exit":
		close(d.stop)
		return nil
	case "p", "press":
		err = d.handleShellPress(args)
	case "r", "release":
		err = d.handleShellRelease(args)
	case "t", "tap":
		err = d.handleShellTap(args)
	default:
		err = errors.Errorf("unknown command: %s", cmd)
	}
	return err
}

func (d *Dupe) handleShellPress(args []string) error {
	for _, name := range args {
		button, err := d.buttonByName(name)
		if err != nil {
			continue
		}
		button.Press()
	}
	return nil
}

func (d *Dupe) handleShellRelease(args []string) error {
	for _, name := range args {
		button, err := d.buttonByName(name)
		if err != nil {
			continue
		}
		button.Press()
	}
	return nil
}

func (d *Dupe) handleShellTap(args []string) error {
	usage := errors.New("usage: tap <button> [millis]")

	if len(args) == 0 {
		return usage
	}
	button, err := d.buttonByName(args[0])
	if err != nil {
		return err
	}

	var millis int
	switch len(args) {
	case 1:
		millis = 100
	case 2:
		millis, err = strconv.Atoi(args[1])
		if err != nil {
			return usage
		}
	default:
		return usage
	}

	button.Press()
	go func() {
		time.Sleep(time.Duration(millis) * time.Millisecond)
		button.Release()
	}()

	return nil
}

func (d *Dupe) buttonByName(name string) (*state.Button, error) {
	var b *state.Button

	switch strings.ToLower(name) {
	case "y":
		b = &d.State().Y
	case "x":
		b = &d.State().X
	case "b":
		b = &d.State().B
	case "a":
		b = &d.State().A
	case "r":
		b = &d.State().R
	case "zr":
		b = &d.State().ZR
	case "l":
		b = &d.State().L
	case "zl":
		b = &d.State().ZL
	case "minus":
		b = &d.State().Minus
	case "plus":
		b = &d.State().Plus
	case "home":
		b = &d.State().Home
	case "capture":
		b = &d.State().Capture
	case "down":
		b = &d.State().Down
	case "up":
		b = &d.State().Up
	case "right":
		b = &d.State().Right
	case "left":
		b = &d.State().Left
	case "leftstick":
		b = &d.State().LeftStick.Button
	case "rightstick":
		b = &d.State().RightStick.Button
	default:
		return nil, errors.New("unknown button")
	}
	return b, nil
}
