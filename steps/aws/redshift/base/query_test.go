package base

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildQuery(t *testing.T) {
	tests := []struct {
		pids     []string
		function string
		expected string
		wantErr  bool
	}{
		{
			pids:     []string{"1"},
			function: "pg_cancel_backend",
			expected: "SELECT pg_cancel_backend(1) as pid_1",
		},
		{
			pids:     []string{"1", "2", "3"},
			function: "pg_cancel_backend",
			expected: "SELECT pg_cancel_backend(1) as pid_1,pg_cancel_backend(2) as pid_2,pg_cancel_backend(3) as pid_3",
		},
		{
			pids:     []string{"", "2", "3"},
			function: "pg_cancel_backend",
			expected: "SELECT pg_cancel_backend(2) as pid_2,pg_cancel_backend(3) as pid_3",
		},
		{
			pids:     []string{"", "2", ""},
			function: "pg_cancel_backend",
			expected: "SELECT pg_cancel_backend(2) as pid_2",
		},
		{
			pids:     []string{"1", "-2", "3"},
			function: "pg_cancel_backend",
			expected: "SELECT pg_cancel_backend(1) as pid_1,pg_cancel_backend(3) as pid_3",
		},
		{
			pids:    []string{"-2"},
			wantErr: true,
		},
		{
			pids:     []string{"1", "0.2", "3"},
			function: "pg_cancel_backend",
			expected: "SELECT pg_cancel_backend(1) as pid_1,pg_cancel_backend(3) as pid_3",
		},
		{
			pids:    nil,
			wantErr: true,
		},
		{
			pids:    []string{},
			wantErr: true,
		},
		{
			pids:    []string{" "},
			wantErr: true,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := BuildPidFuncQuery(test.pids, test.function)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expected, got)
		})

	}
}
