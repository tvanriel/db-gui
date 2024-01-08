package mysqlconnection

import (
	"database/sql"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/tvanriel/db-gui/app/domain"
)

/*
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"columns"`
	Result [][]struct {
		Null  bool   `json:"null"`
		Value string `json:"value"`
	} `json:"result"`
	Error string `json:"error"`
*/

var _ domain.Resultset = &MySQLResult{}

type MySQLResult struct {
	Results []struct {
		ColNames []string
		ColTypes []string
		Result   [][]*string
		Error    string
	}
}

func (m MySQLResult) Get() []struct {
	ColNames []string
	ColTypes []string
	Result   [][]*string
	Error    string
} {
	return m.Results
}

type MySQLResultSet struct {
	ColNames []string
	ColTypes []string
	Result   [][]*string
	Error    string
}

func (m *MySQLResultSet) ColumnNames() []string {
	return m.ColNames
}
func (m *MySQLResultSet) ColumnTypes() []string {
	return m.ColTypes
}

func (m *MySQLResultSet) Rows() [][]*string {
	return m.Result
}

func NewMysqlResultsetWithError(rows *sql.Rows, err error) (domain.Resultset, error) {
	if err != nil {
		return nil, err
	}

	return NewMysqlResultset(rows), nil
}

func NewMysqlResultset(rows *sql.Rows) domain.Resultset {

	result := &MySQLResult{}
	for cont := true; cont; cont = rows.NextResultSet() {
		resultSet := MySQLResultSet{}
		colTypes, err := rows.ColumnTypes()
		if err != nil {
			resultSet.Error = err.Error()
			continue
		}

		// Save the name of the type
		resultSet.ColTypes = make([]string, len(colTypes))
		for i := range colTypes {
			resultSet.ColTypes[i] = colTypes[i].DatabaseTypeName()
		}

		resultSet.ColNames, err = rows.Columns()
		if err != nil {
			resultSet.Error = err.Error()
			continue
		}

		row := make([]any, len(colTypes))

		for i := range colTypes {
			switch colTypes[i].DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				row[i] = new(sql.NullString)
			case "BOOL":
				row[i] = new(sql.NullBool)
			case "INT4":
				row[i] = new(sql.NullInt64)
			case "BLOB", "LONGBLOB":
				row[i] = new(sql.RawBytes)
			default:
				row[i] = new(sql.NullString)
			}
		}

		for rows.Next() {
			err = rows.Scan(row...)

			if err != nil {
				resultSet.Error = err.Error()
				break
			}

			resultRow := make([]*string, len(row))

			for i := range row {
				null, val := sqlValueToString(row[i])

				if null {
					resultRow[i] = nil
				} else {
					resultRow[i] = &val
				}
			}

			resultSet.Result = append(resultSet.Result, resultRow)
		}
		result.Results = append(result.Results, resultSet)
	}

	return result
}

func sqlValueToString(val interface{}) (bool, string) {

	if z, ok := (val).(*sql.NullBool); ok {
		if !z.Valid {
			return !z.Valid, "false"
		}
		return !z.Valid, strconv.FormatBool(z.Bool)
	}

	if z, ok := (val).(*sql.NullString); ok {
		return !z.Valid, z.String
	}

	if z, ok := (val).(*sql.NullFloat64); ok {
		return !z.Valid, strconv.FormatFloat(z.Float64, 'f', 10, 32)
	}

	// Integer values
	if z, ok := (val).(*sql.NullInt32); ok {
		return !z.Valid, strconv.Itoa(int(z.Int32))
	}

	if z, ok := (val).(*sql.NullInt64); ok {
		return !z.Valid, strconv.Itoa(int(z.Int64))
	}

	if z, ok := (val).(*sql.NullInt16); ok {
		return !z.Valid, strconv.Itoa(int(z.Int16))
	}

	if z, ok := (val).(*sql.NullTime); ok {
		return !z.Valid, string(z.Time.Format(time.UnixDate))
	}

	if z, ok := (val).(*sql.NullByte); ok {
		return !z.Valid, strconv.Itoa(int(z.Byte))
	}

	if z, ok := (val).(*sql.RawBytes); ok {
		return false, base64.StdEncoding.EncodeToString([]byte(*z))
	}

	return true, ""
}
