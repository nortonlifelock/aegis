package qualys

import (
	"encoding/xml"
)

// QSAGListOutput holds a list of asset groups, including their relevant appliances and IPs tracked
type QSAGListOutput struct {
	XMLName xml.Name        `xml:"ASSET_GROUP_LIST_OUTPUT"`
	Date    string          `xml:"RESPONSE>DATETIME"`
	Groups  []*QSAssetGroup `xml:"RESPONSE>ASSET_GROUP_LIST>ASSET_GROUP"`
}

// QSAssetGroup is a member of qsAGListOutput and must be exported in order to be marshaled
type QSAssetGroup struct {
	XMLName          xml.Name `xml:"ASSET_GROUP"`
	ID               int      `xml:"ID"`
	Title            *CData   `xml:"TITLE"`
	OwnerID          int      `xml:"OWNER_ID"`
	UnitID           int      `xml:"UNIT_ID"`
	LastUpdated      string   `xml:"LAST_UPDATE"`
	BusinessImpact   string   `xml:"BUSINESS_IMPACT"`
	Appliances       string   `xml:"APPLIANCE_IDS"`
	NetworkID        int      `xml:"NETWORK_ID"`
	Ranges           []string `xml:"IP_SET>IP_RANGE"`
	IPs              []string `xml:"IP_SET>IP"`
	OnlineAppliances []int
}
