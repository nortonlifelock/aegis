package scaffold

import (
	"github.com/benjivesterby/validator"
	"github.com/pkg/errors"
	"strings"
)

type procedureDescriptor struct {
	Name       string
	Parameters map[string]*parameterDescriptor
	PArray     []*parameterDescriptor
	Imports    map[string]bool
	FullProc   string
	Return     *structDescriptor
}

func newProcedureDescriptor(Name string, FullProc string, Return *structDescriptor) (proc *procedureDescriptor, err error) {

	if len(FullProc) > 0 {

		if len(Name) > 0 {

			proc = &procedureDescriptor{
				Name:       Name,
				FullProc:   FullProc,
				Return:     Return,
				Parameters: make(map[string]*parameterDescriptor),
				Imports:    make(map[string]bool),
				PArray:     make([]*parameterDescriptor, 0),
			}

		} else {
			err = errors.New("newProcedureDescriptor - the passed name cannot be empty")
		}

	} else {
		err = errors.New("newProcedureDescriptor - the passed FullProc cannot be empty")
	}

	return proc, err
}

func (proc *procedureDescriptor) addParameter(param *parameterDescriptor) (err error) {

	if validator.IsValid(param) {
		proc.Parameters[param.Name] = param

		// This allows us to track the parameter location for proper
		// procedure call
		if proc.PArray != nil {
			proc.PArray = append(proc.PArray, param)
		}
	}

	return err
}

func (proc *procedureDescriptor) newParameter(name string, dbType string, nullable bool) (err error) {

	if len(name) > 0 {
		if len(dbType) > 0 {

			dbType = strings.ToLower(dbType)

			var size int
			if dbType, size, err = separateDBType(dbType); err == nil {

				goType := ""

				var customType bool
				goType, customType, err = dbTypeToGoType(dbType, size, &proc.Imports, nullable)

				if err == nil {

					param := &parameterDescriptor{
						Name:       name,
						GoType:     goType,
						DBType:     dbType,
						Size:       size,
						Null:       nullable,
						customType: customType,
					}

					err = proc.addParameter(param)
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
