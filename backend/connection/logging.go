package connection

import (
	"database/sql"
	"time"
)

// CreateLog is unlike our other stored procedures. The go code for this one is written by hand (as opposed to being generated)
// This is due to the fact that the dbconn must implement the DatabaseConnection interface, which exists before generation
func (conn dbconn) CreateLog(logType int, text string, logError string, jobHistoryID string, createDate time.Time) (id int, affectedRows int, err error) {

	var etext string
	if len(logError) > 0 {
		etext = logError
	}

	conn.Exec(&Procedure{
		createLog,
		[]interface{}{logType, text, etext, jobHistoryID, createDate},
		func(results interface{}, dberr error) {
			err = dberr

			if result, ok := results.(sql.Result); ok {
				var idOut int64

				// Get the id of the last inserted record
				if idOut, err = result.LastInsertId(); err == nil {
					id = int(idOut)
				}

				// Get the number of affected rows for the execution
				if idOut, err = result.RowsAffected(); ok {
					affectedRows = int(idOut)
				}
			}

		},
	})

	return id, affectedRows, err
}
