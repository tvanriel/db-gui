package domain

type Columns interface {
	Get() []struct {
		Func string
		Name string
	}
}

type Order interface {
	Get() []struct {
		Column    string
		Direction string
	}
}

type WhereClause interface {
	Get() []struct {
		Column   string
		Operator string
		Value    string
	}
}

type Limit interface {
	Offset() int
	Limit() int
}
type Select interface {
	DatabaseName() string
	TableName() string
	Columns() Columns
	OrderBy() Order
	Where() WhereClause

	Limit() Limit
}
