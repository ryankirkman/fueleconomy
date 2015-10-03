package srm

import (
	"bytes"
	"fmt"
	"reflect"
)

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
