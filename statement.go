package qq

import (
	"bytes"
	"database/sql"
)

type TableEntity struct {
	Name  string
	Alias sql.NullString
}

func Table(name string, alias ...string) *TableEntity {
	return &TableEntity{
		Name: name,
		Alias: sql.NullString{
			String: first(alias...),
			Valid:  0 < len(first(alias...)),
		},
	}
}

type Statement struct {
	dialect        dialect
	relationRender Render
	table          *TableEntity
	columns        []ColumnEntity
	predicator     Predicator
	orders         []OrderByExpression
}

func From(table *TableEntity) Statement {
	return Statement{}.From(table)
}

func (stmt Statement) Dialect(dialect dialect) Statement {

	switch dialect {
	case DialectSQLServer:
		stmt.relationRender = new(SQLServerRender)
	default:
		panic("unsupported dialect")
	}

	stmt.dialect = dialect

	return stmt
}

func Dialect(dialect dialect) Statement {
	return Statement{}.Dialect(dialect)
}

func (stmt Statement) From(table *TableEntity) Statement {
	stmt.table = table
	return stmt
}

func WithTable(tb *TableEntity) ColumnOp {
	return func(c *ColumnEntity) {
		c.Table = tb
	}
}

func WithAlias(alias string) ColumnOp {
	return func(c *ColumnEntity) {
		c.Alias.String = alias
		c.Alias.Valid = 0 < len(alias)
	}
}

type ColumnEntity struct {
	Name  string
	Table *TableEntity
	Alias sql.NullString
}

type ColumnOp = func(*ColumnEntity)

func Column(name string, ops ...ColumnOp) ColumnEntity {
	col := ColumnEntity{
		Name: name,
	}

	for _, fn := range ops {
		fn(&col)
	}

	return col
}

func Where(conditions ...Predicator) Statement {
	return Statement{}.Where(conditions...)
}

func (stmt Statement) Where(conditions ...Predicator) Statement {
	stmt.predicator = And(conditions...)
	return stmt
}

func (stmt Statement) Columns(columns ...ColumnEntity) Statement {
	stmt.columns = make([]ColumnEntity, 0, len(columns))
	stmt.columns = append(stmt.columns, columns...)
	return stmt
}

func (stmt Statement) Select() string {

	var buf bytes.Buffer
	stmt.relationRender.Select(&buf, stmt)

	return buf.String()
}

func first[T any](aa ...T) (b T) {
	for _, v := range aa {
		return v
	}
	return
}

type OrderByExpression struct {
	Op     string
	Column ColumnEntity
}

func (stmt Statement) OrderBy(orders ...OrderByExpression) Statement {
	stmt.orders = orders

	return stmt
}

func Asc(column ColumnEntity) OrderByExpression {
	return OrderByExpression{
		Op:     "ascending",
		Column: column,
	}
}

func Desc(column ColumnEntity) OrderByExpression {
	return OrderByExpression{
		Op:     "desending",
		Column: column,
	}
}
