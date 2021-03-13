package main

import (
	"fmt"
	"testing"
)

func Test_buildQuery(t *testing.T) {
	tests := []struct {
		name string
		args Args
		want string
	}{
		{
			name: "Empty",
			args: Args{},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, ""),
		},
		{
			name: "One equal",
			args: Args{ApplicationNameEqual: []string{"first_equal"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (application_name='first_equal')"),
		},
		{
			name: "Two equal",
			args: Args{ApplicationNameEqual: []string{"first_equal", "second_equal"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (application_name='first_equal' OR application_name='second_equal')"),
		},
		{
			name: "One like",
			args: Args{ApplicationNameLike: []string{"first_like%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (application_name LIKE 'first_like%')"),
		},
		{
			name: "Two like",
			args: Args{ApplicationNameLike: []string{"first_like%", "%second_like%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (application_name LIKE 'first_like%' OR application_name LIKE '%second_like%')"),
		},
		{
			name: "One equal one like",
			args: Args{ApplicationNameEqual: []string{"first_equal"}, ApplicationNameLike: []string{"first_like%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (application_name='first_equal' OR application_name LIKE 'first_like%')"),
		},
		{
			name: "Two equal one like",
			args: Args{ApplicationNameEqual: []string{"first_equal", "second_equal"}, ApplicationNameLike: []string{"first_like%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (application_name='first_equal' OR application_name='second_equal' OR application_name LIKE 'first_like%')"),
		},
		{
			name: "One equal two like",
			args: Args{ApplicationNameEqual: []string{"first_equal"}, ApplicationNameLike: []string{"first_like%", "%second_like%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (application_name='first_equal' OR application_name LIKE 'first_like%' OR application_name LIKE '%second_like%')"),
		},
		{
			name: "Two equal two like",
			args: Args{ApplicationNameEqual: []string{"first_equal", "second_equal"}, ApplicationNameLike: []string{"first_like%", "%second_like%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (application_name='first_equal' OR application_name='second_equal' OR application_name LIKE 'first_like%' OR application_name LIKE '%second_like%')"),
		},
		{
			name: "One negative equal",
			args: Args{ExcludeApplicationNameEqual: []string{"first_not_equal"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, " AND (application_name!='first_not_equal')"),
		},
		{
			name: "Two negative equal",
			args: Args{ExcludeApplicationNameEqual: []string{"first_not_equal", "second_not_equal"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, " AND (application_name!='first_not_equal' AND application_name!='second_not_equal')"),
		},
		{
			name: "One negative like",
			args: Args{ExcludeApplicationNameLike: []string{"%first_not_equal%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, " AND (application_name NOT LIKE '%first_not_equal%')"),
		},
		{
			name: "Two negative like",
			args: Args{ExcludeApplicationNameLike: []string{"%first_not_equal%", "second_not_equal%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, " AND (application_name NOT LIKE '%first_not_equal%' AND application_name NOT LIKE 'second_not_equal%')"),
		},
		{
			name: "Two negative equal two negative like",
			args: Args{ExcludeApplicationNameEqual: []string{"first_not_equal", "second_not_equal"}, ExcludeApplicationNameLike: []string{"%first_not_equal%", "second_not_equal%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, " AND (application_name!='first_not_equal' AND application_name!='second_not_equal' AND application_name NOT LIKE '%first_not_equal%' AND application_name NOT LIKE 'second_not_equal%')"),
		},
		{
			name: "Two negative equal two negative like, two equal, two Like",
			args: Args{ApplicationNameEqual: []string{"first_equal", "second_equal"}, ApplicationNameLike: []string{"first_like%", "%second_like%"}, ExcludeApplicationNameEqual: []string{"first_not_equal", "second_not_equal"}, ExcludeApplicationNameLike: []string{"%first_not_equal%", "second_not_equal%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (application_name='first_equal' OR application_name='second_equal' OR application_name LIKE 'first_like%' OR application_name LIKE '%second_like%') AND (application_name!='first_not_equal' AND application_name!='second_not_equal' AND application_name NOT LIKE '%first_not_equal%' AND application_name NOT LIKE 'second_not_equal%')"),
		},
		{
			name: "One pid equal",
			args: Args{PIDsEqual: []string{"123"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (pid='123')"),
		},
		{
			name: "Two pid equal",
			args: Args{PIDsEqual: []string{"123", "456"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (pid='123' OR pid='456')"),
		},
		{
			name: "One negative pid equal",
			args: Args{ExcludePIDs: []string{"123"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, " AND (pid!='123')"),
		},
		{
			name: "Two negative pid equal",
			args: Args{ExcludePIDs: []string{"123", "456"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, " AND (pid!='123' AND pid!='456')"),
		},
		{
			name: "Application name: Two negative equal two negative like, two equal, two Like. PIDs: two equal, two negative equal",
			args: Args{ApplicationNameEqual: []string{"first_equal", "second_equal"}, ApplicationNameLike: []string{"first_like%", "%second_like%"}, ExcludeApplicationNameEqual: []string{"first_not_equal", "second_not_equal"}, ExcludeApplicationNameLike: []string{"%first_not_equal%", "second_not_equal%"}, PIDsEqual: []string{"123", "456"}, ExcludePIDs: []string{"789", "101112"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (application_name='first_equal' OR application_name='second_equal' OR application_name LIKE 'first_like%' OR application_name LIKE '%second_like%' OR pid='123' OR pid='456') AND (application_name!='first_not_equal' AND application_name!='second_not_equal' AND application_name NOT LIKE '%first_not_equal%' AND application_name NOT LIKE 'second_not_equal%' AND pid!='789' AND pid!='101112')"),
		},
		{
			name: "One query like",
			args: Args{QueryLike: []string{"SELECT *%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (query LIKE 'SELECT *%')"),
		},
		{
			name: "Two query like",
			args: Args{QueryLike: []string{"SELECT *%", "%SELECT * FROM%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (query LIKE 'SELECT *%' OR query LIKE '%SELECT * FROM%')"),
		},
		{
			name: "One negative query like",
			args: Args{ExcludeQueryLike: []string{"SELECT *%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, " AND (query NOT LIKE 'SELECT *%')"),
		},
		{
			name: "Two negative query like",
			args: Args{ExcludeQueryLike: []string{"SELECT *%", "%SELECT * FROM%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, " AND (query NOT LIKE 'SELECT *%' AND query NOT LIKE '%SELECT * FROM%')"),
		},
		{
			name: "Application name: Two negative equal two negative like, two equal, two Like. PIDs: two equal, two negative equal. Query: two like, two negative like",
			args: Args{ApplicationNameEqual: []string{"first_equal", "second_equal"}, ApplicationNameLike: []string{"first_like%", "%second_like%"}, ExcludeApplicationNameEqual: []string{"first_not_equal", "second_not_equal"}, ExcludeApplicationNameLike: []string{"%first_not_equal%", "second_not_equal%"}, PIDsEqual: []string{"123", "456"}, ExcludePIDs: []string{"789", "101112"}, QueryLike: []string{"SELECT *%", "%SELECT * FROM%"}, ExcludeQueryLike: []string{"SELECT *%", "%SELECT * FROM%"}},
			want: fmt.Sprintf(longRunningSessionsQuery, 0, "AND (application_name='first_equal' OR application_name='second_equal' OR application_name LIKE 'first_like%' OR application_name LIKE '%second_like%' OR pid='123' OR pid='456' OR query LIKE 'SELECT *%' OR query LIKE '%SELECT * FROM%') AND (application_name!='first_not_equal' AND application_name!='second_not_equal' AND application_name NOT LIKE '%first_not_equal%' AND application_name NOT LIKE 'second_not_equal%' AND pid!='789' AND pid!='101112' AND query NOT LIKE 'SELECT *%' AND query NOT LIKE '%SELECT * FROM%')"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildQuery(tt.args); got != tt.want {
				t.Errorf("buildQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
