package packet

import (
	"github.com/jpas/saddupe/hid"
	"github.com/pkg/errors"
)

type Packet interface {
	Header() hid.ReportHeader
	ID() PacketID
	Encoder
	Decoder
}

type PacketID byte

type packetDecoderKey struct {
	Header hid.ReportHeader
	ID     PacketID
}

var (
	packets = map[packetDecoderKey]Decoder{}
)

func RegisterPacket(p Packet) {
	packets[packetDecoderKey{p.Header(), p.ID()}] = p
}

var ErrUnknownPacket = errors.New("packet: unknown packet")

func DecodeReport(r *hid.Report) (Packet, error) {
	key := packetDecoderKey{r.Header, PacketID(r.Payload[0])}
	target, ok := packets[key]
	if !ok {
		return nil, errors.Wrapf(ErrUnknownPacket, "unknown packet key: %02x %02x", key.Header, key.ID)
	}

	packet, err := decode(r.Payload, target)
	if err != nil {
		return nil, errors.Wrap(err, "decode failed")
	}
	return (packet).(Packet), nil
}

func EncodeReport(p Packet) (*hid.Report, error) {
	payload, err := p.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "encode failed")
	}
	return hid.NewReport(p.Header(), payload)
}
