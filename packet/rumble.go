package packet

import (
	"bytes"

	"github.com/jpas/saddupe/hid"
)

type RumblePacket struct {
	Seqno  uint8
	Rumble [8]byte
}

func init() {
	RegisterPacket(&RumblePacket{})
}

func (p *RumblePacket) Header() hid.ReportHeader {
	return hid.OutputReportHeader
}

func (p *RumblePacket) ID() PacketID {
	return 0x10
}

func (p *RumblePacket) Encode() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(byte(p.ID()))
	buf.WriteByte(p.Seqno & 0xf)
	buf.Write(p.Rumble[:])
	return buf.Bytes(), nil
}

func (p *RumblePacket) Decode(b []byte) error {
	p.Seqno = b[1]
	copy(p.Rumble[:], b[2:10])
	return nil
}
