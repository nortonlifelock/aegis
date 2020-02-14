package database

import (
	"database/sql"
	"github.com/nortonlifelock/aegis/interal/connection"
	"github.com/nortonlifelock/aegis/internal/domain"
	"github.com/pkg/errors"
)

type dbconn struct {
	connection.DatabaseConnection
}

type dbConfig interface {
	DBPath() string
	DBPort() string
	DBUsername() string
	DBPassword() string
	DBSchema() string
}

// NewConnection establishes and returns a connection to the database using the config
func NewConnection(config dbConfig) (domain.DatabaseConnection, error) {
	var db connection.DatabaseConnection
	var err error

	db, err = connection.NewConnection(config)

	var thisDBConn = &dbconn{
		db,
	}

	return thisDBConn, err
}

func (conn *dbconn) getRows(results interface{}, process func(rows *sql.Rows) (err error)) (err error) {
	var rows *sql.Rows
	var ok bool

	if rows, ok = results.(*sql.Rows); ok {
		//noinspection GoUnhandledErrorResult
		defer rows.Close()

		// Check for error from the db
		if rows.Err() == nil {

			// Loop through the results and process them
			for rows.Next() {
				if err = process(rows); err != nil {

					// Break the loop in the event of an error
					break
				}
			}
		} else {
			// Capture the error from the database
			err = rows.Err()
		}
	} else {
		// Expected a row to be returned
		err = errors.Errorf("Invalid return from SQL query. Expected return of *sql.Rows")
	}

	return err
}
