package qq_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/hyhecor/qq"
)

func TestColumn(t *testing.T) {

	tests := []struct {
		name string
		args qq.ColumnEntity
		want string
	}{
		// TODO: Add test cases.
		{"", qq.Column("foo"), "foo"},
		{"", qq.Column("foo", qq.WithAlias("x")), "foo AS x"},
		{"", qq.Column("foo", qq.WithTable(qq.Table("foobar"))), "foobar.foo"},
		{"", qq.Column("foo", qq.WithTable(qq.Table("foobar", "tbl"))), "tbl.foo"},
		{"", qq.Column("foo", qq.WithAlias("x"), qq.WithTable(qq.Table("foobar", "tbl"))), "tbl.foo AS x"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var buf bytes.Buffer
			qq.SQLServerRender{}.Column(&buf, tt.args)

			if got := buf.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Column() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhere(t *testing.T) {

	tests := []struct {
		name string
		args qq.Predicator
		want string
	}{
		// TODO: Add test cases.
		{"", qq.Eq(qq.Column("foo"), "foo"), "foo = @foo"},
		{"Alias 설정 무시", qq.Eq(qq.Column("foo", qq.WithAlias("x")), "foo"), "foo = @foo"},
		{"", qq.Eq(qq.Column("foo", qq.WithTable(qq.Table("foobar"))), "foo"), "foobar.foo = @foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var buf bytes.Buffer
			qq.SQLServerRender{}.Predicator(&buf, tt.args)

			if got := buf.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Where() = %v, want %v", got, tt.want)
			}
		})
	}
}
