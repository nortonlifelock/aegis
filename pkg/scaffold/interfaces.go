package scaffold

import (
	"database/sql"
)

const (
	initiated = iota // Beginning of sproc processing
	returns   = iota // Found the open comment and now looking for the start of returns
	sproc     = iota // Out of the comment now and need to build the sproc call
)

const (
	read    = iota // The Sproc is a Read Sproc
	execute = iota // The Sproc is an Execute Sproc
)

type idb interface {
	Execute(sql string, args ...interface{}) (result sql.Result, err error)
	Migrate(path string) (err error)
}
