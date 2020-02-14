package dal

//**********************************************************
// GENERATED CODE - DO NOT CHANGE
// This file is generated using scaffolding. Any changes to
// this file will be overwritten on the next build
//**********************************************************

import (
	"encoding/json"
	"time"
)

//**********************************************************
// Struct Declaration
//**********************************************************

// AssignmentGroup defines the struct that implements the AssignmentGroup interface
type AssignmentGroup struct {
	DBCreatedDatevar  time.Time
	DBUpdatedDatevar  *time.Time
	GroupNamevar      string
	IPAddressvar      string
	OrganizationIDvar string
	SourceIDvar       int
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myAssignmentGroup AssignmentGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"DBCreatedDate":  myAssignmentGroup.DBCreatedDatevar,
		"DBUpdatedDate":  myAssignmentGroup.DBUpdatedDatevar,
		"GroupName":      myAssignmentGroup.GroupNamevar,
		"IPAddress":      myAssignmentGroup.IPAddressvar,
		"OrganizationID": myAssignmentGroup.OrganizationIDvar,
		"SourceID":       myAssignmentGroup.SourceIDvar,
	})
}

// DBCreatedDate returns the DBCreatedDate parameter from the AssignmentGroup struct
func (myAssignmentGroup *AssignmentGroup) DBCreatedDate() (param time.Time) {
	return myAssignmentGroup.DBCreatedDatevar
}

// DBUpdatedDate returns the DBUpdatedDate parameter from the AssignmentGroup struct
func (myAssignmentGroup *AssignmentGroup) DBUpdatedDate() (param *time.Time) {
	return myAssignmentGroup.DBUpdatedDatevar
}

// GroupName returns the GroupName parameter from the AssignmentGroup struct
func (myAssignmentGroup *AssignmentGroup) GroupName() (param string) {
	return myAssignmentGroup.GroupNamevar
}

// IPAddress returns the IPAddress parameter from the AssignmentGroup struct
func (myAssignmentGroup *AssignmentGroup) IPAddress() (param string) {
	return myAssignmentGroup.IPAddressvar
}

// OrganizationID returns the OrganizationID parameter from the AssignmentGroup struct
func (myAssignmentGroup *AssignmentGroup) OrganizationID() (param string) {
	return myAssignmentGroup.OrganizationIDvar
}

// SourceID returns the SourceID parameter from the AssignmentGroup struct
func (myAssignmentGroup *AssignmentGroup) SourceID() (param int) {
	return myAssignmentGroup.SourceIDvar
}

// SetGroupName sets the GroupName parameter from the AssignmentGroup struct
func (myAssignmentGroup *AssignmentGroup) SetGroupName(val string) {
	myAssignmentGroup.GroupNamevar = val
}

// SetIPAddress sets the IPAddress parameter from the AssignmentGroup struct
func (myAssignmentGroup *AssignmentGroup) SetIPAddress(val string) {
	myAssignmentGroup.IPAddressvar = val
}

// SetOrganizationID sets the OrganizationID parameter from the AssignmentGroup struct
func (myAssignmentGroup *AssignmentGroup) SetOrganizationID(val string) {
	myAssignmentGroup.OrganizationIDvar = val
}
