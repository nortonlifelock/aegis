package qualys

import "encoding/xml"

// ScanSummaryOutput contains a list of scan summary information returned from Qualys, which is useful for discovering
// which hosts were found dead by a scanner.
type ScanSummaryOutput struct {
	XMLName  xml.Name `xml:"SCAN_SUMMARY_OUTPUT"`
	Text     string   `xml:",chardata"`
	Response struct {
		Text            string `xml:",chardata"`
		DateTime        string `xml:"DATETIME"`
		ScanSummaryList struct {
			Text        string        `xml:",chardata"`
			ScanSummary []ScanSummary `xml:"SCAN_SUMMARY"`
		} `xml:"SCAN_SUMMARY_LIST"`
	} `xml:"RESPONSE"`
}

// ScanSummary is a member of ScanSummaryOutput and must be marshaled in order to be exported
type ScanSummary struct {
	Text        string `xml:",chardata"`
	ScanRef     string `xml:"SCAN_REF"`
	ScanDate    string `xml:"SCAN_DATE"`
	HostSummary []struct {
		Text     string `xml:",chardata"`
		Category string `xml:"category,attr"`
		Tracking string `xml:"tracking,attr"`
	} `xml:"HOST_SUMMARY"`
}
