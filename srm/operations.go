package srm

import (
	"bytes"
	"fmt"
	"reflect"
)

// TODO
// Write single query upsert using postgres function and test perfomance
func multiQueryUpsert(db *DbMap, table string, updateOnField string, ptr interface{}) (insertedId int, err error) {
	rowsAffected, err := db.UpdateOne(table, updateOnField, ptr)
	if err != nil {
		return insertedId, err
	}
	if rowsAffected == 0 {
		insertedId, err = db.InsertOne(table, ptr)
		if err != nil {
			return insertedId, err
		}
	}
	return insertedId, nil
}

func deleteall(db *DbMap, table string) error {
	_, err := db.Conn.Exec(fmt.Sprintf("DELETE FROM %s", table))
	if err != nil {
		return err
	}
	return nil
}

func insert(db *DbMap, table string, ptr interface{}) (insertedId int, err error) {
	var (
		queryBuffer  bytes.Buffer
		valuesBuffer bytes.Buffer
		valuesSlice  []interface{}
	)

	queryBuffer.WriteString(fmt.Sprintf("INSERT INTO %s (", table))

	count := 1
	ptrv := reflect.ValueOf(ptr).Elem()
	for i := 0; i < ptrv.NumField(); i++ {
		field := ptrv.Type().Field(i)
		columnName := getColumnForField(field)

		if columnName == "" {
			continue
		}

		if count > 1 {
			queryBuffer.WriteString(",")
			valuesBuffer.WriteString(",")
		}

		queryBuffer.WriteString(columnName)
		valuesBuffer.WriteString(db.Dialect.Placeholder(count))

		fieldValue := getValueForField(ptrv.Field(i))
		valuesSlice = append(valuesSlice, fieldValue)

		count++
	}
	queryBuffer.WriteString(") VALUES (")
	queryBuffer.WriteString(valuesBuffer.String())
	// TODO
	// Dynamic id field
	queryBuffer.WriteString(")")
	queryBuffer.WriteString(db.Dialect.InsertQuerySuffix("id"))

	insertedId, err = db.Dialect.Insert(db, queryBuffer.String(), valuesSlice...)
	if err != nil {
		return insertedId, err
	}

	return insertedId, nil
}

func update(db *DbMap, table string, updateOnColumn string, ptr interface{}) (rowsAffected int64, err error) {
	var (
		queryBuffer   bytes.Buffer
		valuesSlice   []interface{}
		updateOnValue interface{}
	)

	queryBuffer.WriteString(fmt.Sprintf("UPDATE %s SET ", table))

	count := 1
	ptrv := reflect.ValueOf(ptr).Elem()
	for i := 0; i < ptrv.NumField(); i++ {
		field := ptrv.Type().Field(i)
		columnName := getColumnForField(field)

		if columnName == "" {
			continue
		}

		if count > 1 {
			queryBuffer.WriteString(", ")
		}

		queryBuffer.WriteString(fmt.Sprintf("%s = %s", columnName, db.Dialect.Placeholder(count)))

		fieldValue := getValueForField(ptrv.Field(i))
		valuesSlice = append(valuesSlice, fieldValue)

		if columnName == updateOnColumn {
			updateOnValue = fieldValue
		}

		count++
	}

	valuesSlice = append(valuesSlice, updateOnValue)

	queryBuffer.WriteString(fmt.Sprintf(" WHERE %s = %s;", updateOnColumn,
		db.Dialect.Placeholder(count)))

	stmt, err := db.Conn.Prepare(queryBuffer.String())
	if err != nil {
		return rowsAffected, err
	}
	defer stmt.Close()

	r, err := stmt.Exec(valuesSlice...)
	if err != nil {
		return rowsAffected, err
	}

	rowsAffected, err = r.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}
