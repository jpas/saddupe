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
	dev *hid.Device
}

func NewDupe(dev *hid.Device) (*Dupe, error) {
	return &Dupe{dev}, nil
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
			r, err := d.dev.Read()
			if err != nil {
				return errors.Wrap(err, "read failed")
			}
			log.Printf("report received: %+v", r)
		}
	}
}

func (d Dupe) ticker(stop <-chan struct{}) error {
	var p packet.SimpleButtonStatus

	for {
		select {
		case <-stop:
			return nil
		case <-time.After(500 * time.Millisecond):
			r, err := p.Report()
			if err != nil {
				panic(err)
			}

			err = d.dev.Write(r)
			if err != nil {
				panic(err)
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
