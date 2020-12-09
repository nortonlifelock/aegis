package scaffold

import (
	"github.com/benjivesterby/validator"
	"github.com/pkg/errors"
	"strings"
)

type structDescriptor struct {
	Name         string
	Imports      map[string]bool
	Parameters   map[string]*parameterDescriptor
	PArray       []*parameterDescriptor
	Methods      map[string]*methodDescriptor
	SingleReturn bool
	ID           int
}

func newStruct(name string) (str *structDescriptor, err error) {

	if len(name) > 0 {

		str = &structDescriptor{
			Name:       name,
			Parameters: make(map[string]*parameterDescriptor),
			PArray:     make([]*parameterDescriptor, 0),
			Methods:    make(map[string]*methodDescriptor),
			Imports:    make(map[string]bool),
		}
	} else {
		err = errors.New("Empty Struct Name")
	}

	return
}

func copyStruct(in *structDescriptor) (out *structDescriptor, err error) {

	if len(in.Name) > 0 {

		out = &structDescriptor{Name: in.Name}

		// Initialize the maps
		out.Parameters = make(map[string]*parameterDescriptor)
		out.PArray = make([]*parameterDescriptor, 0)

		out.Methods = make(map[string]*methodDescriptor)
		out.Imports = make(map[string]bool)
	} else {
		err = errors.New("empty struct name")
	}

	return
}

func (desc *structDescriptor) addParameter(param *parameterDescriptor) (err error) {

	if validator.IsValid(param) {

		var nullable = param.Null
		// Append the property to the array in the order
		if desc.Parameters[param.Name] == nil {
			desc.PArray = append(desc.PArray, param)
			nullable = param.Null
		} else if param.Null != desc.Parameters[param.Name].Null {
			nullable = true
		}

		desc.Parameters[param.Name] = param
		param.Null = nullable
	}

	return err
}

func (desc *structDescriptor) newParameter(name string, dbType string, dbNull string) (err error) {

	if len(name) > 0 {
		if len(dbType) > 0 {
			var originalDbType = dbType
			dbType = strings.ToLower(dbType)

			// Process desc row as a database column return rather than a
			// mapped object
			var size int
			if dbType, size, err = separateDBType(dbType); err == nil {

				goType := ""
				nullable := false

				// Handle the DB null
				if len(dbNull) > 0 {
					dbNull = strings.ToLower(dbNull)

					if dbNull == "null" {
						nullable = true
					}
				}

				var customType bool
				if goType, customType, err = dbTypeToGoType(originalDbType, size, &desc.Imports, nullable); err == nil {

					err = desc.addParameter(&parameterDescriptor{
						Name:       name,
						GoType:     goType,
						DBType:     dbType,
						Size:       size,
						Null:       nullable,
						customType: customType,
					})
				}
			}
		} else {
			err = errors.New("Database Type cannot be empty")
		}
	} else {
		err = errors.New("Property Name cannot be empty")
	}

	return err
}

// Validate returns a boolean relating to if the structDescriptor has all required fields populated
// TODO: Add setters and getters for the maps and arrays in the struct descriptor so that we can cut down on duplicate
// calls as well as add the ability to subscribe sub objects to the validate method calls.
// This would allow the validate method to fully validate all of the sub objects
func (desc *structDescriptor) Validate() (valid bool) {

	if len(desc.Name) > 0 {
		if desc.Imports != nil {
			// Ensure that the parameters aren't null and that both of the managed
			// structures contain the same number of elements
			if desc.Parameters != nil &&
				len(desc.Parameters) > 0 &&
				desc.PArray != nil &&
				len(desc.Parameters) == len(desc.PArray) {

				//if desc.Methods != nil && len(desc.Methods) > 0 {
				//	if desc.Maps != nil {
				//		valid = true
				//	}
				//}
			}
		}
	}

	return valid
}
