package main

import (
	"fmt"
	"os"

	queryBase "github.com/stackpulse/public-steps/psql/query/base"
	"github.com/stackpulse/steps-sdk-go/step"
)

const QueryEnv = "QUERY"

type PsqlRaw struct {
	psqlQuery *queryBase.PsqlQuery
}

func (p *PsqlRaw) Init() error {
	query := os.Getenv(QueryEnv)
	if query == "" {
		return fmt.Errorf("must specify %s env", QueryEnv)
	}
	psqlQuery, err := queryBase.NewPsqlQuery(query)
	if err != nil {
		return fmt.Errorf("init psql query: %w", err)
	}

	p.psqlQuery = psqlQuery
	return nil
}

func (p *PsqlRaw) Run() (int, []byte, error) {
	output, exitCode, err := p.psqlQuery.RunPsqlQueryCommand(nil)
	if err != nil {
		return exitCode, output, err
	}
	return exitCode, p.psqlQuery.ParseOutput(output), nil
}

func main() {
	step.Run(&PsqlRaw{})
}
