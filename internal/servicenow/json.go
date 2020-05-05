package servicenow

import "fmt"

func (e Err) Error() string {
	//if e.Status == "failure" {
	//	return e.ResError.Detail
	//}
	return fmt.Sprintf("Error: %s, Details: %s", e.ResError.Message, e.ResError.Detail)
}

// Err represents a possible error message that came back from the server
type Err struct {
	ResError struct {
		Detail  string `json:"detail"`
		Message string `json:"message"`
	} `json:"error"`
	Status string `json:"status"`
}

// SvcNowTickets parses the response from the service now server with tickets matching a query
type SvcNowTickets struct {
	Results []*Result `json:"result"`
}

// SvcNowTicketDto is used to parse a single service now ticket
type SvcNowTicketDto struct {
	Result *Result `json:"result"`
}

//type Result struct {
//	Parent                string `json:"parent,omitempty"`
//	UVendorName           string `json:"u_vendor_name,omitempty"`
//	MadeSLA               string `json:"made_sla,omitempty"`
//	WatchList             string `json:"watch_list,omitempty"`
//	UponReject            string `json:"upon_reject,omitempty"`
//	SysUpdatedOn          string `json:"sys_updated_on,omitempty"`
//	UTargetImplementation string `json:"u_target_implementation,omitempty"`
//	ApprovalHistory       string `json:"approval_history,omitempty"`
//	Skills                string `json:"skills,omitempty"`
//	Number                string `json:"number,omitempty"`
//	UBusinessService      string `json:"u_business_service,omitempty"`
//	UTechnicalService     string `json:"u_technical_service,omitempty"`
//	SysUpdatedBy          string `json:"sys_updated_by,omitempty"`
//	PhiHistory            string `json:"phi_history,omitempty"`
//	OpenedBy struct {
//		DisplayValue string `json:"display_value,omitempty"`
//		Link         string `json:"link,omitempty"`
//	} `json:"opened_by,omitempty"`
//	UTargetEndDateTime string `json:"u_target_end_date_time,omitempty"`
//	UserInput          string `json:"user_input,omitempty"`
//	SysCreatedOn       string `json:"sys_created_on,omitempty"`
//	SysDomain struct {
//		DisplayValue string `json:"display_value,omitempty"`
//		Link         string `json:"link,omitempty"`
//	} `json:"sys_domain,omitempty"`
//	ShortDescription       string      `json:"short_description,omitempty"`
//	State                  string      `json:"state,omitempty"`
//	UAdditionalInformation interface{} `json:"u_additional_information,omitempty"`
//	SysCreatedBy           string      `json:"sys_created_by,omitempty"`
//	Knowledge              string      `json:"knowledge,omitempty"`
//	Order                  string      `json:"order,omitempty"`
//	UOnHoldMetric          string      `json:"u_on_hold_metric,omitempty"`
//	UTargetDelivery        string      `json:"u_target_delivery,omitempty"`
//	EncryptedBy            string      `json:"encrypted_by,omitempty"`
//	ClosedAt               string      `json:"closed_at,omitempty"`
//	UChangeNumber          string      `json:"u_change_number,omitempty"`
//	CmdbCi                 string      `json:"cmdb_ci,omitempty"`
//	Impact                 string      `json:"impact,omitempty"`
//	Active                 string      `json:"active,omitempty"`
//	BusinessService        string      `json:"business_service,omitempty"`
//	Priority               string      `json:"priority,omitempty"`
//	TimeWorked             string      `json:"time_worked,omitempty"`
//	ExpectedStart          string      `json:"expected_start,omitempty"`
//	OpenedAt               string      `json:"opened_at,omitempty"`
//	BusinessDuration       string      `json:"business_duration,omitempty"`
//	GroupList              string      `json:"group_list,omitempty"`
//	WorkEnd                string      `json:"work_end,omitempty"`
//	USLAReason             string      `json:"u_sla_reason,omitempty"`
//	CheckEncryption        string      `json:"check_encryption,omitempty"`
//	CorrelationDisplay     string      `json:"correlation_display,omitempty"`
//	WorkStart              string      `json:"work_start,omitempty"`
//	AssignmentGroup        string      `json:"assignment_group,omitempty"`
//	UAssetRedeployable     interface{} `json:"u_asset_redeployable,omitempty"`
//	AdditionalAssigneeList string      `json:"additional_assignee_list,omitempty"`
//	Description            string      `json:"description,omitempty"`
//	CalendarDuration       string      `json:"calendar_duration,omitempty"`
//	UReplicate             string      `json:"u_replicate,omitempty"`
//	CloseNotes             string      `json:"close_notes,omitempty"`
//	ServiceOffering        string      `json:"service_offering,omitempty"`
//	SysClassName           string      `json:"sys_class_name,omitempty"`
//	ClosedBy               string      `json:"closed_by,omitempty"`
//	FollowUp               string      `json:"follow_up,omitempty"`
//	SysID                  string      `json:"sys_id,omitempty"`
//	ContactType            string      `json:"contact_type,omitempty"`
//	Urgency                string      `json:"urgency,omitempty"`
//	Company                string      `json:"company,omitempty"`
//	UCustomerExpectation   string      `json:"u_customer_expectation,omitempty"`
//	ReassignmentCount      string      `json:"reassignment_count,omitempty"`
//	UResponseDate          string      `json:"u_response_date,omitempty"`
//	ActivityDue            string      `json:"activity_due,omitempty"`
//	AssignedTo             string      `json:"assigned_to,omitempty"`
//	UEscalationStatus      interface{} `json:"u_escalation_status,omitempty"`
//	UInStock               interface{} `json:"u_in_stock,omitempty"`
//	Comments               string      `json:"comments,omitempty"`
//	Approval               string      `json:"approval,omitempty"`
//	SLADue                 string      `json:"sla_due,omitempty"`
//	DueDate                string      `json:"due_date,omitempty"`
//	SysModCount            string      `json:"sys_mod_count,omitempty"`
//	SysTags                string      `json:"sys_tags,omitempty"`
//	UMonitoringUser        string      `json:"u_monitoring_user,omitempty"`
//	Requestor              string      `json:"requestor,omitempty"`
//	ToBeEncrypted          string      `json:"to_be_encrypted,omitempty"`
//	UOldNumber             string      `json:"u_old_number,omitempty"`
//	Escalation             string      `json:"escalation,omitempty"`
//	UponApproval           string      `json:"upon_approval,omitempty"`
//	CorrelationID          string      `json:"correlation_id,omitempty"`
//	Location               string      `json:"location,omitempty"`
//	USchedule              string      `json:"u_schedule,omitempty"`
//	WorkNotes              string      `json:"work_notes_list,omitempty"`
//}

//type Result struct {
//	Parent                 string `json:"parent,omitempty"`
//	UVendorName            string `json:"u_vendor_name,omitempty"`
//	PhiHistory             string `json:"phi_history,omitempty"`
//	ChangeApproval         string `json:"change_approval,omitempty"`
//	WatchList              string `json:"watch_list,omitempty"`
//	IgnoreReason           string `json:"ignore_reason,omitempty"`
//	UponReject             string `json:"upon_reject,omitempty"`
//	SysUpdatedOn           string `json:"sys_updated_on,omitempty"`
//	EncryptedDataLog       string `json:"encrypted_data_log,omitempty"`
//	IgnoreExpiration       string `json:"ignore_expiration,omitempty"`
//	Ssl                    string `json:"ssl,omitempty"`
//	ApprovalHistory        string `json:"approval_history,omitempty"`
//	Skills                 string `json:"skills,omitempty"`
//	Number                 string `json:"number,omitempty"`
//	UTechnicalService      string `json:"u_technical_service,omitempty"`
//	FirstFound             string `json:"first_found,omitempty"`
//	AgeClosed              string `json:"age_closed,omitempty"`
//	State                  string `json:"state,omitempty"`
//	SysCreatedBy           string `json:"sys_created_by,omitempty"`
//	Knowledge              string `json:"knowledge,omitempty"`
//	Order                  string `json:"order,omitempty"`
//	UTargetDelivery        string `json:"u_target_delivery,omitempty"`
//	BackupState            string `json:"backup_state,omitempty"`
//	CmdbCi                 string `json:"cmdb_ci,omitempty"`
//	Impact                 string `json:"impact,omitempty"`
//	DNS                    string `json:"dns,omitempty"`
//	Active                 string `json:"active,omitempty"`
//	WorkNotesList          string `json:"work_notes_list,omitempty"`
//	Vulnerability          string `json:"vulnerability,omitempty"`
//	Priority               string `json:"priority,omitempty"`
//	BusinessDuration       string `json:"business_duration,omitempty"`
//	GroupList              string `json:"group_list,omitempty"`
//	ApprovalSet            string `json:"approval_set,omitempty"`
//	Status                 string `json:"status,omitempty"`
//	ShortDescription       string `json:"short_description,omitempty"`
//	CorrelationDisplay     string `json:"correlation_display,omitempty"`
//	WorkStart              string `json:"work_start,omitempty"`
//	AdditionalAssigneeList string `json:"additional_assignee_list,omitempty"`
//	ExternalID             string `json:"external_id,omitempty"`
//	ServiceOffering        string `json:"service_offering,omitempty"`
//	SysClassName           string `json:"sys_class_name,omitempty"`
//	ClosedBy               string `json:"closed_by,omitempty"`
//	FollowUp               string `json:"follow_up,omitempty"`
//	QualysAssigneeEmail    string `json:"qualys_assignee_email,omitempty"`
//	ManagedByVul           string `json:"managed_by_vul,omitempty"`
//	Installation           string `json:"installation,omitempty"`
//	ReassignmentCount      string `json:"reassignment_count,omitempty"`
//	TimesFound             string `json:"times_found,omitempty"`
//	AssignedTo             string `json:"assigned_to,omitempty"`
//	Netbios                string `json:"netbios,omitempty"`
//	BusinessCriticality    string `json:"business_criticality,omitempty"`
//	SLADue                 string `json:"sla_due,omitempty"`
//	IgnoredBy              string `json:"ignored_by,omitempty,omitempty"`
//	CommentsAndWorkNotes   string `json:"comments_and_work_notes,omitempty"`
//	UMonitoringUser        string `json:"u_monitoring_user,omitempty"`
//	Phi                    string `json:"phi,omitempty"`
//	Substate               string `json:"substate,omitempty"`
//	Escalation             string `json:"escalation,omitempty"`
//	UponApproval           string `json:"upon_approval,omitempty"`
//	CorrelationID          string `json:"correlation_id,omitempty"`
//	LastOpened             string `json:"last_opened,omitempty"`
//	USchedule              string `json:"u_schedule,omitempty"`
//	MadeSLA                string `json:"made_sla,omitempty"`
//	Reopened               string `json:"reopened,omitempty"`
//	Source                 string `json:"source,omitempty"`
//	UTargetImplementation  string `json:"u_target_implementation,omitempty"`
//	LastUpdatedByQualys    string `json:"last_updated_by_qualys,omitempty"`
//	Protocol               string `json:"protocol,omitempty"`
//	UBusinessService       string `json:"u_business_service,omitempty"`
//	SysUpdatedBy           string `json:"sys_updated_by,omitempty"`
//	OpenedBy               struct {
//		Link  string `json:"link,omitempty"`
//		Value string `json:"value,omitempty"`
//	} `json:"opened_by,omitempty"`
//	UTargetEndDateTime string `json:"u_target_end_date_time,omitempty"`
//	UserInput          string `json:"user_input,omitempty"`
//	SysCreatedOn       string `json:"sys_created_on,omitempty"`
//	SysDomain          struct {
//		Link  string `json:"link,omitempty"`
//		Value string `json:"value,omitempty"`
//	} `json:"sys_domain"`
//	UAdditionalInformation string `json:"u_additional_information,omitempty"`
//	UOnHoldMetric          string `json:"u_on_hold_metric,omitempty"`
//	EncryptedBy            string `json:"encrypted_by,omitempty"`
//	QualysAssigneeName     string `json:"qualys_assignee_name,omitempty"`
//	ClosedAt               string `json:"closed_at,omitempty"`
//	UChangeNumber          string `json:"u_change_number,omitempty"`
//	QualysTicketState      string `json:"qualys_ticket_state,omitempty"`
//	BusinessService        string `json:"business_service,omitempty"`
//	TimeWorked             string `json:"time_worked,omitempty"`
//	EncryptedData          string `json:"encrypted_data,omitempty"`
//	ExpectedStart          string `json:"expected_start,omitempty"`
//	OpenedAt               string `json:"opened_at,omitempty"`
//	Port                   string `json:"port,omitempty"`
//	QualysSeverity         string `json:"qualys_severity,omitempty"`
//	WorkEnd                string `json:"work_end,omitempty"`
//	USLAReason             string `json:"u_sla_reason,omitempty"`
//	CheckEncryption        string `json:"check_encryption,omitempty"`
//	WorkNotes              string `json:"work_notes,omitempty"`
//	AssignmentGroup       string  `json:"assignment_group,omitempty"`
//	UAssetRedeployable   string `json:"u_asset_redeployable,omitempty"`
//	LastFound            string `json:"last_found,omitempty"`
//	SwVulnerability      string `json:"sw_vulnerability,omitempty"`
//	Description          string `json:"description,omitempty"`
//	CalendarDuration     string `json:"calendar_duration,omitempty"`
//	UReplicate           string `json:"u_replicate,omitempty"`
//	CloseNotes           string `json:"close_notes,omitempty"`
//	SysID                string `json:"sys_id,omitempty"`
//	ContactType          string `json:"contact_type,omitempty"`
//	Urgency              string `json:"urgency,omitempty"`
//	Company              string `json:"company,omitempty"`
//	UCustomerExpectation string `json:"u_customer_expectation,omitempty"`
//	UResponseDate        string `json:"u_response_date,omitempty"`
//	ActivityDue          string `json:"activity_due,omitempty"`
//	UEscalationStatus    string `json:"u_escalation_status,omitempty"`
//	UInStock             string `json:"u_in_stock,omitempty"`
//	Comments             string `json:"comments,omitempty"`
//	QualysTicket         string `json:"qualys_ticket,omitempty"`
//	Approval             string `json:"approval,omitempty"`
//	DueDate              string `json:"due_date,omitempty"`
//	SysModCount          string `json:"sys_mod_count,omitempty"`
//	IPAddress            string `json:"ip_address,omitempty"`
//	SysTags              string `json:"sys_tags,omitempty"`
//	Requestor            string `json:"requestor,omitempty"`
//	ToBeEncrypted        string `json:"to_be_encrypted,omitempty"`
//	UOldNumber           string `json:"u_old_number,omitempty"`
//	Location             string `json:"location,omitempty"`
//	Age                  string `json:"age,omitempty"`
//}

//type Result struct {
//	Parent                 string `json:"parent"`
//	UVendorName            string `json:"u_vendor_name"`
//	PhiHistory             string `json:"phi_history"`
//	ChangeApproval         string `json:"change_approval"`
//	WatchList              string `json:"watch_list"`
//	IgnoreReason           string `json:"ignore_reason"`
//	UponReject             string `json:"upon_reject"`
//	SysUpdatedOn           string `json:"sys_updated_on"`
//	EncryptedDataLog       string `json:"encrypted_data_log"`
//	IgnoreExpiration       string `json:"ignore_expiration"`
//	Ssl                    string `json:"ssl"`
//	ApprovalHistory        string `json:"approval_history"`
//	Skills                 string `json:"skills"`
//	Number                 string `json:"number"`
//	UTechnicalService      string `json:"u_technical_service"`
//	FirstFound             string `json:"first_found"`
//	AgeClosed              string `json:"age_closed"`
//	State                  string `json:"state"`
//	SysCreatedBy           string `json:"sys_created_by"`
//	Knowledge              string `json:"knowledge"`
//	Order                  string `json:"order"`
//	UTargetDelivery        string `json:"u_target_delivery"`
//	BackupState            string `json:"backup_state"`
//	CmdbCi                 string `json:"cmdb_ci"`
//	Impact                 string `json:"impact"`
//	DNS                    string `json:"dns"`
//	Active                 string `json:"active"`
//	WorkNotesList          string `json:"work_notes_list"`
//	Vulnerability          string `json:"vulnerability"`
//	Priority               string `json:"priority"`
//	BusinessDuration       string `json:"business_duration"`
//	GroupList              string `json:"group_list"`
//	ApprovalSet            string `json:"approval_set"`
//	Status                 string `json:"status"`
//	ShortDescription       string `json:"short_description"`
//	CorrelationDisplay     string `json:"correlation_display"`
//	WorkStart              string `json:"work_start"`
//	AdditionalAssigneeList string `json:"additional_assignee_list"`
//	ExternalID             string `json:"external_id"`
//	ServiceOffering        string `json:"service_offering"`
//	SysClassName           string `json:"sys_class_name"`
//	ClosedBy               string `json:"closed_by"`
//	FollowUp               string `json:"follow_up"`
//	QualysAssigneeEmail    string `json:"qualys_assignee_email"`
//	ManagedByVul           string `json:"managed_by_vul"`
//	Installation           struct {
//		DisplayValue string `json:"display_value"`
//		Link         string `json:"link"`
//	} `json:"installation"`
//	ReassignmentCount string `json:"reassignment_count"`
//	TimesFound        string `json:"times_found"`
//	AssignedTo        struct {
//		DisplayValue string `json:"display_value"`
//		Link         string `json:"link"`
//	} `json:"assigned_to"`
//	Netbios               string `json:"netbios"`
//	BusinessCriticality   string `json:"business_criticality"`
//	SLADue                string `json:"sla_due"`
//	IgnoredBy             string `json:"ignored_by"`
//	CommentsAndWorkNotes  string `json:"comments_and_work_notes"`
//	UMonitoringUser       string `json:"u_monitoring_user"`
//	Phi                   string `json:"phi"`
//	Substate              string `json:"substate"`
//	Escalation            string `json:"escalation"`
//	UponApproval          string `json:"upon_approval"`
//	CorrelationID         string `json:"correlation_id"`
//	LastOpened            string `json:"last_opened"`
//	USchedule             string `json:"u_schedule"`
//	MadeSLA               string `json:"made_sla"`
//	Reopened              string `json:"reopened"`
//	Source                string `json:"source"`
//	UTargetImplementation string `json:"u_target_implementation"`
//	LastUpdatedByQualys   string `json:"last_updated_by_qualys"`
//	Protocol              string `json:"protocol"`
//	UBusinessService      string `json:"u_business_service"`
//	SysUpdatedBy          string `json:"sys_updated_by"`
//	OpenedBy              struct {
//		DisplayValue string `json:"display_value"`
//		Link         string `json:"link"`
//	} `json:"opened_by"`
//	UTargetEndDateTime string `json:"u_target_end_date_time"`
//	UserInput          string `json:"user_input"`
//	SysCreatedOn       string `json:"sys_created_on"`
//	SysDomain          struct {
//		DisplayValue string `json:"display_value"`
//		Link         string `json:"link"`
//	} `json:"sys_domain"`
//	UAdditionalInformation string `json:"u_additional_information"`
//	UOnHoldMetric          string `json:"u_on_hold_metric"`
//	EncryptedBy            string `json:"encrypted_by"`
//	QualysAssigneeName     string `json:"qualys_assignee_name"`
//	ClosedAt               string `json:"closed_at"`
//	UChangeNumber          struct {
//		DisplayValue string `json:"display_value"`
//		Link         string `json:"link"`
//	} `json:"u_change_number"`
//	QualysTicketState string `json:"qualys_ticket_state"`
//	BusinessService   struct {
//		DisplayValue string `json:"display_value"`
//		Link         string `json:"link"`
//	} `json:"business_service"`
//	TimeWorked      string `json:"time_worked"`
//	EncryptedData   string `json:"encrypted_data"`
//	ExpectedStart   string `json:"expected_start"`
//	OpenedAt        string `json:"opened_at"`
//	Port            string `json:"port"`
//	QualysSeverity  string `json:"qualys_severity"`
//	WorkEnd         string `json:"work_end"`
//	USLAReason      string `json:"u_sla_reason"`
//	CheckEncryption string `json:"check_encryption"`
//	WorkNotes       string `json:"work_notes"`
//	AssignmentGroup struct {
//		DisplayValue string `json:"display_value"`
//		Link         string `json:"link"`
//	} `json:"assignment_group"`
//	UAssetRedeployable   interface{} `json:"u_asset_redeployable"`
//	LastFound            string      `json:"last_found"`
//	SwVulnerability      string      `json:"sw_vulnerability"`
//	Description          string      `json:"description"`
//	CalendarDuration     string      `json:"calendar_duration"`
//	UReplicate           string      `json:"u_replicate"`
//	CloseNotes           string      `json:"close_notes"`
//	SysID                string      `json:"sys_id"`
//	ContactType          string      `json:"contact_type"`
//	Urgency              string      `json:"urgency"`
//	Company              string      `json:"company"`
//	UCustomerExpectation string      `json:"u_customer_expectation"`
//	UResponseDate        string      `json:"u_response_date"`
//	ActivityDue          string      `json:"activity_due"`
//	UEscalationStatus    interface{} `json:"u_escalation_status"`
//	UInStock             interface{} `json:"u_in_stock"`
//	Comments             string      `json:"comments"`
//	QualysTicket         string      `json:"qualys_ticket"`
//	Approval             string      `json:"approval"`
//	DueDate              string      `json:"due_date"`
//	SysModCount          string      `json:"sys_mod_count"`
//	IPAddress            string      `json:"ip_address"`
//	SysTags              string      `json:"sys_tags"`
//	Requestor            struct {
//		DisplayValue string `json:"display_value"`
//		Link         string `json:"link"`
//	} `json:"requestor"`
//	ToBeEncrypted string `json:"to_be_encrypted"`
//	UOldNumber    string `json:"u_old_number"`
//	Location      string `json:"location"`
//	Age           string `json:"age"`
//}

//type SvcNowRequest struct {
//	SysID                string `json:"sys_id,omitempty"`
//	BusinessService        string `json:"business_service,omitempty"`
//	AssignmentGroup        string `json:"assignment_group,omitempty"`
//	State                  string `json:"state,omitempty"`
//	DNS                    string `json:"dns,omitempty"`
//	Substate               string `json:"substate,omitempty"`
//	Description            string `json:"description,omitempty"`
//	CloseNotes             string `json:"close_notes,omitempty"`
//	AssignedTo             string `json:"assigned_to,omitempty"`
//	Requestor              string `json:"requestor,omitempty"`
//	DueDate                string `json:"due_date,omitempty"`
//	ClosedAt               string `json:"closed_at,omitempty"`
//	QualysSeverity         string `json:"qualys_severity,omitempty"`
//	Port                   string `json:"port,omitempty"`
//	Skills                 string `json:"skills,omitempty"`
//	IPAddress              string `json:"ip_address,omitempty"`
//	UAdditionalInformation string `json:"u_additional_information,omitempty"`
//	UChangeNumber          string `json:"u_change_number,omitempty"`
//	GroupList              string `json:"group_list,omitempty"`
//	CmdbCi                 string `json:"cmdb_ci,omitempty"`
//	QualysTicket           string `json:"qualys_ticket,omitempty"`
//	IgnoreReason           string `json:"ignore_reason,omitempty"`
//	Ssl                    string `json:"ssl,omitempty"`
//	ActivityDue            string `json:"activity_due,omitempty"`
//	UserInput              string `json:"user_input,omitempty"`
//	WatchList              string `json:"watch_list,omitempty"`
//	BusinessCriticality    string `json:"business_criticality,omitempty"`
//	Protocol               string `json:"protocol,omitempty"`
//	WorkNotes              string `json:"work_notes,omitempty"`
//	QualysAssigneeEmail    string `json:"qualys_assignee_email,omitempty"`
//	QualysAssigneeName     string `json:"qualys_assignee_name,omitempty"`
//	QualysTicketState      string `json:"qualys_ticket_state,omitempty"`
//	Source                 string `json:"source,omitempty"`
//	Active                 string `json:"active,omitempty"`
//	WorkNotesList          string `json:"work_notes_list,omitempty"`
//	Vulnerability          string `json:"vulnerability,omitempty"`
//	Priority               string `json:"priority,omitempty"`
//	ManagedByVul           string `json:"managed_by_vul,omitempty"`
//	Installation           string `json:"installation,omitempty"`
//	SLADue                 string `json:"sla_due,omitempty"`
//	MadeSLA                string `json:"made_sla,omitempty"`
//	USLAReason             string `json:"u_sla_reason,omitempty"`
//	SwVulnerability        string `json:"sw_vulnerability,omitempty"`
//	Comments               string `json:"comments,omitempty"`
//	CorrelationDisplay     string `json:"correlation_display,omitempty"`
//	CorrelationID          string `json:"correlation_id,omitempty"`
//	AdditionalAssigneeList string `json:"additional_assignee_list,omitempty"`
//}

// SvcNowRequest holds the information required for the API to process a request
type SvcNowRequest struct {
	Priority               string `json:"priority,omitempty"`
	DueDate                string `json:"due_date,omitempty"`
	UAdditionalInformation string `json:"u_additional_information,omitempty"`
	AssignmentGroup        string `json:"assignment_group,omitempty"`
	State                  string `json:"state,omitempty"`
	Substate               string `json:"substate,omitempty"`
	Port                   string `json:"port,omitempty"`
	IPAddress              string `json:"ip_address,omitempty"`
	DNS                    string `json:"dns,omitempty"`
	ActivityDue            string `json:"activity_due,omitempty"`
	QualysSeverity         string `json:"qualys_severity,omitempty"`
	Skills                 string `json:"skills,omitempty"`
	Requestor              string `json:"requestor,omitempty"`
	AdditionalAssigneeList string `json:"additional_assignee_list,omitempty"`
	WatchList              string `json:"watch_list,omitempty"`
	QualysTicket           string `json:"qualys_ticket,omitempty"`
	GroupList              string `json:"group_list,omitempty"`
	UserInput              string `json:"user_input,omitempty"`
	Description            string `json:"description,omitempty"`
	CloseNotes             string `json:"close_notes,omitempty"`
	CorrelationDisplay     string `json:"correlation_display,omitempty"`
	CorrelationID          string `json:"correlation_id,omitempty"`
	WorkNotes              string `json:"work_notes,omitempty"`
	Active                 string `json:"active,omitempty"`
	ClosedAt               string `json:"closed_at,omitempty"`
	AssignedTo             string `json:"assigned_to,omitempty"`
	QualysAssigneeEmail    string `json:"qualys_assignee_email,omitempty"`
	Vulnerability          string `json:"vulnerability,omitempty"`
	CmdbCi                 string `json:"cmdb_ci,omitempty"`
}

// Result parses the particular fields of a service now response
type Result struct {
	Parent                 string `json:"parent,omitempty"`
	UVendorName            string `json:"u_vendor_name,omitempty"`
	PhiHistory             string `json:"phi_history,omitempty"`
	ChangeApproval         string `json:"change_approval,omitempty"`
	WatchList              string `json:"watch_list,omitempty"`
	IgnoreReason           string `json:"ignore_reason,omitempty"`
	UponReject             string `json:"upon_reject,omitempty"`
	SysUpdatedOn           string `json:"sys_updated_on,omitempty"`
	EncryptedDataLog       string `json:"encrypted_data_log,omitempty"`
	IgnoreExpiration       string `json:"ignore_expiration,omitempty"`
	Ssl                    string `json:"ssl,omitempty"`
	ApprovalHistory        string `json:"approval_history,omitempty"`
	Skills                 string `json:"skills,omitempty"`
	Number                 string `json:"number,omitempty"`
	UTechnicalService      string `json:"u_technical_service,omitempty"`
	FirstFound             string `json:"first_found,omitempty"`
	AgeClosed              string `json:"age_closed,omitempty"`
	State                  string `json:"state,omitempty"`
	SysCreatedBy           string `json:"sys_created_by,omitempty"`
	Knowledge              string `json:"knowledge,omitempty"`
	Order                  string `json:"order,omitempty"`
	UTargetDelivery        string `json:"u_target_delivery,omitempty"`
	BackupState            string `json:"backup_state,omitempty"`
	Impact                 string `json:"impact,omitempty"`
	DNS                    string `json:"dns,omitempty"`
	Active                 string `json:"active,omitempty"`
	WorkNotesList          string `json:"work_notes_list,omitempty"`
	Priority               string `json:"priority,omitempty"`
	BusinessDuration       string `json:"business_duration,omitempty"`
	GroupList              string `json:"group_list,omitempty"`
	ApprovalSet            string `json:"approval_set,omitempty"`
	Status                 string `json:"status,omitempty"`
	ShortDescription       string `json:"short_description,omitempty"`
	CorrelationDisplay     string `json:"correlation_display,omitempty"`
	WorkStart              string `json:"work_start,omitempty"`
	AdditionalAssigneeList string `json:"additional_assignee_list,omitempty"`
	ExternalID             string `json:"external_id,omitempty"`
	ServiceOffering        string `json:"service_offering,omitempty"`
	SysClassName           string `json:"sys_class_name,omitempty"`
	ClosedBy               string `json:"closed_by,omitempty"`
	FollowUp               string `json:"follow_up,omitempty"`
	QualysAssigneeEmail    string `json:"qualys_assignee_email,omitempty"`
	ManagedByVul           string `json:"managed_by_vul,omitempty"`
	Installation           string `json:"installation,omitempty"`
	ReassignmentCount      string `json:"reassignment_count,omitempty"`
	TimesFound             string `json:"times_found,omitempty"`
	AssignedTo             struct {
		DisplayValue string `json:"display_value,omitempty"`
		Link         string `json:"link,omitempty"`
	} `json:"assigned_to,omitempty"`
	Netbios               string `json:"netbios,omitempty"`
	BusinessCriticality   string `json:"business_criticality,omitempty"`
	SLADue                string `json:"sla_due,omitempty"`
	IgnoredBy             string `json:"ignored_by,omitempty"`
	CommentsAndWorkNotes  string `json:"comments_and_work_notes,omitempty"`
	UMonitoringUser       string `json:"u_monitoring_user,omitempty"`
	Phi                   string `json:"phi,omitempty"`
	Substate              string `json:"substate,omitempty"`
	Escalation            string `json:"escalation,omitempty"`
	UponApproval          string `json:"upon_approval,omitempty"`
	CorrelationID         string `json:"correlation_id,omitempty"`
	LastOpened            string `json:"last_opened,omitempty"`
	USchedule             string `json:"u_schedule,omitempty"`
	MadeSLA               string `json:"made_sla,omitempty"`
	Reopened              string `json:"reopened,omitempty"`
	Source                string `json:"source,omitempty"`
	UTargetImplementation string `json:"u_target_implementation,omitempty"`
	LastUpdatedByQualys   string `json:"last_updated_by_qualys,omitempty"`
	Protocol              string `json:"protocol,omitempty"`
	UBusinessService      string `json:"u_business_service,omitempty"`
	SysUpdatedBy          string `json:"sys_updated_by,omitempty"`
	OpenedBy              struct {
		DisplayValue string `json:"display_value,omitempty"`
		Link         string `json:"link,omitempty"`
	} `json:"opened_by,omitempty"`
	UTargetEndDateTime string `json:"u_target_end_date_time,omitempty"`
	UserInput          string `json:"user_input,omitempty"`
	SysCreatedOn       string `json:"sys_created_on,omitempty"`
	SysDomain          struct {
		DisplayValue string `json:"display_value,omitempty"`
		Link         string `json:"link,omitempty"`
	} `json:"sys_domain,omitempty"`
	UAdditionalInformation string `json:"u_additional_information,omitempty"`
	UOnHoldMetric          string `json:"u_on_hold_metric,omitempty"`
	EncryptedBy            string `json:"encrypted_by,omitempty"`
	QualysAssigneeName     string `json:"qualys_assignee_name,omitempty"`
	ClosedAt               string `json:"closed_at,omitempty"`
	UChangeNumber          struct {
		DisplayValue string `json:"display_value,omitempty"`
		Link         string `json:"link,omitempty"`
	} `json:"u_change_number,omitempty"`
	QualysTicketState string `json:"qualys_ticket_state,omitempty"`
	BusinessService   string `json:"business_service,omitempty"`
	TimeWorked        string `json:"time_worked,omitempty"`
	EncryptedData     string `json:"encrypted_data,omitempty"`
	ExpectedStart     string `json:"expected_start,omitempty"`
	OpenedAt          string `json:"opened_at,omitempty"`
	Port              string `json:"port,omitempty"`
	QualysSeverity    string `json:"qualys_severity,omitempty"`
	WorkEnd           string `json:"work_end,omitempty"`
	USLAReason        string `json:"u_sla_reason,omitempty"`
	CheckEncryption   string `json:"check_encryption,omitempty"`
	WorkNotes         string `json:"work_notes,omitempty"`
	AssignmentGroup   struct {
		DisplayValue string `json:"display_value,omitempty"`
		Link         string `json:"link,omitempty"`
	} `json:"assignment_group,omitempty"`
	Vulnerability struct {
		DisplayValue string `json:"display_value,omitempty"`
		Link         string `json:"link,omitempty"`
	} `json:"vulnerability,omitempty"`
	CmdbCi struct {
		DisplayValue string `json:"display_value,omitempty"`
		Link         string `json:"link,omitempty"`
	} `json:"cmdb_ci,omitempty"`
	UAssetRedeployable   interface{} `json:"u_asset_redeployable,omitempty"`
	LastFound            string      `json:"last_found,omitempty"`
	SwVulnerability      string      `json:"sw_vulnerability,omitempty"`
	Description          string      `json:"description,omitempty"`
	CalendarDuration     string      `json:"calendar_duration,omitempty"`
	UReplicate           string      `json:"u_replicate,omitempty"`
	CloseNotes           string      `json:"close_notes,omitempty"`
	SysID                string      `json:"sys_id,omitempty"`
	ContactType          string      `json:"contact_type,omitempty"`
	Urgency              string      `json:"urgency,omitempty"`
	Company              string      `json:"company,omitempty"`
	UCustomerExpectation string      `json:"u_customer_expectation,omitempty"`
	UResponseDate        string      `json:"u_response_date,omitempty"`
	ActivityDue          string      `json:"activity_due,omitempty"`
	UEscalationStatus    interface{} `json:"u_escalation_status,omitempty"`
	UInStock             interface{} `json:"u_in_stock,omitempty"`
	Comments             string      `json:"comments,omitempty"`
	QualysTicket         string      `json:"qualys_ticket,omitempty"`
	Approval             string      `json:"approval,omitempty"`
	DueDate              string      `json:"due_date,omitempty"`
	SysModCount          string      `json:"sys_mod_count,omitempty"`
	IPAddress            string      `json:"ip_address,omitempty"`
	SysTags              string      `json:"sys_tags,omitempty"`
	Requestor            struct {
		DisplayValue string `json:"display_value,omitempty"`
		Link         string `json:"link,omitempty"`
	} `json:"requestor,omitempty"`
	ToBeEncrypted string `json:"to_be_encrypted,omitempty"`
	UOldNumber    string `json:"u_old_number,omitempty"`
	Location      string `json:"location,omitempty"`
	Age           string `json:"age,omitempty"`
}

//type Result struct {
//		CorrelationDisplay     string `json:"correlation_display,omitempty"`
//		WatchList              string `json:"watch_list,omitempty"`
//		AdditionalAssigneeList string `json:"additional_assignee_list,omitempty"`
//		Description            string `json:"description,omitempty"`
//		CloseNotes             string `json:"close_notes,omitempty"`
//		Skills                 string `json:"skills,omitempty"`
//		Number                 string `json:"number,omitempty"`
//		SysID                  string `json:"sys_id,omitempty"`
//		QualysAssigneeEmail    string `json:"qualys_assignee_email,omitempty"`
//		UserInput              string `json:"user_input,omitempty"`
//		State                  string `json:"state,omitempty"`
//		UAdditionalInformation string `json:"u_additional_information,omitempty"`
//		ActivityDue            string `json:"activity_due,omitempty"`
//		QualysTicket           string `json:"qualys_ticket,omitempty"`
//		ClosedAt               string `json:"closed_at,omitempty"`
//		DueDate                string `json:"due_date,omitempty"`
//		DNS                    string `json:"dns,omitempty"`
//		Active                 string `json:"active,omitempty"`
//		IPAddress              string `json:"ip_address,omitempty"`
//		Priority               string `json:"priority,omitempty"`
//		Requestor              struct {
//			DisplayValue string `json:"display_value,omitempty"`
//			Link         string `json:"link,omitempty"`
//		} `json:"requestor,omitempty"`
//		GroupList      string `json:"group_list,omitempty"`
//		Substate       string `json:"substate,omitempty"`
//		Port           string `json:"port,omitempty"`
//		QualysSeverity string `json:"qualys_severity,omitempty"`
//		CorrelationID  string `json:"correlation_id,omitempty"`
//		WorkNotes      string `json:"work_notes,omitempty"`
//	}
