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

	err := encodeState(b[1:], &p.State)
	if err != nil {
		return nil, errors.Wrap(err, "status encode faild")
	}

	ret, err := p.Ret.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "ret encode failed")
	}
	copy(b[13:], ret)
	return b[:], nil
}

func (p *RetPacket) Decode(b []byte) error {
	op := OpCode(b[14])
	target, ok := rets[op]
	if !ok {
		return errors.Wrapf(ErrUnknownPacket, "ret unknown opcode: %02x", op)
	}

	ret, err := decode(b[13:], target)
	if err != nil {
		return errors.Wrap(err, "ret decode failed")
	}
	p.Ret = (ret).(Ret)

	err = decodeState(b[1:], &p.State)
	if err != nil {
		return errors.Wrap(err, "state decode faild")
	}

	return nil
}

type Ret interface {
	Op() OpCode
	Ack() bool
	Type() byte
	EncodeDecoder
}

func putRetHeader(b []byte, r Ret) error {
	if r.Ack() {
		b[0] = 0x80 | byte(r.Type())&0x7f
	} else {
		b[0] = 0x00
	}
	b[1] = byte(r.Op())
	return nil
}

var (
	rets = map[OpCode]Decoder{}
)

func RegisterRet(r Ret) {
	rets[r.Op()] = r
}

type RetAck struct {
	op  OpCode
	ack bool
}

func NewRetAck(op OpCode, ack bool) *RetAck {
	return &RetAck{op, ack}
}

func (r *RetAck) Op() OpCode {
	return r.op
}

func (r *RetAck) Ack() bool {
	return r.ack
}

func (r *RetAck) Type() byte {
	return 0x00
}

func (r *RetAck) Encode() ([]byte, error) {
	b := make([]byte, 2)
	err := putRetHeader(b, r)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *RetAck) Decode(b []byte) error {
	r.op = OpCode(b[0])
	r.ack = bitIsSet(b[1], 7)
	return nil
}
