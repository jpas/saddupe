package l2

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

// Addr implements net.Addr for L2CAP Addresses
type Addr struct {
	MAC [6]byte
	PSM uint16
}

// NewAddr returns an Addr for L2CAP
func NewAddr(mac string, psm uint16) (*Addr, error) {
	var a Addr

	if len(mac) != 17 {
		return nil, errors.Errorf("invalid MAC address: %s", mac)
	}
	hw, err := net.ParseMAC(mac)
	if err != nil {
		return nil, err
	}
	for i := range a.MAC {
		a.MAC[i] = hw[i]
	}

	a.PSM = psm

	return &a, nil
}

func (a *Addr) String() string {
	hw := net.HardwareAddr(a.MAC[:]) // does propper string formatting for us
	return fmt.Sprintf("[%s]:%d", hw, a.PSM)
}

// Network returns the "network" for L2CAP connections
func (a *Addr) Network() string {
	return "l2cap"
}
