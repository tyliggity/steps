package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKeyPairEnv(t *testing.T) {
	tests := []struct {
		name                string
		inputEnv            string
		shouldErr           bool
		expectedParsedValue map[string]string
	}{
		{"empty value", "", false, map[string]string{}},
		{"valid key pair", "key=value", false, map[string]string{"key": "value"}},
		{"valid multiple key pair", "key=value,key2=value2", false, map[string]string{"key": "value", "key2": "value2"}},
		{"valid key pair, extra additional ',' ", "key=value,", false, map[string]string{"key": "value"}},
		{"valid key pair, empty value", "key=", false, map[string]string{"key": ""}},
		{"key pair - empty key", "=value", false, map[string]string{"": "value"}},
		{"invalid key pair - no delimiter", "justkey", true, map[string]string{}},
		{"invalid key pair - multiple values, no delimiter", "key=valuekey2=value2", true, map[string]string{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parsedValue, err := ParseKeyValueEnv(test.inputEnv)
			if test.shouldErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.expectedParsedValue, parsedValue)
			}
		})
	}
}
