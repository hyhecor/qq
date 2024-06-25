package qq_test

import (
	"testing"

	qq "github.com/hyhecor/qq/opop"
)

func Test_Select(t *testing.T) {

	t1 := qq.Table("foobar", qq.TableAlias("t1"))

	c1 := qq.Column("foo", qq.ColumnTable(t1))
	c2 := qq.Column("bar", qq.ColumnTable(t1))

	sql := qq.Select(
		qq.Dialect(qq.DialectSQLServer),
		qq.From(t1),
		qq.Columns(c1, c2),
		qq.Where(qq.Eq(c1, "foo"), qq.Gt(c2, "bar")),
		qq.Orders(qq.Asc(c1), qq.Desc(c2)),
	).Render()

	t.Log(sql)
	_ = sql
}
