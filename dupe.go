package main

import (
	"log"
	"reflect"
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
	tg    TaskGroup
}

func NewDupe(dev *hid.Device) (*Dupe, error) {
	d := &Dupe{
		dev:   dev,
		state: state.NewState(),
	}
	if err := d.start(); err != nil {
		dev.Close()
		return nil, err
	}
	return d, nil
}

func NewBtDupe(console string) (*Dupe, error) {
	var dupe *Dupe
	err := Retry(3, 500*time.Millisecond, func() error {
		dev, err := hid.Dial("bt", console)
		if err != nil {
			return errors.Wrap(err, "unable to connect to console")
		}

		dupe, err = NewDupe(dev)
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

func (d *Dupe) start() error {
	d.tg.Add(d.ticker)
	d.tg.Add(d.pusher)
	d.tg.Add(d.receiver)
	d.tg.Start()
	// we might have early failure so wait a little bit to make sure everything is okay

	select {
	case err := <-d.Exited():
		return errors.Wrap(err, "early failure")
	case <-time.After(1 * time.Second):
		return nil
	}
}

func (d *Dupe) State() *state.State {
	return d.state
}

func (d *Dupe) Stop() error {
	d.tg.Stop()
	return d.Wait()
}

func (d *Dupe) Wait() error {
	d.tg.Wait()
	return d.tg.Err()
}

func (d *Dupe) Exited() <-chan error {
	err := make(chan error)
	go func() {
		<-d.tg.Done()
		err <- d.tg.Err()
	}()
	return err
}

func (d *Dupe) ticker() error {
	tick := time.NewTicker(time.Second / 240)
	for {
		select {
		case <-d.tg.Done():
			return nil
		case <-tick.C:
			d.state.Tick += 1
		}
	}
}

func (d *Dupe) pusher() error {
	push := time.NewTicker(time.Second / 60)
	for {
		select {
		case <-d.tg.Done():
			return nil
		case <-push.C:
			p := &packet.StatePacket{State: *d.state}
			if err := d.send(p); err != nil {
				log.Println(errors.Wrap(err, "send failed"))
				continue
			}
		}
	}
}

func (d *Dupe) receiver() error {
	recv := make(chan packet.Packet)
	errc := make(chan error, 1)
	go func() {
		for {
			select {
			case <-d.tg.Done():
				return
			default:
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
		}
	}()

	for {
		select {
		case <-d.tg.Done():
			return nil
		case err := <-errc:
			return err
		case p := <-recv:
			err := d.handlePacket(p)
			if err != nil {
				log.Println(errors.Wrap(err, "packet handler failed"))
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

func (d *Dupe) handlePacket(p packet.Packet) error {
	switch p := p.(type) {
	case *packet.CmdPacket:
		return d.handleCmd(p.Cmd)
	case *packet.RumblePacket:
		return nil
	default:
		return errors.Errorf("no packet handler for %s", reflect.TypeOf(p))
	}
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

	p := &packet.RetPacket{State: *d.State(), Ret: r}
	if err := d.send(p); err != nil {
		log.Println(errors.Wrap(err, "send failed"))
		return err
	}
	return nil
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
