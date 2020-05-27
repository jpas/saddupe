package main

import (
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jpas/saddupe/hid"
	"github.com/jpas/saddupe/packet"
	"github.com/jpas/saddupe/state"
	"github.com/pkg/errors"
)

type Dupe struct {
	dev     *hid.Device
	state   *state.State
	returns chan packet.Ret
}

func NewDupe(dev *hid.Device) (*Dupe, error) {
	d := &Dupe{
		dev:     dev,
		state:   state.NewState(),
		returns: make(chan packet.Ret),
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

func (d *Dupe) Run() error {
	defer d.dev.Close()
	return waitAll(
		d.reader,
		d.writer,
	)
}

func (d *Dupe) State() *state.State {
	return d.state
}

func (d *Dupe) reader(stop <-chan struct{}) error {
	for {
		select {
		case <-stop:
			return nil
		default:
			p, err := d.recv()
			if errors.Is(err, packet.ErrUnknownPacket) {
				log.Printf("recv failed: %s", err)
				continue
			}

			if err != nil {
				return err
			}

			err = d.handlePacket(p)
			if err != nil {
				log.Printf("handler failed: %s", err)
				continue
			}
		}
	}
}

func (d *Dupe) recv() (packet.Packet, error) {
	r, err := d.dev.Read()
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}

	if r.Header != hid.OutputReportHeader {
		return nil, errors.New("not an output report")
	}

	p, err := packet.DecodeReport(r)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (d *Dupe) writer(stop <-chan struct{}) error {
	for {
		var p packet.Packet
		d.state.Tick += 1

		select {
		case <-stop:
			return nil
		case ret := <-d.returns:
			p = &packet.RetPacket{State: *d.state, Ret: ret}
		case <-time.After(15 * time.Millisecond):
			p = &packet.StatePacket{State: *d.state}
		}
		if err := d.send(p); err != nil {
			return errors.Wrap(err, "send failed")
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

func (d *Dupe) handlePacket(p packet.Packet) error {
	switch p := p.(type) {
	case *packet.CmdPacket:
		return d.handleCmd(p.Cmd)
	}
	return nil
}

func (d *Dupe) handleCmd(c packet.Cmd) error {
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

	d.returns <- r
	return err
}

func (d *Dupe) handleCmdAuxSetConfig(c packet.Cmd) (packet.Ret, error) {
	log.Printf("cmd: %s", spew.Sdump(c))
	return &packet.RetAuxSetConfig{}, nil
}

func (d *Dupe) handleCmdButtonTime(c packet.Cmd) (packet.Ret, error) {
	return &packet.RetButtonTime{L: 20, R: 20}, nil
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
