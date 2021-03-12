package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

const (
	defaultBlockingSeconds = 60
)

//go:embed query.tpl
var templateQuery string

type query struct {
	onlyUsers       []string
	ignoredUsers    []string
	blockingSeconds int
}

func (q *query) validate() error {
	if len(q.onlyUsers) > 0 && len(q.ignoredUsers) > 0 {
		return fmt.Errorf("either only-users or ignored-users can be set")
	}

	if q.blockingSeconds <= 0 {
		q.blockingSeconds = defaultBlockingSeconds
	}

	return nil
}

func (q *query) quoteList(values []string) []string {
	var list []string
	for _, v := range values {
		list = append(list, "'"+v+"'")
	}
	return list
}

func (q *query) whereClause() (string, error) {
	if err := q.validate(); err != nil {
		return "", err
	}

	var where strings.Builder
	if len(q.onlyUsers) > 0 {
		users := strings.Join(q.quoteList(q.onlyUsers), ", ")
		where.WriteString(fmt.Sprintf("AND username IN (%s) ", users))
	}

	if len(q.ignoredUsers) > 0 {
		users := strings.Join(q.quoteList(q.ignoredUsers), ", ")
		where.WriteString(fmt.Sprintf("AND username NOT IN (%s) ", users))
	}

	where.WriteString(fmt.Sprintf("AND block_sec>%d ", q.blockingSeconds))

	return where.String(), nil
}

func (q *query) generate() (string, error) {
	if err := q.validate(); err != nil {
		return "", err
	}

	where, err := q.whereClause()
	if err != nil {
		return "", fmt.Errorf("generate where clause: %w", err)
	}

	tmpl, err := template.New("query").Parse(templateQuery)
	if err != nil {
		return "", fmt.Errorf("generate template: %w", err)
	}

	out := &bytes.Buffer{}
	err = tmpl.Execute(out, struct {
		Conditions string
	}{Conditions: where})
	if err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return out.String(), nil
}
