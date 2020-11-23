package qualys

import "encoding/xml"

//----------------------------------------------
// Appliance API
//----------------------------------------------

// QAppliances holds a list of information pertaining to scanner appliances, which are used to
// assess the security of internal network systems, devices and web applications
type QAppliances struct {
	XMLName    xml.Name     `xml:"APPLIANCE_LIST_OUTPUT"`
	Appliances []QAppliance `xml:"RESPONSE>APPLIANCE_LIST>APPLIANCE"`
}

// QAppliance is a member of QAppliances and must be exported in order to be marshaled
type QAppliance struct {
	XMLName         xml.Name `xml:"APPLIANCE"`
	ID              int      `xml:"ID"`
	UUID            string   `xml:"UUID"`
	Name            string   `xml:"NAME"`
	NetworkID       int      `xml:"NETWORK_ID"`
	SoftwareVersion float32  `xml:"SOFTWARE_VERSION"`
	RunningSlices   int      `xml:"RUNNING_SLICES_COUNT"`
	RunningScans    int      `xml:"RUNNING_SCAN_COUNT"`
	Status          string   `xml:"STATUS"`
}
