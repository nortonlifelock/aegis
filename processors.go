package scaffold

import (
	"fmt"
	"github.com/nortonlifelock/files"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ProcessSprocs generates classes for every stored procedure that reads information back from the database, as well as generates the go code
// required to execute the sproc against the database and parse the response
func ProcessSprocs(dbConn idb, sprocPath string, domainEngPath, sprocGenPath, templatePath string, generateFiles bool) (err error) {

	var methodMap = make(map[string]string)
	var includesMap = make(map[string]bool)

	var builder *builder
	// Create a builder for use with the scaffolding
	if builder, err = newBuilder(dbConn, templatePath); err == nil {

		// Recurse the directory tree and execute the database changes that haven't
		// yet been applied
		err = files.ExecuteThroughDirectory(sprocPath, true, func(path string, file os.FileInfo) (err error) {

			// Verify that we're only executing sql files against the database
			if filepath.Ext(path) == ".sql" {

				// Process the sproc and add it to the maps for the builder
				err = builder.processSproc(path)

				if err != nil {
					fmt.Println(err)
				}
			}

			return err
		})

		if err == nil {

			// Create domain objects in the database package
			path := sprocGenPath + "/dal/"

			if generateFiles {
				if _, err = os.Stat(path); err != nil {
					// this error will overwrite the error that brought us to this block.
					// if MkdirAll is successful, the error should be nil and the next block should enter
					err = os.MkdirAll(path, 0775)
				}
			}

			if err == nil {

				if generateFiles {
					path = path + "generated.%s.go"

					// For each class build the file for that class
					for _, class := range builder.classes {

						var genClass string
						if genClass, err = builder.generateClass(class); err == nil {

							filePath := fmt.Sprintf(path, strings.ToLower(class.Name))

							// Overwrite the file in the folder structure
							var formattedTemplate []byte
							formattedTemplate, err = format.Source([]byte(genClass))
							if err == nil {
								err = ioutil.WriteFile(filePath, formattedTemplate, 0644)
							}

							if err != nil {
								fmt.Println(err)
							}
						}
					}
				}

				var signatures string
				var interfaceImports map[string]bool
				signatures, interfaceImports, err = createsProceduresInDBAndGenerateSignatures(builder, methodMap, includesMap, dbConn)

				if generateFiles {
					path = domainEngPath + "/generated.interface.go"
					err = generateDatabaseInterface(signatures, templatePath, interfaceImports, path)
					if err != nil {
						fmt.Println(err.Error())
					}
				}

				// Build the methods files for each database
				if err == nil && generateFiles {
					path = sprocGenPath + "/generated.procedures.go"
					err = generateGoCodeThatExecutesSprocs(includesMap, methodMap, templatePath, path)
				}
			}
		}
	}

	return err
}

func generateGoCodeThatExecutesSprocs(includesMap map[string]bool, methodMap map[string]string, templatePath string, path string) (err error) {
	var methods string
	var imports string
	for imp := range includesMap {
		imports += "\t\"" + imp + "\"\n"
	}
	alphabetizedMethods := getOrderedKeyListFromStringMap(methodMap)
	for _, methodName := range alphabetizedMethods {
		if len(methodMap[methodName]) > 0 {
			methods += methodMap[methodName]
		} else {
			fmt.Println("Failed to generate method for", methodName)
		}
	}
	var dbsprocTemplate *Template
	dbsprocTemplate, err = NewTemplate(templatePath, "dbsproc")
	if err == nil {

		// Copy the methods into the template string
		dbsprocTemplate.Repl("%methods", methods)
		dbsprocTemplate.Repl("%imports", imports) // Add the imports to the template

		var formattedTemplate []byte
		formattedTemplate, err = format.Source([]byte(dbsprocTemplate.Get()))
		if err == nil {
			err = ioutil.WriteFile(path, formattedTemplate, 0644)
		} else {
			err = ioutil.WriteFile(path, []byte(dbsprocTemplate.Get()), 0644)
		}
	}
	return err
}

func getOrderedKeyListFromStringMap(in map[string]string) (alphabetized []string) {
	alphabetized = make([]string, 0)
	for key := range in {
		alphabetized = append(alphabetized, key)
	}
	sort.Strings(alphabetized)

	return alphabetized
}

func getOrderedKeyListFromMethodMap(in map[string]*methodDescriptor) (alphabetized []string) {
	alphabetized = make([]string, 0)
	for key := range in {
		alphabetized = append(alphabetized, key)
	}
	sort.Strings(alphabetized)

	return alphabetized
}

func getOrderedListFromBoolMap(in map[string]bool) (alphabetized []string) {
	alphabetized = make([]string, 0)
	for key := range in {
		alphabetized = append(alphabetized, key)
	}
	sort.Strings(alphabetized)

	return alphabetized
}

func getOrderedListFromParameterMap(in map[string]*parameterDescriptor) (alphabetized []string) {
	alphabetized = make([]string, 0)
	for key := range in {
		alphabetized = append(alphabetized, key)
	}
	sort.Strings(alphabetized)

	return alphabetized
}

func getOrderedListFromProcedureMap(in map[string]*procedureDescriptor) (alphabetized []string) {
	alphabetized = make([]string, 0)
	for key := range in {
		alphabetized = append(alphabetized, key)
	}
	sort.Strings(alphabetized)

	return alphabetized
}

func generateDatabaseInterface(signatures string, templatePath string, interfaceImports map[string]bool, path string) (err error) {
	if len(signatures) > 0 {
		var interfaceTemplate *Template
		interfaceTemplate, err = NewTemplate(templatePath, "dbsproc_interface")

		if err == nil {

			var imports string
			for k := range interfaceImports {
				imports = fmt.Sprintf("%s\n\t\"%s\"", imports, k)
			}

			interfaceTemplate.
				Repl("%imports", imports). // Add the imports to the template
				Repl("%methods", strings.Replace(signatures, "domain.", "", -1))

			var formattedTemplate []byte
			formattedTemplate, err = format.Source([]byte(interfaceTemplate.Get()))
			if err == nil {
				err = ioutil.WriteFile(path, formattedTemplate, 0644)
			}
		}
	}

	return err
}

func createsProceduresInDBAndGenerateSignatures(builder *builder, methodMap map[string]string, includesMap map[string]bool, dbConn idb) (string, map[string]bool, error) {
	var err error

	var signatures string
	var interfaceImports = make(map[string]bool)
	// For each sproc build the code file for accessing the file
	alphabetizedProcedures := getOrderedListFromProcedureMap(builder.procs)
	for _, index := range alphabetizedProcedures {

		var genMethod string
		var signature string

		if genMethod, signature, err = builder.generateStoredProcedure(builder.procs[index]); err == nil {

			// Don't want to duplicate interface method CreateLog
			if strings.Index(signature, "CreateLog") < 0 {
				signatures = fmt.Sprintf("%s\n%s", signatures, signature)
			}

			// Add this sproc method to this database's generated methods
			methodMap[builder.procs[index].Name] = genMethod

			// Add includes for this sproc
			for k := range builder.procs[index].Imports {
				interfaceImports[k] = true
				includesMap[k] = true
			}
			// Execute the sproc create against the db
			_, err = dbConn.Execute(builder.procs[index].FullProc)

			if err != nil { //the error is overwritten unless it was from the last function call
				fmt.Println(builder.procs[index].Name, err.Error())
			}
		} else {
			fmt.Println(builder.procs[index].Name, err.Error())
		}
	}
	return signatures, interfaceImports, err
}
