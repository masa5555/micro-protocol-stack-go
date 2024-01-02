package main

import (
	"fmt"
	"strings"
)

// max length of interface name
const IFNAMSIZ = 16

// max length of mac address
const NET_DEVICE_ADDR_LEN = 6

const NET_DEVICE_TYPE_NULL uint16 = 0

const NET_DEVICE_FLAG_UP uint16 = 0x0001

type NetDevice struct {
	Next  *NetDevice
	Index uint
	Name  [IFNAMSIZ]byte
	Type  uint16
	/* maximum transmission unit */
	Mtu   uint16
	Flags uint16
	/* header length */
	Hlen uint16
	/* address length */
	Alen uint16
	Addr [NET_DEVICE_ADDR_LEN]byte
	// hardware address of device
	Peer      [NET_DEVICE_ADDR_LEN]byte
	Broadcast [NET_DEVICE_ADDR_LEN]byte
	Ops       NetDeviceOps
	Priv      func()
}

// デバイスドライバに実装されている関数へのポインタを格納する構造体
type NetDeviceOps struct {
	Open     func(dev *NetDevice) error
	Close    func(dev *NetDevice) error
	Transmit func(dev *NetDevice, typ uint16, data *[]byte, len uint16, dst func()) error
	Poll     func(*NetDevice) error
}

func (dev *NetDevice) IsUp() bool {
	return (dev.Flags & NET_DEVICE_FLAG_UP) != 0
}

func (dev *NetDevice) DebugNetDevice() string {
	sprint := fmt.Sprintf("dev=%s type=%d mtu=%d flags=%d hlen=%d alen=%d", dev.Name, dev.Type, dev.Mtu, dev.Flags, dev.Hlen, dev.Alen)
	// null を除外
	return strings.Replace(sprint, "\x00", "", -1)
}
