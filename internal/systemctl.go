package internal

import (
	"github.com/pkg/errors"
	"os/exec"
	"time"
)

func RestartBluetooth() error {
	if err := exec.Command("systemctl", "restart", "bluetooth.service").Run(); err != nil {
		return errors.Wrap(err, "systemctl restart failed")
	}

	time.Sleep(1 * time.Second)

	return nil
}
