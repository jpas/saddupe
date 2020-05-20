package main

import (
	"io"
	"log"
	"time"

	"github.com/jpas/saddupe/packet"

	"github.com/jpas/saddupe/l2"
	"github.com/pkg/errors"
)

type Dupe struct {
	ctrl io.Closer
	intr io.ReadWriteCloser
}

func NewDupe(console string) (*Dupe, error) {
	ctrl, err := l2.NewConn(console, 17)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open control channel")
	}

	intr, err := l2.NewConn(console, 19)
	if err != nil {
		ctrl.Close()
		return nil, errors.Wrap(err, "unable to open interrupt channel")
	}

	return &Dupe{io.Closer(ctrl), io.ReadWriteCloser(intr)}, nil
}

func (d *Dupe) Run() {
	go d.tick()

	var buf [1024]byte
	for {
		n, err := d.intr.Read(buf[:])
		if err != nil {
			panic(err)
		}
		log.Printf("recv: %v", buf[:n])
	}
}

func (d *Dupe) tick() {
	var p packet.SimpleButtonStatus

	for {
		b, err := p.Pack()
		if err != nil {
			panic(err)
		}

		_, err = d.intr.Write(b)
		if err != nil {
			panic(err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
