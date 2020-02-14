package connection

import (
	"database/sql"
	"github.com/pkg/errors"
)

// Execute executes raw sql strings against the database
func (conn *dbconn) Execute(sql string, args ...interface{}) (result sql.Result, err error) {

	// Verify that the sql statement is not empty
	if len(sql) > 0 {

		transaction, err2 := conn.db.Begin()
		if err2 == nil {

			if result, err = transaction.Exec(sql, args...); err == nil {
				err = transaction.Commit()
			} else {
				_ = transaction.Rollback()
			}
		} else {
			err = err2
		}

	} else {
		err = errors.New("empty SQL string")
	}
	return result, err
}
