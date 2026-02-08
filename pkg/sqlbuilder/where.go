package sqlbuilder

import (
	"bytes"
	"strings"
)

type valueBuilder struct {
	args []any
	sql  *bytes.Buffer
}

type whereBuilder struct {
	parent     sqlBuilder
	args       []any
	sql        string
	sqlGroupBy string
	sqlOrderBy string
}

func (w *whereBuilder) SetPaginate(limit, offset int64) WhereBuilder {
	w.parent.dialect.SetPagination()
	w.args = append(w.args, offset)
	w.args = append(w.args, limit)
	return w
}

func (w *whereBuilder) Build() (string, []any) {
	query := bytes.Buffer{}
	query.WriteString(w.parent.sql)
	query.WriteString(w.sql)

	if strings.TrimSpace(w.sqlGroupBy) != "" {
		query.WriteString(w.sqlGroupBy)
	}

	if strings.TrimSpace(w.sqlOrderBy) != "" {
		query.WriteString(w.sqlOrderBy)
	}

	return w.parent.dialect.ToSQL(query.String()), w.args
}

func (v *valueBuilder) Add(sql string, args ...any) {
	v.sql.WriteString(" ")
	v.sql.WriteString(strings.TrimSpace(sql))

	v.args = append(v.args, args...)
}

func (w sqlBuilder) Build() (string, []any) {
	query := bytes.Buffer{}
	query.WriteString(w.sql)
	return query.String(), nil
}

func (w *whereBuilder) OrderBy(sql string) WhereBuilder {
	w.sqlOrderBy = " " + strings.TrimSpace(sql)
	return w
}

func (w *whereBuilder) GroupBy(sql string) WhereBuilder {
	w.sqlGroupBy = " " + strings.TrimSpace(sql)
	return w
}
