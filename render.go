package qq

import (
	"fmt"
	"io"
)

type RelationRender interface {
	Select(w io.Writer, stmt statement)
	From(w io.Writer, stmt statement)
	Columns(w io.Writer, stmt statement)
	Column(w io.Writer, stmt statement, col column)
	Where(w io.Writer, stmt statement)
	Predicator(w io.Writer, stmt statement, cond Predicator)
	LogicalOp(w io.Writer, stmt statement, cond LogicalOp)
	ComparisonOp(w io.Writer, stmt statement, cond ComparisonOp)
}

type SQLServerRender struct{}

func (render SQLServerRender) Select(w io.Writer, stmt statement) {
	fmt.Fprint(w, "SELECT ")
	render.Columns(w, stmt)
	render.From(w, stmt)
	render.Where(w, stmt)
}

func (render SQLServerRender) From(w io.Writer, stmt statement) {
	fmt.Fprint(w, " FROM ")
	fmt.Fprint(w, stmt.table.Name)

	if stmt.table.Alias.Valid {
		fmt.Fprint(w, " AS ")
		fmt.Fprint(w, stmt.table.Alias.String)
	}
}

func (render SQLServerRender) Columns(w io.Writer, stmt statement) {
	for i := range stmt.columns {
		if 0 < i {
			fmt.Fprint(w, ", ")
		}
		render.Column(w, stmt, stmt.columns[i])
	}
}

func (render SQLServerRender) Column(w io.Writer, stmt statement, col column) {

	switch {

	case col.Table != nil && col.Table.Alias.Valid:
		fmt.Fprint(w, col.Table.Alias.String)
		fmt.Fprint(w, ".")
		fmt.Fprint(w, col.Name)

	case col.Table != nil && !col.Table.Alias.Valid:
		fmt.Fprint(w, col.Table.Name)
		fmt.Fprint(w, ".")
		fmt.Fprint(w, col.Name)

	case col.Table == nil:
		fmt.Fprint(w, col.Name)
	}

	switch {
	case col.Alias.Valid:
		fmt.Fprint(w, " AS ")
		fmt.Fprint(w, col.Alias.String)
	}

}

func (render SQLServerRender) Where(w io.Writer, stmt statement) {
	fmt.Fprint(w, " WHERE ")

	render.Predicator(w, stmt, stmt.conditions)
}

func (render SQLServerRender) Predicator(w io.Writer, stmt statement, predicator Predicator) {

	switch v := predicator.(type) {
	case ComparisonOp:
		render.ComparisonOp(w, stmt, v)

	case LogicalOp:
		render.LogicalOp(w, stmt, v)
	}
}

func (render SQLServerRender) LogicalOp(w io.Writer, stmt statement, cond LogicalOp) {

	switch cond.Op {
	case "and":
		for i := range cond.Predicators {
			if 0 < i {
				fmt.Fprint(w, " AND ")
			}

			render.Predicator(w, stmt, cond.Predicators[i])
		}
	case "or":
		for i := range cond.Predicators {
			if 0 < i {
				fmt.Fprint(w, " OR ")
			}

			render.Predicator(w, stmt, cond.Predicators[i])
		}
	case "not":
		if len(cond.Predicators) == 0 {
			panic(error(ErrEmpty{Op: cond.Op}))
		}

		fmt.Fprint(w, "NOT ")
		render.Predicator(w, stmt, cond.Predicators[0])

	}
}
func (render SQLServerRender) ComparisonOp(w io.Writer, stmt statement, cond ComparisonOp) {
	switch cond.Op {
	case "equal":
		if len(cond.Param) != 1 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 1))
		}

		render.Column(w, stmt, cond.Column)
		fmt.Fprintf(w, " = @%s", cond.Param[0])

	case "greater_than":
		if len(cond.Param) != 1 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 1))
		}

		render.Column(w, stmt, cond.Column)
		fmt.Fprintf(w, " > @%s", cond.Param[0])

	case "less_than":
		if len(cond.Param) != 1 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 1))
		}

		render.Column(w, stmt, cond.Column)
		fmt.Fprintf(w, " < @%s", cond.Param[0])

	case "greater_than_or_equal":
		if len(cond.Param) != 1 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 1))
		}

		render.Column(w, stmt, cond.Column)
		fmt.Fprintf(w, " >= @%s", cond.Param[0])

	case "less_than_or_equal":
		if len(cond.Param) != 1 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 1))
		}

		render.Column(w, stmt, cond.Column)
		fmt.Fprintf(w, " <= @%s", cond.Param[0])

	case "between":
		if len(cond.Param) != 2 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 2))
		}

		render.Column(w, stmt, cond.Column)
		fmt.Fprintf(w, " BETWEEN @%s AND @%s", cond.Param[0], cond.Param[1])

	case "in":
		if len(cond.Param) == 0 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 1))
		}

		render.Column(w, stmt, cond.Column)
		fmt.Fprint(w, " IN (")
		for i := range cond.Param {
			if 0 < i {
				fmt.Fprint(w, ", ")
			}

			fmt.Fprintf(w, "@%s", cond.Param[i])
		}

		fmt.Fprint(w, ")")
	}
}
