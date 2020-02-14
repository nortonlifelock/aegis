package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nortonlifelock/aegis/backend/config"
	"github.com/nortonlifelock/aegis/backend/connection"
	"github.com/nortonlifelock/aegis/backend/files"
	"github.com/nortonlifelock/scaffold"
)

const (
	migrateMode = 1 << iota
	deleteMode
	procedureMode
	generateMode
)

// Setting up config arguments for starting the job runner
// configFile - The filename of the config to load
// configPath - The directory path of the config to load

// When the generate flag is not present, the generate paths are not required (domainGenPath, dalGenPath). If you're pulling the code from Github, the generated files are included
// domainGenPath - The path where the domain interfaces are generated (/path-to/domain)
// dalGenPath - The path where the dal objects implementing the domain are generated (/path-to/database)

// sprocPath - The path to where the stored procedures waiting for generation are located (/path-to/aegis-scaffold/procedures)
// schemaMigrationPath - The path where the migrate files are located (/path-to/aegis-scaffold/migrations)
// templatePath - The path where the 'templates' directory is located (/path-to/aegis-scaffold/templates)

// migrateFlag - One of four flags [m|d|p|g] that control what parts of the scaffolding execute. If m is present among the flags, the migrations will be run. If none of the 4 flags are present, all steps of the scaffolding execute
// deleteFlag - One of four flags [m|d|p|g] that control what parts of the scaffolding execute. If d is present among the flags, the old generated files will be deleted. If none of the 4 flags are present, all steps of the scaffolding execute
// procedureFlag - One of four flags [m|d|p|g] that control what parts of the scaffolding execute. If p is present among the flags, the stored procedures will be processed. If none of the 4 flags are present, all steps of the scaffolding execute
// generateFlag - One of four flags [m|d|p|g] that control what parts of the scaffolding execute. If p is present among the flags, the stored procedures files will be generated. If none of the 4 flags are present, all steps of the scaffolding execute
func RunScaffold(configFile, configPath, domainGenPath, dalGenPath, sprocPath, schemaMigrationPath, templatePath string, migrateFlag, deleteFlag, procedureFlag, generateFlag bool) {
	var err error

	if len(configFile) == 0 {
		configFile = "app.json"
	}

	// executionMode controls which steps of the scaffolding execute based on which flags were present
	// if none of the flags are present, all steps are enabled by default
	var executionMode = getExecutionMode(migrateFlag, deleteFlag, procedureFlag, generateFlag)

	needGeneratePaths := deleteFlag || generateFlag

	appConfig, domainPath, dalPath, schemaMigration, err := processFlags(sprocPath, domainGenPath, dalGenPath, schemaMigrationPath, configFile, configPath, templatePath, needGeneratePaths)

	if err == nil {
		var dbConn connection.DatabaseConnection
		if dbConn, err = connection.NewConnection(appConfig); err == nil {
			err = runScaffolding(dbConn, executionMode, schemaMigration, domainPath, dalPath, sprocPath, templatePath)
		} else {
			err = fmt.Errorf("error while connecting to database - %s", err.Error())
		}
	}

	if err == nil {
		fmt.Println("\nScaffolding finished without error")
	} else {
		fmt.Println(err)
	}
}

func runScaffolding(dbConn connection.DatabaseConnection, executionMode int, schemaMigration string, domainPath string, dalPath string, sprocPath string, templatePath string) (err error) {
	if hasFlag(executionMode, migrateMode) {
		fmt.Println("BEGINNING processing of database changes")
		if err = dbConn.Migrate(schemaMigration); err == nil {
			fmt.Println("FINISHED processing of database changes")
		} else {
			err = fmt.Errorf("error while performing migration - %s", err.Error())
		}
	}

	if err == nil {
		if hasFlag(executionMode, deleteMode) {
			fmt.Println("DELETING old generated interfaces and structures")
			if err = cleanOutOldGeneratedFiles(domainPath, dalPath); err == nil {
				fmt.Println("FINISHED deleting old generated interfaces and structures")
			} else {
				err = fmt.Errorf("error while cleaning out old generated files - %s", err.Error())
			}
		}
	}

	if err == nil {
		if hasFlag(executionMode, procedureMode) {
			fmt.Println("BEGINNING processing of stored procedures")
			if err = scaffold.ProcessSprocs(dbConn, sprocPath, domainPath, dalPath, templatePath, hasFlag(executionMode, generateMode)); err == nil {
				fmt.Println("FINISHED processing of stored procedures")
			} else {
				err = fmt.Errorf("error while processing stored procedures - %s", err.Error())
			}
		}
	}

	return err
}

func processFlags(sprocPath string, domainGenPath string, dalGenPath string, schemaMigrationPath string, configFile string, configPath string, templatePath string, needGeneratePaths bool) (appConfig config.AppConfig, domainPath string, dalPath string, schemaMigration string, err error) {
	err = validateFlags(sprocPath, domainGenPath, dalGenPath, schemaMigrationPath, configFile, configPath, templatePath, needGeneratePaths)
	if err == nil {
		if appConfig, err = config.LoadConfig(configPath, configFile); err == nil {

			domainPath = domainGenPath
			dalPath = dalGenPath
			schemaMigration = schemaMigrationPath

			// Generate the directory to hold the generated interfaces if it does not exist
			if _, err := os.Stat(domainPath); err != nil {
				_ = os.MkdirAll(domainPath, 0775)
			}

			// Generate the directory to hold the generated structs if it does not exist
			if _, err := os.Stat(dalPath); err != nil {
				_ = os.MkdirAll(dalPath, 0775)
			}

			if len(templatePath) == 0 || !files.ValidDir(templatePath) {
				templatePath = appConfig.AegisPath()
			}

		} else {
			err = fmt.Errorf("errors while loading config | %s", err.Error())
		}
	}

	return appConfig, domainPath, dalPath, schemaMigration, err
}

func validateFlags(sprocPath string, domainGenPath string, dalGenPath string, schemaMigrationPath string, configFile string, configPath string, templatePath string, needGeneratePaths bool) (err error) {
	if files.ValidDir(sprocPath) {
		if !needGeneratePaths || files.ValidDir(domainGenPath) {
			if !needGeneratePaths || files.ValidDir(dalGenPath) {

				if files.ValidDir(schemaMigrationPath) {

					if len(configFile) > 0 && len(configPath) > 0 && len(templatePath) > 0 {

					} else {
						err = fmt.Errorf("need to provide the -config, -cpath, and -tpath")
					}
				} else {
					err = fmt.Errorf("cannot have an empty schema migration path")
				}
			} else {
				err = fmt.Errorf("cannot have an empty dal generation path")
			}
		} else {
			err = fmt.Errorf("cannot have an empty domain generation path")
		}
	} else {
		err = fmt.Errorf("cannot have an empty path to stored procedures")
	}

	return err
}

func cleanOutOldGeneratedFiles(domainpath string, sprocpath string) (err error) {
	err = files.ExecuteThroughDirectory(domainpath, true, func(fpath string, file os.FileInfo) (err error) {

		filename := filepath.Base(fpath)
		if strings.Index(filename, "generated") >= 0 {

			fmt.Printf("Removing File [%s]\n", fpath)
			err = os.Remove(fpath)
		}

		return err
	})

	if err == nil {
		err = files.ExecuteThroughDirectory(sprocpath, true, func(fpath string, file os.FileInfo) (err error) {

			filename := filepath.Base(fpath)
			if strings.Index(filename, "generated") >= 0 {

				fmt.Printf("Removing File [%s]\n", fpath)
				err = os.Remove(fpath)
			}

			return err
		})
	}

	return err
}

func getExecutionMode(migrateFlag bool, deleteFlag bool, procedureFlag bool, generateFlag bool) (mode int) {
	if !migrateFlag && !deleteFlag && !procedureFlag && !generateFlag {
		// if none of the flags were set, we do all processes by default
		mode = migrateMode | deleteMode | procedureMode | generateMode
	} else {
		// otherwise we enable each mode bit by performing a logical or-equals if the corresponding flag is present
		if migrateFlag {
			mode |= migrateMode
		}
		if deleteFlag {
			mode |= deleteMode
		}
		if procedureFlag {
			mode |= procedureMode
		}
		if generateFlag {
			mode |= generateMode
		}
	}

	return mode
}

func hasFlag(mode int, flag int) bool {
	return mode&flag == flag
}
