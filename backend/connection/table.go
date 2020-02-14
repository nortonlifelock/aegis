package connection

import (
	"fmt"
	"github.com/pkg/errors"
)

type table struct {
	Name    string
	columns []*column
}

func (table *table) addColumn(column *column) (err error) {

	if column != nil {
		table.columns = append(table.columns, column)
	} else {
		err = errors.New("invalid column object")
	}

	return err
}

func (table *table) valid() bool {
	valid := false

	if table != nil && len(table.Name) > 0 && len(table.columns) > 0 {
		valid = true
	}

	return valid
}

// Builds the create statement for adding a table to the mysql database
func (table *table) build(database string) (statement string, err error) {

	// Check for valid table object
	if table != nil && table.valid() {
		// Table is valid so start building the statement
		statement = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s ( ", database, table.Name)

		columnCount := len(table.columns)

		for i := 0; i < columnCount; i++ {

			column := table.columns[i]

			// Check for valid column object
			if column != nil {

				statement += column.CName + " " + column.CType + " "

				// Define whether or not the field is nullable in the table
				if column.CNull {
					statement += "NULL "
				} else {
					statement += "NOT NULL "
				}

				if column.CPrimary {
					statement += "PRIMARY KEY "
				}

				// Define a default if one exists
				if len(column.CDefault) > 0 {
					statement += "DEFAULT " + column.CDefault
				}

				// Add commas
				if i != (columnCount - 1) {
					statement += ", "
				}
			} else {
				err = errors.New("Invalid null column while attempting to build create table sql string")
			}
		}

		// Close the create table statement
		statement += ")"

	} else {
		err = errors.New("Invalid Table Object")
	}

	return statement, err
}
