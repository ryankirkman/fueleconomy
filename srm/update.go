package srm

import (
	"bytes"
	"fmt"
	"reflect"
)

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
