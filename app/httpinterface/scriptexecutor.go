package httpinterface

import (
	"time"

	"github.com/tvanriel/db-gui/app/domain"
	"github.com/tvanriel/db-gui/app/mysqlconnection/queries"
)

func executeSqlScript(
	statements []queries.QueryOrStatement,
	communicator domain.ScriptExecutionCommunicator,
	q domain.Queriable,
) {
	amountOfStatements := len(statements)

	communicator.ExecutionProgress(0, amountOfStatements)

	for i := range statements {

		sql := statements[i].Sql

		communicator.StatementPlan(sql)

		if statements[i].Query {
			start := time.Now()
			resultSet, err := q.Query(
				sql,
			)
			end := time.Now()

			if err != nil {
				communicator.StatementError(
					err,
					sql,
				)
				goto next
			}

			communicator.QueryResult(
				sql,
				resultSet,
				end.UnixMicro()-start.UnixMicro(),
			)

		} else {
			start := time.Now()
			sqlResult, err := q.Exec(sql)
			end := time.Now()

			if err != nil {
				communicator.StatementError(
					err,
					sql,
				)
				goto next
			}

			rowsAffected, err := sqlResult.RowsAffected()

			if err != nil {
				communicator.StatementError(
					err,
					sql,
				)
				goto next
			}

			lastInsertId, err := sqlResult.LastInsertId()

			if err != nil {
				communicator.StatementError(
					err,
					sql,
				)
				goto next
			}

			communicator.ExecResult(
				rowsAffected,
				lastInsertId,
				sql,
				end.UnixMicro()-start.UnixMicro(),
			)

		}
	next:
		communicator.ExecutionProgress(i+1, amountOfStatements)

	}
}
