package main

import (
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/jpas/saddupe/hid"
	. "github.com/jpas/saddupe/internal"
	"github.com/jpas/saddupe/packet"
	"github.com/jpas/saddupe/state"
	"github.com/pkg/errors"
)

type Dupe struct {
	dev   *hid.Device
	state *state.State
	done  chan struct{}

	starter sync.Once
	started chan struct{}
}

func NewDupe(state *state.State, dev *hid.Device) (*Dupe, error) {
	d := &Dupe{
		dev:     dev,
		state:   state,
		done:    make(chan struct{}),
		started: make(chan struct{}),
	}
	return d, nil
}

func NewBtDupe(s *state.State, console string) (*Dupe, error) {
	var dupe *Dupe
	err := Retry(3, 500*time.Millisecond, func() error {
		dev, err := hid.Dial("bt", console)
		if err != nil {
			return errors.Wrap(err, "unable to connect to console")
		}

		dupe, err = NewDupe(s, dev)
		if err != nil {
			return errors.Wrap(err, "failed to initialize device")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return dupe, nil
}

func (d *Dupe) State() *state.State {
	return d.state
}

func (d *Dupe) Run() error {
	defer close(d.done)

	go d.ticker()
	go d.pusher()

	for {
		p, err := d.recv()
		if errors.Is(err, hid.ErrDeviceClosed) {
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "recv failed")
		}

		if p.Header() != hid.OutputReportHeader {
			continue
		}

		p, err = d.handlePacket(p)
		if err != nil {
			log.Println("dupe:", errors.Wrap(err, "packet handler failed"))
			continue
		}

		if p == nil {
			continue
		}
		err = d.send(p)
		if errors.Is(err, hid.ErrDeviceClosed) {
			return nil
		}
	}
}

func (d *Dupe) Close() error {
	return d.dev.Close()
}

func (d *Dupe) Started() {
	<-d.started
}

func (d *Dupe) ticker() {
	tick := time.NewTicker(time.Second / 240)
	for {
		select {
		case <-d.done:
			return
		case <-tick.C:
			d.state.Tick += 1
		}
	}
}

func (d *Dupe) pusher() {
	push := time.NewTicker(time.Second / 60)
	for {
		select {
		case <-d.done:
			return
		case <-push.C:
			p := &packet.StatePacket{State: *d.state}
			if err := d.send(p); err != nil {
				if !errors.Is(err, hid.ErrDeviceClosed) {
					log.Println(errors.Wrap(err, "send failed"))
				}
				return
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
		d.starter.Do(func() {
			close(d.started)
		})
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
		Kind:     d.state.Kind(),
		MAC:      d.dev.LocalAddr(),
		HasColor: d.state.HasColor(),
	}
	return ret, nil
}

func (d *Dupe) handleCmdFlashRead(c packet.Cmd) (packet.Ret, error) {
	cmd := c.(*packet.CmdFlashRead)
	data := make([]byte, cmd.Len)
	if err := d.state.Read(data, cmd.Addr, cmd.Len); err != nil {
		return nil, errors.Wrap(err, "unable to read flash")
	}
	return &packet.RetFlashRead{Addr: cmd.Addr, Data: data}, nil
}

func (d *Dupe) handleCmdModeSet(c packet.Cmd) (packet.Ret, error) {
	cmd := c.(*packet.CmdModeSet)
	d.state.Mode = cmd.Mode
	return packet.NewRetAck(cmd.Op(), true), nil
}
