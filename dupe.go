package main

import (
	"log"
	"reflect"
	"time"

	"github.com/jpas/saddupe/hid"
	"github.com/jpas/saddupe/packet"
	"github.com/jpas/saddupe/state"
	"github.com/pkg/errors"
)

type Dupe struct {
	dev     *hid.Device
	state   state.State
	flash   *Flash
	returns chan packet.Ret
}

func NewDupe(dev *hid.Device) (*Dupe, error) {
	d := &Dupe{
		dev:     dev,
		flash:   NewFlash(),
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

	if r.Header != hid.OutputReportHeader {
		return nil, errors.New("not an output report")
	}

	p, err := packet.DecodeReport(r)
	if err != nil {
		return nil, err
	}

	switch t := p.(type) {
	case *packet.RumblePacket:
		// do nothing, just spammy
	case *packet.CmdPacket:
		log.Printf("recv: %08x %s %s", d.state.Tick, reflect.TypeOf(p), reflect.TypeOf(t.Cmd))
	default:
		log.Printf("recv: %08x %s", d.state.Tick, reflect.TypeOf(p))
	}

	return p, nil
}

func (d *Dupe) handlePacket(p packet.Packet) error {
	switch p := p.(type) {
	case *packet.CmdPacket:
		return d.handleCmd(p.Cmd)
	}
	return nil
}

func (d *Dupe) handleCmd(c packet.Cmd) error {
	switch c := c.(type) {
	case *packet.CmdDeviceInfo:
		_ = c
		ret := &packet.RetDeviceInfo{
			Kind:     0x03, // hard coded pro controller
			MAC:      d.dev.MAC,
			HasColor: false,
		}
		d.returns <- ret
	case *packet.CmdShipmentState, *packet.CmdEnableVibration, *packet.CmdEnableSixaxis:
		d.returns <- packet.NewRetAck(c.Op(), true)
	case *packet.CmdFlashRead:
		data := make([]byte, c.Len)
		if err := d.flash.Read(data, c.Addr, c.Len); err != nil {
			d.returns <- packet.NewRetAck(c.Op(), false)
			break
		}
		d.returns <- &packet.RetFlashRead{Addr: c.Addr, Data: data}
	case *packet.CmdSetMode:
		d.state.Mode = c.Mode
		d.returns <- packet.NewRetAck(c.Op(), c.Mode == state.FullMode)
	case *packet.CmdButtonTime:
		d.returns <- &packet.RetButtonTime{L: 20, R: 20}
	}

	return nil
}

func (d *Dupe) Run() error {
	defer d.dev.Close()
	return waitAll(
		d.reader,
		d.writer,
	)
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

func (d *Dupe) writer(stop <-chan struct{}) error {
	for {
		var p packet.Packet
		d.state.Tick += 1

		select {
		case <-stop:
			return nil
		case ret := <-d.returns:
			p = &packet.RetPacket{State: d.state, Ret: ret}
		case <-time.After(15 * time.Millisecond):
			p = &packet.FullStatePacket{State: d.state}
		}
		if err := d.send(p); err != nil {
			return errors.Wrap(err, "send failed")
		}
	}
}
