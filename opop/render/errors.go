package render

import (
	"fmt"

	"github.com/hyhecor/qq/opop/models"
)

type ErrEmpty struct {
	Op   models.Operator
	Name string
}

func (err ErrEmpty) Error() string {
	if len(err.Name) == 0 {
		return fmt.Sprintf("%s: param is emtpy", err.Op)
	}
	return fmt.Sprintf("%s: %q param is emtpy", err.Name, err.Op)
}

func MustParamCountOfComparisonExpression(cond models.ComparisonExpression, expected int) {
	if len(cond.Param) == expected {
		return
	}

	panic(error(ErrParamCountNotMatched{Op: cond.Op, Name: cond.Column.Name, ParamCount: expected}))
}

func MustNotParamCountOfComparisonExpression(cond models.ComparisonExpression, expected int) {
	if len(cond.Param) != expected {
		return
	}

	panic(error(ErrParamCountNotMatched{Op: cond.Op, Name: cond.Column.Name, ParamCount: expected}))
}

type ErrParamCountNotMatched struct {
	Op         models.Operator
	Name       string
	ParamCount int
}

func (err ErrParamCountNotMatched) Error() string {
	if len(err.Name) == 0 {
		return fmt.Sprintf("%s: expected param count is %d", err.Op, err.ParamCount)
	}
	return fmt.Sprintf("%s: %q expected param count is %d", err.Name, err.Op, err.ParamCount)
}
