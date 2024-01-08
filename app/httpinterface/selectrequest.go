package httpinterface

import (
	"github.com/tvanriel/db-gui/app/domain"
	"github.com/tvanriel/db-gui/app/formatting"
)

type SelectRequestConfigOrder struct {
	Column    string `json:"column"`
	Direction string `json:"direction"`
}

type SelectRequestConfigColumn struct {
	FieldName string `json:"name"`
	FieldFunc string `json:"func"`
}

type SelectRequestConfigLimitOffset struct {
	FieldLimit  int `json:"limit"`
	FieldOffset int `json:"offset"`
}

type SelectRequestConfigWhereClause struct {
	Column   string `json:"columnName"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}
type SelectRequestConfig struct {
	Config struct {
		FieldOrder   []SelectRequestConfigOrder       `json:"order"`
		FieldColumns []SelectRequestConfigColumn      `json:"columns"`
		FieldLimit   SelectRequestConfigLimitOffset   `json:"limit" binding:"required"`
		FieldWhere   []SelectRequestConfigWhereClause `json:"where"`
	} `json:"config"`

	FieldDatabaseName string `json:"databaseName" binding:"required"`
	FieldTableName    string `json:"tableName" binding:"required"`
}

type SelectResponse struct {
	Sql     string                     `json:"sql"`
	Results []formatting.ResultSetJson `json:"results"`
}

var _ domain.Select = &SelectRequestConfig{}

type RequestColumns []SelectRequestConfigColumn

func (r RequestColumns) Get() []struct {
	Func string
	Name string
} {
	result := make([]struct {
		Func string
		Name string
	}, len(r))
	for i := range r {
		result[i] = struct {
			Func string
			Name string
		}{
			Func: r[i].FieldFunc,
			Name: r[i].FieldName,
		}
	}
	return result
}

func (c *SelectRequestConfig) Columns() domain.Columns {
	return RequestColumns(c.Config.FieldColumns)
}

func (s SelectRequestConfigLimitOffset) Limit() int {
	return s.FieldLimit
}

func (s SelectRequestConfigLimitOffset) Offset() int {
	return s.FieldOffset
}

func (c *SelectRequestConfig) Limit() domain.Limit {
	return c.Config.FieldLimit
}

type RequestWhere []SelectRequestConfigWhereClause

func (r RequestWhere) Get() []struct {
	Column   string
	Operator string
	Value    string
} {
	result := make([]struct {
		Column   string
		Operator string
		Value    string
	}, len(r))
	for i := range r {
		result[i] = struct {
			Column   string
			Operator string
			Value    string
		}{
			Column:   r[i].Column,
			Operator: r[i].Operator,
			Value:    r[i].Value,
		}
	}
	return result
}

func (c *SelectRequestConfig) Where() domain.WhereClause {
	return RequestWhere(c.Config.FieldWhere)
}

type RequestOrder []SelectRequestConfigOrder

func (r RequestOrder) Get() []struct {
	Column    string
	Direction string
} {
	result := make([]struct {
		Column    string
		Direction string
	}, len(r))

	for i := range r {
		result[i] = struct {
			Column    string
			Direction string
		}{
			Column:    r[i].Column,
			Direction: r[i].Direction,
		}
	}

	return result
}

func (c *SelectRequestConfig) OrderBy() domain.Order {
	return RequestOrder(c.Config.FieldOrder)
}

func (c *SelectRequestConfig) DatabaseName() string {
	return c.FieldDatabaseName
}

func (c *SelectRequestConfig) TableName() string {
	return c.FieldTableName
}
