package hid

import (
	"github.com/pkg/errors"
)

type Report struct {
	Header  ReportHeader
	Payload []byte
}

type ReportHeader byte

// As per Bluetooth HID Profile, Table 3.12 (pp. 35)
const (
	InputReportHeader   ReportHeader = 0xa1
	OutputReportHeader  ReportHeader = 0xa2
	FeatureReportHeader ReportHeader = 0xa3
)

func NewReport(header ReportHeader, payload []byte) (*Report, error) {
	switch header {
	case InputReportHeader, OutputReportHeader, FeatureReportHeader:
		return &Report{header, payload}, nil
	default:
		return nil, errors.New("invalid report header")
	}
}

func (r Report) Bytes() []byte {
	return append([]byte{byte(r.Header)}, r.Payload...)
}
