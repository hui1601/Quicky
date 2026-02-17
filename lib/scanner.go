package quicky

import (
	"github.com/hui1601/Quicky/internal/discovery"
	"github.com/hui1601/Quicky/internal/product"
	"github.com/hui1601/Quicky/internal/utils"
	"tinygo.org/x/bluetooth"
)

type AdvertisementInfo = utils.AdvertisementInfo
type Product = product.Product

type ScanResult struct {
	Address       bluetooth.Address
	RSSI          int16
	Name          string
	Advertisement *AdvertisementInfo
}

func (r ScanResult) GetProductInfo() (*Product, bool) {
	if r.Advertisement == nil {
		return nil, false
	}
	return product.Lookup(r.Advertisement.VendorID)
}

type Scanner struct {
	s *discovery.Scanner
}

func NewScanner() *Scanner {
	return &Scanner{s: discovery.NewScanner(bluetooth.DefaultAdapter)}
}

func NewScannerWithAdapter(adapter *bluetooth.Adapter) *Scanner {
	return &Scanner{s: discovery.NewScanner(adapter)}
}

func (s *Scanner) Scan(callback func(ScanResult)) error {
	return s.s.Scan(func(result discovery.ScanResult) {
		callback(ScanResult{
			Address:       result.Address,
			RSSI:          result.RSSI,
			Name:          result.Name,
			Advertisement: result.Advertisement,
		})
	})
}

func (s *Scanner) StopScan() error {
	return s.s.StopScan()
}
