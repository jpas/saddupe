package report

const (
	InputReportID  byte = 0xA1 // TODO: double check value
	OutputReportID      = 0xA2 // TODO: double check value
)

type Report interface {
	Report() ([]byte, error)
}

type ReportHeader struct {
	ReportType byte
	ReportID   byte
}

type OtherReport struct {
	*ReportHeader
	Data []byte
}

type RumbleData [6]byte

type CommandCallHeader struct {
	*ReportHeader
	*RumbleData
	CommandID byte
}

type OutputHeader struct {
	Seqno uint8
}

// Remote procedure call with piggyback rumble.
type RawOutput0x01 struct {
	*CommandCallHeader
}

// NFC/IR MCU Firmware Update packet
type RawOutput0x03 struct {
	*OutputHeader
}

// Send rumble only.
type RawOutput0x10 struct {
	*OutputHeader
	*RumbleData
}

// Request data from NFC/IR MCU.
type RawOutput0x11 struct {
	*OutputHeader
	*RumbleData // TODO: maybe "Can also send rumble"
	// more unknown data...
}

// Unknown, does the same thing as Command0x28
type RawOutput0x12 struct {
	*OutputHeader
	// more unknown data...
}

type InputHeader struct {
	*ReportHeader
}

type RawInput0x21 struct {
	*InputHeader
}

type RawInput0x23 struct {
	*InputHeader
}

type RawInput0x31 struct {
	*InputHeader
}

type RawInput0x32 struct {
	*InputHeader
}

type RawInput0x33 struct {
	*InputHeader
}

type RawInput0x3F struct {
	*InputHeader
	Buttons  [2]byte
	StickHat byte
	Filler   [8]byte
}
