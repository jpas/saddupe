package main

import (
	"log"
)

type Dupe struct {
	ctrl *L2Socket
	intr *L2Socket
}

func NewDupe(addr *BtAddr) (*Dupe, error) {
	ctrl, err := NewL2Socket()
	if err != nil {
		return nil, err
	}

	if err := ctrl.Connect(&L2Addr{*addr, 17}); err != nil {
		return nil, err
	}

	intr, err := NewL2Socket()
	if err != nil {
		return nil, err
	}

	if err := intr.Connect(&L2Addr{*addr, 19}); err != nil {
		return nil, err
	}

	return &Dupe{ctrl, intr}, nil
}

func (d *Dupe) Run() {
	var buf [1024]byte
	for {
		n, err := d.intr.Read(buf[:])
		if err != nil {
			panic(err)
		}
		log.Printf("recv: %v", buf[:n])
	}
}
