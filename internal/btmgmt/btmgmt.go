package btmgmt

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

// Btmgmt is a btmgmt interface for a specific Bluetooth controller
type Btmgmt struct {
	index string
}

// New constructs a btmgmt interface for a controller with a given address
func New(addr string) (*Btmgmt, error) {
	var bt Btmgmt

	if addr == "" {
		return &bt, nil
	}

	out, err := bt.Run("info")
	if err != nil {
		return nil, err
	}

	index, err := findAddrIndex(addr, out)
	if err != nil {
		return nil, errors.Wrap(err, `cannot find controller index`)
	}

	return &Btmgmt{index}, nil
}

// Run executes the given command with btmgmt
func (bt Btmgmt) Run(args ...string) (string, error) {
	a := args
	if bt.index != "" {
		a = append([]string{"--index", bt.index}, a...)
	}
	cmd := exec.Command("btmgmt", a...)

	// btmgmt does not output anything to stdout unless stdin is a file that supports polling.
	_, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	// errors are reported on stdout but with red asni escape codes,
	// the exit status is not consistent so it is ignored
	out, _ := cmd.Output()

	var buf bytes.Buffer
	s := bufio.NewScanner(io.Reader(bytes.NewBuffer(out)))
	for s.Scan() {
		line := s.Bytes()

		err := parseError(line)
		if err != nil {
			return "", errors.Wrapf(err, "btmgmt %v failed", args[0])
		}

		buf.Write(stripANSI(line))
		buf.WriteRune('\n')
	}

	return buf.String(), nil
}

var addrIndexPattern = regexp.MustCompile(`hci(\d+?):.+?\n\taddr ([:[:xdigit:]]+)`)

func findAddrIndex(addr, str string) (string, error) {
	matches := addrIndexPattern.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		if addr == match[2] {
			return match[1], nil
		}
	}

	return "", errors.Errorf("no controller with address: %v", addr)
}

var ansiEscapes = regexp.MustCompile(`[\\u001B\\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[-a-zA-Z\\d\\/#&.:=?%@~_]*)*)?\\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PR-TZcf-ntqry=><~]))`)

func stripANSI(b []byte) []byte {
	return ansiEscapes.ReplaceAllLiteral(b, nil)
}

var errPattern = regexp.MustCompile(`status 0x([[:xdigit:]]{2})`)

func parseError(b []byte) error {
	matches := errPattern.FindSubmatchIndex(b)
	if len(matches) == 0 {
		return nil
	}

	start, end := matches[2], matches[3]
	str := string(b[start:end])
	code, err := strconv.ParseUint(str, 16, 8)
	if err != nil {
		return err
	}

	switch code {
	case 0x00:
		return nil
	case 0x01:
		return ErrUnknownCommand
	case 0x02:
		return ErrNotConnected
	case 0x03:
		return ErrFailed
	case 0x04:
		return ErrConnectFailed
	case 0x05:
		return ErrAuthenticationFailed
	case 0x06:
		return ErrNotPaired
	case 0x07:
		return ErrNoResources
	case 0x08:
		return ErrTimeout
	case 0x09:
		return ErrAlreadyConnected
	case 0x0a:
		return ErrBusy
	case 0x0b:
		return ErrRejected
	case 0x0c:
		return ErrNotSupported
	case 0x0d:
		return ErrInvalidParameters
	case 0x0e:
		return ErrDisconnected
	case 0x0f:
		return ErrNotPowered
	case 0x10:
		return ErrCancelled
	case 0x11:
		return ErrInvalidIndex
	case 0x12:
		return ErrRFKilled
	case 0x13:
		return ErrAlreadyPaired
	case 0x14:
		return ErrPermissionDenied
	default:
		return ErrUnknownErrorCode
	}
}

// Errors for each btmgmt error
var (
	ErrUnknownCommand       = errors.New("unknown command")
	ErrNotConnected         = errors.New("not connected")
	ErrFailed               = errors.New("failed")
	ErrConnectFailed        = errors.New("connect failed")
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrNotPaired            = errors.New("not paired")
	ErrNoResources          = errors.New("no resources")
	ErrTimeout              = errors.New("timeout")
	ErrAlreadyConnected     = errors.New("already connected")
	ErrBusy                 = errors.New("busy")
	ErrRejected             = errors.New("rejected")
	ErrNotSupported         = errors.New("not supported")
	ErrInvalidParameters    = errors.New("invalid parameters")
	ErrDisconnected         = errors.New("disconnected")
	ErrNotPowered           = errors.New("not powered")
	ErrCancelled            = errors.New("cancelled")
	ErrInvalidIndex         = errors.New("invalid index")
	ErrRFKilled             = errors.New("rfkilled")
	ErrAlreadyPaired        = errors.New("already paired")
	ErrPermissionDenied     = errors.New("permission denied")
	ErrUnknownErrorCode     = errors.New("unknown error code")
)
