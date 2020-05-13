package internal

import (
	"github.com/muka/go-bluetooth/bluez/profile/profile"
	"github.com/godbus/dbus"
)

type Profile struct {
	m *profile.ProfileManager1
}

func RegisterProfile(path, uuid, record string) (*Profile, error) {
	m, err := profile.NewProfileManager1()
	if err != nil {
		return nil, err
	}

	err = m.RegisterProfile(
		dbus.ObjectPath(path),
		uuid,
		map[string]interface{}{
			"Role":                  "server",
			"RequireAuthentication": false,
			"RequireAuthorization":  false,
			"ServiceRecord":         record,
		},
	)
	if err != nil {
		m.Close()
		return nil, err
	}

	return &Profile{m}, nil
}

func (p *Profile) Unregister() {
	p.m.Close()
}
