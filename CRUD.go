package qq

import (
	"database/sql"
)

type Foobar struct {
	Foo string
	Bar string
}

func QueryFoobar(db *sql.DB, foo string) error {

	rows, err := db.Query("SELECT foo, bar FROM foobar WHERE foo = ?", foo)
	if err != nil {
		return err
	}

	type Foobar struct {
		Foo string
		Bar string
	}

	var foobar Foobar
	for rows.Next() {
		if err := rows.Scan(&foobar.Foo, &foobar.Bar); err != nil {
			return err
		}
	}

	return nil
}

func OneFoobar(db *sql.DB, foo string, bar string) error {

	rows, err := db.Query("SELECT foo, bar FROM foobar WHERE foo = ? AND bar = ?", foo, bar)
	if err != nil {
		return err
	}

	var foobar Foobar
	for rows.Next() {
		if err := rows.Scan(&foobar.Foo, &foobar.Bar); err != nil {
			return err
		}
	}

	return nil
}
