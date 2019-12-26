package scaffold

type parameterDescriptor struct {
	Name   string
	GoType string
	DBType string
	Size   int  // This will include checks for size in the generated class
	Null   bool // This will include checks for null in the generated class
	Index  int

	IsMapType bool // Is a mapping parameter
	MapSingle bool // Many or One indicator

	customType bool
}

// Validate checks whether parameterDescriptor is populated with all the required fields
func (param *parameterDescriptor) Validate() (valid bool) {

	if param.IsMapType {
		if len(param.Name) > 0 {
			if len(param.GoType) > 0 {
				valid = true
			}
		}

	} else if len(param.Name) > 0 {
		if len(param.GoType) > 0 {
			if len(param.DBType) > 0 {
				if param.Size != 0 {
					valid = true
				}
			}
		}
	}

	return valid
}
