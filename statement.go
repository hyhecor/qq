package qq

import (
	"bytes"
	"database/sql"
)

type table struct {
	Name  string
	Alias sql.NullString
}

type statement struct {
	relationRender RelationRender
	table          *table
	columns        []column
	conditions     Predicator
}

func (stmt statement) Dialect(dialect Dialect) statement {

	switch dialect {
	case DialectSQLServer:
		stmt.relationRender = new(SQLServerRender)
	default:
		panic("unsupported dialect")
	}

	return stmt
}

func Table(name string, alias ...string) *table {
	return &table{
		Name: name,
		Alias: sql.NullString{
			String: first(alias...),
			Valid:  0 < len(first(alias...)),
		},
	}
}

func (stmt statement) From(table *table) statement {
	stmt.table = table
	return stmt
}

func From(table *table) statement {
	return statement{}.From(table)
}

type ColumnOp func(*column)

func WithTable(tb *table) ColumnOp {
	return func(c *column) {
		c.Table = tb
	}
}

func WithAlias(alias string) ColumnOp {
	return func(c *column) {
		c.Alias.String = alias
		c.Alias.Valid = 0 < len(alias)
	}
}

func Column(name string, ops ...ColumnOp) column {
	col := column{
		Name: name,
	}

	for _, fn := range ops {
		fn(&col)
	}

	return col
}

type column struct {
	Name  string
	Table *table
	Alias sql.NullString
}

func (stmt statement) Where(conditions ...Predicator) statement {
	stmt.conditions = And(conditions...)
	return stmt
}

func Where(conditions ...Predicator) statement {
	return statement{}.Where(conditions...)
}

func (stmt statement) Columns(columns ...column) statement {
	stmt.columns = make([]column, 0, len(columns))
	stmt.columns = append(stmt.columns, columns...)
	return stmt
}

func (stmt statement) Select() string {

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
