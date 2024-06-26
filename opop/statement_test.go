package qq_test

import (
	"testing"

	qq "github.com/hyhecor/qq/opop"
	"github.com/stretchr/testify/assert"
)

func Test_Select(t *testing.T) {

	{
		T1 := qq.Table("foobar")

		C1 := qq.Column("foo")
		C2 := qq.Column("bar")

		sql := qq.Select(
			qq.Dialect(qq.DialectSQLServer),
			qq.From(T1),
			qq.Columns(C1, C2),
		).Render()

		assert.Equal(t, sql, "SELECT foo, bar FROM foobar")
	}
	{
		T1 := qq.Table("foobar")

		C1 := qq.Column("foo", qq.ColumnTable(T1))
		C2 := qq.Column("bar", qq.ColumnTable(T1))

		sql := qq.Select(
			qq.Dialect(qq.DialectSQLServer),
			qq.From(T1),
			qq.Columns(C1, C2),
		).Render()

		assert.Equal(t, sql, "SELECT foobar.foo, foobar.bar FROM foobar")
	}
	{
		T1 := qq.Table("foobar", qq.TableAlias("T"))

		C1 := qq.Column("foo", qq.ColumnTable(T1))
		C2 := qq.Column("bar", qq.ColumnTable(T1))

		sql := qq.Select(
			qq.Dialect(qq.DialectSQLServer),
			qq.From(T1),
			qq.Columns(C1, C2),
		).Render()

		assert.Equal(t, sql, "SELECT T.foo, T.bar FROM foobar AS T")
	}
	{
		T1 := qq.Table("foobar", qq.TableAlias("T"))

		C1 := qq.Column("foo", qq.ColumnTable(T1), qq.ColunmAlias("f"))
		C2 := qq.Column("bar", qq.ColumnTable(T1), qq.ColunmAlias("b"))

		sql := qq.Select(
			qq.Dialect(qq.DialectSQLServer),
			qq.From(T1),
			qq.Columns(C1, C2),
		).Render()

		assert.Equal(t, sql, "SELECT T.foo AS f, T.bar AS b FROM foobar AS T")
	}
	{
		T1 := qq.Table("foobar", qq.TableAlias("T"))

		C1 := qq.Column("foo", qq.ColumnTable(T1), qq.ColunmAlias("f"))
		C2 := qq.Column("bar", qq.ColumnTable(T1), qq.ColunmAlias("b"))

		sql := qq.Select(
			qq.Dialect(qq.DialectSQLServer),
			qq.From(T1),
			qq.Columns(C1, C2),
			qq.Where(qq.And(qq.In(C1, "one", "two"))),
		).Render()

		assert.Equal(t, sql, "SELECT T.foo AS f, T.bar AS b FROM foobar AS T WHERE T.foo IN (@one, @two)")
	}
	{
		T1 := qq.Table("foobar", qq.TableAlias("T"))

		C1 := qq.Column("foo", qq.ColumnTable(T1), qq.ColunmAlias("f"))
		C2 := qq.Column("bar", qq.ColumnTable(T1), qq.ColunmAlias("b"))

		sql := qq.Select(
			qq.Dialect(qq.DialectSQLServer),
			qq.From(T1),
			qq.Columns(C1, C2),
			qq.Where(qq.And(qq.Between(C1, "one", "two"))),
		).Render()

		assert.Equal(t, sql, "SELECT T.foo AS f, T.bar AS b FROM foobar AS T WHERE T.foo BETWEEN @one AND @two")
	}

	tE := qq.Table("EMPLOYEE", qq.TableAlias("E"))

	empId := qq.Column("empId", qq.ColumnTable(tE))
	name := qq.Column("name", qq.ColumnTable(tE))
	dept := qq.Column("dept", qq.ColumnTable(tE))

	sql := qq.Select(
		qq.Dialect(qq.DialectSQLServer),
		qq.From(tE),
		qq.Columns(empId, name, dept, qq.Literal("getdate()")),
		qq.Where(qq.And(qq.Gt(empId, "empId"))),
		qq.Orders(qq.Asc(empId)),
	).Render()

	assert.Equal(t, sql, "SELECT E.empId, E.name, E.dept, getdate() FROM EMPLOYEE AS E WHERE E.empId > @empId ORDER BY E.empId ASC")

}

func Test_Insert(t *testing.T) {

	T1 := qq.Table("foobar")

	C1 := qq.Column("foo")
	C2 := qq.Column("bar")

	sql := qq.Insert(
		qq.Dialect(qq.DialectSQLServer),
		qq.Into(T1),
		qq.Columns(C1, C2),
		qq.Values("C1", "C2"),
	).Render()

	assert.Equal(t, sql, "INSERT INTO foobar ( foo, bar ) VALUES ( @C1, @C2 )")

}
