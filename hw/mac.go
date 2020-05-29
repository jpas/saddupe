package hw

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type MAC [6]byte

func ParseMAC(mac string) (*MAC, error) {
	hw, err := net.ParseMAC(mac)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse MAC address")
	}

	var m MAC
	if len(hw) != len(m) {
		return nil, errors.Errorf("invalid MAC address: %s", mac)
	}
	for i := 0; i < len(m); i++ {
		m[i] = hw[5-i]
	}

	return &m, nil
}

func (m MAC) Bytes() []byte {
	return m[:]
}

func (m MAC) String() string {
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", m[5], m[4], m[3], m[2], m[1], m[0])
}
