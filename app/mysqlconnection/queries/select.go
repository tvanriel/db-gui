package queries

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/tvanriel/db-gui/app/domain"
	"vitess.io/vitess/go/sqlescape"
	"vitess.io/vitess/go/vt/sqlparser"
)

func escapeString(value string) string {
	str, _ := json.Marshal(value)
	return string(str)
}

func makeWhereClause(w domain.WhereClause) string {
	where := w.Get()
	if len(where) == 0 {
		return ""
	}

	stringClauses := make([]string, len(where))

	for i := range where {
		setOrEmpty := func(predicate, val string) string {
			if predicate == "" {
				return ""
			}
			return val
		}
		opAndVal := ""

		switch where[i].Operator {

		case "SQL":
			opAndVal = where[i].Value
		case "IN":
			opAndVal = strings.Join([]string{
				"IN (",
				where[i].Value,
				")",
			}, "")
		case "NOT IN":
			opAndVal = strings.Join([]string{
				"NOT IN (",
				where[i].Value,
				")",
			}, "")
		case "LIKE %%":
			opAndVal = strings.Join([]string{
				"LIKE '%",
				where[i].Value,
				"%'",
			}, "")
		case "NOT LIKE %%":
			opAndVal = strings.Join([]string{
				"NOT LIKE %",
				where[i].Value,
				"%",
			}, "")
		default:
			opAndVal = strings.Join([]string{
				where[i].Operator,
				setOrEmpty(where[i].Value, " "),
				escapeString(where[i].Value),
			}, "")
		}
		stringClauses[i] = strings.Join([]string{
			sqlescape.EscapeID(where[i].Column),
			" ",
			opAndVal,
		}, "")
	}

	return strings.Join([]string{
		" WHERE ",
		strings.Join(stringClauses, " AND "),
	}, "")

}

func makeCols(c domain.Columns) string {
	col := c.Get()
	if len(col) == 0 {
		return "*"
	}

	s := make([]string, len(col))

	for i := range col {
		switch col[i].Func {
		case "":
			s[i] = sqlescape.EscapeID(col[i].Name)
		case "count distinct":
			s[i] = strings.Join([]string{
				"COUNT(DISTINCT ",
				sqlescape.EscapeID(col[i].Name),
				")",
			}, "")
		case "distinct":
			s[i] = strings.Join([]string{
				"DISTINCT ",
				sqlescape.EscapeID(col[i].Name),
			}, "")
		default:
			s[i] = strings.Join([]string{col[i].Func, "(", sqlescape.EscapeID(col[i].Name), ")"}, "")
		}

	}
	return strings.Join(s, ", ")
}

func makeSource(database string, table string) string {
	return strings.Join([]string{sqlescape.EscapeID(database), ".", sqlescape.EscapeID(table)}, "")
}

type LimitOffset struct {
	Limit  int
	Offset int
}

func makeLimitClause(limit domain.Limit) string {
	return strings.Join([]string{
		" LIMIT",
		strconv.Itoa(limit.Limit()),
		"OFFSET",
		strconv.Itoa(limit.Offset()),
	}, " ")
}

func makeOrderClause(o domain.Order) string {
	order := o.Get()
	if len(order) == 0 {
		return ""
	}
	orderClauses := make([]string, len(order))

	for i := range order {
		dir := "ASC"
		if order[i].Direction == "DESC" {
			dir = "DESC"
		}

		orderClauses[i] = strings.Join([]string{
			sqlescape.EscapeID(order[i].Column),
			dir,
		}, "")
	}

	return strings.Join([]string{" ORDER BY", strings.Join(orderClauses, ", ")}, " ")
}

func Select(s domain.Select) string {
	whereClause := makeWhereClause(s.Where())
	cols := makeCols(s.Columns())
	source := makeSource(s.DatabaseName(), s.TableName())
	limitClause := makeLimitClause(s.Limit())
	orderClause := makeOrderClause(s.OrderBy())

	query := strings.Join([]string{
		"SELECT ",
		cols,
		" FROM ",
		source,
		whereClause,
		orderClause,
		limitClause,
	}, "")

	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return ""
	}

	return stringifyStatement(stmt)
}
