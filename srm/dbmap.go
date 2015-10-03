package srm

import (
	"database/sql"
)

// Struct Relational Mapper
type DbMap struct {
	Conn    *sql.DB
	Dialect Dialect
}

func (db *DbMap) DeleteAll(table string) (err error) {
	err = deleteall(db, table)
	return err
}

func (db *DbMap) InsertMany(table string, list ...interface{}) (insertedIds []int, err error) {
	for _, ptr := range list {
		insertedId, err := insert(db, table, ptr)
		if err != nil {
			return insertedIds, err
		}
		insertedIds = append(insertedIds, insertedId)
	}
	return insertedIds, nil
}

func (db *DbMap) InsertOne(table string, ptr interface{}) (insertedId int, err error) {
	insertedId, err = insert(db, table, ptr)
	return insertedId, err
}

func (db *DbMap) SelectOne(ptr interface{}, query string, args ...interface{}) (err error) {
	err = selectone(db, ptr, query, args...)
	return err
}

func (db *DbMap) SelectMany(ptr interface{}, query string, args ...interface{}) (err error) {
	err = selectmany(db, ptr, query, args...)
	return err
}

func (db *DbMap) UpdateOne(table string, updateOnField string, ptr interface{}) (rowsAffected int64, err error) {
	rowsAffected, err = update(db, table, updateOnField, ptr)
	return rowsAffected, err
}

func (db *DbMap) UpsertOne(table string, updateOnField string, ptr interface{}) (insertedId int, err error) {
	insertedId, err = multiQueryUpsert(db, table, updateOnField, ptr)
	return insertedId, err
}

func (db *DbMap) UpsertMany(table string, updateOnField string, list ...interface{}) (insertedIds []int, err error) {
	for _, ptr := range list {
		insertedId, err := multiQueryUpsert(db, table, updateOnField, ptr)
		if err != nil {
			return insertedIds, err
		}
		insertedIds = append(insertedIds, insertedId)
	}
	return insertedIds, nil
}
