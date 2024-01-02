package main

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/augustoroman/hexdump"
)

func (devices *NetDevice) NetRun() error {
	slog.Info("Open all devices...")
	// 登録済みの全デバイスを開く
	for devices != nil {
		err := devices.NetDeviceOpen()
		if err != nil {
			slog.Error("NetDeviceOpen failed: %s", err)
		}
		devices = devices.Next
	}
	slog.Info("running...")
	return nil
}

func (devices *NetDevice) NetShutdown() error {
	slog.Info("Close all devices...")
	// 登録済みの全デバイスを閉じる
	for devices != nil {
		err := devices.NetDeviceClose()
		if err != nil {
			slog.Error("NetDeviceClose failed: %s", err)
		}
		devices = devices.Next
	}
	slog.Info("shutdown.")
	return nil
}

// デバイス構造体を登録する
func NetDeviceRegister(devices *NetDevice, dev *NetDevice) (newDevices *NetDevice, err error) {
	// 登録するときに名前を生成する
	// net0 -> net1 -> net2 -> net3 -> net4 -> net5
	dev.Name = [IFNAMSIZ]byte{'n', 'e', 't', byte(fmt.Sprintf("%d", dev.Index)[0])}

	dev.Next = *(&devices)
	newDevices = *(&dev)
	slog.Info("Registered", "NetDevice=", dev.DebugNetDevice())
	return newDevices, nil
}

func (dev *NetDevice) NetDeviceOpen() error {
	if dev.IsUp() {
		var errMsg = fmt.Sprintf("already up: %s", dev.Name)
		return errors.New(errMsg)
	}
	if dev.Ops.Open != nil {
		err := dev.Ops.Open(dev)
		if err != nil {
			return err
		}
	}

	// upフラグを立てる
	dev.Flags |= NET_DEVICE_FLAG_UP
	slog.Info("Device Opened", "NetDevice=", dev.DebugNetDevice())
	return nil
}

func (dev *NetDevice) NetDeviceClose() error {
	if !dev.IsUp() {
		var errMsg = fmt.Sprintf("already down: %s", dev.Name)
		return errors.New(errMsg)
	}
	if dev.Ops.Close != nil {
		err := dev.Ops.Close(dev)
		if err != nil {
			return err
		}
	}
	// upフラグを落とす
	dev.Flags &= ^NET_DEVICE_FLAG_UP
	slog.Info("Device Closed", "NetDevice=", dev.DebugNetDevice())
	return nil
}

func (dev *NetDevice) NetDeviceOutput(data *[]byte, len uint16, dst func()) error {
	if !dev.IsUp() {
		return errors.New("device is down")
	}
	if len > dev.Mtu {
		return errors.New("too long packet")
	}
	hexdump.Dump(*data)
	if dev.Ops.Transmit != nil {
		err := dev.Ops.Transmit(dev, dev.Type, data, len, dst)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dev *NetDevice) NetInputHandler(typ uint16, data *[]byte, len uint16, dst func()) error {
	slog.Info("InputHandler", "NetDevice=", dev.DebugNetDevice())
	slog.Info(hexdump.Dump(*data))
	return nil
}
