package main

import (
	"log"
	"sync"
	"time"

	"github.com/jpas/saddupe/hid"
	"github.com/jpas/saddupe/packet"

	"github.com/pkg/errors"
)

type Dupe struct {
	dev     *hid.Device
	returns chan packet.Ret
}

func NewDupe(dev *hid.Device) (*Dupe, error) {
	return &Dupe{dev, make(chan packet.Ret)}, nil
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
	log.Printf("send: %#+v", p)
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

	log.Printf("recv: %#+v", p)

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
	case *packet.CmdShipmentState:
		d.returns <- &packet.RetAck{c.Op()}
	}
	return nil
}

func (d *Dupe) Run() error {
	defer d.dev.Close()

	return runAll(
		d.reader,
		d.ticker,
	)
}

func (d Dupe) reader(stop <-chan struct{}) error {
	for {
		select {
		case <-stop:
			return nil
		default:
			p, err := d.recv()
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

func (d Dupe) ticker(stop <-chan struct{}) error {
	for {
		select {
		case <-stop:
			return nil
		case ret := <-d.returns:
			p := &packet.RetPacket{State: nil, Ret: ret}
			if err := d.send(p); err != nil {
				return err
			}
		case <-time.After(1 * time.Second):
			p := &packet.BasicStatePacket{}
			if err := d.send(p); err != nil {
				return err
			}
		}
	}
}

func runAll(fn ...func(<-chan struct{}) error) error {
	stop := make(chan struct{})
	result := make(chan error, len(fn))

	if len(fn) < 1 {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(fn))

	for _, f := range fn {
		f := f
		go func() {
			defer wg.Done()
			result <- f(stop)
		}()
	}

	// defer to let others finish in the background
	defer wg.Wait()
	defer close(stop)

	return <-result
}
