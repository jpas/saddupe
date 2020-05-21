package main

import (
	"log"
	"time"

	"github.com/jpas/saddupe/hid"
	"github.com/jpas/saddupe/packet"

	"github.com/pkg/errors"
)

type Dupe struct {
	dev *hid.Device
}

func NewDupe(dev *hid.Device) (*Dupe, error) {
	return &Dupe{dev}, nil
}

func NewDupeBluetooth(console string) (*Dupe, error) {
	dev, err := hid.Dial("bt", console)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to console")
	}
	return NewDupe(dev)
}

func (d *Dupe) Run() {
	go d.tick()

	for {
		r, err := d.dev.Read()
		if err != nil {
			panic(err)
		}
		log.Printf("report received: %+v", r)
	}
}

func (d *Dupe) tick() {
	var p packet.SimpleButtonStatus

	for {
		r, err := p.Report()
		if err != nil {
			panic(err)
		}

		err = d.dev.Write(r)
		if err != nil {
			panic(err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
