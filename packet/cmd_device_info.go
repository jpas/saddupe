package packet

type CmdDeviceInfo struct{}

const CmdDeviceInfoOp OpCode = 0x02

func init() {
	RegisterCmd(&CmdDeviceInfo{})
}

func (c *CmdDeviceInfo) Op() OpCode {
	return CmdDeviceInfoOp
}

func (c *CmdDeviceInfo) Encode() ([]byte, error) {
	return []byte{byte(c.Op())}, nil
}

func (c *CmdDeviceInfo) Decode([]byte) error {
	return nil
}

type RetDeviceInfo struct {
	Firmware [2]byte
	MAC      [6]byte
	Kind     byte
	HasColor bool
}

func init() {
	RegisterRet(&RetDeviceInfo{})
}

func (r *RetDeviceInfo) Op() OpCode {
	return CmdDeviceInfoOp
}

func (r *RetDeviceInfo) Ack() bool {
	return true
}

func (r *RetDeviceInfo) Type() byte {
	return byte(CmdDeviceInfoOp)
}

func (r *RetDeviceInfo) Encode() ([]byte, error) {
	b := make([]byte, 14)
	putRetHeader(b, r)
	b[2], b[3] = 0x04, 0x06
	b[4] = r.Kind
	b[5] = 0x02
	copy(b[6:], r.MAC[:])
	b[12] = 0x01
	b[13] = byte(boolToBit(r.HasColor, 0))
	return b, nil
}

func (r *RetDeviceInfo) Decode(b []byte) error {
	copy(r.Firmware[:], b[2:])
	r.Kind = b[4]
	copy(r.MAC[:], b[6:])
	r.HasColor = bitIsSet(b[13], 0)
	return nil
}
