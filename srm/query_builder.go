package srm

import (
	"bytes"
	"fmt"
	"strings"
)

type QueryBuilder struct {
	Db         *DbMap
	Table      string
	Limit      int
	Offset     int
	WhereExact map[string]interface{}
	WhereFuzzy map[string]string
}

func (qb *QueryBuilder) BuildCount() (string, []interface{}) {
	var sqlArgs []interface{}
	sqlQuery := bytes.Buffer{}
	sqlQuery.WriteString("SELECT COUNT(*) FROM ")
	sqlQuery.WriteString(qb.Table)
	first := true
	count := 1

	sqlQuery.WriteString(qb.buildWhere(&count, &first, &sqlArgs))

	return sqlQuery.String(), sqlArgs
}

func (qb *QueryBuilder) BuildSelect() (string, []interface{}) {
	var sqlArgs []interface{}
	sqlQuery := bytes.Buffer{}
	sqlQuery.WriteString("SELECT * FROM ")
	sqlQuery.WriteString(qb.Table)
	first := true
	count := 1
	sqlQuery.WriteString(qb.buildWhere(&count, &first, &sqlArgs))

	if qb.Limit > 0 {
		sqlQuery.WriteString(fmt.Sprintf(" LIMIT %s", qb.Db.Dialect.Placeholder(count)))
		sqlArgs = append(sqlArgs, qb.Limit)
		count++
	}

	if qb.Offset > 0 {
		sqlQuery.WriteString(fmt.Sprintf(" OFFSET %s", qb.Db.Dialect.Placeholder(count)))
		sqlArgs = append(sqlArgs, qb.Offset)
	}

	return sqlQuery.String(), sqlArgs
}

func (qb *QueryBuilder) buildWhere(count *int, first *bool, args *[]interface{}) string {
	buff := bytes.Buffer{}
	for col, val := range qb.WhereExact {
		if *first {
			buff.WriteString(" WHERE ")
		}
		if !*first {
			buff.WriteString(" AND ")
		}
		buff.WriteString(fmt.Sprintf("%s = %s", col, qb.Db.Dialect.Placeholder(*count)))
		*args = append(*args, val)
		*first = false
		*count++
	}
	for col, val := range qb.WhereFuzzy {
		if *first {
			buff.WriteString(" WHERE ")
		}
		if !*first {
			buff.WriteString(" AND ")
		}
		buff.WriteString(fmt.Sprintf("lower(%s) LIKE %s", col, qb.Db.Dialect.Placeholder(*count)))
		*args = append(*args, trailingPercent(strings.ToLower(val)))
		*first = false
		*count++
	}
	return buff.String()
}

func trailingPercent(str string) string {
	buff := bytes.Buffer{}
	buff.WriteString(str)
	buff.WriteString("%")

	return buff.String()
}
