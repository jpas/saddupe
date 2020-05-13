package main

import (
	"github.com/jpas/saddupe/internal"
	"github.com/spf13/cobra"
	"log"
)

var pairCmd = &cobra.Command{
	Use:   "pair",
	Short: "Pairs with a alternating device over Bluetooth",
	Run:   pairRun,
}

func pairRun(cmd *cobra.Command, args []string) {
	host, err := NewBtAddr("80:32:53:37:22:19")
	if err != nil {
		log.Fatal(err)
	}

	if err := Pair(host); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(pairCmd)
}

func Pair(host *BtAddr) error {
	ctrl, err := NewL2Socket()
	if err != nil {
		return err
	}
	defer ctrl.Close()

	intr, err := NewL2Socket()
	if err != nil {
		return err
	}
	defer intr.Close()

	if err := listenOnSockets(ctrl, intr); err != nil {
		return err
	}

	spoof, err := StartSpoofing(host)
	if err != nil {
		return err
	}
	defer spoof.Stop()
	log.Println("started controller spoofing")

	c, err := ctrl.Accept()
	if err != nil {
		return err
	}
	defer c.Close()
	log.Println("accepted control channel")

	i, err := intr.Accept()
	if err != nil {
		return err
	}
	defer i.Close()
	log.Println("accepted interrupt channel")

	log.Printf("paired with: %s", c.RemoteAddr().Addr)

	// make socket on port 17, then 19
	// probably bind to the one we want, then to the switch
	// expect binding locally to not work, restart bluetooth and try again
	// restarting bluetooth is a hack, can disable input to do it, but we
	// only need this for pairing

	return nil
}

func listenOnSockets(ctrl, intr *L2Socket) error {
	log.Println("restarting bluetooth service")
	if err := internal.RestartBluetooth(); err != nil {
		return err
	}

	if err := ctrl.Bind(&L2Addr{BtAddrAny, 17}); err != nil {
		return err
	}

	if err := intr.Bind(&L2Addr{BtAddrAny, 19}); err != nil {
		return err
	}

	if err := ctrl.Listen(1); err != nil {
		return err
	}

	if err := intr.Listen(1); err != nil {
		return err
	}

	return nil
}

type Spoof struct {
	bt      *internal.Btmgmt
	profile *internal.Profile
}

const (
	profilePath = `/bluez/saddupe/hid`
	profileUUID = `00001124-0000-1000-8000-00805f9b34fb`

	// exported with sdptool
	profileServiceRecord = `<?xml version="1.0" encoding="UTF-8" ?><record><attribute id="0x0000"><uint32 value="0x00010000" /></attribute><attribute id="0x0001"><sequence><uuid value="0x1124" /></sequence></attribute><attribute id="0x0004"><sequence><sequence><uuid value="0x0100" /><uint16 value="0x0011" /></sequence><sequence><uuid value="0x0011" /></sequence></sequence></attribute><attribute id="0x0005"><sequence><uuid value="0x1002" /></sequence></attribute><attribute id="0x0006"><sequence><uint16 value="0x656e" /><uint16 value="0x006a" /><uint16 value="0x0100" /></sequence></attribute><attribute id="0x0009"><sequence><sequence><uuid value="0x1124" /><uint16 value="0x0101" /></sequence></sequence></attribute><attribute id="0x000d"><sequence><sequence><sequence><uuid value="0x0100" /><uint16 value="0x0013" /></sequence><sequence><uuid value="0x0011" /></sequence></sequence></sequence></attribute><attribute id="0x0100"><text value="Wireless Gamepad" /></attribute><attribute id="0x0101"><text value="Gamepad" /></attribute><attribute id="0x0102"><text value="Nintendo" /></attribute><attribute id="0x0201"><uint16 value="0x0111" /></attribute><attribute id="0x0202"><uint8 value="0x08" /></attribute><attribute id="0x0203"><uint8 value="0x21" /></attribute><attribute id="0x0204"><boolean value="true" /></attribute><attribute id="0x0205"><boolean value="true" /></attribute><attribute id="0x0206"><sequence><sequence><uint8 value="0x22" /><text encoding="hex" value="05010905a1010601ff8521092175089530810285300930750895308102853109317508966901810285320932750896690181028533093375089669018102853f05091901291015002501750195108102050109391500250775049501814205097504950181010501093009310933093416000027ffff00007510950481020601ff85010901750895309102851009107508953091028511091175089530910285120912750895309102c0" /></sequence></sequence></attribute><attribute id="0x0207"><sequence><sequence><uint16 value="0x0409" /><uint16 value="0x0100" /></sequence></sequence></attribute><attribute id="0x0209"><boolean value="true" /></attribute><attribute id="0x020a"><boolean value="true" /></attribute><attribute id="0x020c"><uint16 value="0x0c80" /></attribute><attribute id="0x020d"><boolean value="false" /></attribute><attribute id="0x020e"><boolean value="false" /></attribute></record>`
)

func StartSpoofing(host *BtAddr) (*Spoof, error) {
	bt, err := internal.NewBtmgmt(host.String())
	if err != nil {
		return nil, err
	}

	profile, err := internal.RegisterProfile(profilePath, profileUUID, profileServiceRecord)
	if err != nil {
		return nil, err
	}

	cmds := [][]string{
		[]string{"power", "off"},
		[]string{"le", "off"},
		[]string{"name", "Pro Controller"},
		[]string{"linksec", "off"},
		[]string{"class", "5", "8"},
		[]string{"pairable", "on"},
		[]string{"connectable", "on"},
		[]string{"discov", "off"},
		[]string{"power", "on"},
		[]string{"clr-uuids"},
		[]string{"discov", "limited", "60"},
	}

	for _, cmd := range cmds {
		if _, err := bt.Run(cmd...); err != nil {
			return nil, err
		}
	}

	return &Spoof{bt, profile}, nil
}

func (s *Spoof) Stop() {
	s.profile.Unregister()
	internal.RestartBluetooth()
	// TODO(jpas) we _could_ restart the bluetooth service here to reset controller state.
}
