package connection

type column struct {
	CName    string
	CType    string
	CNull    bool
	CDefault string
	CPrimary bool
}

func newColumn(cname string, ctype string, cnull bool, cdefault string, cprimary bool) (col *column, err error) {
	newColumn := &column{
		cname,
		ctype,
		cnull,
		cdefault,
		cprimary,
	}

	if newColumn.valid() {
		col = newColumn
	}

	return col, err
}

func (column *column) valid() (valid bool) {

	if column != nil && len(column.CName) > 0 && len(column.CType) > 0 {
		valid = true
	}

	return valid
}
