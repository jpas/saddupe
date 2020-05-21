package internal

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez/profile/profile"
)

// Profile manages a Service Discovery Protocol profile
type Profile struct {
	m *profile.ProfileManager1
}

// RegisterProfile registers a Service Discovery Protocol profile
func RegisterProfile(path string, uuid string, options map[string]interface{}) (*Profile, error) {
	m, err := profile.NewProfileManager1()
	if err != nil {
		return nil, err
	}

	err = m.RegisterProfile(dbus.ObjectPath(path), uuid, options)
	if err != nil {
		m.Close()
		return nil, err
	}

	return &Profile{m}, nil
}

// Unregister removes the profile from registration
func (p *Profile) Unregister() {
	p.m.Close()
}
