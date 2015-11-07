package srm

import (
	"bytes"
	"fmt"
	"strings"
)

type QueryBuilder struct {
	Db    *DbMap
	Table string
}

func (qb QueryBuilder) BuildCount(whereExact map[string]interface{},
	whereFuzzy map[string]string) (string, []interface{}) {

	var sqlArgs []interface{}
	sqlQuery := bytes.Buffer{}
	sqlQuery.WriteString("SELECT COUNT(*) FROM ")
	sqlQuery.WriteString(qb.Table)
	first := true
	count := 1

	for col, val := range whereExact {
		if first {
			sqlQuery.WriteString(" WHERE ")
		}
		if !first {
			sqlQuery.WriteString(" AND ")
		}
		sqlQuery.WriteString(fmt.Sprintf("%s = %s", col, qb.Db.Dialect.Placeholder(count)))
		sqlArgs = append(sqlArgs, val)
		first = false
		count++
	}

	for col, val := range whereFuzzy {
		if first {
			sqlQuery.WriteString(" WHERE ")
		}
		if !first {
			sqlQuery.WriteString(" AND ")
		}
		sqlQuery.WriteString(fmt.Sprintf("lower(%s) LIKE %s", col, qb.Db.Dialect.Placeholder(count)))
		sqlArgs = append(sqlArgs, trailingPercent(strings.ToLower(val)))
		first = false
		count++
	}

	return sqlQuery.String(), sqlArgs
}

func (qb QueryBuilder) BuildSelect(limit int, offset int, whereExact map[string]interface{},
	whereFuzzy map[string]string) (string, []interface{}) {

	var sqlArgs []interface{}
	sqlQuery := bytes.Buffer{}
	sqlQuery.WriteString("SELECT * FROM ")
	sqlQuery.WriteString(qb.Table)
	first := true
	count := 1

	for col, val := range whereExact {
		if first {
			sqlQuery.WriteString(" WHERE ")
		}
		if !first {
			sqlQuery.WriteString(" AND ")
		}
		sqlQuery.WriteString(fmt.Sprintf("%s = %s", col, qb.Db.Dialect.Placeholder(count)))
		sqlArgs = append(sqlArgs, val)
		first = false
		count++
	}

	for col, val := range whereFuzzy {
		if first {
			sqlQuery.WriteString(" WHERE ")
		}
		if !first {
			sqlQuery.WriteString(" AND ")
		}
		sqlQuery.WriteString(fmt.Sprintf("lower(%s) LIKE %s", col, qb.Db.Dialect.Placeholder(count)))
		sqlArgs = append(sqlArgs, trailingPercent(strings.ToLower(val)))
		first = false
		count++
	}

	if limit > 0 {
		sqlQuery.WriteString(fmt.Sprintf(" LIMIT %s", qb.Db.Dialect.Placeholder(count)))
		sqlArgs = append(sqlArgs, limit)
		count++
	}

	if offset > 0 {
		sqlQuery.WriteString(fmt.Sprintf(" OFFSET %s", qb.Db.Dialect.Placeholder(count)))
		sqlArgs = append(sqlArgs, offset)
	}

	return sqlQuery.String(), sqlArgs
}

func trailingPercent(str string) string {
	buff := bytes.Buffer{}
	buff.WriteString(str)
	buff.WriteString("%")

	return buff.String()
}
