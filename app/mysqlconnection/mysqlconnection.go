package mysqlconnection

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/tvanriel/db-gui/app/domain"
	"github.com/tvanriel/db-gui/app/mysqlconnection/queries"
)

var _ domain.Connection = &MySQLConnection{}

type MySQLConnection struct {
	Host     string
	Username string
	Password string
	db       *sql.DB
}

// Get a new MySQL connection with the current username and password.
func NewMySQLConnection(username, password, host string) domain.Connection {

	conn := &MySQLConnection{
		Host:     host,
		Username: username,
		Password: password,
	}

	return conn
}

// Get a mysql.Config from the password and username that were given in the constructor
func (c *MySQLConnection) getMySQLConfig() *mysql.Config {
	config := mysql.NewConfig()

	// Set the authentication parameters
	config.User = c.Username
	config.Passwd = c.Password
	config.Net = "tcp"
	config.Addr = c.Host

	// Allow this for GCP CloudSQLProxy
	config.AllowCleartextPasswords = true
	config.MultiStatements = true

	return config
}

// Ping the server to see if it is connected.
func (c *MySQLConnection) Ping() error {
	if c.db == nil {
		return errors.New("ping: database is not open")
	}

	return c.db.Ping()
}

// Get the list of databases in this connection.
func (c *MySQLConnection) GetDatabases() ([]string, error) {
	result, err := c.Query(queries.GetSchemas())
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

// Get a new connection config with the database pre-selected.
func (c *MySQLConnection) GetDatabase(name string) domain.Database {
	if name == "" {
		return nil
	}

	return &MySQLDatabase{
		Host:     c.Host,
		Username: c.Username,
		Password: c.Password,
		DBName:   name,
	}
}

// Open a connection to the given database server.
func (c *MySQLConnection) Connect() error {
	db, err := sql.Open("mysql", c.getMySQLConfig().FormatDSN())
	if err != nil {
		return err
	}
	c.db = db
	return nil
}

// Close the connection to the database
func (c *MySQLConnection) Close() error {
	if c.db == nil {
		return nil
	}
	return c.db.Close()
}

// Execute a query on the database and return the result.
func (c *MySQLConnection) Query(statement string) (domain.Resultset, error) {
	return NewMysqlResultsetWithError(c.db.Query(statement))
}

func (c *MySQLConnection) Exec(statement string) (sql.Result, error) {
	return c.db.Exec(statement)
}

func (c *MySQLConnection) Connection() *sql.DB {
	return c.db
}

func (c *MySQLConnection) AsQueriable() domain.Queriable {
	return c
}
