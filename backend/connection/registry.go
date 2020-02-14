package connection

import (
	"database/sql"
	"time"
)

// Procedures
const (
	createLog = "CreateLog"
)

// DatabaseConnection defines the methods that are required to fulfil the needs of the application for interacting with a database
type DatabaseConnection interface {
	CreateLog(logType int, text string, logError string, jobHistoryID string, createTime time.Time) (id int, affectedRows int, err error)
	Execute(sql string, args ...interface{}) (result sql.Result, err error)
	Exec(procedure *Procedure)
	Read(procedure *Procedure)
	Migrate(path string) (err error)
}
