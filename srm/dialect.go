package srm

import (
	"fmt"
)

type Dialect interface {
	InsertQuerySuffix(string) string
	Insert(*DbMap, string, ...interface{}) (int, error)
	Placeholder(int) string
}

type PostgresDialect struct{}

func (p PostgresDialect) InsertQuerySuffix(pkName string) string {
	return fmt.Sprintf(" RETURNING %s;", pkName)
}

func (p PostgresDialect) Insert(db *DbMap, sqlString string, params ...interface{}) (insertedId int, err error) {
	stmt, err := db.Conn.Prepare(sqlString)
	if err != nil {
		return insertedId, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(params...).Scan(&insertedId)
	if err != nil {
		return insertedId, err
	}

	return insertedId, err
}

func (p PostgresDialect) Placeholder(count int) string {
	return fmt.Sprintf("$%d", count)
}

type Sqlite3Dialect struct{}

func (s Sqlite3Dialect) InsertQuerySuffix(pkName string) string {
	return ";"
}

func (p Sqlite3Dialect) Insert(db *DbMap, sqlString string, params ...interface{}) (insertedId int, err error) {
	r, err := db.Conn.Exec(sqlString, params...)
	if err != nil {
		return insertedId, err
	}

	insertedId64, err := r.LastInsertId()
	if err != nil {
		return insertedId, err
	}

	return int(insertedId64), err
}

func (s Sqlite3Dialect) Placeholder(count int) string {
	return "?"
}
