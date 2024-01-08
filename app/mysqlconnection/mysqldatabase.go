package mysqlconnection

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/tvanriel/db-gui/app/domain"
	"github.com/tvanriel/db-gui/app/mysqlconnection/mysqldump"
	"github.com/tvanriel/db-gui/app/mysqlconnection/queries"
)

var _ domain.Database = &MySQLDatabase{}

type MySQLDatabase struct {
	Host     string
	Username string
	Password string
	DBName   string

	db *sql.DB
}

func (db *MySQLDatabase) getConfig() *mysql.Config {
	config := mysql.NewConfig()

	config.User = db.Username
	config.Passwd = db.Password
	config.Net = "tcp"
	config.Addr = db.Host
	config.DBName = db.DBName

	config.AllowCleartextPasswords = true
	config.MultiStatements = true

	return config
}

func (db *MySQLDatabase) Connect() error {
	conn, err := sql.Open("mysql", db.getConfig().FormatDSN())
	db.db = conn

	return err
}

func (db *MySQLDatabase) Close() error {
	if db.db == nil {
		return nil
	}

	return db.db.Close()
}

func (db *MySQLDatabase) Query(statement string) (domain.Resultset, error) {
	return NewMysqlResultsetWithError(db.db.Query(statement))
}

func (c *MySQLDatabase) Exec(statement string) (sql.Result, error) {
	return c.db.Exec(statement)
}

func (db *MySQLDatabase) GetTableNames() ([]string, error) {
	result, err := db.Query(queries.GetTableNames())

	if err != nil {
		return nil, err
	}

	resultSet := result.Get()
	if len(resultSet) != 1 {
		return nil, errors.New("Unexpected number of results when listing schemas")
	}

	res := []string{}
	for i := range resultSet[0].Result {
		maybeName := resultSet[0].Result[i][0]
		if maybeName == nil {
			continue
		}
		res = append(res, *maybeName)
	}
	return res, nil
}

func (db *MySQLDatabase) DescribeTable(name string) (domain.Table, error) {
	stmt, err := db.db.Prepare(queries.DescribeTable())

	if err != nil {
		return nil, err
	}

	describeResultSet, err := stmt.Query(name, db.Name())

	if err != nil {
		return nil, err
	}

	columns := []*MySQLColumn{}

	for describeResultSet.Next() {
		var field, colType, comment, null, def sql.NullString

		err = describeResultSet.Scan(&field, &colType, &null, &comment, &def)
		if err != nil {
			return nil, err
		}
		columns = append(columns, &MySQLColumn{
			field:      field.String,
			colType:    colType.String,
			null:       null.String,
			comment:    comment.String,
			colDefault: def.String,
		})
	}
	err = stmt.Close()
	if err != nil {
		return nil, err
	}

	indeces := []*MySQLIndex{}

	return &MySQLTable{
		name:    name,
		columns: columns,
		indeces: indeces,
	}, nil
}

func (db *MySQLDatabase) Name() string {
	return db.DBName
}

func (db *MySQLDatabase) Connection() *sql.DB {
	return db.db
}

func (db *MySQLDatabase) AsQueriable() domain.Queriable {
	return db
}

func (db *MySQLDatabase) Dump(strategy string) (string, error) {
	com := &mysqldump.SummaryCommunicator{}
	err := mysqldump.DumpDatabase(mysqldump.DumpParameters{
		DB:                  db.db,
		Name:                db.DBName,
		TableCreateStrategy: strategy,
		DumpCommunicator:    com,
	})
	if err != nil {
		return "", err
	}
	return com.GetDumpFilename(), nil
}

func (db *MySQLDatabase) Select(s domain.Select) (domain.Resultset, string, error) {
	sql := queries.Select(s)
	result, err := db.Query(sql)
	return result, sql, err
}
