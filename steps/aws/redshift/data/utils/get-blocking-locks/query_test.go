package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryWhereClause(t *testing.T) {
	tests := []struct {
		query    *query
		expected string
		wantErr  bool
	}{
		{
			query:    &query{},
			expected: "AND block_sec>60 ",
		},
		{
			query: &query{
				onlyUsers: []string{"user1"},
			},
			expected: `AND username IN ('user1') AND block_sec>60 `,
		},
		{
			query: &query{
				onlyUsers: []string{"user1", "user2"},
			},
			expected: `AND username IN ('user1', 'user2') AND block_sec>60 `,
		},
		{
			query: &query{
				onlyUsers: []string{"user1", "user2"},
			},
			expected: `AND username IN ('user1', 'user2') AND block_sec>60 `,
		},
		{
			query: &query{
				onlyUsers:    []string{"user1", "user2"},
				ignoredUsers: []string{"user1", "user2"},
			},
			wantErr: true,
		},
		{
			query: &query{
				onlyUsers:       []string{"user1", "user2"},
				blockingSeconds: 1500,
			},
			expected: `AND username IN ('user1', 'user2') AND block_sec>1500 `,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := test.query.whereClause()
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expected, got)
		})
	}
}

func TestQueryGenerate(t *testing.T) {
	tests := []struct {
		query    *query
		expected string
		wantErr  bool
	}{
		{
			query: &query{
				onlyUsers:       []string{"datauser"},
				blockingSeconds: 1500,
			},
			expected: `SELECT
    a.pid,
    a.xid,
    a.pidlist,
    a.username,
    a.block_sec,
    a.max_sec_blocking,
    a.num_blocking,
    b.querytxt
FROM admin.v_get_blocking_locks a
LEFT JOIN stl_query b on b.xid=a.xid
WHERE num_blocking>0 AND username IN ('datauser') AND block_sec>1500 
GROUP BY 1,2,3,4,5,6,7,8
ORDER BY a.block_sec asc`,
		},
		{
			query: &query{
				ignoredUsers:    []string{"datauser"},
				blockingSeconds: 1500,
			},
			expected: `SELECT
    a.pid,
    a.xid,
    a.pidlist,
    a.username,
    a.block_sec,
    a.max_sec_blocking,
    a.num_blocking,
    b.querytxt
FROM admin.v_get_blocking_locks a
LEFT JOIN stl_query b on b.xid=a.xid
WHERE num_blocking>0 AND username NOT IN ('datauser') AND block_sec>1500 
GROUP BY 1,2,3,4,5,6,7,8
ORDER BY a.block_sec asc`,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := test.query.generate()
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expected, got)
		})
	}

}
