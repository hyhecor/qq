package models

import (
	"bytes"
	"io"
)

type Dialect string

type Mode int

const (
	ModeSelect Mode = iota
	ModeInsert Mode = iota
)

type Statement struct {
	Mode    Mode
	Dialect Dialect
	Painter Painter

	// select
	From       *Table
	Columns    []ColumnIndentifier
	Conditions ConditionIndentifier
	Orders     []OrderByExpression

	// insert
	Into          *Table
	InsertColumns []ColumnIndentifier
	Values        []string
}

type StatementOp = func(*Statement)

type Painter interface {
	Select(w io.Writer, stmt Statement)
	Insert(w io.Writer, stmt Statement)
}

func (stmt Statement) Render() string {
	var buf bytes.Buffer

	switch stmt.Mode {
	case ModeSelect:
		stmt.Painter.Select(&buf, stmt)
	case ModeInsert:
		stmt.Painter.Insert(&buf, stmt)
	}

	return buf.String()
}

type Alias *string

type Table struct {
	Name string
	// Alias sql.NullString
	Alias Alias
}

type ColumnIndentifier interface {
	_column()
}

type Column struct {
	Name  string
	Table *Table
	// Alias sql.NullString
	Alias Alias
}

func (Column) _column() {}

type Literal struct {
	Expression string
}

func (Literal) _column() {}

type ConditionIndentifier interface {
	_condition()
}

type ComparisonExpression struct {
	Op     Operator
	Column Column
	Param  []string
}

func (ComparisonExpression) _condition() {}

type LogicalExpression struct {
	Op          Operator
	Predicators []ConditionIndentifier
}

func (LogicalExpression) _condition() {}

type OrderByExpression struct {
	Op     Operator
	Column Column
}

type Operator string

const (
	Equal              Operator = "equal"
	GreaterThan        Operator = "greater_than"
	LessThan           Operator = "less_than"
	GreaterThanOrEqual Operator = "greater_than_or_equal"
	LessThanOrEqual    Operator = "less_than_or_equal"
	Between            Operator = "between"
	In                 Operator = "in"
	IsNull             Operator = "is_null"

	Not       Operator = "not"
	And       Operator = "and"
	Or        Operator = "or"
	Ascending Operator = "ascending"
	Desending Operator = "desending"
)
