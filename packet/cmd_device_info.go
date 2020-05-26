package packet

type CmdDeviceGetInfo struct{}

const CmdDeviceGetInfoOp OpCode = 0x02

func init() {
	RegisterCmd(&CmdDeviceGetInfo{})
}

func (c *CmdDeviceGetInfo) Op() OpCode {
	return CmdDeviceGetInfoOp
}

func (c *CmdDeviceGetInfo) Encode() ([]byte, error) {
	return []byte{byte(c.Op())}, nil
}

func (c *CmdDeviceGetInfo) Decode([]byte) error {
	return nil
}

type RetDeviceGetInfo struct {
	Firmware [2]byte
	MAC      [6]byte
	Kind     byte
	HasColor bool
}

func init() {
	RegisterRet(&RetDeviceGetInfo{})
}

func (r *RetDeviceGetInfo) Op() OpCode {
	return CmdDeviceGetInfoOp
}

func (r *RetDeviceGetInfo) Ack() bool {
	return true
}

func (r *RetDeviceGetInfo) Type() byte {
	return byte(CmdDeviceGetInfoOp)
}

func (r *RetDeviceGetInfo) Encode() ([]byte, error) {
	b := make([]byte, 14)
	if err := putRetHeader(b, r); err != nil {
		return nil, err
	}
	b[2], b[3] = 0x04, 0x06
	b[4] = r.Kind
	b[5] = 0x02
	copy(b[6:], r.MAC[:])
	b[12] = 0x01
	b[13] = byte(boolToBit(r.HasColor, 0))
	return b, nil
}

func (r *RetDeviceGetInfo) Decode(b []byte) error {
	copy(r.Firmware[:], b[2:])
	r.Kind = b[4]
	copy(r.MAC[:], b[6:])
	r.HasColor = bitIsSet(b[13], 0)
	return nil
}
