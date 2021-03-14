package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitLines(t *testing.T) {
	testText := `this
is
test
`
	splitText := SplitLines(testText)
	assert.Len(t, splitText, 3)
	assert.Equal(t, []string{"this", "is", "test"}, splitText)
}

func TestParseTable(t *testing.T) {

	table := `NAME                                                   CDS        LDS        EDS        RDS          ISTIOD                     VERSION
details-v1-79c697d759-fh8qj.default                    SYNCED     SYNCED     SYNCED     SYNCED       istiod-d7464f9db-2htf7     1.7.2
istio-egressgateway-c999fffd6-xfjqn.istio-system       SYNCED     SYNCED     SYNCED     NOT SENT     istiod-d7464f9db-2htf7     1.7.2
istio-ingressgateway-7c56f5b4b7-d4fxb.istio-system     SYNCED     SYNCED     SYNCED     NOT SENT     istiod-d7464f9db-2htf7     1.7.2
productpage-v1-65576bb7bf-22hpg.default                SYNCED     SYNCED     SYNCED     SYNCED       istiod-d7464f9db-2htf7     1.7.2
ratings-v1-7d99676f7f-bgvhb.default                    SYNCED     SYNCED     SYNCED     SYNCED       istiod-d7464f9db-2htf7     1.7.2
reviews-v1-987d495c-2nz2h.default                      SYNCED     SYNCED     SYNCED     SYNCED       istiod-d7464f9db-2htf7     1.7.2
reviews-v2-6c5bf657cf-zxmp7.default                    SYNCED     SYNCED     SYNCED     SYNCED       istiod-d7464f9db-2htf7     1.7.2
reviews-v3-5f7b9f4f77-mmh49.default                    SYNCED     SYNCED     SYNCED     SYNCED       istiod-d7464f9db-2htf7     1.7.2`

	tbl, err := ParseTable(table)
	assert.NoError(t, err)
	assert.Len(t, tbl, 8)

	expectedNames := []string{
		"details-v1-79c697d759-fh8qj.default",
		"istio-egressgateway-c999fffd6-xfjqn.istio-system",
		"istio-ingressgateway-7c56f5b4b7-d4fxb.istio-system",
		"productpage-v1-65576bb7bf-22hpg.default",
		"ratings-v1-7d99676f7f-bgvhb.default",
		"reviews-v1-987d495c-2nz2h.default",
		"reviews-v2-6c5bf657cf-zxmp7.default",
		"reviews-v3-5f7b9f4f77-mmh49.default",
	}

	for i, name := range expectedNames {
		assert.Equal(t, name, tbl[i]["NAME"])
	}

}
