# QuiCkY

[English](README.md) | [한국어](README.ko.md)

BLE GATT를 통해 리눅스에서 QCY 블루투스 이어폰을 제어하는 Go 라이브러리.

> **개념 증명 단계입니다.** 정상 작동을 보장하지 않습니다. 사용에 따른 책임은 본인에게 있습니다.

## 설치

```bash
go get github.com/hui1601/Quicky
```

리눅스에서 BlueZ가 필요합니다.

## 사용법

### 디바이스 스캔

```go
package main

import (
	"fmt"

	quicky "github.com/hui1601/Quicky/lib"
)

func main() {
	scanner := quicky.NewScanner()
	scanner.Scan(func(result quicky.ScanResult) {
		fmt.Printf("발견: %s (RSSI: %d)\n", result.Address.String(), result.RSSI)
		if result.Advertisement != nil {
			fmt.Printf("  배터리: 좌=%d%% 우=%d%% 케이스=%d%%\n",
				result.Advertisement.LeftBattery,
				result.Advertisement.RightBattery,
				result.Advertisement.BoxBattery)
			fmt.Printf("  제어 MAC: %s\n", result.Advertisement.ControlMAC)
		}
	})
}
```

### 연결 및 제어

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
	fmt.Printf("배터리: 좌=%d%% 우=%d%% 케이스=%d%%\n",
		battery.Left.Level, battery.Right.Level, battery.Box.Level)

	client.SetNoiseCancelMode(quicky.NoiseCancelANC)
	client.SetLowLatency(true)
	client.SetVolume(80, 80)
}
```

### 이벤트 수신

```go
go func() {
	for ev := range client.Events() {
		fmt.Printf("이벤트: type=%d cmdID=0x%02x\n", ev.Type, ev.CmdID)
	}
}()
```

## 기능

- **디바이스 탐색** — BLE 제조사 데이터(CompanyID `0x521c`)로 QCY 기기 스캔, 광고 패킷에서 배터리 잔량, 충전 상태, MAC 주소 파싱
- **40개 이상의 명령** — 노이즈 캔슬링, EQ (v1/v2/채널별), 키 기능 매핑, LED 효과, 공간 음향, 착용 감지, 알람, 음악 제어 등
- **이벤트 스트림** — 타입이 지정된 이벤트 파싱을 통한 비동기 알림 처리
- **직접 읽기** — 특성(Characteristic) 읽기를 통한 배터리 및 펌웨어 버전 확인

## 문서

- [프로토콜 레퍼런스](docs/protocol.md)
- [GATT 서비스 UUID](docs/service.md)