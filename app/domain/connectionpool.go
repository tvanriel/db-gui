package domain

import (
	"database/sql"
)

type Column interface {
	Field() string
	Type() string
	Null() bool
	Default() string
	Comment() string
}

type Index interface {
	Table() string
	NonUnique() bool
	KeyName() string
	SeqInIndex() int
	ColumnName() string
	Collation() string
	Cardinality() string
	SubPart() string
	Packed() string
	Null() bool
	IndexType() string
	Comment() string
	IndexComment() string
	Visible() bool
	Expression() string
}

type Table interface {
	Columns() []string
	GetColumn(string) Column
	Name() string
	IndexSequence() []int
	GetIndex(int) Index
}

type Resultset interface {
	Get() []struct {
		ColNames []string
		ColTypes []string
		Result   [][]*string
		Error    string
	}
}

type Queriable interface {
	Query(stmt string) (Resultset, error)
	Exec(stmt string) (sql.Result, error)
	Close() error
}

type Database interface {
	GetTableNames() ([]string, error)
	DescribeTable(string) (Table, error)
	Name() string
	Connect() error
	Close() error
	Dump(strategy string) (string, error)
	Select(Select) (Resultset, string, error)
	AsQueriable() Queriable
}

type Connection interface {
	Ping() (err error)
	GetDatabases() (databases []string, err error)
	GetDatabase(databaseName string) Database
	Connect() (err error)
	Close() (err error)
	AsQueriable() Queriable
}

type ConnectionConfig struct {
	Script string
}

func NewConnectionConfig(Script string) *ConnectionConfig {
	return &ConnectionConfig{
		Script: Script,
	}
}

type ConnectionPool interface {
	NewToken() (token string)
	Get(token string) Connection
	Add(*ConnectionConfig, string) error
	Delete(token string)
}

type ScriptExecutionCommunicator interface {
	ParsingProgress(i int)
	ExecutionProgress(executed, total int)

	StatementPlan(sql string)
	QueryResult(sql string, results Resultset, Time int64)
	ExecResult(rowsAffected, lastInsertId int64, sql string, Time int64)

	ShouldContinueFetching(fetched int) bool

	StatementError(err error, sql string)
	GenericError(err error)

	Close() error
}

type QueryOrExecuteSubject interface {
	ShouldQuery() bool
	Statement() string
}

type StatementList struct {
	Items []QueryOrExecuteSubject
}
