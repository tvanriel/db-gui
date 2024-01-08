package mysqlconnection

import "github.com/tvanriel/db-gui/app/domain"

var _ domain.Column = &MySQLColumn{}

type MySQLColumn struct {
	field      string
	colType    string
	null       string
	comment    string
	colDefault string
}

func (c *MySQLColumn) Field() string {
	return c.field
}

func (c *MySQLColumn) Type() string {
	return c.colType
}

func (c *MySQLColumn) Null() bool {
	return c.null == "YES"
}

func (c *MySQLColumn) Default() string {
	return c.colDefault
}

func (c *MySQLColumn) Comment() string {
	return c.comment
}
