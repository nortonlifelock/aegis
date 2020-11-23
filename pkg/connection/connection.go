package connection

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/nortonlifelock/aegis/pkg/files"

	"sync"

	// Adding in the mysql driver using a blank import to register the driver in the database/sql package
	_ "github.com/go-sql-driver/mysql"

	"github.com/pkg/errors"
)

type dbconn struct {
	db      *sql.DB
	execute chan *Procedure
	read    chan *Procedure
	schema  string
}

// Procedure holds the information required for call a stored procedure
type Procedure struct {
	Proc       string
	Parameters []interface{}
	Callback   func(results interface{}, dberr error)
}

type dbConfig interface {
	DBPath() string
	DBPort() string
	DBUsername() string
	DBPassword() string
	DBSchema() string
}

// NewConnection returns a connection to the sql database
func NewConnection(config dbConfig) (DatabaseConnection, error) {
	var proto = "tcp"
	return connectionBuilder(buildConnString(
		config.DBPath(),
		config.DBPort(),
		config.DBSchema(),
		config.DBUsername(),
		config.DBPassword(),
		proto,
		true,
	), config.DBSchema())
}

func connectionBuilder(connectionString string, schema string) (DatabaseConnection, error) {
	var db *dbconn
	var err error

	var dbInstance *sql.DB
	if dbInstance, err = sql.Open("mysql", connectionString); err == nil {

		db = &dbconn{
			db:      dbInstance,
			execute: make(chan *Procedure),
			read:    make(chan *Procedure),
			schema:  schema,
		}

		go db.listenForDatabaseRequests(context.Background())
	}

	return db, err
}

// Migrate executes db changes recursively through the directory passed in the path argument
// Only files containing '.up.' are executed. Each '.up.' change file should have a corresponding '.down.' file that undoes the change
func (conn *dbconn) Migrate(path string) (err error) {
	var exists bool
	if exists, err = conn.tableExists("schema_migrations"); err == nil {
		if !exists {
			_, err = conn.Execute("CREATE TABLE `schema_migrations` (`Name` VARCHAR(300) NOT NULL, `Date` DATETIME NOT NULL DEFAULT NOW());")
		}
	}

	if err == nil {
		err = files.ExecuteThroughDirectory(path, false, func(fpath string, file os.FileInfo) (err error) {

			if strings.Contains(file.Name(), ".up.") {
				var exists bool
				if exists, err = conn.dbChangeExists(file.Name()); !exists && err == nil {

					fmt.Printf("Executing %v\n", file.Name())

					var command string
					// Get the sql string from the file
					if command, err = files.GetStringFromFile(fpath); err == nil {

						// Execute the file against the database
						_, err = conn.Execute(command)

						if err == nil {
							_, err = conn.Execute(fmt.Sprintf("INSERT INTO schema_migrations (Name) VALUE ('%v');", file.Name()))
						} else {
							fmt.Printf("Error in file [%s/%s] : %s\n", path, file.Name(), err.Error())
						}
					}
				}
			}

			return err
		})
	}

	return err
}

// This method checks the db changes table for whether or not this change has been applied
func (conn *dbconn) dbChangeExists(filename string) (exists bool, err error) {
	if len(filename) > 0 {
		command := fmt.Sprintf("SELECT EXISTS(SELECT Name FROM schema_migrations WHERE Name = \"%s\" limit 1)", filename)
		exists, err = conn.existsCheck(command)
	}

	return exists, err
}

// Executes a query against the database that is passed to it. It expects a single row to be returned
// with the first column containing either a 1 or 0 indicating whether something exists or not.
// This function is primarily used for determining the existence of database Schema objects such as tables
// or stored procedures
func (conn *dbconn) existsCheck(sql string) (exists bool, err error) {
	if len(sql) > 0 {
		// Execute the sql against the database
		if rows, err2 := conn.db.Query(sql); err2 == nil && rows != nil {

			if err = rows.Err(); err == nil {
				defer func() {
					_ = rows.Close()
				}()

				for err == nil && rows.Next() {
					err = rows.Scan(&exists)
				}
			}
		} else {
			err = err2
		}
	} else {
		err = errors.New("empty SQL String")
	}

	return exists, err
}

func (conn *dbconn) tableExists(table string) (exists bool, err error) {
	if len(table) > 0 {

		// TODO: Validate table name to ensure that there is not sql injection here
		var command string
		if len(conn.schema) > 0 {
			command = fmt.Sprintf("select EXISTS(select * from information_schema.tables where table_name = '%s'  and table_schema = '%s' limit 1) as result;", table, conn.schema)
		} else {
			command = fmt.Sprintf("select EXISTS(select * from information_schema.tables where table_name = '%s' limit 1) as result;", table)
		}

		// Run exists check
		exists, err = conn.existsCheck(command)
	} else {
		err = errors.New("empty table name")
	}

	return exists, err
}

// Exec executes a procedure against the database and calls a function to process the results
func (conn *dbconn) Exec(proc *Procedure) {
	var wg = sync.WaitGroup{}

	wg.Add(1)
	var originalCallback = proc.Callback

	// wrap the original callback function with a wait group
	proc.Callback = func(results interface{}, dberr error) {
		defer wg.Done()
		originalCallback(results, dberr)
	}

	conn.execute <- proc

	wg.Wait()
}

// Read reads information from the database and calls a function to process the results
func (conn *dbconn) Read(proc *Procedure) {
	var wg = sync.WaitGroup{}

	wg.Add(1)
	var originalCallback = proc.Callback

	// wrap the original callback function with a wait group
	proc.Callback = func(results interface{}, dberr error) {
		defer wg.Done()
		originalCallback(results, dberr)
	}

	conn.read <- proc

	wg.Wait()
}

func (conn *dbconn) listenForDatabaseRequests(ctx context.Context) {
	defer func() {
		select {
		case <-ctx.Done():
			return
		default:
			conn.panicRecover()

			go conn.listenForDatabaseRequests(ctx)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case proc := <-conn.execute:
			proc.Callback(conn.exec(ctx, proc.Proc, proc.Parameters))
		case procedure := <-conn.read:
			if procedure != nil {
				procedure.Callback(conn.call(ctx, procedure.Proc, procedure.Parameters))
			}
		}
	}
}

// Calls a stored Procedure expecting a set of records to be returned in the form of a mysql.Row
func (conn *dbconn) call(ctx context.Context, sproc string, params []interface{}) (rows *sql.Rows, err error) {

	// TODO: check Parameters for injection here
	// TODO: update the params parameter to use an interface type? So that Parameters to the sproc can have the types validated?

	// Build the command
	var command string
	if command, err = conn.buildSprocCall(sproc, params); err == nil {

		// Verify that a valid command was created
		if len(command) > 0 {
			var stmt *sql.Stmt
			if stmt, err = conn.db.Prepare(command); err == nil {
				defer func() {
					_ = stmt.Close()
				}()

				// Return rows and error from the query
				rows, err = stmt.Query(params...)
			}
		} else {
			err = errors.Errorf("Invalid command creation for stored Procedure %s", sproc)
		}
	} else {
		err = errors.Errorf("Error while building execution call for stored Procedure %s: %v", sproc, err)
	}

	return rows, err
}

// Calls a stored Procedure expecting a set of records to be returned in the form of a mysql.Row
func (conn *dbconn) exec(ctx context.Context, name string, params []interface{}) (result sql.Result, err error) {

	// TODO: check Parameters for injection here
	// TODO: update the params parameter to use an interface type? So that Parameters to the sproc can have the types validated?

	// Build the command
	var command string
	if command, err = conn.buildSprocCall(name, params); err == nil {

		// Verify that a valid command was created
		if len(command) > 0 {
			result, err = conn.db.ExecContext(ctx, command, params...)
		} else {
			err = errors.Errorf("Invalid command creation for stored Procedure %s", name)
		}
	} else {
		err = errors.Errorf("Error while building execution call for stored Procedure %s: %v", name, err)
	}
	return result, err
}

func buildConnString(server string, port string, schema string, user string, password string, protocol string, multi bool) (connString string) {
	connString = fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true", user, password, protocol, server, port, schema)
	if multi {
		connString += "&multiStatements=true"
	}

	return connString
}

// Builds a call for a stored Procedure given a name and a list of parameter values in the
// order that they need for calling the sproc
func (conn *dbconn) buildSprocCall(name string, params []interface{}) (command string, err error) {
	if len(name) > 0 {
		// MySQL command for calling an sproc and sproc name
		command = fmt.Sprintf("CALL `%s`(", name)

		if params != nil {
			var plength = len(params)

			if plength > 0 {
				// Build the Parameters going into the sproc
				for index := range params {
					command += "?"

					// if it is not the last field then add a comma for the next parameter
					if index != (plength - 1) {
						command += ","
					}
				}
			}
		}

		command += ");"
	} else {
		err = errors.New("stored Procedure name cannot be empty")
	}

	return command, err
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (conn *dbconn) panicRecover() {

	if retVal := recover(); retVal != nil {

		var stackWrapper = errors.New(retVal.(string))
		err, ok := errors.Cause(stackWrapper).(stackTracer)

		var newLog = retVal.(string)

		if ok {
			stackTrace := err.StackTrace()
			newLog = fmt.Sprintf("%s\n%+v\n", retVal, stackTrace[3:]) // top two frames
		}

		fmt.Println(newLog)
	}
}
