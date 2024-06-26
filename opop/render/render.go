package render

import (
	"fmt"
	"io"

	"github.com/hyhecor/qq/opop/models"
)

type SQLServerPainter struct{}

func (render SQLServerPainter) Select(w io.Writer, stmt models.Statement) {
	fmt.Fprint(w, "SELECT ")
	render.Columns(w, stmt)
	render.From(w, stmt)
	render.Where(w, stmt)
	render.OrderBy(w, stmt)
}

func (render SQLServerPainter) From(w io.Writer, stmt models.Statement) {
	fmt.Fprint(w, " FROM ")
	fmt.Fprint(w, stmt.From.Name)

	if stmt.From.Alias != nil {
		fmt.Fprint(w, " AS ")
		fmt.Fprint(w, string(*stmt.From.Alias))
	}
}

func (render SQLServerPainter) Columns(w io.Writer, stmt models.Statement) {

	if len(stmt.Columns) == 0 {
		for _, table := range []*models.Table{stmt.From} {
			stmt.Columns = _map(_ColumnToInterface)([]models.Column{{Name: "*", Table: table}})
		}
	}

	for i := range stmt.Columns {
		if 0 < i {
			fmt.Fprint(w, ", ")
		}
		render.Column(w, stmt.Columns[i])
		render.ColumnAlias(w, stmt.Columns[i])
	}
}

func (render SQLServerPainter) Column(w io.Writer, col models.ColumnIndentifier) {

	switch col := col.(type) {
	case models.Column:
		switch {
		case col.Table != nil && col.Table.Alias != nil:
			fmt.Fprint(w, string(*col.Table.Alias))
			fmt.Fprint(w, ".")
			fmt.Fprint(w, col.Name)

		case col.Table != nil && col.Table.Alias == nil:
			fmt.Fprint(w, col.Table.Name)
			fmt.Fprint(w, ".")
			fmt.Fprint(w, col.Name)

		case col.Table == nil:
			fmt.Fprint(w, col.Name)
		}

	case models.Literal:
		fmt.Fprint(w, col.Expression)
	}
}

func (render SQLServerPainter) ColumnAlias(w io.Writer, col models.ColumnIndentifier) {

	switch col := col.(type) {
	case models.Column:
		switch {
		case col.Alias != nil:
			fmt.Fprint(w, " AS ")
			fmt.Fprint(w, string(*col.Alias))
		}

	case models.Literal:
	}
}

func (render SQLServerPainter) Where(w io.Writer, stmt models.Statement) {
	if stmt.Conditions == nil {
		return
	}

	fmt.Fprint(w, " WHERE ")
	render.Predicator(w, stmt.Conditions)
}

func (render SQLServerPainter) Predicator(w io.Writer, predicator models.ConditionIndentifier) {

	switch v := predicator.(type) {
	case models.ComparisonExpression:
		render.ComparisonOp(w, v)

	case models.LogicalExpression:
		render.LogicalOp(w, v)
	}
}

func (render SQLServerPainter) LogicalOp(w io.Writer, cond models.LogicalExpression) {

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

func (render SQLServerPainter) ComparisonOp(w io.Writer, cond models.ComparisonExpression) {
	switch cond.Op {
	case models.Equal:
		MustParamCountOfComparisonExpression(cond, 1)

		render.Column(w, cond.Column)
		fmt.Fprintf(w, " = @%s", cond.Param[0])

	case models.GreaterThan:
		MustParamCountOfComparisonExpression(cond, 1)

		render.Column(w, cond.Column)
		fmt.Fprintf(w, " > @%s", cond.Param[0])

	case models.LessThan:
		MustParamCountOfComparisonExpression(cond, 1)

		render.Column(w, cond.Column)
		fmt.Fprintf(w, " < @%s", cond.Param[0])

	case models.GreaterThanOrEqual:
		MustParamCountOfComparisonExpression(cond, 1)

		render.Column(w, cond.Column)
		fmt.Fprintf(w, " >= @%s", cond.Param[0])

	case models.LessThanOrEqual:
		MustParamCountOfComparisonExpression(cond, 1)

		render.Column(w, cond.Column)
		fmt.Fprintf(w, " <= @%s", cond.Param[0])

	case models.Between:
		MustParamCountOfComparisonExpression(cond, 2)

		render.Column(w, cond.Column)
		fmt.Fprintf(w, " BETWEEN @%s AND @%s", cond.Param[0], cond.Param[1])

	case models.In:
		MustNotParamCountOfComparisonExpression(cond, 0)

		render.Column(w, cond.Column)
		fmt.Fprint(w, " IN (")
		for i := range cond.Param {
			if 0 < i {
				fmt.Fprint(w, ", ")
			}

			fmt.Fprintf(w, "@%s", cond.Param[i])
		}

		fmt.Fprint(w, ")")

	case models.IsNull:
		MustParamCountOfComparisonExpression(cond, 0)

		render.Column(w, cond.Column)
		fmt.Fprintf(w, " IS NULL")
	}
}

func (render SQLServerPainter) OrderBy(w io.Writer, stmt models.Statement) {

	for i := range stmt.Orders {
		if i == 0 {
			fmt.Fprint(w, " ORDER BY ")
		}
		if 0 < i {
			fmt.Fprint(w, ", ")
		}

		item := stmt.Orders[i]

		switch item.Op {
		case models.Ascending:
			render.Column(w, item.Column)
			fmt.Fprint(w, " ASC")

		case models.Desending:
			render.Column(w, item.Column)
			fmt.Fprint(w, " DESC")
		}
	}
}

func (render SQLServerPainter) Insert(w io.Writer, stmt models.Statement) {
	fmt.Fprint(w, "INSERT ")
	render.Into(w, stmt)
	render.InsertColumns(w, stmt)
	render.Values(w, stmt)
}

func (render SQLServerPainter) Into(w io.Writer, stmt models.Statement) {
	fmt.Fprint(w, "INTO ")
	fmt.Fprint(w, stmt.Into.Name)
}

func (render SQLServerPainter) InsertColumns(w io.Writer, stmt models.Statement) {

	if len(stmt.Columns) == 0 {
		return
	}

	fmt.Fprint(w, " ( ")

	for i := range stmt.Columns {
		if 0 < i {
			fmt.Fprint(w, ", ")
		}
		render.InsertColumn(w, stmt.Columns[i])
	}

	fmt.Fprint(w, " )")
}

func (render SQLServerPainter) InsertColumn(w io.Writer, col models.ColumnIndentifier) {

	switch col := col.(type) {
	case models.Column:
		fmt.Fprint(w, col.Name)

	case models.Literal:
		fmt.Fprint(w, col.Expression)
	}
}

func (render SQLServerPainter) Values(w io.Writer, stmt models.Statement) {

	if len(stmt.Columns) == 0 {
		return
	}

	fmt.Fprint(w, " VALUES ")

	fmt.Fprint(w, "( ")

	for i := range stmt.Values {
		if 0 < i {
			fmt.Fprint(w, ", ")
		}

		fmt.Fprintf(w, "@%s", stmt.Values[i])
	}

	fmt.Fprint(w, " )")
}
