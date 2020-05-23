package packet

type CmdButtonTime struct{}

func init() {
	RegisterCmd(&CmdButtonTime{})
}

func (c *CmdButtonTime) Op() OpCode {
	return 0x04
}

func (c *CmdButtonTime) Encode() ([]byte, error) {
	b := []byte{byte(c.Op())}
	return b, nil
}

func (c *CmdButtonTime) Decode(b []byte) error {
	return nil
}
