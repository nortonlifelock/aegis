package scaffold

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nortonlifelock/files"
	"github.com/pkg/errors"
)

// builder holds a list of all the procedures, classes, and methods and is used to hold the information required to generate code
type builder struct {
	dbconn        idb
	classes       map[string]*structDescriptor
	procs         map[string]*procedureDescriptor
	accessMethods map[string]string
	templatePath  string
}

func newBuilder(dbconn idb, templatePath string) (b *builder, err error) {

	b = &builder{
		classes:       make(map[string]*structDescriptor),
		accessMethods: make(map[string]string),
		procs:         make(map[string]*procedureDescriptor),
		dbconn:        dbconn,
		templatePath:  templatePath,
	}

	return b, err
}

func (builder *builder) processSproc(path string) (err error) {
	status := initiated

	sprocType := execute

	var declarationStartIndex int
	var declarationEndIndex int
	var sprocCreateIndex int

	var fstring string

	// Load the SPROC file into memory
	if fstring, err = files.GetStringFromFile(path); err == nil {

		var lines []string
		// Split the file up into an array of strings at the line breaks
		if lines, err = stringToLines(fstring); err == nil {

			// Identify the important lines of the file
			for index := range lines {

				var fields []string
				fields, err = getFields(lines[index])

				if err == nil {

					// Find the comments, return and create statement to pull information from for the class
					if len(fields) > 0 {
						switch strings.ToLower(fields[0]) {
						case "/*":
							// Only list if status is initiated because we only want
							// builder comment if it's at the beginning of the file
							if status == initiated {

								sprocType = read

								// Begin class declaration
								// Update status to reflect the change
								status = returns
								declarationStartIndex = index + 1
							}
						case "*/":
							// Only update the status if builder is part of a declaration
							if status == returns {
								// End class declaration
								declarationEndIndex = index // no minus 1 because the slice takes that into account
								status = sproc
							}

						case "create":
							status = sproc
							// Indicate the line where the sproc create begins
							sprocCreateIndex = index
						}
					}
				}
			}

			var returnVal *structDescriptor

			var createString string
			var createProcedure = true
			createString, err = builder.getSprocCreateString(fstring, lines[sprocCreateIndex])

			if err == nil {
				if sprocType == read {

					// Slice the lines out that include the returns only
					returnVal, createProcedure, err = builder.handleDeclarations(lines[declarationStartIndex:declarationEndIndex])
				}

				if createProcedure {
					err = builder.createProcedureDescriptor(createString, fstring, returnVal)
				}
			}
		}
	}

	return err
}

func (builder *builder) createProcedureDescriptor(create string, fullproc string, returnVal *structDescriptor) (err error) {

	// Get the index of the first '('
	indexOfParamStart := strings.Index(create, "(")
	indexOfParamEnd := strings.LastIndex(create, ")")

	var sprocName string
	if sprocName, err = getSprocName(create[:indexOfParamStart]); err == nil {

		var proc *procedureDescriptor
		if proc, err = newProcedureDescriptor(sprocName, fullproc, returnVal); err == nil {

			// Add builder to the builders sprocs
			builder.procs[sprocName] = proc

			// ***********************************************************
			// Stored Procedure Parameter Parsing
			// ***********************************************************

			// If indexOfParamStart + 1 is not < indexOfParamEnd, there are no parameters to process
			if indexOfParamStart+1 < indexOfParamEnd {
				// Separate the parameters by comma
				params := strings.Split(create[indexOfParamStart+1:indexOfParamEnd], ",")

				// Evaluate each of the parameters
				for index := range params {

					// Remove any excess whitespace at the beginning or end
					param := strings.TrimSpace(params[index])

					// Get the fields of each parameter
					paramFields := strings.Fields(param)

					// Verify the param broke up properly
					if len(paramFields) == 2 {

						// Break out of the loop if there is an issue assigning a parameter
						if err = proc.newParameter(paramFields[0], paramFields[1], false); err != nil {
							break
						}
					} else {
						err = errors.Errorf("createProcedureDescriptor - paramFields did not have 2 fields, it had %d", len(paramFields))
					}
				}
			}
		}

	}

	return err
}

func getSprocName(createString string) (dbsprocname string, err error) {

	if len(createString) > 0 {

		// Clean up newlines and ticks
		createString = strings.Replace(createString, "\n", "", -1)
		createString = strings.Replace(createString, "\r", "", -1)
		createString = strings.Replace(createString, "`", "", -1)
		createString = strings.Replace(createString, "`", "", -1)

		createFields := strings.Fields(createString)

		if len(createFields) == 3 {
			// Get the third field which should include the sproc db and name
			dbsprocname = createFields[2]
		} else {
			err = errors.New("getSprocName - createFields could not be split into three")
		}
	} else {
		err = errors.New("getSprocName - createString cannot be empty")
	}

	return dbsprocname, err
}

func (builder *builder) getSprocCreateString(fullFileString string, createStatementStart string) (fullCreateStatement string, err error) {

	if len(fullFileString) > 0 {
		if len(createStatementStart) > 0 {

			// Get the index of the create statement in the full file
			indexOfCreateInOriginalFile := strings.Index(fullFileString, createStatementStart)

			// Create a substring of the file starting at the create statement of the stored procedure
			subst := fullFileString[indexOfCreateInOriginalFile:]

			// Get the index of the Begin statement that should follow the create statement
			indexOfBeginStatement := strings.Index(subst, "#BEGIN#")

			// Get the stubstring that contains the full stored procedure create statement
			// so that the full thing can be parsed for a proper method call
			fullCreateStatement = subst[:indexOfBeginStatement]

			// Clean up any newlines
			fullCreateStatement = strings.Replace(fullCreateStatement, "\n", "", -1)
			fullCreateStatement = strings.Replace(fullCreateStatement, "\r", "", -1)

		} else {
			err = errors.New("getSprocCreateString - createStatementStart cannot be empty")
		}
	} else {
		err = errors.New("getSprocCreateString - the fullFileString cannot be empty")
	}

	return fullCreateStatement, err
}

func (builder *builder) getOrBuildClass(fields []string) (myClass *structDescriptor, err error) {
	if className := fields[1]; len(className) > 0 {

		// Only create the procClass if it's NIL
		myClass = builder.classes[className]
		if myClass == nil {

			if myClass, err = newStruct(className); err == nil {
				builder.classes[className] = myClass
			}
		}
	}

	return myClass, err
}

func (builder *builder) handleDeclarations(returns []string) (currentReturn *structDescriptor, createProcedure bool, err error) {
	createProcedure = true

	if returns != nil && len(returns) > 0 {

		// Current class that's being worked on
		var existingClass *structDescriptor

		for index := range returns {

			var line = returns[index]

			// Ensure that builder is not an empty line
			if len(strings.TrimSpace(returns[index])) > 0 {
				var fields []string

				// Break up the line by field
				if fields, err = getFields(line); err == nil {

					if len(fields) >= 1 {
						switch strings.ToLower(fields[0]) {
						case "gen":
							createProcedure = false
							fallthrough
						case "return":

							if len(fields) >= 2 {
								existingClass, err = builder.getOrBuildClass(fields)
								if err == nil {
									currentReturn, err = copyStruct(existingClass)
								}
							} else {
								err = errors.Errorf("handleDeclarations - Return is not valid for line %s", line)
							}
						default:
							if len(fields) >= 3 {
								// Add new parameters to the procClass
								err = currentReturn.newParameter(fields[0], fields[1], fields[2])
								if err == nil {
									err = existingClass.newParameter(fields[0], fields[1], fields[2])
								}
							} else {
								err = errors.Errorf("handleDeclarations - The line %s could not be split into three fields to create a parameter", line)
							}
						}
					} else {
						err = errors.Errorf("handleDeclarations - could not retrieve fields from the line %s", line)
					}
				}
			} else {
				err = errors.New("handleDeclarations - line must have a length greater than zero")
			}
		}

	}

	return currentReturn, createProcedure, err
}

func (builder *builder) generateClass(class *structDescriptor) (genClass string, err error) {
	if class != nil {
		var genClassTemplate *Template
		if genClassTemplate, err = NewTemplate(builder.templatePath, "struct"); err == nil {

			var mapParams = make(map[string]*parameterDescriptor)

			for k, p := range class.Parameters {
				if p.IsMapType {
					mapParams[k] = p
				}
			}

			var params string
			var jsonParams string
			if params, jsonParams, err = builder.generateParameters(class, class.Parameters); err == nil {

				var methods string
				alphabetizedMethods := getOrderedKeyListFromMethodMap(class.Methods)
				for _, index := range alphabetizedMethods {
					methods += "\n\n" + class.Methods[index].Method
				}

				var imports string
				alphabetizedImports := getOrderedListFromBoolMap(class.Imports)
				for _, imp := range alphabetizedImports {
					imports += "\t\"" + imp + "\"\n"
				}

				genClassTemplate.Repl("%class", class.Name).
					Repl("%imports", imports).
					Repl("%json", jsonParams).
					Repl("%parameters", params).
					Repl("%methods", methods)

				genClass = genClassTemplate.Get()

			}
		}

	}

	return genClass, err
}

func (builder *builder) generateParameters(class *structDescriptor, parameters map[string]*parameterDescriptor) (genParams string, jsonParams string, err error) {

	if parameters != nil {

		var count = 0

		alphabetizedParameters := getOrderedListFromParameterMap(parameters)
		for _, index := range alphabetizedParameters {
			var param = parameters[index]

			// Force all generated parameters to be private so that they have to use
			// the generated getters and setters, or the new class method
			pname := param.Name

			var isSlice = strings.Index(param.GoType, "[]") >= 0

			var p string
			if param.customType && strings.Index(param.GoType, "domain") < 0 {
				p = pname + "var domain." + param.GoType
			} else {
				p = pname + "var " + param.GoType
			}

			jsonParams += fmt.Sprintf("\t\t\"%s\":my%s.%svar,\n", pname, class.Name, pname)

			genParams += fmt.Sprintf("\t%s\n", p)

			err = builder.generateParamGetters(param, class)
			if err != nil {
				break
			}

			err = builder.generateParamSetters(param, class, isSlice)
			if err != nil {
				break
			}

			if count < len(parameters)-1 {
				count++
			}
		}
	} else {
		err = errors.New("generateParameters - parameters cannot be passed empty")
	}

	// Prevents domain.*DatabaseType from being in the struct definition
	genParams = strings.Replace(genParams, "domain.*", "domain.", -1)

	return genParams, jsonParams, err
}

func (builder *builder) generateParamSetters(param *parameterDescriptor, class *structDescriptor, isSlice bool) (err error) {
	var set string

	if param.Size > 0 {

		if param.Null {
			if set, err = builder.updateParamTemplate("setwnullable", class.Name, param, true); err == nil {
				mname := "set-" + param.Name
				class.Methods[mname] = newMethod(mname, set)
			}
		} else {
			if set, err = builder.updateParamTemplate("set", class.Name, param, true); err == nil {
				mname := "set-" + param.Name
				class.Methods[mname] = newMethod(mname, set)
			}
		}
	}

	return err
}

func (builder *builder) generateParamGetters(param *parameterDescriptor, class *structDescriptor) (err error) {
	var get string

	if param.Null && !param.customType {
		if get, err = builder.updateParamTemplate("getwnullable", class.Name, param, true); err == nil {
			mname := "get-" + param.Name
			class.Methods[mname] = newMethod(mname, get)
		}
	} else {
		if get, err = builder.updateParamTemplate("get", class.Name, param, true); err == nil {
			mname := "get-" + param.Name
			class.Methods[mname] = newMethod(mname, get)
		}
	}
	return err
}

func (builder *builder) updateParamTemplate(templatePath string, class string, param *parameterDescriptor, pullPointer bool) (t string, err error) {

	var template *Template
	template, err = NewTemplate(builder.templatePath, templatePath)

	if err == nil {

		var pname = param.Name

		template.Repl("%class", class).
			Repl("%parameter", pname)

		// Private Parameter name
		template.Repl("%pparameter", pname)

		var goType = param.GoType

		if pullPointer {
			// PULL THE * OFF THE GO TYPE BECAUSE WE DON'T WANT THEM PASSING POINTERS TO THE OBJECT WE WANT A NEW POINTER
			goType = strings.Replace(param.GoType, "*", "", -1)
		}

		if param.customType && strings.Index(goType, "domain") < 0 {
			goType = "domain." + goType
		}

		template.Repl("%type", goType)

		template.Repl("%nullable", strconv.FormatBool(param.Null)).
			Repl("%size", strconv.Itoa(param.Size))

		t = template.Get()
	}

	return t, err
}

func appearsToBeSingleReturn(sprocString string) (res bool) {
	sprocString = strings.ToLower(sprocString)

	splitByNewline := strings.Split(sprocString, "\n")
	for index := range splitByNewline {
		var returnIndex = strings.Index(splitByNewline[index], "return")
		var singleIndex = strings.Index(splitByNewline[index], "single")

		if returnIndex >= 0 {
			res = singleIndex >= 0
			break
		}
	}

	return res
}

func (builder *builder) generateStoredProcedure(procedure *procedureDescriptor) (sproc string, signature string, err error) {

	// Differentiate between a read and execute stored procedure call
	var template = "sproc.exec"
	if procedure.Return != nil {
		template = "sproc.read"
	}

	var sprocTemplate *Template
	sprocTemplate, err = NewTemplate(builder.templatePath, template)
	if err == nil && len(sprocTemplate.Get()) > 0 {

		// Replace the %procedure with the name of the sproc\
		sprocTemplate.Repl("%pname", strings.Title(procedure.Name)). // %procedure
										Repl("%procedure", procedure.Name) // %procedure

		var methodParameters string // %parameters // Method Parameters
		var returnParameters string // %returns // Return Parameters
		var dbMapping string        // %dbparams // Mapping of db parameters to object properties

		// ***********************************************************
		// Build the Stored Procedure Method CallA Parameters
		// ***********************************************************

		for i := 0; i < len(procedure.PArray); i++ {
			pname := procedure.PArray[i].Name

			methodParameters += fmt.Sprintf("%s %s,", pname, procedure.PArray[i].GoType)

			// TODO: Make sure builder is correct
			dbMapping += fmt.Sprintf("%s,", pname)
		}

		// Update the template with the sproc parameter information
		sprocTemplate.Repl("%parameters", methodParameters). // %parameters
									Repl("%dbparams", dbMapping) // %dbparams

		// ***********************************************************
		// Build the Stored Procedure Method CallA Parameters
		// ***********************************************************

		if procedure.Return != nil {

			var readTemplate string
			var rets string
			var initReturn string

			if appearsToBeSingleReturn(procedure.FullProc) {
				rets = "domain.%s,"
				readTemplate = "sproc.read.classmap.single"
				initReturn = fmt.Sprintf("var ret%s domain.%s", procedure.Return.Name, procedure.Return.Name)
			} else {
				rets = "[]domain.%s,"
				readTemplate = "sproc.read.classmap"
				initReturn = fmt.Sprintf("var ret%s = make([]domain.%s, 0)", procedure.Return.Name, procedure.Return.Name)
			}

			returnParameters += fmt.Sprintf(rets, procedure.Return.Name)

			var classmap string
			if classmap, err = builder.buildClassMap(procedure.Return, &procedure.Imports, readTemplate); err == nil {

				classmap = strings.Replace(classmap, "%pnum", strconv.Itoa(procedure.Return.ID), -1)

				// Add the RetMaps to the code
				sprocTemplate.
					Repl("%retmap", classmap). // %retmap
					Repl("%initreturn", initReturn).
					Repl("%returns", returnParameters). // %returns
					Repl("%rname", fmt.Sprintf("ret%s", procedure.Return.Name))

				sproc = sprocTemplate.Get()
			}

			signature = fmt.Sprintf("\t%s(%s) (%s error)", procedure.Name, methodParameters, returnParameters)
		} else {
			sproc = sprocTemplate.Get()
			signature = fmt.Sprintf("\t%s(%s) (id int, affectedRows int, err error)", procedure.Name, methodParameters)
		}
	}

	return sproc, signature, err
}

func (builder *builder) buildClassMap(class *structDescriptor, imports *map[string]bool, template string) (classmap string, err error) {
	var classMapTemplate *Template
	classMapTemplate, err = NewTemplate(builder.templatePath, template)
	if err == nil {
		classmap = classMapTemplate.Get()
		if classmap, err = builder.assignClassMap(class, imports, classmap); err == nil {

		}
	}
	return classmap, err
}

func (builder *builder) assignClassMap(class *structDescriptor, imports *map[string]bool, templateval string) (classmap string, err error) {
	var cmap string
	var variables string
	var vmap string
	var classMapTemplate *Template
	classMapTemplate = NewTemplateEmpty()
	classMapTemplate.UpdateBase(templateval)
	classMapTemplate.Repl("%class", class.Name) // %dbparams

	for i := 0; i < len(class.PArray); i++ {

		if class.PArray[i].GoType == "time.Time" || class.PArray[i].GoType == "time.Duration" {

			(*imports)["time"] = true
		}

		// This checks to see if the struct representation has a null property while the procedure representation
		// does not have a null property to ensure that the address of the return is properly stored in the struct
		// upon creation
		var pname = class.PArray[i].Name

		cmap = builder.calculateClassMap(class, i, cmap, pname)

		var gotype = class.PArray[i].GoType

		if class.PArray[i].GoType == "bool" {
			gotype = "[]uint8"
		} else if class.PArray[i].GoType == "*bool" {
			gotype = "*[]uint8"
		}

		if !class.PArray[i].IsMapType {
			if strings.Index(gotype, "domain") >= 0 {
				variables += fmt.Sprintf("\t\t\t\t\t\tvar my%s %s\n", class.PArray[i].Name, strings.Replace(gotype, "domain", "*dal", 1))
			} else if class.PArray[i].customType {
				variables += fmt.Sprintf("\t\t\t\t\t\tvar my%s *dal.%s\n", class.PArray[i].Name, gotype)
			} else {
				variables += fmt.Sprintf("\t\t\t\t\t\tvar my%s %s\n", class.PArray[i].Name, gotype)
			}

			vmap += fmt.Sprintf("\t\t\t\t\t\t\t\t\t&my%s,\n", class.PArray[i].Name)
		}
	}

	classMapTemplate.Repl("%map", cmap). // %dbparams -- new class
						Repl("%variables", variables). // %dbparams   -- next
						Repl("%vparams", vmap)         // %dbparams           -- scan

	classmap = classMapTemplate.Get()
	return classmap, err
}

func (builder *builder) calculateClassMap(class *structDescriptor, i int, cmap string, pname string) string {
	if class.PArray[i].GoType == "bool" || class.PArray[i].GoType == "*bool" {
		cmap += fmt.Sprintf("%svar: my%s[0] > 0 && my%s[0] != 48, // converts uint8 to bool (48 is ASCII code for 0, which is reserved for false) \n", pname, pname, pname)
	} else if class.PArray[i].IsMapType {
		if class.PArray[i].MapSingle {
			cmap += fmt.Sprintf("%svar: nil, \n", pname)
		} else {
			cmap += fmt.Sprintf("%svar: make(%s, 0), \n", pname, class.PArray[i].GoType)
		}
	} else if class.PArray[i].customType && strings.Index(class.PArray[i].GoType, "[]") >= 0 {
		var structType = strings.Replace(class.PArray[i].GoType, "domain", "*dal", 1)

		cmap += fmt.Sprintf(`%s: func(structType %s) (interfaceType %s) {for index := range structType {interfaceType[index] = structType[index]}treturn interfaceType} (my%s),`+"\n", pname, structType, class.PArray[i].GoType, pname)
		//cmap = strings.Replace(cmap, "\\t", "\t", -1)
	} else {
		if builder.classes[class.Name].Parameters[class.PArray[i].Name].Null && !class.PArray[i].Null {
			cmap += fmt.Sprintf("%svar: &my%s, \n", pname, pname)
		} else {
			cmap += fmt.Sprintf("%svar: my%s, \n", pname, pname)
		}

	}
	return cmap
}

// stringToLines was pulled from https://siongui.github.io/2016/04/06/go-readlines-from-file-or-string/
func stringToLines(s string) (lines []string, err error) {

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return lines, err
}

// getFields retrieves a list of the fields in a string and strips off any space at the beginning and end of
// each of those fields
func getFields(s string) (fields []string, err error) {

	if len(s) > 0 {
		// Parse out the fields from the string
		fields = strings.Fields(s)

		if len(fields) > 0 {

			// Trim off the empty strings before and after each field
			for i, field := range fields {
				fields[i] = strings.TrimSpace(field)
			}
		} else {
			err = errors.New("Invalid Fields Array returned")
		}
	} else {
		err = errors.New("Cannot get fields of empty string")
	}

	return fields, err
}
