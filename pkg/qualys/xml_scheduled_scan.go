package qualys

import "encoding/xml"

type ScheduleScanListOutput struct {
	XMLName  xml.Name `xml:"SCHEDULE_SCAN_LIST_OUTPUT"`
	Text     string   `xml:",chardata"`
	Response struct {
		Text             string `xml:",chardata"`
		ScheduleScanList struct {
			Text string `xml:",chardata"`
			Scan []struct {
				Text         string `xml:",chardata"`
				ID           string `xml:"ID"`
				Active       string `xml:"ACTIVE"`
				Title        string `xml:"TITLE"`
				UserLogin    string `xml:"USER_LOGIN"`
				Target       string `xml:"TARGET"`
				NetworkID    string `xml:"NETWORK_ID"`
				IScannerName string `xml:"ISCANNER_NAME"`
				EC2Instance  struct {
					Text           string `xml:",chardata"`
					ConnectorUUID  string `xml:"CONNECTOR_UUID"`
					EC2Endpoint    string `xml:"EC2_ENDPOINT"`
					EC2OnlyClassic string `xml:"EC2_ONLY_CLASSIC"`
				} `xml:"EC2_INSTANCE"`
				AssetTags struct {
					Text                string `xml:",chardata"`
					TagsIncludeSelector string `xml:"TAG_INCLUDE_SELECTOR"`
					TagSetInclude       string `xml:"TAG_SET_INCLUDE"`
					UseIPNTRangeTags    string `xml:"USE_IP_NT_RANGE_TAGS"`
				} `xml:"ASSET_TAGS"`
				OptionProfile struct {
					Text        string `xml:",chardata"`
					Title       string `xml:"TITLE"`
					DefaultFlag string `xml:"DEFAULT_FLAG"`
				} `xml:"OPTION_PROFILE"`
				ProcessingPriority string `xml:"PROCESSING_PRIORITY"`
				Schedule           struct {
					Text  string `xml:",chardata"`
					Daily struct {
						Text          string `xml:",chardata"`
						FrequencyDays string `xml:"frequency_days,attr"`
					} `xml:"DAILY"`
					StartDateUTC  string `xml:"START_DATE_UTC"`
					StartHour     string `xml:"START_HOUR"`
					StartMinute   string `xml:"START_MINUTE"`
					NextLaunchUTC string `xml:"NEXTLAUNCH_UTC"`
					TimeZone      struct {
						Text            string `xml:",chardata"`
						TimeZoneCode    string `xml:"TIME_ZONE_CODE"`
						TimeZoneDetails string `xml:"TIME_ZONE_DETAILS"`
					} `xml:"TIME_ZONE"`
					DSTSelected string `xml:"DST_SELECTED"`
				} `xml:"SCHEDULE"`
			} `xml:"SCAN"`
		} `xml:"SCHEDULE_SCAN_LIST"`
	} `xml:"RESPONSE"`
}
