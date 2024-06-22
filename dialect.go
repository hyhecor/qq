package qq

const (
	DialectSQLServer Dialect = "sqlserver"
	DialectSQLite3   Dialect = "sqlite3" /* implemented yet */
	DialectMySQL     Dialect = "mysql"   /* implemented yet */
)

type Dialect string
