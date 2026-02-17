# QuiCkY

[English](README.md) | [한국어](README.ko.md)

Go library for controlling QCY Bluetooth earphones on Linux via BLE GATT.

> **Proof of concept.** No guarantee it will work properly. Use at your own risk.

## Install

```bash
go get github.com/hui1601/Quicky
```

Requires BlueZ on Linux.

## Usage

### Scanning for Devices

```go
package main

import (
	"fmt"

	quicky "github.com/hui1601/Quicky/lib"
)

func main() {
	scanner := quicky.NewScanner()
	scanner.Scan(func(result quicky.ScanResult) {
		fmt.Printf("Found: %s (RSSI: %d)\n", result.Address.String(), result.RSSI)
		if result.Advertisement != nil {
			fmt.Printf("  Battery: L=%d%% R=%d%% Box=%d%%\n",
				result.Advertisement.LeftBattery,
				result.Advertisement.RightBattery,
				result.Advertisement.BoxBattery)
			fmt.Printf("  Control MAC: %s\n", result.Advertisement.ControlMAC)
			
			if product, ok := result.GetProductInfo(); ok {
				fmt.Printf("  Model: %s\n", product.Title)
				if product.Features.ANC != nil {
					fmt.Printf("    ANC: %d modes\n", len(product.Features.ANC.Modes))
				}
			}
		}
	})
}
```

### Connecting and Controlling

```go
package main

import (
	"context"
	"fmt"
	"time"

	quicky "github.com/hui1601/Quicky/lib"
)

func main() {
	client, _ := quicky.New("AA:BB:CC:DD:EE:FF")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Connect(ctx)
	defer client.Disconnect()

	battery, _ := client.ReadBattery()
	fmt.Printf("Battery: L=%d%% R=%d%% Box=%d%%\n",
		battery.Left.Level, battery.Right.Level, battery.Box.Level)

	client.SetNoiseCancelMode(quicky.NoiseCancelANC)
	client.SetLowLatency(true)
	client.SetVolume(80, 80)
}
```

### Receiving Events

```go
go func() {
	for ev := range client.Events() {
		fmt.Printf("Event: type=%d cmdID=0x%02x\n", ev.Type, ev.CmdID)
	}
}()
```

## Features

- **Discovery** — Scan for QCY devices via BLE manufacturer data (CompanyID `0x521c`), parse battery levels, charging state, and MAC addresses from advertisements
- **Model identification** — Embedded product database (199 models) maps vendorId to model name and supported features (ANC, EQ, etc.)
- **40+ commands** — Noise cancellation, EQ (v1/v2/per-channel), key function mapping, LED effects, spatial audio, wearing detection, alarms, music control, and more
- **Event stream** — Async notification handling with typed event parsing
- **Direct reads** — Battery and firmware version via characteristic reads

## Documentation

- [Protocol Reference](docs/protocol.md)
- [GATT Service UUIDs](docs/service.md)