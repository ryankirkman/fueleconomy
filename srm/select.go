package srm

import (
	"reflect"
)

func selectone(db *DbMap, ptr interface{}, query string, args ...interface{}) error {
	structVal := reflect.Indirect(reflect.ValueOf(ptr))

	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	colToFieldIndex := makeColToFieldIndex(structVal.Type(), cols)
	dest := make([]interface{}, len(cols))

	for x := range cols {
		idx := colToFieldIndex[x]
		if idx == nil {
			var dummy dummyField
			dest[x] = &dummy
			continue
		}

		f := structVal.FieldByIndex(idx)
		target := f.Addr().Interface()
		dest[x] = target
	}

	rows.Next()
	err = rows.Scan(dest...)
	if err != nil {
		return err
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func selectmany(db *DbMap, ptr interface{}, query string, args ...interface{}) error {
	sliceVal := reflect.Indirect(reflect.ValueOf(ptr))
	structType := reflect.TypeOf(ptr).Elem().Elem()

	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	colToFieldIndex := makeColToFieldIndex(structType, cols)
	for {
		if !rows.Next() {
			if rows.Err() != nil {
				return rows.Err()
			}
			break
		}

		structVal := reflect.New(structType)
		dest := make([]interface{}, len(cols))

		for x := range cols {
			el := structVal.Elem()
			idx := colToFieldIndex[x]
			if idx == nil {
				var dummy dummyField
				dest[x] = &dummy
				continue
			}

			f := el.FieldByIndex(idx)
			target := f.Addr().Interface()
			dest[x] = target
		}

		err := rows.Scan(dest...)
		if err != nil {
			return err
		}

		sliceVal.Set(reflect.Append(sliceVal, structVal.Elem()))
	}

	return nil
}
