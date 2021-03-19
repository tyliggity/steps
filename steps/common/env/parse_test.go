package env

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ItemCheck interface {
	Check(t *testing.T, res interface{})
}

type stringItem struct {
	Input []string `env:"INPUT,required"`
}

func (i *stringItem) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type intItem struct {
	Input []int `env:"INPUT,required"`
}

func (i *intItem) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type int8Item struct {
	Input []int8 `env:"INPUT,required"`
}

func (i *int8Item) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type int16Item struct {
	Input []int16 `env:"INPUT,required"`
}

func (i *int16Item) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type int32Item struct {
	Input []int32 `env:"INPUT,required"`
}

func (i *int32Item) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type int64Item struct {
	Input []int64 `env:"INPUT,required"`
}

func (i *int64Item) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type uintItem struct {
	Input []uint `env:"INPUT,required"`
}

func (i *uintItem) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type uint8Item struct {
	Input []uint8 `env:"INPUT,required"`
}

func (i *uint8Item) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type uint16Item struct {
	Input []uint16 `env:"INPUT,required"`
}

func (i *uint16Item) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type uint32Item struct {
	Input []uint32 `env:"INPUT,required"`
}

func (i *uint32Item) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type uint64Item struct {
	Input []uint64 `env:"INPUT,required"`
}

func (i *uint64Item) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type float32Item struct {
	Input []float32 `env:"INPUT,required"`
}

func (i *float32Item) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type float64Item struct {
	Input []float64 `env:"INPUT,required"`
}

func (i *float64Item) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type boolItem struct {
	Input []bool `env:"INPUT,required"`
}

func (i *boolItem) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type durationItem struct {
	Input []time.Duration `env:"INPUT,required"`
}

func (i *durationItem) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

type urlItem struct {
	Input []url.URL `env:"INPUT,required"`
}

func (i *urlItem) Check(t *testing.T, res interface{}) {
	assert.Equal(t, res, i.Input)
}

func ValidateError(t *testing.T, wantErr error, err error) bool {
	t.Helper()

	if wantErr != nil {
		require.Error(t, err)
		assert.Contains(t, err.Error(), wantErr.Error())
		return false
	}
	require.NoError(t, err)
	return true
}

func TestParse(t *testing.T) {
	sp := "https://stackpulse.io"
	spURL, err := url.Parse(sp)
	require.NoError(t, err)

	google := "http://google.com"
	googleURL, err := url.Parse(google)
	require.NoError(t, err)

	tests := []struct {
		name    string
		item    ItemCheck
		input   string
		res     interface{}
		wantErr error
	}{
		{
			name:  "Srting array list",
			input: "1,2, abc",
			res:   []string{"1", "2", " abc"},
			item:  &stringItem{},
		},
		{
			name:  "Srting array json",
			input: `["1","2","abc"]`,
			res:   []string{"1", "2", "abc"},
			item:  &stringItem{},
		},
		{
			name:  "Int array json",
			input: `[1,2,3]`,
			res:   []int{1, 2, 3},
			item:  &intItem{},
		},
		{
			name:  "Int array list",
			input: `1,2, 3`,
			res:   []int{1, 2, 3},
			item:  &intItem{},
		},
		{
			name:  "Int single json",
			input: `[1]`,
			res:   []int{1},
			item:  &intItem{},
		},
		{
			name:  "Int single not json",
			input: `1`,
			res:   []int{1},
			item:  &intItem{},
		},
		{
			name:    "Int array json error",
			input:   `["1",2,3]`,
			res:     nil,
			item:    &intItem{},
			wantErr: fmt.Errorf(`[]int": wrong type`),
		},
		{
			name:    "Int array list error",
			input:   `1,2,a`,
			res:     nil,
			item:    &intItem{},
			wantErr: fmt.Errorf("strconv.ParseInt"),
		},
		{
			name:  "Int8 array json",
			input: `[1,2,3]`,
			res:   []int8{1, 2, 3},
			item:  &int8Item{},
		},
		{
			name:  "Int8 array list",
			input: `1,2, 3`,
			res:   []int8{1, 2, 3},
			item:  &int8Item{},
		},
		{
			name:  "Int16 array json",
			input: `[1,2, 3]`,
			res:   []int16{1, 2, 3},
			item:  &int16Item{},
		},
		{
			name:  "Int16 array list",
			input: `1,2, 3`,
			res:   []int16{1, 2, 3},
			item:  &int16Item{},
		},
		{
			name:  "Int32 array json",
			input: `[1,2,3]`,
			res:   []int32{1, 2, 3},
			item:  &int32Item{},
		},
		{
			name:  "Int32 array list",
			input: `1,2, 3`,
			res:   []int32{1, 2, 3},
			item:  &int32Item{},
		},
		{
			name:  "Int64 array json",
			input: `[1,2,3]`,
			res:   []int64{1, 2, 3},
			item:  &int64Item{},
		},
		{
			name:  "Int64 array list",
			input: `1,2, 3`,
			res:   []int64{1, 2, 3},
			item:  &int64Item{},
		},
		{
			name:  "Uint array list",
			input: `1,2, 3`,
			res:   []uint{1, 2, 3},
			item:  &uintItem{},
		},
		{
			name:    "Uint array json error",
			input:   `["1",2,3]`,
			res:     nil,
			item:    &uintItem{},
			wantErr: fmt.Errorf(`[]uint": wrong type`),
		},
		{
			name:    "Uint array list error",
			input:   `1,2,a`,
			res:     nil,
			item:    &uintItem{},
			wantErr: fmt.Errorf("strconv.ParseUint"),
		},
		{
			name:  "Uint8 array json",
			input: `[1,2,3]`,
			res:   []uint8{1, 2, 3},
			item:  &uint8Item{},
		},
		{
			name:  "Uint8 array list",
			input: `1,2, 3`,
			res:   []uint8{1, 2, 3},
			item:  &uint8Item{},
		},
		{
			name:  "Uint16 array json",
			input: `[1,2, 3]`,
			res:   []uint16{1, 2, 3},
			item:  &uint16Item{},
		},
		{
			name:  "Uint16 array list",
			input: `1,2, 3`,
			res:   []uint16{1, 2, 3},
			item:  &uint16Item{},
		},
		{
			name:  "Uint32 array json",
			input: `[1,2,3]`,
			res:   []uint32{1, 2, 3},
			item:  &uint32Item{},
		},
		{
			name:  "uint32 array list",
			input: `1,2, 3`,
			res:   []uint32{1, 2, 3},
			item:  &uint32Item{},
		},
		{
			name:  "uint64 array json",
			input: `[1,2,3]`,
			res:   []uint64{1, 2, 3},
			item:  &uint64Item{},
		},
		{
			name:  "Uint64 array list",
			input: `1,2, 3`,
			res:   []uint64{1, 2, 3},
			item:  &uint64Item{},
		},
		{
			name:  "Float32 array json",
			input: `[1.2,2.3,3.4]`,
			res:   []float32{1.2, 2.3, 3.4},
			item:  &float32Item{},
		},
		{
			name:  "Float32 array list",
			input: `1.2,2.3, 3.4`,
			res:   []float32{1.2, 2.3, 3.4},
			item:  &float32Item{},
		},
		{
			name:    "Float32 array json error",
			input:   `["1.2",2.3,3.4]`,
			res:     nil,
			item:    &float32Item{},
			wantErr: fmt.Errorf(`[]float32": wrong type`),
		},
		{
			name:    "Float32 array list error",
			input:   `1.2,2.3,a`,
			res:     nil,
			item:    &float32Item{},
			wantErr: fmt.Errorf("strconv.ParseFloat"),
		},
		{
			name:  "Float64 array json",
			input: `[1.2,2.3,3.4]`,
			res:   []float64{1.2, 2.3, 3.4},
			item:  &float64Item{},
		},
		{
			name:  "Float64 array list",
			input: `1.2,2.3, 3.4`,
			res:   []float64{1.2, 2.3, 3.4},
			item:  &float64Item{},
		},
		{
			name:  "Bool array json",
			input: `[true,false]`,
			res:   []bool{true, false},
			item:  &boolItem{},
		},
		{
			name:  "Bool array list",
			input: `true, false`,
			res:   []bool{true, false},
			item:  &boolItem{},
		},
		{
			name:    "Bool array json error",
			input:   `["true",false]`,
			res:     nil,
			item:    &boolItem{},
			wantErr: fmt.Errorf(`[]bool": wrong type`),
		},
		{
			name:    "Bool array list error",
			input:   `foo, false`,
			res:     nil,
			item:    &boolItem{},
			wantErr: fmt.Errorf("strconv.ParseBool"),
		},
		{
			name:  "Duration array json",
			input: `["5s", "4h"]`,
			res:   []time.Duration{5 * time.Second, 4 * time.Hour},
			item:  &durationItem{},
		},
		{
			name:  "Duration array list",
			input: `5s, 4h`,
			res:   []time.Duration{5 * time.Second, 4 * time.Hour},
			item:  &durationItem{},
		},
		{
			name:    "Duration array json error",
			input:   `["5a","5s"]`,
			res:     nil,
			item:    &durationItem{},
			wantErr: fmt.Errorf(`[]time.Duration": parsing type`),
		},
		{
			name:    "Duration array list error",
			input:   `5a, 5s`,
			res:     nil,
			item:    &durationItem{},
			wantErr: fmt.Errorf("unable to parse duration"),
		},
		{
			name:  "URL array json",
			input: fmt.Sprintf(`["%s", "%s"]`, sp, google),
			res:   []url.URL{*spURL, *googleURL},
			item:  &urlItem{},
		},
		{
			name:  "URL array list",
			input: fmt.Sprintf(`%s, %s`, sp, google),
			res:   []url.URL{*spURL, *googleURL},
			item:  &urlItem{},
		},
		{
			name:    "URL array json error",
			input:   fmt.Sprintf(`["%s", "%s"]`, "\r", google),
			res:     nil,
			item:    &urlItem{},
			wantErr: fmt.Errorf("unable to parse URL"),
		},
		{
			name:    "URL array list error",
			input:   fmt.Sprintf(`%s, %s`, "\a", google),
			res:     nil,
			item:    &urlItem{},
			wantErr: fmt.Errorf("unable to parse URL"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.Setenv("INPUT", tc.input))

			err := Parse(tc.item)
			if !ValidateError(t, tc.wantErr, err) {
				return
			}

			tc.item.Check(t, tc.res)
		})
	}
}
