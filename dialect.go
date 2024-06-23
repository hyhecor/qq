package qq

const (
	DialectSQLServer dialect = "sqlserver"
	DialectSQLite3   dialect = "sqlite3" /* implemented yet */
	DialectMySQL     dialect = "mysql"   /* implemented yet */
)

type dialect string
