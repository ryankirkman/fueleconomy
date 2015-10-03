package srm

import (
	"reflect"
	"strings"
	"time"
)

// for fields that exist in DB table, but don't exist in struct
type dummyField struct{}

// Scan implements the Scanner interface.
func (nt *dummyField) Scan(value interface{}) error {
	return nil
}

func getColumnForField(v reflect.StructField) string {
	tags := strings.Split(v.Tag.Get("db"), ", ")
	colName := tags[0]
	if colName == "-" {
		return ""
	}
	if len(tags) > 1 && (tags[1] == "primaryKey" || tags[1] == "autoSet") {
		return ""
	}
	if colName != "" {
		return colName
	}

	return strings.ToLower(v.Name)
}

func getValueForField(v reflect.Value) interface{} {
	if t, ok := v.Interface().(time.Time); ok {
		return t
	}

	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Float64:
		return v.Float()
	case reflect.Int:
		return v.Int()
	case reflect.Bool:
		return v.Bool()
	default:
		return nil
	}
}

func makeColToFieldIndex(t reflect.Type, cols []string) [][]int {
	var colToFieldIndex = make([][]int, len(cols))

	for x := range cols {
		colName := strings.ToLower(cols[x])
		field, found := t.FieldByNameFunc(func(fieldName string) bool {
			field, _ := t.FieldByName(fieldName)
			tags := strings.Split(field.Tag.Get("db"), ", ")
			dbNameFromTag := tags[0]
			if dbNameFromTag == "" {
				return strings.ToLower(field.Name) == colName
			}
			return colName == dbNameFromTag
		})
		if found {
			colToFieldIndex[x] = field.Index
		}
	}

	return colToFieldIndex
}
