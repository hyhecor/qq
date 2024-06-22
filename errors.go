package qq

import "fmt"

type ErrEmpty struct {
	Op   string
	Name string
}

func (err ErrEmpty) Error() string {
	if len(err.Name) == 0 {
		return fmt.Sprintf("%s: param is emtpy", err.Op)
	}
	return fmt.Sprintf("%s: %q param is emtpy", err.Name, err.Op)
}

type ErrParamCountNotMatched struct {
	Op         string
	Name       string
	ParamCount int
}

func (err ErrParamCountNotMatched) Error() string {
	if len(err.Name) == 0 {
		return fmt.Sprintf("%s: expected param count is %d", err.Op, err.ParamCount)
	}
	return fmt.Sprintf("%s: %q expected param count is %d", err.Name, err.Op, err.ParamCount)
}
