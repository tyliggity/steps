package json

import (
	"fmt"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stretchr/testify/assert"

	"os"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    string
		wantErr bool
	}{
		{
			name:    "map json in start",
			data:    `{"a":"b", "c":"d"}`,
			want:    env.EndMarker() + `{"a":"b", "c":"d"}`,
			wantErr: false,
		},
		{
			name:    "array json in start",
			data:    `[{"a":"b"}, {"c":"d"}]`,
			want:    env.EndMarker() + `{"output":[{"a":"b"}, {"c":"d"}]}`,
			wantErr: false,
		},
		{
			name:    "map json not in start",
			data:    `foo bar {"a":"b", "c":"d"}`,
			want:    env.EndMarker() + `{"a":"b", "c":"d"}`,
			wantErr: false,
		},
		{
			name:    "map json with not parsable json string ahead",
			data:    `foo {bar} {"a":"b", "c":"d"}`,
			want:    env.EndMarker() + `{"a":"b", "c":"d"}`,
			wantErr: false,
		},
		{
			name:    "array json with not parsable json strings ahead",
			data:    `{foo} {bar} [{"a":"b"}, {"c":"d"}]`,
			want:    env.EndMarker() + `{"output":[{"a":"b"}, {"c":"d"}]}`,
			wantErr: false,
		},
		{
			name:    "no json",
			data:    `foo bar`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "no parsable json",
			data:    `foo {bar}`,
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Format() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatAsJSONOutput(t *testing.T) {
	output := "foo"
	x, err := FormatAsJSONOutput([]byte(output))
	if err != nil {
		t.Error("FormatAsJSONOutput failed but it shouldnt")
		return
	}

	assert.Equal(t, fmt.Sprintf(`{"%s":"%s"}`, outputKey, output), string(x))
}

func TestMain(m *testing.M) {
	os.Setenv("SP_DEBUG", "true")
	m.Run()
}
