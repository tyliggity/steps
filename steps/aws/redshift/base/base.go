package base

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/avast/retry-go"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"
	"github.com/olekukonko/tablewriter"
	"github.com/stackpulse/public-steps/common/log"
)

var ErrTimeout = fmt.Errorf("timeout")

// QueryRunError represents the query execution failure
type QueryRunError struct {
	Reason string
	Status string
}

func (e *QueryRunError) Error() string {
	return fmt.Sprintf("query run failed due to '%s' status: %s", e.Reason, e.Status)
}

// Using output struct here is a bit problematic because the output should return
// as JSON array of unknown fields (according to the given query)
const (
	resultsOutputName     = "results"
	delayBeforeFirstFetch = 5 * time.Second
)

type Args struct {
	Cluster   string        `env:"CLUSTER,required"`
	Region    string        `env:"REGION,required"`
	Database  string        `env:"DATABASE,required"`
	Timeout   time.Duration `env:"TIMEOUT" envDefault:"0s"`
	User      string        `env:"USER"`
	SecretArn string        `env:"SECRET_ARN"`
}

func (a Args) Validate() error {
	if (a.User == "" && a.SecretArn == "") || (a.User != "" && a.SecretArn != "") {
		return fmt.Errorf("either user or secret arn should be supplied")
	}

	return nil
}

type RedshiftAWSRunner struct {
	Args
	service *redshiftdata.Client
}

type RedshiftQueryResults struct {
	Headers []string
	Values  [][]interface{}
}

func (r *RedshiftAWSRunner) ExecuteStatement(ctx context.Context, sql string) (string, error) {
	var secret *string
	if r.SecretArn != "" {
		secret = aws.String(r.SecretArn)
	}

	executeOut, err := r.service.ExecuteStatement(ctx, &redshiftdata.ExecuteStatementInput{
		ClusterIdentifier: aws.String(r.Cluster),
		Database:          aws.String(r.Database),
		DbUser:            aws.String(r.User),
		SecretArn:         secret,
		Sql:               aws.String(sql),
	})
	if err != nil {
		return "", fmt.Errorf("execute statement: %w", err)
	}

	if executeOut.Id == nil || *executeOut.Id == "" {
		return "", fmt.Errorf("no execution id")
	}

	return *executeOut.Id, nil
}

func (r *RedshiftAWSRunner) FetchResult(ctx context.Context, id string) (*RedshiftQueryResults, error) {
	const (
		attempts = 5
		delay    = 3 * time.Second
	)
	var zeroDuration time.Duration
	var timeoutAt time.Time
	if r.Timeout > zeroDuration {
		timeoutAt = time.Now().Add(r.Timeout)
	}

	var statementResultsOut *redshiftdata.GetStatementResultOutput
	err := retry.Do(func() error {
		var err error
		statementResultsOut, err = r.service.GetStatementResult(ctx, &redshiftdata.GetStatementResultInput{
			Id: aws.String(id),
		})

		if !timeoutAt.IsZero() && time.Now().After(timeoutAt) {
			return ErrTimeout
		}
		if err == nil {
			return nil
		}

		// catch failed executions to abort redundant retries
		out, descErr := r.service.DescribeStatement(ctx, &redshiftdata.DescribeStatementInput{Id: aws.String(id)})
		if descErr == nil {
			if out.Status == types.StatusStringAborted || out.Status == types.StatusStringFailed {
				var reason string
				if out.Error != nil {
					reason = *out.Error
				}
				return &QueryRunError{
					Reason: reason,
					Status: string(out.Status),
				}
			}
		}

		return err
	},
		retry.RetryIf(func(err error) bool {
			if errors.Is(err, ErrTimeout) {
				return false
			}

			var qerr *QueryRunError
			if errors.As(err, &qerr) {
				return false
			}

			return true
		}),
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(delay),
		retry.Attempts(attempts),
		retry.OnRetry(func(n uint, err error) {
			log.Debug("Retry get statement result (%d/%d) error: %v", n, attempts, err)
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("get statement error: %w", err)
	}
	if statementResultsOut == nil {
		return nil, fmt.Errorf("get statement output is nil")
	}

	numColumns := len(statementResultsOut.ColumnMetadata)
	results := &RedshiftQueryResults{}
	results.Headers = make([]string, numColumns)
	for i, c := range statementResultsOut.ColumnMetadata {
		results.Headers[i] = *c.Name
	}

	values := make([][]interface{}, statementResultsOut.TotalNumRows)

	// Note: the 'Value' field is not accessible via the Field interface, but every field has 'Value' member.
	out, err := json.Marshal(statementResultsOut.Records)
	if err != nil {
		return nil, fmt.Errorf("marshal redshift records: %w", err)
	}

	var rsValues [][]struct {
		Value interface{} `json:"Value"`
	}

	err = json.Unmarshal(out, &rsValues)
	if err != nil {
		return nil, fmt.Errorf("unmarshal redshift records: %w", err)
	}

	for i, r := range rsValues {
		values[i] = make([]interface{}, len(r))
		for j, field := range r {
			value := strings.Trim(fmt.Sprintf("%v", field.Value), " ")
			values[i][j] = value
		}
	}

	results.Values = values

	return results, nil
}

func (r *RedshiftAWSRunner) RunQuery(sql string) (*RedshiftQueryResults, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(r.Region))
	if err != nil {
		return nil, fmt.Errorf("load default config for region '%s': %w", r.Region, err)
	}

	r.service = redshiftdata.NewFromConfig(cfg)
	id, err := r.ExecuteStatement(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("execute statement: %w", err)
	}

	log.Debug("Execute id: %s", id)

	// most of the times, immediate request to the get result would fail so we give it a bit of slack before first try
	time.Sleep(delayBeforeFirstFetch)

	results, err := r.FetchResult(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetch result: %w", err)
	}

	return results, nil
}

func (q *RedshiftQueryResults) JSON() (*gabs.Container, error) {
	root, _ := gabs.New().Array()

	for j, row := range q.Values {
		obj := gabs.New()
		for i, value := range row {
			_, err := obj.Set(value, q.Headers[i])
			if err != nil {
				log.Debug("Failed to set json result key '%s' to '%v'", q.Headers[i], value)
				continue
			}
		}
		err := root.ArrayAppend(obj)
		if err != nil {
			log.Debug("Failed to append object to result array line: #%d", j)
		}
	}

	return root, nil
}

func (q *RedshiftQueryResults) Text() string {
	builder := &strings.Builder{}
	table := tablewriter.NewWriter(builder)
	table.SetHeader(q.Headers)

	for _, v := range q.Values {
		row := make([]string, len(v))
		for i := range row {
			row[i] = fmt.Sprint(v[i])
		}
		table.Append(row)
	}

	table.Render()

	return builder.String()
}

func (q *RedshiftQueryResults) ResultsOutput() ([]byte, error) {
	results, err := q.JSON()
	if err != nil {
		return nil, fmt.Errorf("json results: %w", err)
	}

	output := gabs.New()
	_, err = output.Set(results, resultsOutputName)
	if err != nil {
		return nil, fmt.Errorf("set results json array key: %w", err)
	}

	return output.Bytes(), nil
}
