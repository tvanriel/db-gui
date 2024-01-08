package mysqlconnection

import "github.com/tvanriel/db-gui/app/domain"

var _ domain.Table = &MySQLTable{}

type MySQLTable struct {
	indeces []*MySQLIndex
	columns []*MySQLColumn
	name    string
}

func (t *MySQLTable) Columns() []string {
	names := make([]string, len(t.columns))

	for i, col := range t.columns {
		names[i] = col.field
	}

	return names
}
func (t *MySQLTable) GetColumn(name string) domain.Column {
	for _, col := range t.columns {

		if col.field == name {
			return col
		}

	}

	return nil
}

func (t *MySQLTable) Name() string {
	return t.name
}

func (t *MySQLTable) IndexSequence() []int {
	sequence := []int{}

	for _, index := range t.indeces {
		sequence = append(sequence, index.seqInIndex)
	}

	return sequence
}

func (t *MySQLTable) GetIndex(seq int) domain.Index {
	for _, index := range t.indeces {

		if index.seqInIndex == seq {
			return index
		}

	}

	return nil
}
