package models

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func parseStringToOutStruct(inValue reflect.Value, inField reflect.StructField, outValue reflect.Value) error {
	var outField reflect.Value
	outEl := outValue.Elem()

	if boolTag := inField.Tag.Get("parseBool"); boolTag != "" {
		split := strings.Split(boolTag, ",")
		boolOutName := split[0]
		boolTrueVals := split[1:]
		outToSet := parseStringToBool(inValue.String(), boolTrueVals...)
		outField = outEl.FieldByName(boolOutName)
		outField.Set(reflect.ValueOf(outToSet))
		return nil
	}

	outField = outEl.FieldByName(inField.Name)

	switch outField.Kind() {
	case reflect.String:
		outField.Set(inValue)
	case reflect.Int:
		outField.SetInt(parseStringToInt(inValue.String()))
	case reflect.Float64:
		outField.SetFloat(parseStringToFloat(inValue.String()))
	default:
		if _, ok := outField.Interface().(time.Time); ok {
			outField.Set(reflect.ValueOf(parseStringToTime(inValue.String())))
		} else {
			return errors.New(fmt.Sprintf("Unsupported: %s", outField.Kind().String()))
		}
	}

	return nil
}

func parseStringToBool(val string, trueVals ...string) bool {
	for i := range trueVals {
		if val == trueVals[i] {
			return true
		}
	}
	return false
}

func parseStringToFloat(val string) (f float64) {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		f = 0.0
	}
	if f == -1.0 {
		f = 0.0
	}
	return f
}

func parseStringToInt(val string) (i int64) {
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		i = 0
	}
	if i == -1 {
		i = 0
	}
	return i
}

func parseStringToTime(val string) time.Time {
	t, err := time.Parse("2006-01-02T00:00:00-07:00", val)
	if err != nil {
		return time.Time{}
	}
	return t
}
