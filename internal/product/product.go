package product

import (
	_ "embed"
	"encoding/json"
	"strconv"
)

//go:embed products.json
var productsJSON []byte

type Product struct {
	VendorId uint16   `json:"vendorId"`
	Title    string   `json:"title"`
	SubTitle string   `json:"subTitle"`
	Category string   `json:"category"`
	Features Features `json:"features,omitempty"`
}

type Features struct {
	ANC            *ANCFeature      `json:"anc,omitempty"`
	EQ             *EQFeature       `json:"eq,omitempty"`
	KeyFunction    *KeyFuncFeature  `json:"key_function,omitempty"`
	ChannelBalance bool             `json:"channel_balance,omitempty"`
	FindEarphone   bool             `json:"find_earphone,omitempty"`
	DeviceName     bool             `json:"device_name,omitempty"`
	AutoOffTimer   *AutoOffFeature  `json:"auto_off_timer,omitempty"`
	Settings       []SettingItem    `json:"settings,omitempty"`
}

type ANCFeature struct {
	Modes []ANCMode `json:"modes"`
}

type ANCMode struct {
	Name       string     `json:"name"`
	StartCmdID int        `json:"startcmdid"`
	EndCmdID   int        `json:"endcmdid"`
	DefaultCmd int        `json:"defaultcmd"`
	ViewType   int        `json:"viewtype"`
	Items      []ANCItem  `json:"items,omitempty"`
}

type ANCItem struct {
	Name       string `json:"name"`
	StartCmdID int    `json:"startcmdid"`
	EndCmdID   int    `json:"endcmdid"`
}

type EQFeature struct {
	Bands          int      `json:"bands"`
	MinDB          int      `json:"mindb"`
	MaxDB          int      `json:"maxdb"`
	Freq           string   `json:"freq"`
	Characteristic string   `json:"characteristic"`
	Presets        []string `json:"presets"`
}

type KeyFuncFeature struct {
	Events []KeyEvent `json:"events"`
}

type KeyEvent struct {
	Name      string   `json:"name"`
	Functions []string `json:"functions"`
}

type AutoOffFeature struct {
	CmdID  int `json:"cmdid"`
	Repeat int `json:"repeat"`
}

type SettingItem struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	CmdID *int   `json:"cmdid,omitempty"`
	Cmd   string `json:"cmd,omitempty"`
}

var productDatabase map[string]Product

func init() {
	productDatabase = make(map[string]Product)
	if err := json.Unmarshal(productsJSON, &productDatabase); err != nil {
		panic("failed to load product database: " + err.Error())
	}
}

func Lookup(vendorId uint16) (*Product, bool) {
	key := strconv.Itoa(int(vendorId))
	p, ok := productDatabase[key]
	if !ok {
		return nil, false
	}
	return &p, true
}

func LookupByString(vendorIdStr string) (*Product, bool) {
	p, ok := productDatabase[vendorIdStr]
	if !ok {
		return nil, false
	}
	return &p, true
}

func Count() int {
	return len(productDatabase)
}
