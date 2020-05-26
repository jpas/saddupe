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
	b := []byte{
		0x04, 0x00, // firmware version
		r.Kind,
		0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // MAC
		0x01,
		byte(boolToBit(r.HasColor, 0)),
	}
	copy(b[4:], r.MAC[:])
	return b, nil
}

func (r *RetDeviceInfo) Decode(b []byte) error {
	return nil
}
