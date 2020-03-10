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

// TicketSummary defines the struct that implements the TicketSummary interface
type TicketSummary struct {
	CreatedDatevar    *time.Time
	DetectionIDvar    string
	DueDatevar        time.Time
	OrganizationIDvar string
	ResolutionDatevar *time.Time
	Statusvar         string
	Titlevar          string
	UpdatedDatevar    *time.Time
}

//**********************************************************
// Struct Methods
//**********************************************************

// MarshalJSON marshals the struct by converting it to a map
func (myTicketSummary TicketSummary) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"CreatedDate":    myTicketSummary.CreatedDatevar,
		"DetectionID":    myTicketSummary.DetectionIDvar,
		"DueDate":        myTicketSummary.DueDatevar,
		"OrganizationID": myTicketSummary.OrganizationIDvar,
		"ResolutionDate": myTicketSummary.ResolutionDatevar,
		"Status":         myTicketSummary.Statusvar,
		"Title":          myTicketSummary.Titlevar,
		"UpdatedDate":    myTicketSummary.UpdatedDatevar,
	})
}

// CreatedDate returns the CreatedDate parameter from the TicketSummary struct
func (myTicketSummary *TicketSummary) CreatedDate() (param *time.Time) {
	return myTicketSummary.CreatedDatevar
}

// DetectionID returns the DetectionID parameter from the TicketSummary struct
func (myTicketSummary *TicketSummary) DetectionID() (param string) {
	return myTicketSummary.DetectionIDvar
}

// DueDate returns the DueDate parameter from the TicketSummary struct
func (myTicketSummary *TicketSummary) DueDate() (param time.Time) {
	return myTicketSummary.DueDatevar
}

// OrganizationID returns the OrganizationID parameter from the TicketSummary struct
func (myTicketSummary *TicketSummary) OrganizationID() (param string) {
	return myTicketSummary.OrganizationIDvar
}

// ResolutionDate returns the ResolutionDate parameter from the TicketSummary struct
func (myTicketSummary *TicketSummary) ResolutionDate() (param *time.Time) {
	return myTicketSummary.ResolutionDatevar
}

// Status returns the Status parameter from the TicketSummary struct
func (myTicketSummary *TicketSummary) Status() (param string) {
	return myTicketSummary.Statusvar
}

// Title returns the Title parameter from the TicketSummary struct
func (myTicketSummary *TicketSummary) Title() (param string) {
	return myTicketSummary.Titlevar
}

// UpdatedDate returns the UpdatedDate parameter from the TicketSummary struct
func (myTicketSummary *TicketSummary) UpdatedDate() (param *time.Time) {
	return myTicketSummary.UpdatedDatevar
}

// SetDetectionID sets the DetectionID parameter from the TicketSummary struct
func (myTicketSummary *TicketSummary) SetDetectionID(val string) {
	myTicketSummary.DetectionIDvar = val
}

// SetOrganizationID sets the OrganizationID parameter from the TicketSummary struct
func (myTicketSummary *TicketSummary) SetOrganizationID(val string) {
	myTicketSummary.OrganizationIDvar = val
}

// SetStatus sets the Status parameter from the TicketSummary struct
func (myTicketSummary *TicketSummary) SetStatus(val string) {
	myTicketSummary.Statusvar = val
}

// SetTitle sets the Title parameter from the TicketSummary struct
func (myTicketSummary *TicketSummary) SetTitle(val string) {
	myTicketSummary.Titlevar = val
}
