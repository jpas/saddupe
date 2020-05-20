package packet

import (
	"github.com/pkg/errors"
)

var (
	ErrTooShort  = errors.New("packet: buffer too short")
	ErrMalformed = errors.New("packet: packet malformed")
)

type Packet interface {
	Pack() ([]byte, error)
}

type PacketID uint16

func unpackPacketID(p []byte) (PacketID, error) {
	if len(p) < 2 {
		return 0, ErrTooShort
	}
	return PacketID(p[0])<<8 | PacketID(p[1]), nil
}

func (i PacketID) packInto(p []byte) error {
	if len(p) < 2 {
		return ErrTooShort
	}
	p[0] = byte(i >> 8)
	p[1] = byte(i)
	return nil
}

const (
	InputPacketID  PacketID = 0xA100
	OutputPacketID PacketID = 0xA200
)

type Unpacker func([]byte) (Packet, error)

var (
	Unpackers = make(map[PacketID]Unpacker)
)

func Unpack(p []byte) (Packet, error) {
	id, err := unpackPacketID(p)
	if err != nil {
		return nil, err
	}

	unpacker, ok := Unpackers[id]
	if !ok {
		return nil, errors.Errorf("packet: unknown packet id: %04x", id)
	}
	return unpacker(p)
}
