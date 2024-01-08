package main

import (
	"github.com/tvanriel/db-gui/app/queries"
)

func main() {
	sql := queries.Select("information_schema", "COLUMNS",
		[]queries.SelectColumn{
			{Name: "COLUMN_NAME", Func: "count"},
		},
		[]queries.WhereClause{
			{
				ColumnName: "TABLE_NAME",
				Operator:   "=",
				Value:      "cities",
			},
		},
		queries.LimitOffset{
			Limit:  50,
			Offset: 0,
		},
		[]queries.OrderBy{
			{
				Direction: "ASC",
				By:        "Test",
			},
		},
	)

}
