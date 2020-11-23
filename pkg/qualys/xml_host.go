package qualys

import "encoding/xml"

// HostListOutput is the struct for mapping Host List data from the Qualys API
type HostListOutput struct {
	XMLName  xml.Name `xml:"HOST_LIST_OUTPUT"`
	Text     string   `xml:",chardata"`
	Response struct {
		Text     string `xml:",chardata"`
		HostList struct {
			Text string `xml:",chardata"`
			Host []struct {
				Text           string `xml:",chardata"`
				ID             string `xml:"ID"`
				IP             string `xml:"IP"`
				NetworkID      string `xml:"NETWORK_ID"`
				AssetGroupIDs  string `xml:"ASSET_GROUP_IDS"`
				TrackingMethod string `xml:"TRACKING_METHOD"` // e.g. IP/EC2
			} `xml:"HOST"`
		} `xml:"HOST_LIST"`
	} `xml:"RESPONSE"`
}
