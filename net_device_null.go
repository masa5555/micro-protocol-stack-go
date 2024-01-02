package main

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/augustoroman/hexdump"
)

const NULL_MTU = math.MaxInt16

func InitNullDevice(index uint) (dev *NetDevice, err error) {
	dev = &NetDevice{
		Index: index,
		Type:  NET_DEVICE_TYPE_NULL,
		Mtu:   NULL_MTU,
		Flags: 0x0000,
		Hlen:  0, // ヘッダーは存在しない
		Alen:  0, // アドレスは存在しない
		Ops: NetDeviceOps{
			nil,
			nil,
			NullDeviceTransmit,
			nil,
		},
	}

	slog.Info(fmt.Sprintf("Initialized null device: %s", dev.Name))
	return dev, nil
}

func NullDeviceTransmit(dev *NetDevice, typ uint16, data *[]byte, len uint16, dst func()) error {
	slog.Info("NullDeviceTransmit", "NetDevice=", dev.DebugNetDevice())
	slog.Info(hexdump.Dump(*data))
	/* drop data */
	return nil
}
