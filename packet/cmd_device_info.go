package packet

type CmdDeviceInfo struct{}

func init() {
	RegisterCmd(&CmdDeviceInfo{})
}

func (c *CmdDeviceInfo) Op() OpCode {
	return 0x02
}

func (c *CmdDeviceInfo) Encode() ([]byte, error) {
	return []byte{byte(c.Op())}, nil
}

func (c *CmdDeviceInfo) Decode([]byte) error {
	return nil
}
