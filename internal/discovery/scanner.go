package discovery

import (
	"github.com/hui1601/Quicky/internal/constant"
	"github.com/hui1601/Quicky/internal/utils"
	"tinygo.org/x/bluetooth"
)

type ScanResult struct {
	Address         bluetooth.Address
	RSSI            int16
	Name            string
	Advertisement   *utils.AdvertisementInfo
}

type Scanner struct {
	adapter *bluetooth.Adapter
}

func NewScanner(adapter *bluetooth.Adapter) *Scanner {
	return &Scanner{adapter: adapter}
}

func (s *Scanner) Scan(callback func(ScanResult)) error {
	return s.adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		for _, mfr := range result.ManufacturerData() {
			if mfr.CompanyID != constant.QCYCompanyID {
				continue
			}
			info, err := utils.ParseManufacturerData(mfr.Data)
			if err != nil {
				continue
			}
			callback(ScanResult{
				Address:       result.Address,
				RSSI:          result.RSSI,
				Name:          result.LocalName(),
				Advertisement: info,
			})
			return
		}
	})
}

func (s *Scanner) StopScan() error {
	return s.adapter.StopScan()
}
