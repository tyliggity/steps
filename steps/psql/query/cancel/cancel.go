package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/lib/pq"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
	psqlBase "github.com/stackpulse/steps/psql/base"
	queryBase "github.com/stackpulse/steps/psql/query/base"
)

const cancelQuery = `
SELECT pg_cancel_backend(pid) FROM pg_stat_activity WHERE pid=%s;
`

type Args struct {
	PIDs []int `env:"PIDS,required"`
}

type PsqlCancel struct {
	psqlQuery *queryBase.PsqlQuery
}

func (p *PsqlCancel) Init() error {
	args := &Args{}
	if err := env.Parse(args); err != nil {
		return fmt.Errorf("parsing arguments: %w", err)
	}
	if len(args.PIDs) == 0 {
		return fmt.Errorf("must specify at least one pid")
	}

	pidsStr := make([]string, len(args.PIDs))
	for i, pid := range args.PIDs {
		pidsStr[i] = pq.QuoteLiteral(strconv.Itoa(pid))
	}
	psqlQuery, err := queryBase.NewPsqlQuery(fmt.Sprintf(cancelQuery, strings.Join(pidsStr, " OR pid=")))
	if err != nil {
		return fmt.Errorf("init psql query: %w", err)
	}

	p.psqlQuery = psqlQuery
	return nil
}

func (p *PsqlCancel) Run() (int, []byte, error) {
	output, exitCode, err := p.psqlQuery.RunPsqlQueryCommand(nil)
	if err != nil {
		return exitCode, output, err
	}

	parsed := p.psqlQuery.BasePsql.ParseOutputJSON(output)
	canceledCount := len(parsed.S(psqlBase.DataRootJSONName).Children())

	retGc := gabs.New()
	_, _ = retGc.Set(canceledCount, "canceled")
	return exitCode, retGc.Bytes(), nil
}

func main() {
	step.Run(&PsqlCancel{})
}
