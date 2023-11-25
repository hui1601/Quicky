package utils

import (
	"strings"
)

func PathToMac(path string) string {
	if !strings.HasPrefix(path, "/org/bluez/hci") {
		return ""
	}
	if !strings.Contains(path, "/dev_") {
		return ""
	}
	res := strings.ReplaceAll(path, "_", ":")[strings.LastIndex(path, "/")+1:]
	res = strings.ReplaceAll(res, "dev:", "")
	if len(res) != 17 {
		return ""
	}
	return res
}

func MacToPath(mac string) string {
	return strings.ReplaceAll(mac, ":", "_")
}
