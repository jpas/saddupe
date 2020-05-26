package packet

import (
	"bytes"

	"github.com/jpas/saddupe/hid"
	"github.com/pkg/errors"
)

type OpCode byte

type CmdPacket struct {
	Seqno  uint8
	Rumble [8]byte
	Cmd    Cmd
}

func init() {
	RegisterPacket(&CmdPacket{})
}

func (p *CmdPacket) Header() hid.ReportHeader {
	return hid.OutputReportHeader
}

func (p *CmdPacket) ID() PacketID {
	return 0x01
}

func (p *CmdPacket) Encode() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(byte(p.ID()))
	buf.WriteByte(p.Seqno & 0xf)
	buf.Write(p.Rumble[:])

	cmd, err := p.Cmd.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "cmd encode faild")
	}
	buf.Write(cmd)

	return buf.Bytes(), nil
}

func (p *CmdPacket) Decode(b []byte) error {
	op := OpCode(b[10])
	target, ok := cmds[op]
	if !ok {
		return errors.Wrapf(ErrUnknownPacket, "cmd unknown opcode: %02x", op)
	}

	cmd, err := decode(b[10:], target)
	if err != nil {
		return errors.Wrap(err, "decode failed")
	}
	p.Cmd = (cmd).(Cmd)

	p.Seqno = b[1]
	copy(p.Rumble[:], b[2:10])

	return nil
}

type Cmd interface {
	Op() OpCode
	EncodeDecoder
}

var (
	cmds = map[OpCode]Decoder{}
)

func RegisterCmd(c Cmd) {
	cmds[c.Op()] = c
}
