package queries

import (
	"io"

	"github.com/tvanriel/db-gui/app/domain"
	"vitess.io/vitess/go/vt/sqlparser"
)

var _ domain.QueryOrExecuteSubject = QueryOrStatement{}

type QueryOrStatement struct {
	Sql   string
	Query bool
}

func (q QueryOrStatement) ShouldQuery() bool {
	return q.Query
}
func (q QueryOrStatement) Statement() string {
	return q.Sql
}

func ParseScript(sql string, report func(int)) ([]QueryOrStatement, error) {
	statements := []sqlparser.Statement{}

	t := sqlparser.NewStringTokenizer(sql)
	t.AllowComments = true
	i := 0
	for {
		next, err := sqlparser.ParseNext(t)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		i++
		if i%5 == 0 {
			report(t.Pos)
		}
		statements = append(statements, next)
	}
	report(t.Pos)
	result := make([]QueryOrStatement, i)

	for i := range statements {

		shouldQuery := false
		switch sqlparser.ASTToStatementType(statements[i]) {
		case sqlparser.StmtSelect, sqlparser.StmtExplain, sqlparser.StmtStream, sqlparser.StmtShow:
			shouldQuery = true
		}

		result[i] = QueryOrStatement{
			Sql:   stringifyStatement(statements[i]),
			Query: shouldQuery,
		}
	}

	return result, nil
}

func stringifyStatement(node sqlparser.SQLNode) string {
	if node == nil {
		return ""
	}

	return sqlparser.CanonicalString(node)
}
