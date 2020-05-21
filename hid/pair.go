package hid

import (
	"os/exec"
	"time"

	. "github.com/jpas/saddupe/hid/internal"
	"github.com/jpas/saddupe/hid/internal/l2"
)

func Pair(host string) (string, error) {
	_ = exec.Command("systemctl", "stop", "bluetooth").Run()
	time.Sleep(1 * time.Second)

	ctrl, err := l2.NewListener("00:00:00:00:00:00", 17)
	if err != nil {
		return "", err
	}
	defer ctrl.Close()

	intr, err := l2.NewListener("00:00:00:00:00:00", 19)
	if err != nil {
		return "", err
	}
	defer intr.Close()

	_ = exec.Command("systemctl", "start", "bluetooth").Run()
	time.Sleep(1 * time.Second)

	spoof, err := NewSpoofer(host)
	if err != nil {
		return "", err
	}
	defer spoof.Stop()

	c, err := ctrl.Accept()
	if err != nil {
		return "", err
	}
	defer c.Close()

	i, err := intr.Accept()
	if err != nil {
		return "", err
	}
	defer i.Close()

	return c.RemoteAddr().String()[1:17], nil
}

type Spoofer struct {
	bt      *Btmgmt
	profile *Profile
}

const (
	profilePath = `/bluez/saddupe/hid`
	profileUUID = `00001124-0000-1000-8000-00805f9b34fb`

	// exported with sdptool
	profileServiceRecord = `<?xml version="1.0" encoding="UTF-8" ?><record><attribute id="0x0000"><uint32 value="0x00010000" /></attribute><attribute id="0x0001"><sequence><uuid value="0x1124" /></sequence></attribute><attribute id="0x0004"><sequence><sequence><uuid value="0x0100" /><uint16 value="0x0011" /></sequence><sequence><uuid value="0x0011" /></sequence></sequence></attribute><attribute id="0x0005"><sequence><uuid value="0x1002" /></sequence></attribute><attribute id="0x0006"><sequence><uint16 value="0x656e" /><uint16 value="0x006a" /><uint16 value="0x0100" /></sequence></attribute><attribute id="0x0009"><sequence><sequence><uuid value="0x1124" /><uint16 value="0x0101" /></sequence></sequence></attribute><attribute id="0x000d"><sequence><sequence><sequence><uuid value="0x0100" /><uint16 value="0x0013" /></sequence><sequence><uuid value="0x0011" /></sequence></sequence></sequence></attribute><attribute id="0x0100"><text value="Wireless Gamepad" /></attribute><attribute id="0x0101"><text value="Gamepad" /></attribute><attribute id="0x0102"><text value="Nintendo" /></attribute><attribute id="0x0201"><uint16 value="0x0111" /></attribute><attribute id="0x0202"><uint8 value="0x08" /></attribute><attribute id="0x0203"><uint8 value="0x21" /></attribute><attribute id="0x0204"><boolean value="true" /></attribute><attribute id="0x0205"><boolean value="true" /></attribute><attribute id="0x0206"><sequence><sequence><uint8 value="0x22" /><text encoding="hex" value="05010905a1010601ff8521092175089530810285300930750895308102853109317508966901810285320932750896690181028533093375089669018102853f05091901291015002501750195108102050109391500250775049501814205097504950181010501093009310933093416000027ffff00007510950481020601ff85010901750895309102851009107508953091028511091175089530910285120912750895309102c0" /></sequence></sequence></attribute><attribute id="0x0207"><sequence><sequence><uint16 value="0x0409" /><uint16 value="0x0100" /></sequence></sequence></attribute><attribute id="0x0209"><boolean value="true" /></attribute><attribute id="0x020a"><boolean value="true" /></attribute><attribute id="0x020c"><uint16 value="0x0c80" /></attribute><attribute id="0x020d"><boolean value="false" /></attribute><attribute id="0x020e"><boolean value="false" /></attribute></record>`
)

func NewSpoofer(host string) (*Spoofer, error) {
	bt, err := NewBtmgmt(host)
	if err != nil {
		return nil, err
	}

	options := map[string]interface{}{
		"Role":                  "server",
		"RequireAuthentication": false,
		"RequireAuthorization":  false,
		"ServiceRecord":         profileServiceRecord,
	}

	profile, err := RegisterProfile(profilePath, profileUUID, options)
	if err != nil {
		return nil, err
	}

	// TODO(jpas) Should I disable the agent?

	cmds := [][]string{
		{"power", "off"},
		{"name", "Pro Controller"},
		{"class", "5", "8"},
		{"pairable", "on"},
		{"connectable", "on"},
		{"discov", "off"},
		{"power", "on"},
		{"clr-uuids"},
		{"discov", "limited", "60"},
	}

	for _, cmd := range cmds {
		if _, err := bt.Run(cmd...); err != nil {
			return nil, err
		}
	}

	return &Spoofer{bt, profile}, nil
}

func (s Spoofer) Stop() {
	s.profile.Unregister()
	_ = exec.Command("systemctl", "restart", "bluetooth").Run()
}
