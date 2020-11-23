package jira

import (
	"fmt"
	"strings"
)

// Query is the struct that is used to execute a JQL against JIRA
type Query struct {
	JQL    string
	Size   int
	Fields map[string]bool
}

func (connector *ConnectorJira) queryStart() *Query {
	return NewQuery().
		equals(connector.GetFieldMap(backendProject), connector.project)
}

// NewQuery initializes the Query object for querying JIRA's JQL
func NewQuery() (q *Query) {

	q = &Query{
		JQL:    "",
		Fields: make(map[string]bool),
		Size:   1000,
	}

	return q
}

func (q *Query) getKey(in interface{}) (key string) {
	if field, ok := in.(*Field); ok {
		key = field.getQueryID()
	} else if field, ok := in.(string); ok {
		key = field
	}

	return key
}

func (q *Query) and() *Query {
	q.JQL += " AND "

	return q
}

func (q *Query) or() *Query {
	q.JQL += " OR "

	return q
}

func (q *Query) beginGroup() *Query {
	q.JQL += "("

	return q
}

func (q *Query) endGroup() *Query {
	q.JQL += ")"

	return q
}

func (q *Query) equals(in interface{}, value string) *Query {
	q.JQL += fmt.Sprintf("\"%s\" = \"%s\"", q.getKey(in), value)

	return q
}

func (q *Query) lessOrEquals(in interface{}, value string) *Query {
	q.JQL += fmt.Sprintf("\"%s\" <= %s", q.getKey(in), value)

	return q
}

func (q *Query) greaterOrEquals(in interface{}, value string) *Query {

	q.JQL += fmt.Sprintf("\"%s\" >= %s", q.getKey(in), value)

	return q
}

func (q *Query) lessThan(in interface{}, value string) *Query {
	q.JQL += fmt.Sprintf("\"%s\" < %s", q.getKey(in), value)

	return q
}

func (q *Query) greaterThan(in interface{}, value string) *Query {
	q.JQL += fmt.Sprintf("\"%s\" > %s", q.getKey(in), value)

	return q
}

func (q *Query) notEquals(in interface{}, value string) *Query {
	q.JQL += fmt.Sprintf("\"%s\" != \"%s\"", q.getKey(in), value)

	return q
}

func (q *Query) isEmpty(in interface{}) *Query {
	q.JQL += fmt.Sprintf("\"%s\" IS EMPTY", q.getKey(in))

	return q
}

func (q *Query) notEmpty(in interface{}) *Query {
	q.JQL += fmt.Sprintf("\"%s\" IS NOT EMPTY", q.getKey(in))

	return q
}

func (q *Query) contains(in interface{}, value string) *Query {
	q.JQL += fmt.Sprintf("\"%s\" ~ \"(%s)\"", q.getKey(in), cleanString(value))

	return q
}

func (q *Query) doesNotContain(in interface{}, value string) *Query {
	q.JQL += fmt.Sprintf("\"%s\" !~ %s", q.getKey(in), cleanString(value))

	return q
}

func (q *Query) in(in interface{}, value []string) *Query {
	q.JQL += fmt.Sprintf("\"%s\" in (%s)", q.getKey(in), strings.Join(value, ","))

	return q
}

func (q *Query) notIn(in interface{}, value []string) *Query {
	q.JQL += fmt.Sprintf("\"%s\" not in (%s)", q.getKey(in), strings.Join(value, ","))

	return q
}

func (q *Query) orderByDescend(orderBy string) *Query { //ASC
	q.JQL += fmt.Sprintf(" ORDER BY %s DESC ", orderBy)

	return q
}

func (q *Query) orderByAscend(orderBy string) *Query { //ASC
	q.JQL += fmt.Sprintf(" ORDER BY %s ASC ", orderBy)

	return q
}

func cleanString(in string) (out string) {
	if len(in) > 0 {
		out = strings.Replace(in, "\"", "", -1)
		out = strings.Replace(out, "(", "", -1)
		out = strings.Replace(out, ")", "", -1)
	}

	return out
}
