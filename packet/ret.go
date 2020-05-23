package packet

import (
	"github.com/jpas/saddupe/hid"
	"github.com/jpas/saddupe/state"
	"github.com/pkg/errors"
)

type RetPacket struct {
	State state.State
	Ret   Ret
}

func init() {
	RegisterPacket(&RetPacket{})
}

func (p *RetPacket) Header() hid.ReportHeader {
	return hid.InputReportHeader
}

func (p *RetPacket) ID() PacketID {
	return 0x21
}

func (p *RetPacket) Encode() ([]byte, error) {
	var b [50]byte

	b[0] = byte(p.ID())

	state, err := p.State.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "status encode faild")
	}
	copy(b[1:], state)

	ret, err := p.Ret.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "ret encode failed")
	}
	copy(b[15:], ret)

	if p.Ret.Ack() {
		b[13] = 0x80 | p.Ret.Type()&0x7f
	} else {
		b[13] = 0x00
	}

	b[14] = byte(p.Ret.Op())

	return b[:], nil
}

func (p *RetPacket) Decode(b []byte) error {
	op := OpCode(b[10])
	target, ok := rets[op]
	if !ok {
		return errors.Wrapf(ErrUnknownPacket, "unknown opcode: %02x", op)
	}

	ret, err := decode(b[10:], target)
	if err != nil {
		return errors.Wrap(err, "ret decode failed")
	}
	p.Ret = (ret).(Ret)

	err = p.State.Decode(b[1:])
	if err != nil {
		return errors.Wrap(err, "status decode faild")
	}

	return nil
}

type Ret interface {
	Op() OpCode
	Ack() bool
	Type() byte
	EncodeDecoder
}

var (
	rets = map[OpCode]Decoder{}
)

func RegisterRet(r Ret) {
	rets[r.Op()] = r
}

type RetAck struct {
	TheOp OpCode
}

func (r *RetAck) Op() OpCode {
	return r.TheOp
}

func (r *RetAck) Ack() bool {
	return true
}

func (r *RetAck) Type() byte {
	return 0x00
}

func (r *RetAck) Encode() ([]byte, error) {
	return nil, nil
}

func (r *RetAck) Decode(b []byte) error {
	r.TheOp = OpCode(b[0])
	return nil
}
