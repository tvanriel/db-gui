package mysqldump

import (
	"compress/gzip"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"vitess.io/vitess/go/sqlescape"
)

const header = `-- Go SQL Dump
--
-- ------------------------------------------------------
-- Server version	{{ .ServerVersion }}
-- Server version	{{ .Date }}

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
`

type tmplData struct {
	Date          string
	ServerVersion string
}

type DumpParameters struct {
	DB                  *sql.DB
	Name                string
	TableCreateStrategy string
	DumpCommunicator    DumpCommunicator
}

// The DumpDatabase mnethod takes a set of parameters describing how to reverse engineer
// the database, writes it all into a sql.gz file and returns the filename of the compressed file.
func DumpDatabase(params DumpParameters) error {

	db := params.DB
	gzipFile, err := os.CreateTemp(os.TempDir(), "mysql-dump-*.sql.gz")
	if err != nil {
		return err
	}
	sqlFile, err := os.CreateTemp(os.TempDir(), "mysql-dump-*.sql")
	if err != nil {
		return err
	}

	defer sqlFile.Close()
	defer gzipFile.Close()

	writeHeader(sqlFile, db)

	// Get a list of all tables in the database
	tables, err := getTables(db)
	if err != nil {
		return err
	}

	// Get a list of all views in the database
	views, err := getViews(db)
	if err != nil {
		return err
	}

	params.DumpCommunicator.Start(
		len(tables),
		len(views),
	)

	// Iterate through the tables and generate the CREATE TABLE statements
	for i, table := range tables {
		params.DumpCommunicator.DumpingTable(table, i)
		// Get the CREATE TABLE statement for the table
		createTable, err := getCreateTableStatement(db, table)
		if err != nil {
			return err
		}

		// Append the CREATE TABLE statement to the SQL dump
		sqlFile.WriteString(createTable)
		sqlFile.WriteString("\n\n")

		// Get the data for the table
		data, err := getTableData(db, table)
		if err != nil {
			return err
		}

		// Append the INSERT INTO statements to the SQL dump
		for _, row := range data {
			sqlFile.WriteString(row)
			sqlFile.WriteString("\n")
		}
		sqlFile.WriteString("\n")
	}

	// Iterate through the views and generate the CREATE VIEW statements
	for i, view := range views {
		params.DumpCommunicator.DumpingView(view, i)

		// Get the CREATE VIEW statement for the view
		createView, err := getCreateViewStatement(db, view)
		if err != nil {
			return err
		}

		// Append the CREATE VIEW statement to the SQL dump
		_, err = sqlFile.WriteString(createView)
		if err != nil {
			return err
		}
		_, err = sqlFile.WriteString("\n\n")
		if err != nil {
			return err
		}
	}

	gz, err := gzip.NewWriterLevel(gzipFile, gzip.BestCompression)
	if err != nil {
		return err
	}
	_, err = io.Copy(gz, sqlFile)
	if err != nil {
		return err
	}
	gz.Close()
	params.DumpCommunicator.Complete(gzipFile.Name())

	return nil
}

// writeHeader writes the MySQL dump header to the given io.Writer.
// It retrieves the MySQL server version using a query and the current date using the time package.
// It then uses a template and the tmplData struct to generate the header and write it to the io.Writer.
func writeHeader(f io.Writer, db *sql.DB) error {
	t := time.Now().Format("2006-01-02 15:04:05")

	// Query the database for the MySQL server version
	r := db.QueryRow("SELECT version()")

	// Scan the result of the query into a string variable
	var version string
	err := r.Scan(&version)
	if err != nil {
		return err
	}

	// Create a tmplData struct with the server version and current date
	data := tmplData{
		ServerVersion: version,
		Date:          t,
	}

	// Parse the header template
	templ, err := template.New("mysqldump_header").Parse(header)
	if err != nil {
		return err
	}

	// Execute the template with the tmplData struct and write the result to the io.Writer
	return templ.Execute(f, data)
}

// getTables returns a list of all tables in the database
func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW FULL TABLES WHERE Table_Type = 'BASE TABLE'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		var tableType string // Discard, we know the answer, but not selecting it is a pain
		if err := rows.Scan(&table, &tableType); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

// getCreateTableStatement returns the CREATE TABLE statement for the given table
func getCreateTableStatement(db *sql.DB, table string) (string, error) {
	// Query the database for the CREATE TABLE statement
	row := db.QueryRow(fmt.Sprintf("SHOW CREATE TABLE %s", sqlescape.EscapeID(table)))

	// Scan the result into a string
	var tableName, createTable string
	if err := row.Scan(&tableName, &createTable); err != nil {
		return "", err
	}

	// Return the CREATE TABLE statement
	return createTable, nil
}

// getTableData returns the data for the given table as a slice of INSERT INTO statements
func getTableData(db *sql.DB, table string) ([]string, error) {
	// Query the database for the data in the table
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", sqlescape.EscapeID(table)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get the column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	valuePtrs := make([]interface{}, len(columns))

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	for i := range columnTypes {
		switch columnTypes[i].DatabaseTypeName() {
		case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
			valuePtrs[i] = new(sql.NullString)
		case "BOOL":
			valuePtrs[i] = new(sql.NullBool)
		case "INT4":
			valuePtrs[i] = new(sql.NullInt64)
		case "BLOB", "LONGBLOB":
			valuePtrs[i] = new(sql.RawBytes)
		default:
			valuePtrs[i] = new(sql.NullString)
		}
	}

	// Create a slice to hold the data rows
	var data []string

	// Scan the rows into the data slice
	for rows.Next() {
		// Create a slice to hold the column values for this row

		// Scan the row into the value pointers slice
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Build the INSERT INTO statement for the row
		var row strings.Builder
		row.WriteString(fmt.Sprintf("INSERT INTO %s (", sqlescape.EscapeID(table)))
		row.WriteString(strings.Join(columns, ", "))
		row.WriteString(") VALUES (")
		for i, value := range valuePtrs {
			// Convert the value to a string
			str, err := valueToString(value)
			if err != nil {
				return nil, err
			}

			// Append the string value to the INSERT INTO statement
			row.WriteString(str)
			if i < len(valuePtrs)-1 {
				row.WriteString(", ")
			}
		}
		row.WriteString(");")

		// Append the INSERT INTO statement to the data slice
		data = append(data, row.String())
	}
	return data, nil
}

func addQuotes(s string) string {
	return strings.Join([]string{
		"\"",
		strings.ReplaceAll(s, "'", "''"),
		"\"",
	}, "")
}

// valueToString converts a value to a string representation suitable for use in an SQL statement
func valueToString(val interface{}) (string, error) {

	if z, ok := (val).(*sql.NullBool); ok {
		if !z.Valid {
			return "NULL", nil
		}
		return addQuotes(strconv.FormatBool(z.Bool)), nil
	}

	if z, ok := (val).(*sql.NullString); ok {
		return addQuotes(z.String), nil
	}

	if z, ok := (val).(*sql.NullFloat64); ok {
		if !z.Valid {
			return "NULL", nil
		}
		return strconv.FormatFloat(z.Float64, 'f', 10, 32), nil
	}

	// Integer values
	if z, ok := (val).(*sql.NullInt32); ok {
		if !z.Valid {
			return "NULL", nil
		}
		return strconv.Itoa(int(z.Int32)), nil
	}
	if z, ok := (val).(*sql.NullInt64); ok {
		if !z.Valid {
			return "NULL", nil
		}
		return strconv.Itoa(int(z.Int64)), nil
	}
	if z, ok := (val).(*sql.NullInt16); ok {
		if !z.Valid {
			return "NULL", nil
		}
		return strconv.Itoa(int(z.Int16)), nil
	}

	if z, ok := (val).(*sql.NullTime); ok {
		if !z.Valid {
			return "NULL", nil
		}
		return addQuotes(string(z.Time.Format("2006-01-02 15:04:05"))), nil
	}
	if z, ok := (val).(*sql.NullByte); ok {
		if !z.Valid {
			return "NULL", nil
		}
		return strconv.Itoa(int(z.Byte)), nil
	}
	if z, ok := (val).(*sql.RawBytes); ok {
		return addQuotes(base64.StdEncoding.EncodeToString([]byte(*z))), nil
	}
	if z, ok := (val).([]uint8); ok {
		return addQuotes(base64.StdEncoding.EncodeToString([]byte(z))), nil
	}
	if val == nil {
		return "NULL", nil
	}
	return "", fmt.Errorf("unable to format: %T", val)
}

// getCreateView Statement returns the CREATE VIEW statement for the given view
func getCreateViewStatement(db *sql.DB, view string) (string, error) {
	// Query the database for the CREATE VIEW statement
	row := db.QueryRow(fmt.Sprintf("SHOW CREATE VIEW %s", sqlescape.EscapeID(view)))

	// Scan the result into a string
	viewName := ""
	createView := ""
	characterSetClient := ""
	collationConnection := ""

	if err := row.Scan(&viewName, &createView, &characterSetClient, &collationConnection); err != nil {
		return "", err
	}

	// Return the CREATE VIEW statement
	return createView, nil
}

// getViews returns a list of all views in the database
func getViews(db *sql.DB) ([]string, error) {
	// Query the database for the list of views
	rows, err := db.Query("SHOW FULL TABLES WHERE Table_Type = 'VIEW'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a slice to hold the view names
	views := []string{}

	// Scan the rows and append the view names to the slice
	for rows.Next() {
		tableName := ""
		tableType := ""

		if err := rows.Scan(&tableName, &tableType); err != nil {
			return nil, err
		}

		views = append(views, tableName)
	}
	return views, nil
}
