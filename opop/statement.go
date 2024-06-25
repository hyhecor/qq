package qq

import (
	"bytes"

	"github.com/hyhecor/qq/opop/models"
)

const (
	DialectSQLServer models.Dialect = "sqlserver"
	DialectSQLite3   models.Dialect = "sqlite3" /* implemented yet */
	DialectMySQL     models.Dialect = "mysql"   /* implemented yet */
)

type Statement struct {
	Mode       Mode
	Dialect    models.Dialect
	Painter    Painter
	Table      *models.Table
	Columns    []models.Column
	Predicator models.Predicator
	Orders     []models.OrderByExpression
}

func (stmt Statement) Render() string {
	var buf bytes.Buffer

	switch stmt.Mode {
	case ModeSelect:
		stmt.Painter.Select(&buf, stmt)
	}

	return buf.String()
}

type Mode int

const (
	ModeSelect Mode = iota
)

type StatementOp = func(*Statement)

func Select(ops ...StatementOp) Statement {
	stmt := Statement{
		Mode: ModeSelect,
	}

	for _, op := range ops {
		op(&stmt)
	}
	return stmt
}

func Dialect(dialect models.Dialect) StatementOp {
	return func(s *Statement) {

		switch dialect {
		case DialectSQLServer:
			s.Painter = new(SQLServerRender)
		default:
			panic("unsupported dialect")
		}

		s.Dialect = dialect
	}
}

func From(table models.Table) StatementOp {
	return func(s *Statement) {
		s.Table = &table
	}
}

func Table(name string, ops ...TableOp) models.Table {

	tbl := models.Table{
		Name: name,
	}

	for _, op := range ops {
		op(&tbl)
	}
	return tbl
}

type TableOp = func(*models.Table)

func TableAlias(alias string) TableOp {
	return func(c *models.Table) {
		c.Alias = &alias
	}
}

func Columns(columns ...models.Column) StatementOp {
	return func(s *Statement) {
		s.Columns = make([]models.Column, 0, len(columns))
		s.Columns = append(s.Columns, columns...)
	}
}

func Column(name string, ops ...ColumnOp) models.Column {
	col := models.Column{
		Name: name,
	}

	for _, fn := range ops {
		fn(&col)
	}
	return col
}

type ColumnOp = func(*models.Column)

func ColumnTable(tb models.Table) ColumnOp {
	return func(c *models.Column) {
		c.Table = &tb
	}
}

func ColunmAlias(alias string) ColumnOp {
	return func(c *models.Column) {
		c.Alias = &alias
	}
}

func Where(conditions ...models.Predicator) StatementOp {
	return func(s *Statement) {
		s.Predicator = And(conditions...)
	}
}

var (
	Eq = Equal
	Gt = GreaterThan
	Lt = LessThan
	Ge = GreaterThanOrEqual
	Le = LessThanOrEqual
)

func Equal(column models.Column, param string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     "equal",
		Column: column,
		Param:  []string{param},
	}
}

func GreaterThan(column models.Column, param string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     "greater_than",
		Column: column,
		Param:  []string{param},
	}
}

func LessThan(column models.Column, param string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     "less_than",
		Column: column,
		Param:  []string{param},
	}
}

func GreaterThanOrEqual(column models.Column, param string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     "greater_than_or_equal",
		Column: column,
		Param:  []string{param},
	}
}

func LessThanOrEqual(column models.Column, param string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     "less_than_or_equal",
		Column: column,
		Param:  []string{param},
	}
}

func Between(column models.Column, from string, to string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     "between",
		Column: column,
		Param:  []string{from, to},
	}
}

func In(column models.Column, param ...string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     "in",
		Column: column,
		Param:  param,
	}
}

func Not(value models.Predicator) models.LogicalExpression {
	return models.LogicalExpression{
		Op:          "not",
		Predicators: []models.Predicator{value},
	}
}

func And(condition ...models.Predicator) models.LogicalExpression {
	return models.LogicalExpression{
		Op:          "and",
		Predicators: condition,
	}
}

func Or(condition ...models.Predicator) models.LogicalExpression {
	return models.LogicalExpression{
		Op:          "or",
		Predicators: condition,
	}
}

func Orders(orders ...models.OrderByExpression) StatementOp {
	return func(s *Statement) {
		s.Orders = orders
	}
}

func Asc(column models.Column) models.OrderByExpression {
	return models.OrderByExpression{
		Op:     "ascending",
		Column: column,
	}
}

func Desc(column models.Column) models.OrderByExpression {
	return models.OrderByExpression{
		Op:     "desending",
		Column: column,
	}
}
