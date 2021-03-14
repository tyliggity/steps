package main

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/step"
	queryBase "github.com/stackpulse/public-steps/psql/query/base"
	"strings"
)

// This const should be processed by fmt.Sprintf
const longRunningSessionsQuery = `
SELECT  pid, pg_blocking_pids(pid) as blocked_by,
        state, application_name, usename,
        now() - xact_start as xact_age ,
        now() - query_start as age,
        wait_event_type, wait_event, client_port, client_addr, query
FROM    pg_stat_activity
WHERE   state != 'idle' AND EXTRACT(EPOCH FROM (now() - query_start) ) > %d AND 
        query NOT ILIKE '%%pg_stat_activity%%' %s
ORDER BY query_start asc;
`

type Args struct {
	QueryLike                   []string `env:"QUERY_LIKE"`
	ApplicationNameLike         []string `env:"APPLICATION_NAME_LIKE"`
	ApplicationNameEqual        []string `env:"APPLICATION_NAME_EQUAL"`
	ExcludeQueryLike            []string `env:"EXCLUDE_QUERY_LIKE"`
	ExcludeApplicationNameLike  []string `env:"EXCLUDE_APPLICATION_NAME_LIKE"`
	ExcludeApplicationNameEqual []string `env:"EXCLUDE_APPLICATION_NAME_EQUAL"`
	PIDsEqual                   []string `env:"PIDS_EQUAL"`
	ExcludePIDs                 []string `env:"EXCLUDE_PIDS"`
	LongDurationSeconds         int      `env:"LONG_DURATION_SECONDS" envDefault:"30"`
}

type PsqlLongRunning struct {
	psqlQuery *queryBase.PsqlQuery
}

func buildLike(sb *strings.Builder, negative bool, values []string, columnName string) {
	if len(values) == 0 {
		return
	}

	likeOperator := "LIKE"
	logicalOperator := "OR"
	if negative {
		likeOperator = "NOT LIKE"
		logicalOperator = "AND"
	}

	if sb.Len() > 0 {
		sb.WriteString(fmt.Sprintf(" %s ", logicalOperator))
	}
	sb.WriteString(fmt.Sprintf("%s %s ", columnName, likeOperator))
	sb.WriteString(strings.Join(values, fmt.Sprintf(" %s %s %s ", logicalOperator, columnName, likeOperator)))
}

func buildEqual(sb *strings.Builder, negative bool, values []string, columnName string) {
	if len(values) == 0 {
		return
	}

	equalOperator := "="
	logicalOperator := "OR"
	if negative {
		equalOperator = "!="
		logicalOperator = "AND"
	}

	if sb.Len() > 0 {
		sb.WriteString(fmt.Sprintf(" %s ", logicalOperator))
	}

	sb.WriteString(columnName + equalOperator)
	sb.WriteString(strings.Join(values, fmt.Sprintf(" %s %s", logicalOperator, columnName+equalOperator)))
}

func quoteFilters(args Args) {
	for i, v := range args.ApplicationNameEqual {
		args.ApplicationNameEqual[i] = pq.QuoteLiteral(v)
	}
	for i, v := range args.ApplicationNameLike {
		args.ApplicationNameLike[i] = pq.QuoteLiteral(v)
	}
	for i, v := range args.ExcludeApplicationNameEqual {
		args.ExcludeApplicationNameEqual[i] = pq.QuoteLiteral(v)
	}
	for i, v := range args.ExcludeApplicationNameLike {
		args.ExcludeApplicationNameLike[i] = pq.QuoteLiteral(v)
	}
	for i, v := range args.PIDsEqual {
		args.PIDsEqual[i] = pq.QuoteLiteral(v)
	}
	for i, v := range args.ExcludePIDs {
		args.ExcludePIDs[i] = pq.QuoteLiteral(v)
	}
	for i, v := range args.QueryLike {
		args.QueryLike[i] = pq.QuoteLiteral(v)
	}
	for i, v := range args.ExcludeQueryLike {
		args.ExcludeQueryLike[i] = pq.QuoteLiteral(v)
	}
}

func buildQuery(args Args) string {
	quoteFilters(args)

	var positiveWhere strings.Builder
	buildEqual(&positiveWhere, false, args.ApplicationNameEqual, "application_name")
	buildLike(&positiveWhere, false, args.ApplicationNameLike, "application_name")
	buildEqual(&positiveWhere, false, args.PIDsEqual, "pid")
	buildLike(&positiveWhere, false, args.QueryLike, "query")

	var negativeWhere strings.Builder
	buildEqual(&negativeWhere, true, args.ExcludeApplicationNameEqual, "application_name")
	buildLike(&negativeWhere, true, args.ExcludeApplicationNameLike, "application_name")
	buildEqual(&negativeWhere, true, args.ExcludePIDs, "pid")
	buildLike(&negativeWhere, true, args.ExcludeQueryLike, "query")

	where := ""
	if positiveWhere.Len() > 0 {
		where = fmt.Sprintf("AND (%s)", positiveWhere.String())
	}

	if negativeWhere.Len() > 0 {
		where = fmt.Sprintf("%s AND (%s)", where, negativeWhere.String())
	}

	return fmt.Sprintf(longRunningSessionsQuery, args.LongDurationSeconds, where)
}

func (p *PsqlLongRunning) Init() error {
	args := Args{}
	if err := env.Parse(&args); err != nil {
		return fmt.Errorf("parsing arguments: %w", err)
	}

	psqlQuery, err := queryBase.NewPsqlQuery(buildQuery(args))
	if err != nil {
		return fmt.Errorf("init psql query: %w", err)
	}

	p.psqlQuery = psqlQuery
	return nil
}

func (p *PsqlLongRunning) Run() (int, []byte, error) {
	output, exitCode, err := p.psqlQuery.RunPsqlQueryCommand(nil)
	if err != nil {
		return exitCode, output, err
	}
	return exitCode, p.psqlQuery.ParseOutput(output), nil
}

func main() {
	step.Run(&PsqlLongRunning{})
}
