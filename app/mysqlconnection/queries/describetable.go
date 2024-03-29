package queries

func DescribeTable() string {
	return `SELECT
	COLUMNS.COLUMN_NAME,
	COLUMNS.COLUMN_TYPE,
	COLUMNS.IS_NULLABLE,
	COLUMNS.COLUMN_COMMENT,
	COLUMNS.COLUMN_DEFAULT
FROM information_schema.COLUMNS
WHERE
	COLUMNS.TABLE_NAME = ?
	AND COLUMNS.TABLE_SCHEMA = ?`
}
