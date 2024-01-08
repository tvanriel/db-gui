package mysqlconnection

import "github.com/tvanriel/db-gui/app/domain"

var _ domain.Index = &MySQLIndex{}

type MySQLIndex struct {
	table        string
	nonUnique    bool
	keyName      string
	seqInIndex   int
	columnName   string
	collation    string
	cardinality  string
	subPart      string
	packed       string
	null         bool
	indexType    string
	comment      string
	indexComment string
	visible      bool
	expression   string
}

func (i *MySQLIndex) Table() string {
	return i.table
}

func (i *MySQLIndex) NonUnique() bool {
	return i.nonUnique
}

func (i *MySQLIndex) KeyName() string {
	return i.keyName
}

func (i *MySQLIndex) SeqInIndex() int {
	return i.seqInIndex
}

func (i *MySQLIndex) ColumnName() string {
	return i.columnName
}

func (i *MySQLIndex) Collation() string {
	return i.collation
}

func (i *MySQLIndex) Cardinality() string {
	return i.cardinality
}

func (i *MySQLIndex) SubPart() string {
	return i.subPart
}

func (i *MySQLIndex) Packed() string {
	return i.packed
}

func (i *MySQLIndex) Null() bool {
	return i.null
}

func (i *MySQLIndex) IndexType() string {
	return i.indexType
}

func (i *MySQLIndex) Comment() string {
	return i.comment
}

func (i *MySQLIndex) IndexComment() string {
	return i.indexComment
}

func (i *MySQLIndex) Visible() bool {
	return i.visible
}

func (i *MySQLIndex) Expression() string {
	return i.expression
}
