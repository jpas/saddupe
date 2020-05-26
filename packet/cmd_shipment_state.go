package packet

type CmdShipmentStateSet struct {
	state bool
}

const CmdShipmentStateSetOp OpCode = 0x08

func init() {
	RegisterCmd(&CmdShipmentStateSet{})
}

func (c *CmdShipmentStateSet) Op() OpCode {
	return CmdShipmentStateSetOp
}

func (c *CmdShipmentStateSet) Encode() ([]byte, error) {
	return []byte{byte(boolToBit(c.state, 0))}, nil
}

func (c *CmdShipmentStateSet) Decode(b []byte) error {
	c.state = b[1] != 0
	return nil
}

func init() {
	RegisterRet(&RetAck{op: CmdShipmentStateSetOp})
}
