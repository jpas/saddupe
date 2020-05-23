package packet

import (
	"github.com/jpas/saddupe/hid"
	"github.com/jpas/saddupe/state"
	"github.com/pkg/errors"
)

type FullStatePacket struct {
	State state.State
}

func init() {
	RegisterPacket(&FullStatePacket{})
}

func (p *FullStatePacket) Header() hid.ReportHeader {
	return hid.InputReportHeader
}

func (p *FullStatePacket) ID() PacketID {
	return 0x30
}

func (p *FullStatePacket) Encode() ([]byte, error) {
	var b [50]byte

	b[0] = byte(p.ID())

	state, err := p.State.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "status encode faild")
	}
	copy(b[1:], state)

	return b[:], nil
}

func (p *FullStatePacket) Decode(b []byte) error {
	err := p.State.Decode(b[1:])
	if err != nil {
		return errors.Wrap(err, "status decode faild")
	}
	return nil
}
