package hid

import (
	"bytes"

	"github.com/pkg/errors"
)

type Report struct {
	Header  ReportHeader
	Payload []byte
}

type ReportHeader byte

// As per Bluetooth HID Profile, Table 3.12 (pp. 35)
const (
	InputReport   ReportHeader = 0xa1
	OutputReport  ReportHeader = 0xa2
	FeatureReport ReportHeader = 0xa3
)

func NewInputReport(payload []byte) *Report {
	return &Report{InputReport, payload}
}

func NewOutputReport(payload []byte) *Report {
	return &Report{OutputReport, payload}
}

func NewFeatureReport(payload []byte) *Report {
	return &Report{FeatureReport, payload}
}

func NewReport(buf []byte) (*Report, error) {
	if len(buf) < 1 {
		return nil, errors.New("buf too short")
	}
	header := ReportHeader(buf[0])
	switch header {
	case InputReport, OutputReport, FeatureReport:
		return &Report{header, buf[1:]}, nil
	default:
		return nil, errors.New("invalid report header")
	}
}

func (r Report) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(r.Header))
	buf.Write(r.Payload)
	return buf.Bytes()
}
