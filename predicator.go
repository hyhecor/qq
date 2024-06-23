package qq

var (
	Eq = Equal
	Gt = GreaterThan
	Lt = LessThan
	Ge = GreaterThanOrEqual
	Le = LessThanOrEqual
)

type Predicator interface {
	Predicate()
}

type ComparisonExpression struct {
	Op     string
	Column ColumnEntity
	Param  []string
}

func (ComparisonExpression) Predicate() {}

func Equal(column ColumnEntity, param string) ComparisonExpression {
	return ComparisonExpression{
		Op:     "equal",
		Column: column,
		Param:  []string{param},
	}
}

func GreaterThan(column ColumnEntity, param string) ComparisonExpression {
	return ComparisonExpression{
		Op:     "greater_than",
		Column: column,
		Param:  []string{param},
	}
}

func LessThan(column ColumnEntity, param string) ComparisonExpression {
	return ComparisonExpression{
		Op:     "less_than",
		Column: column,
		Param:  []string{param},
	}
}

func GreaterThanOrEqual(column ColumnEntity, param string) ComparisonExpression {
	return ComparisonExpression{
		Op:     "greater_than_or_equal",
		Column: column,
		Param:  []string{param},
	}
}

func LessThanOrEqual(column ColumnEntity, param string) ComparisonExpression {
	return ComparisonExpression{
		Op:     "less_than_or_equal",
		Column: column,
		Param:  []string{param},
	}
}

func Between(column ColumnEntity, from string, to string) ComparisonExpression {
	return ComparisonExpression{
		Op:     "between",
		Column: column,
		Param:  []string{from, to},
	}
}

func In(column ColumnEntity, param ...string) ComparisonExpression {
	return ComparisonExpression{
		Op:     "in",
		Column: column,
		Param:  param,
	}
}

type LogicalExpression struct {
	Op          string
	Predicators []Predicator
}

func (LogicalExpression) Predicate() {}

func Not(value Predicator) LogicalExpression {
	return LogicalExpression{
		Op:          "not",
		Predicators: []Predicator{value},
	}
}

func And(condition ...Predicator) LogicalExpression {
	return LogicalExpression{
		Op:          "and",
		Predicators: condition,
	}
}

func Or(condition ...Predicator) LogicalExpression {
	return LogicalExpression{
		Op:          "or",
		Predicators: condition,
	}
}
