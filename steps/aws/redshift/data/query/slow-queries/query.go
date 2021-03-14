package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed query.tpl
var templateQuery string

func generate(limit int) (string, error) {
	if limit < 0 {
		return "", fmt.Errorf("invalid limit: %d", limit)
	}

	tmpl, err := template.New("query").Parse(templateQuery)
	if err != nil {
		return "", fmt.Errorf("generate template: %w", err)
	}

	out := &bytes.Buffer{}
	err = tmpl.Execute(out, struct {
		Limit int
	}{Limit: limit})
	if err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return out.String(), nil
}
