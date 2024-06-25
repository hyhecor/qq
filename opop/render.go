package qq

import (
	"fmt"
	"io"

	"github.com/hyhecor/qq/opop/models"
)

type Painter interface {
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
	fmt.Fprint(w, stmt.Table.Name)

	if stmt.Table.Alias != nil {
		fmt.Fprint(w, " AS ")
		fmt.Fprint(w, string(*stmt.Table.Alias))
	}
}

func (render SQLServerRender) Columns(w io.Writer, stmt Statement) {

	if len(stmt.Columns) == 0 {
		stmt.Columns = []models.Column{{Name: "*"}}
	}

	for i := range stmt.Columns {
		if 0 < i {
			fmt.Fprint(w, ", ")
		}
		render.Column(w, stmt.Columns[i])
		render.ColumnAlias(w, stmt.Columns[i])
	}
}

func (render SQLServerRender) Column(w io.Writer, col models.Column) {

	switch {
	case 0 < len(col.Table.Name) && col.Table.Alias != nil:
		fmt.Fprint(w, string(*col.Table.Alias))
		fmt.Fprint(w, ".")
		fmt.Fprint(w, col.Name)

	case 0 < len(col.Table.Name) && col.Table.Alias == nil:
		fmt.Fprint(w, col.Table.Name)
		fmt.Fprint(w, ".")
		fmt.Fprint(w, col.Name)

	case 0 < len(col.Table.Name):
		fmt.Fprint(w, col.Name)
	}
}

func (render SQLServerRender) ColumnAlias(w io.Writer, col models.Column) {

	switch {
	case col.Alias != nil:
		fmt.Fprint(w, " AS ")
		fmt.Fprint(w, string(*col.Alias))
	}
}

func (render SQLServerRender) Where(w io.Writer, stmt Statement) {
	if stmt.Predicator == nil {
		return
	}

	fmt.Fprint(w, " WHERE ")
	render.Predicator(w, stmt.Predicator)
}

func (render SQLServerRender) Predicator(w io.Writer, predicator models.Predicator) {

	switch v := predicator.(type) {
	case models.ComparisonExpression:
		render.ComparisonOp(w, v)

	case models.LogicalExpression:
		render.LogicalOp(w, v)
	}
}

func (render SQLServerRender) LogicalOp(w io.Writer, cond models.LogicalExpression) {

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

func (render SQLServerRender) ComparisonOp(w io.Writer, cond models.ComparisonExpression) {
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

	for i := range stmt.Orders {
		if i == 0 {
			fmt.Fprint(w, " ORDER BY ")
		}

		if 0 < i {
			fmt.Fprint(w, ", ")
		}

		item := stmt.Orders[i]

		switch item.Op {
		case "ascending":
			render.Column(w, stmt.Columns[i])
			fmt.Fprint(w, " ASC")

		case "desending":
			render.Column(w, stmt.Columns[i])
			fmt.Fprint(w, " DESC")
		}
	}
}
