package query

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps/psql/base"
)

const QueryFile = "/tmp/query.txt"

type Args struct {
	QueryUsingFile bool `env:"QUERY_USING_FILE", envDefault:"false"`
}

type PsqlQuery struct {
	BasePsql *base.PsqlCommand
	Query    string
	args     Args
}

func NewPsqlQuery(query string) (*PsqlQuery, error) {
	args := &Args{}
	if err := env.Parse(args); err != nil {
		return nil, fmt.Errorf("parse args: %w", err)
	}

	basePsql, err := base.NewPsqlCommand()
	if err != nil {
		return nil, err
	}

	if args.QueryUsingFile {
		if err := ioutil.WriteFile(QueryFile, []byte(query), os.ModePerm); err != nil {
			return nil, fmt.Errorf("can't dump query to file: %w", err)
		}
	}

	return &PsqlQuery{
		BasePsql: basePsql,
		Query:    query,
		args:     *args,
	}, nil
}

func (p *PsqlQuery) RunPsqlQueryCommand(queryExtraArgs []string) ([]byte, int, error) {
	var extraArgs []string
	if p.args.QueryUsingFile {
		extraArgs = []string{"-f", QueryFile}
	} else {
		query := strings.Replace(strings.Replace(p.Query, "\r", "", -1), "\n", " ", -1)
		extraArgs = []string{"-c", query}
	}

	extraArgs = append(extraArgs, queryExtraArgs...)

	return p.BasePsql.RunPsqlCommand(extraArgs)
}

func (p *PsqlQuery) ParseOutput(output []byte) []byte {
	return p.BasePsql.ParseOutput(output)
}
