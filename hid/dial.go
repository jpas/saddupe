package hid

import (
	"github.com/jpas/saddupe/hid/internal/l2"
	"github.com/pkg/errors"
)

func Dial(proto string, addr string) (*Device, error) {
	switch proto {
	case "bt":
		return BtDial(addr)
	default:
		return nil, errors.New("unsupported proto")
	}
}

func BtDial(mac string) (*Device, error) {
	control, err := l2.NewConn(mac, 17)
	if err != nil {
		return nil, errors.Wrap(err, "unable to control channel")
	}

	interrupt, err := l2.NewConn(mac, 19)
	if err != nil {
		return nil, errors.Wrap(err, "unable to interrupt channel")
	}

	lmac := (interrupt.LocalAddr()).(*l2.Addr).MAC

	return &Device{control, interrupt, *lmac}, nil
}
