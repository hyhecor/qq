package models

type Dialect string

type Alias *string

type Table struct {
	Name string
	// Alias sql.NullString
	Alias Alias
}

type Column struct {
	Name  string
	Table *Table
	// Alias sql.NullString
	Alias Alias
}

type Predicator interface {
	Predicate()
}

type ComparisonExpression struct {
	Op     string
	Column Column
	Param  []string
}

func (ComparisonExpression) Predicate() {}

type LogicalExpression struct {
	Op          string
	Predicators []Predicator
}

func (LogicalExpression) Predicate() {}

type OrderByExpression struct {
	Op     string
	Column Column
}
