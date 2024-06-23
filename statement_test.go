package qq_test

import (
	"reflect"
	"testing"

	"github.com/hyhecor/qq"
)

func TestStatement(t *testing.T) {
	dialect := qq.Dialect(qq.DialectSQLServer)
	table := qq.Table("foobar")
	tabelWithAlias := qq.Table("foobar", "a")
	withTable := qq.WithTable(table)
	withTableWithAlias := qq.WithTable(tabelWithAlias)

	tests := []struct {
		name string
		args string
		want string
	}{
		// TODO: Add test cases.
		{"", dialect.From(table).Select(),
			"SELECT * FROM foobar"},
		{"", dialect.From(table).Columns(qq.Column("foo"), qq.Column("bar")).Select(),
			"SELECT foo, bar FROM foobar"},
		{"", dialect.From(table).Columns(qq.Column("foo"), qq.Column("bar")).OrderBy(qq.Asc(qq.Column("foo"))).Select(),
			"SELECT foo, bar FROM foobar ORDER BY foo ASC"},
		{"", dialect.From(table).Columns(qq.Column("foo"), qq.Column("bar")).OrderBy(qq.Asc(qq.Column("foo")), qq.Desc(qq.Column("bar"))).Select(),
			"SELECT foo, bar FROM foobar ORDER BY foo ASC, bar DESC"},

		{"", dialect.From(table).Columns(qq.Column("foo", withTable), qq.Column("bar", withTable)).Select(),
			"SELECT foobar.foo, foobar.bar FROM foobar"},
		{"", dialect.From(table).Columns(qq.Column("foo", withTable), qq.Column("bar", withTable)).OrderBy(qq.Asc(qq.Column("foo", withTable))).Select(),
			"SELECT foobar.foo, foobar.bar FROM foobar ORDER BY foobar.foo ASC"},

		{"", dialect.From(tabelWithAlias).Select(),
			"SELECT * FROM foobar AS a"},
		{"", dialect.From(tabelWithAlias).Columns(qq.Column("foo", withTableWithAlias), qq.Column("bar", withTableWithAlias)).Select(),
			"SELECT a.foo, a.bar FROM foobar AS a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.args; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Column() = %v, want %v", got, tt.want)
			}
		})
	}
}
