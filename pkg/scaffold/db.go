package scaffold

import (
	"fmt"
	"strconv"
	"strings"
)

func separateDBType(dbType string) (dbtype string, size int, err error) {
	// Look for the first paren of the db type if it exists
	oparen := strings.Index(dbType, "(")
	cparen := strings.Index(dbType, ")")
	size = -1

	if oparen > 0 {
		dbtype = dbType[0:oparen]

		// Pull out the size of the field from the string
		if cparen >= 0 {
			size, err = strconv.Atoi(dbType[oparen+1 : cparen])
		}
	} else {
		dbtype = dbType
	}

	return dbtype, size, err
}

// TODO:
//TINYINT  -->  int8
//UNSIGNED TINYINT  -->  uint8
//SMALLINT  -->  int16
//UNSIGNED SMALLINT  -->  uint16
//MEDIUMINT, INT  -->  int32
//UNSIGNED MEDIUMINT, UNSIGNED INT  -->  uint32
//BIGINT  -->  int64
//UNSIGNED BIGINT  -->  uint64
//FLOAT  -->  float32
//DOUBLE  -->  float64
//DECIMAL  -->  float64
//TIMESTAMP, DATETIME  -->  time.Time
//DATE  -->  mysql.Date
//TIME  -->  time.Duration
//YEAR  -->  int16
//CHAR, VARCHAR, BINARY, VARBINARY  -->  []byte
//TEXT, TINYTEXT, MEDIUMTEXT, LONGTEX  -->  []byte
//BLOB, TINYBLOB, MEDIUMBLOB, LONGBLOB  -->  []byte
//BIT  -->  []byte
//SET, ENUM  -->  []byte
//NULL  -->  nil

var types = map[string]string{
	"BIT":       "bool",
	"BOOL":      "bool",
	"BOOLEAN":   "bool",
	"SMALLINT":  "int16",
	"MEDIUMINT": "int",
	"INT":       "int",
	"INTEGER":   "int",
	"BIGINT":    "int",
	"DECIMAL":   "float32",
	"DEC":       "float32",
	"NUMERIC":   "float32",
	"FIXED":     "float32",
	"FLOAT":     "float32",
	"DOUBLE":    "float32",
	"REAL":      "float32",
	"DATE":      "time.Time",
	"DATETIME":  "time.Time",
	"TIMESTAMP": "time.Time",
	"TIME":      "time.Duration",
	"NVARCHAR":  "string",
	"VARCHAR":   "string",
}

func dbTypeToGoType(DBType string, size int, imports *map[string]bool, nullable bool) (goType string, customType bool, err error) {
	if len(DBType) > 0 {
		if strings.Contains(strings.ToUpper(DBType), "VARCHAR") || strings.Contains(strings.ToUpper(DBType), "TEXT") {
			DBType = "NVARCHAR"
		} else if strings.Contains(strings.ToUpper(DBType), "VARCHAR") {
			DBType = "VARCHAR"
		}

		goType = types[strings.ToUpper(DBType)]
		if len(goType) == 0 {
			goType = DBType
			(*imports)["github.com/nortonlifelock/aegis/pkg/domain"] = true
			customType = true

			arrayIndex := strings.Index(DBType, "[]")
			if arrayIndex >= 0 {
				goType = DBType[:arrayIndex+2] + "domain." + DBType[arrayIndex+2:]
			}
		}

		if goType == "time.Time" || goType == "time.Duration" {
			(*imports)["time"] = true
		}

		if nullable {
			goType = fmt.Sprintf("*%s", goType)
		}
	} else {
		err = fmt.Errorf("empty db type provided")
	}

	return goType, customType, err
}
