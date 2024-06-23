package qq

import (
	"fmt"
	"io"
)

type Render interface {
	Select(w io.Writer, stmt Statement)
}

type SQLServerRender struct{}

func (render SQLServerRender) Select(w io.Writer, stmt Statement) {
	fmt.Fprint(w, "SELECT ")
	render.Columns(w, stmt)
	render.From(w, stmt)
	render.Where(w, stmt)
	render.OrderBy(w, stmt)
}

func (render SQLServerRender) From(w io.Writer, stmt Statement) {
	fmt.Fprint(w, " FROM ")
	fmt.Fprint(w, stmt.table.Name)

	if stmt.table.Alias.Valid {
		fmt.Fprint(w, " AS ")
		fmt.Fprint(w, stmt.table.Alias.String)
	}
}

func (render SQLServerRender) Columns(w io.Writer, stmt Statement) {

	if len(stmt.columns) == 0 {
		stmt.columns = []ColumnEntity{{Name: "*"}}
	}

	for i := range stmt.columns {
		if 0 < i {
			fmt.Fprint(w, ", ")
		}
		render.Column(w, stmt.columns[i])
		render.ColumnAlias(w, stmt.columns[i])
	}
}

func (render SQLServerRender) Column(w io.Writer, col ColumnEntity) {

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
}

func (render SQLServerRender) ColumnAlias(w io.Writer, col ColumnEntity) {

	switch {
	case col.Alias.Valid:
		fmt.Fprint(w, " AS ")
		fmt.Fprint(w, col.Alias.String)
	}
}

func (render SQLServerRender) Where(w io.Writer, stmt Statement) {
	if stmt.predicator == nil {
		return
	}

	fmt.Fprint(w, " WHERE ")
	render.Predicator(w, stmt.predicator)
}

func (render SQLServerRender) Predicator(w io.Writer, predicator Predicator) {

	switch v := predicator.(type) {
	case ComparisonExpression:
		render.ComparisonOp(w, v)

	case LogicalExpression:
		render.LogicalOp(w, v)
	}
}

func (render SQLServerRender) LogicalOp(w io.Writer, cond LogicalExpression) {

	switch cond.Op {
	case "and":
		for i := range cond.Predicators {
			if 0 < i {
				fmt.Fprint(w, " AND ")
			}

			render.Predicator(w, cond.Predicators[i])
		}
	case "or":
		for i := range cond.Predicators {
			if 0 < i {
				fmt.Fprint(w, " OR ")
			}

			render.Predicator(w, cond.Predicators[i])
		}
	case "not":
		if len(cond.Predicators) == 0 {
			panic(error(ErrEmpty{Op: cond.Op}))
		}

		fmt.Fprint(w, "NOT ")
		render.Predicator(w, cond.Predicators[0])

	}
}

func (render SQLServerRender) ComparisonOp(w io.Writer, cond ComparisonExpression) {
	switch cond.Op {
	case "equal":
		if len(cond.Param) != 1 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 1))
		}

		render.Column(w, cond.Column)
		fmt.Fprintf(w, " = @%s", cond.Param[0])

	case "greater_than":
		if len(cond.Param) != 1 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 1))
		}

		render.Column(w, cond.Column)
		fmt.Fprintf(w, " > @%s", cond.Param[0])

	case "less_than":
		if len(cond.Param) != 1 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 1))
		}

		render.Column(w, cond.Column)
		fmt.Fprintf(w, " < @%s", cond.Param[0])

	case "greater_than_or_equal":
		if len(cond.Param) != 1 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 1))
		}

		render.Column(w, cond.Column)
		fmt.Fprintf(w, " >= @%s", cond.Param[0])

	case "less_than_or_equal":
		if len(cond.Param) != 1 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 1))
		}

		render.Column(w, cond.Column)
		fmt.Fprintf(w, " <= @%s", cond.Param[0])

	case "between":
		if len(cond.Param) != 2 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 2))
		}

		render.Column(w, cond.Column)
		fmt.Fprintf(w, " BETWEEN @%s AND @%s", cond.Param[0], cond.Param[1])

	case "in":
		if len(cond.Param) == 0 {
			panic(fmt.Errorf("%s: expected param count is %d", cond.Op, 1))
		}

		render.Column(w, cond.Column)
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

func (render SQLServerRender) OrderBy(w io.Writer, stmt Statement) {

	for i := range stmt.orders {
		if i == 0 {
			fmt.Fprint(w, " ORDER BY ")
		}

		if 0 < i {
			fmt.Fprint(w, ", ")
		}

		item := stmt.orders[i]

		switch item.Op {
		case "ascending":
			render.Column(w, stmt.columns[i])
			fmt.Fprint(w, " ASC")

		case "desending":
			render.Column(w, stmt.columns[i])
			fmt.Fprint(w, " DESC")
		}
	}
}
