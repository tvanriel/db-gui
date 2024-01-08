package queries

func GetSchemas() string {
	return `SELECT
	SCHEMA_NAME
FROM information_schema.SCHEMATA
ORDER BY SCHEMA_NAME
`
}
