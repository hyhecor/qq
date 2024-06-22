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

type ComparisonOp struct {
	Op     string
	Column column
	Param  []string
}

func (ComparisonOp) Predicate() {}

func Equal(column column, param string) ComparisonOp {
	return ComparisonOp{
		Op:     "equal",
		Column: column,
		Param:  []string{param},
	}
}

func GreaterThan(column column, param string) ComparisonOp {
	return ComparisonOp{
		Op:     "greater_than",
		Column: column,
		Param:  []string{param},
	}
}

func LessThan(column column, param string) ComparisonOp {
	return ComparisonOp{
		Op:     "less_than",
		Column: column,
		Param:  []string{param},
	}
}

func GreaterThanOrEqual(column column, param string) ComparisonOp {
	return ComparisonOp{
		Op:     "greater_than_or_equal",
		Column: column,
		Param:  []string{param},
	}
}

func LessThanOrEqual(column column, param string) ComparisonOp {
	return ComparisonOp{
		Op:     "less_than_or_equal",
		Column: column,
		Param:  []string{param},
	}
}

func Between(column column, from string, to string) ComparisonOp {
	return ComparisonOp{
		Op:     "between",
		Column: column,
		Param:  []string{from, to},
	}
}

func In(column column, param ...string) ComparisonOp {
	return ComparisonOp{
		Op:     "in",
		Column: column,
		Param:  param,
	}
}

type LogicalOp struct {
	Op          string
	Predicators []Predicator
}

func (LogicalOp) Predicate() {}

func Not(value Predicator) LogicalOp {
	return LogicalOp{
		Op:          "not",
		Predicators: []Predicator{value},
	}
}

func And(condition ...Predicator) LogicalOp {
	return LogicalOp{
		Op:          "and",
		Predicators: condition,
	}
}

func Or(condition ...Predicator) LogicalOp {
	return LogicalOp{
		Op:          "or",
		Predicators: condition,
	}
}
