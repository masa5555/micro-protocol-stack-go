package main

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/augustoroman/hexdump"
)

const LOOPBACK_MTU = math.MaxInt16

func InitLoopbackDevice(index uint) (dev *NetDevice, err error) {
	dev = &NetDevice{
		Index: index,
		Type:  NET_DEVICE_TYPE_LOOPBACK,
		Mtu:   LOOPBACK_MTU,
		Flags: NET_DEVICE_FLAG_UP,
		Hlen:  0, // ヘッダーは存在しない
		Alen:  0, // アドレスは存在しない
		Ops: NetDeviceOps{
			nil,
			nil,
			LoopbackDeviceTransmit,
			nil,
		},
	}

	slog.Info(fmt.Sprintf("Initialized loopback device: %s", dev.Name))
	return dev, nil
}

func LoopbackDeviceTransmit(dev *NetDevice, typ uint16, data *[]byte, len uint16, dst func()) error {
	slog.Info("LoopbackDeviceTransmit", "NetDevice=", dev.DebugNetDevice())
	slog.Info(hexdump.Dump(*data))

	return dev.NetInputHandler(typ, data, len, dst)
}
