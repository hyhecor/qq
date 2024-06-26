package qq

import (
	"github.com/hyhecor/qq/opop/models"
	"github.com/hyhecor/qq/opop/render"
)

const (
	DialectSQLServer models.Dialect = "sqlserver"
	DialectSQLite3   models.Dialect = "sqlite3" /* implemented yet */
	DialectMySQL     models.Dialect = "mysql"   /* implemented yet */
)

func Select(ops ...models.StatementOp) models.Statement {
	stmt := models.Statement{
		Mode: models.ModeSelect,
	}

	for _, op := range ops {
		op(&stmt)
	}
	return stmt
}

func Dialect(dialect models.Dialect) models.StatementOp {
	return func(s *models.Statement) {

		switch dialect {
		case DialectSQLServer:
			s.Painter = new(render.SQLServerPainter)
		default:
			panic("unsupported dialect")
		}

		s.Dialect = dialect
	}
}

func From(table models.Table) models.StatementOp {
	return func(s *models.Statement) {
		s.From = &table
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

func Columns(columns ...models.ColumnIndentifier) models.StatementOp {
	return func(s *models.Statement) {
		s.Columns = columns
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

func Literal(expression string) models.Literal {
	literal := models.Literal{
		Expression: expression,
	}

	return literal
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

func Where(conditions ...models.ConditionIndentifier) models.StatementOp {
	return func(s *models.Statement) {
		s.Conditions = And(conditions...)
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
		Op:     models.Equal,
		Column: column,
		Param:  []string{param},
	}
}

func GreaterThan(column models.Column, param string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     models.GreaterThan,
		Column: column,
		Param:  []string{param},
	}
}

func LessThan(column models.Column, param string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     models.LessThan,
		Column: column,
		Param:  []string{param},
	}
}

func GreaterThanOrEqual(column models.Column, param string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     models.GreaterThanOrEqual,
		Column: column,
		Param:  []string{param},
	}
}

func LessThanOrEqual(column models.Column, param string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     models.LessThanOrEqual,
		Column: column,
		Param:  []string{param},
	}
}

func Between(column models.Column, from string, to string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     models.Between,
		Column: column,
		Param:  []string{from, to},
	}
}

func In(column models.Column, param ...string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     models.In,
		Column: column,
		Param:  param,
	}
}

func IsNull(column models.Column, param ...string) models.ComparisonExpression {
	return models.ComparisonExpression{
		Op:     models.IsNull,
		Column: column,
		Param:  param,
	}
}

func Not(value models.ConditionIndentifier) models.LogicalExpression {
	return models.LogicalExpression{
		Op:          models.Not,
		Predicators: []models.ConditionIndentifier{value},
	}
}

func And(condition ...models.ConditionIndentifier) models.LogicalExpression {
	return models.LogicalExpression{
		Op:          models.And,
		Predicators: condition,
	}
}

func Or(condition ...models.ConditionIndentifier) models.LogicalExpression {
	return models.LogicalExpression{
		Op:          models.Or,
		Predicators: condition,
	}
}

func Orders(orders ...models.OrderByExpression) models.StatementOp {
	return func(s *models.Statement) {
		s.Orders = orders
	}
}

func Asc(column models.Column) models.OrderByExpression {
	return models.OrderByExpression{
		Op:     models.Ascending,
		Column: column,
	}
}

func Desc(column models.Column) models.OrderByExpression {
	return models.OrderByExpression{
		Op:     models.Desending,
		Column: column,
	}
}

func Insert(ops ...models.StatementOp) models.Statement {
	stmt := models.Statement{
		Mode: models.ModeInsert,
	}

	for _, op := range ops {
		op(&stmt)
	}
	return stmt
}

func Into(table models.Table) models.StatementOp {
	return func(s *models.Statement) {
		s.Into = &table
	}
}

func Values(columns ...string) models.StatementOp {
	return func(s *models.Statement) {
		s.Values = columns
	}
}
