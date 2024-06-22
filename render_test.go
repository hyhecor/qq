package qq_test

import (
	"testing"

	"github.com/hyhecor/qq"
)

func TestMain(t *testing.T) {
	tbl := qq.Table("foobar", "t")

	q := qq.From(tbl).
		Dialect(qq.DialectSQLServer).
		Columns(
			qq.Column("foo"),
			qq.Column("foo", qq.WithTable(tbl)),
			qq.Column("foo", qq.WithAlias("x")),
			qq.Column("foo", qq.WithTable(tbl), qq.WithAlias("x")),
			qq.Column("bar"),
		).
		Where(
			qq.In(qq.Column("foo", qq.WithTable(tbl)), "enum1", "enum2", "enum3"),
			qq.Or(
				qq.Eq(qq.Column("foo", qq.WithTable(tbl)), "foo1"),
				qq.Eq(qq.Column("foo", qq.WithTable(tbl)), "foo2"),
			),

			qq.Between(qq.Column("bar", qq.WithTable(tbl)), "enum1", "enum2"),
		).
		Select()

	println(q)
}
