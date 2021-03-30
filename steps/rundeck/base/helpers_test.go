package base

import (
	"fmt"
	"testing"

	"github.com/stackpulse/steps-sdk-go/testers"
	"github.com/stretchr/testify/assert"
)

func Test_jsonKeysCaseSerializer(t *testing.T) {
	tests := []struct {
		name           string
		inputJson      string
		expectedErr    error
		expectedOutput string
	}{
		{"empty obj", `{}`, nil, `{}`},
		{"empty Array", `[]`, nil, `[]`},
		{"no changes", `{"hello":"holla"}`, nil, `{"hello":"holla"}`},
		{"no changes 2 words", `{"hello_shalom":"holla"}`, nil, `{"hello_shalom":"holla"}`},
		{"camel case 1 word", `{"Hello":"holla"}`, nil, `{"hello":"holla"}`},
		{"camel case 2 words", `{"HelloShalom":"Holla"}`, nil, `{"hello_shalom":"Holla"}`},
		{"kebab case 2 words", `{"helloShalom":"Holla"}`, nil, `{"hello_shalom":"Holla"}`},
		{"obj with empty value", `{"helloShalom":""}`, nil, `{"hello_shalom":""}`},
		{"obj with multiple diffrent keys", `{"fun_fun_fun":"ok","HelloShalom":["gello","dwa"],"shalomHello":"Holla"}`, nil, `{"fun_fun_fun":"ok","hello_shalom":["gello","dwa"],"shalom_hello":"Holla"}`},
		{"obj in array in obj", `{"HelloShalom":[{"helloShalom":"Holla"},{"helloShalom":"Holla"},{"hello_shalom":"Holla"}]}`, nil, `{"hello_shalom":[{"hello_shalom":"Holla"},{"hello_shalom":"Holla"},{"hello_shalom":"Holla"}]}`},
		{"array in obj in array", `[{"hello_shalom":[{"HelloShalom":"Holla"}]}]`, nil, `[{"hello_shalom":[{"hello_shalom":"Holla"}]}]`},
		{"invalid json", "key=value,key2=value2", fmt.Errorf("looking for beginning of value"), "key=value,key2=value2"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parsedValue, err := jsonKeysCaseSerializer([]byte(test.inputJson))
			testers.ValidateError(t, test.expectedErr, err)
			assert.EqualValues(t, test.expectedOutput, string(parsedValue))
		})
	}
}
