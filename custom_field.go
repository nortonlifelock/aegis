package jira

import (
	"fmt"
)

// GetFieldMap returns the custom project field equivalent to the backend field name
func (connector *ConnectorJira) GetFieldMap(in string) *Field {
	return connector.Fields[connector.GetFieldMapName(in)]
}

// GetFieldMapName returns the custom project field name equivalent to the backend field name
func (connector *ConnectorJira) GetFieldMapName(in string) (out string) {
	out = in

	if connector.payload.FieldMap != nil {
		if len(connector.payload.FieldMap[in]) > 0 {
			out = connector.payload.FieldMap[in]
		}
	}

	return out
}

func (ji *Issue) getCF(connector *ConnectorJira, cf string) (obj interface{}, err error) {
	cf = ji.connector.GetFieldMapName(cf)

	var ret interface{}

	if ji.Issue != nil {
		if ji.Issue.Fields != nil {
			if ji.Issue.Fields.Unknowns != nil && len(ji.Issue.Fields.Unknowns) > 0 {
				if connector != nil {
					if connector.Fields != nil && len(connector.Fields) > 0 {
						fvalue := connector.Fields[cf]

						if fvalue != nil {
							ret = ji.Issue.Fields.Unknowns[fvalue.getCreateID()]

							if ret != nil {

								if v, ok := ret.(map[string]interface{}); ok {
									obj = v
								} else {

									obj = ret
								}
							}
						} else {
							err = fmt.Errorf("nil field value [%s] retrieved from connector", cf)
						}
					} else {
						err = fmt.Errorf("connector did not properly load fields")
					}
				} else {
					err = fmt.Errorf("nil connector passed to getCF")
				}
			} else {
				err = fmt.Errorf("jira issue did not properly load field unknowns")
			}
		} else {
			err = fmt.Errorf("jira issue did not properly load fields")
		}

	} else {
		err = fmt.Errorf("jira issue passed nil to getCF")
	}

	return obj, err
}

func (ji *Issue) setCF(connector *ConnectorJira, cf string, value interface{}) (err error) {
	cf = ji.connector.GetFieldMapName(cf)

	if ji.Issue != nil {
		if ji.Issue.Fields != nil {
			if ji.Issue.Fields.Unknowns != nil {
				if connector != nil {
					if connector.Fields != nil && len(connector.Fields) > 0 {
						fValue := connector.Fields[cf]

						if fValue != nil {

							if fValue.Schema.Type == "string" ||
								fValue.Schema.Type == "date" {

								ji.Issue.Fields.Unknowns[fValue.getCreateID()] = value
							} else {
								if fValue.Schema.Type == "group" {

									// Set the custom field to the Value that was passed in
									ji.Issue.Fields.Unknowns[fValue.getCreateID()] = CF{Name: value}
								} else {

									// Set the custom field to the Value that was passed in
									ji.Issue.Fields.Unknowns[fValue.getCreateID()] = CF{Value: value}
								}
							}
						} else {
							err = fmt.Errorf("nil field value retrieved from connector")
						}
					} else {
						err = fmt.Errorf("connector did not properly load fields")
					}
				} else {
					err = fmt.Errorf("nil connector passed to setCF")
				}
			} else {
				err = fmt.Errorf("jira issue did not properly load field unknowns")
			}
		} else {
			err = fmt.Errorf("jira issue did not properly load fields")
		}

	} else {
		err = fmt.Errorf("jira issue passed nil to setCF")
	}

	return err
}

func (connector *ConnectorJira) getField(issue *Issue, field string) (fieldValue interface{}, err error) {
	var cf interface{}
	cf, err = issue.getCF(connector, field)
	if cf != nil {

		if v, ok := cf.(map[string]interface{}); ok {
			if v["value"] != nil {
				fieldValue = v["value"]
			} else {
				fieldValue = v["name"]
			}
		} else {

			fieldValue = cf
		}

	} else {
		// TODO:
	}

	return fieldValue, err
}

func (cf *Field) getQueryID() (id string) {
	if cf != nil {
		if cf.Schema.CustomID == 0 {
			id = cf.ID
		} else {
			id = fmt.Sprintf("CF[%v]", cf.Schema.CustomID)
		}
	} else {
		fmt.Println("NIL CUSTOM FIELD IN GET QUERY ID")
	}

	return id
}

func (cf *Field) getCreateID() (id string) {
	if cf != nil {
		id = cf.ID
	} else {
		fmt.Println("NIL CUSTOM FIELD in GET CREATE ID")
	}

	return id
}
