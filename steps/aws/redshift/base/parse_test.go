package base

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePidQueryResults(t *testing.T) {
	tests := []struct {
		results RedshiftQueryResults
		success []string
		failed  []string
		wantErr bool
	}{
		{
			results: RedshiftQueryResults{
				Headers: []string{"pid_16125", "pid_16126"},
				Values:  [][]interface{}{{"0", "1"}},
			},
			success: []string{"16126"},
			failed:  []string{"16125"},
		},
		{
			results: RedshiftQueryResults{
				Headers: []string{"pid_16125", "pid_16126"},
				Values:  [][]interface{}{{"1", "1"}},
			},
			success: []string{"16125", "16126"},
			failed:  []string{},
		},
		{
			results: RedshiftQueryResults{
				Headers: []string{"pid_16125", "pid_16126"},
				Values:  [][]interface{}{{"0", "0"}},
			},
			success: []string{},
			failed:  []string{"16125", "16126"},
		},
		{
			results: RedshiftQueryResults{
				Headers: []string{"pid_16125", "pid_16126", "pid_16126"},
				Values:  [][]interface{}{{"0", "0", "0"}},
			},
			wantErr: true,
		},
		{
			results: RedshiftQueryResults{
				Headers: []string{"pid_16125", "pid_16126", "pid_16126"},
				Values:  [][]interface{}{{"0", "0"}},
			},
			wantErr: true,
		},
		{
			results: RedshiftQueryResults{
				Headers: []string{"pid_16125", "pid_16126"},
				Values:  [][]interface{}{{"0", "0", "0"}},
			},
			wantErr: true,
		},
		{
			results: RedshiftQueryResults{
				Headers: []string{"pid_16125", "pid_16126"},
				Values:  [][]interface{}{{"1", "non-1-value"}},
			},
			success: []string{"16125"},
			failed:  []string{"16126"},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			success, failed, err := ParsePidQueryResults(&test.results)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.EqualValues(t, test.success, success)
			assert.EqualValues(t, test.failed, failed)
		})
	}

}
