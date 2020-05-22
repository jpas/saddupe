package packet

type CmdShipmentState struct {
	state bool
}

func init() {
	RegisterCmd(&CmdShipmentState{})
}

func (c *CmdShipmentState) Op() OpCode {
	return 0x08
}

func (c *CmdShipmentState) Encode() ([]byte, error) {
	return []byte{byte(boolToBit(c.state, 0))}, nil
}

func (c *CmdShipmentState) Decode(b []byte) error {
	c.state = b[1] != 0
	return nil
}
