package main

import (
	"os"
	"os/signal"
	"time"
)

var test_data = []byte{
	0x45, 0x00, 0x00, 0x30,
	0x00, 0x80, 0x00, 0x00,
	0xff, 0x01, 0xbd, 0x4a,
	0x7f, 0x00, 0x00, 0x01,
	0x7f, 0x00, 0x00, 0x01,
	0x08, 0x00, 0x35, 0x64,
	0x00, 0x80, 0x00, 0x01,
	0x31, 0x32, 0x33, 0x34,
	0x35, 0x36, 0x37, 0x38,
	0x39, 0x30, 0x21, 0x40,
	0x23, 0x24, 0x25, 0x5e,
	0x26, 0x2a, 0x28, 0x29,
}

func main() {
	var devices *NetDevice

	var i uint = 0
	for i < 2 {
		dev, err := InitLoopbackDevice(i)
		if err != nil {
			panic(err)
		}
		devices, err = NetDeviceRegister(devices, dev)
		if err != nil {
			panic(err)
		}
		i++
	}

	// プロトコルスタックの起動
	err := devices.NetRun()
	if err != nil {
		panic(err)
	}

	// Ctrl+Cが押されるまで待機
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt)

	ticker := time.NewTicker(1 * time.Second)
loop:
	for {
		select {
		case <-ticker.C:
			err := devices.NetDeviceOutput(&test_data, uint16(len(test_data)), nil)
			if err != nil {
				panic(err)
			}
		case <-sigs:
			break loop
		}
	}

	// プロトコルスタックの終了
	err = devices.NetShutdown()
	if err != nil {
		panic(err)
	}
}

func GetHello() string {
	return "Hello, world!"
}
