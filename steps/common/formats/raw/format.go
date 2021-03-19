package raw

import (
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/stackpulse/steps-sdk-go/env"
)

func Format(dataB []byte) (string, error) {
	ret := strings.Builder{}
	ret.WriteString(env.EndMarker())

	gc := gabs.New()
	_, _ = gc.Set(string(dataB), "output")

	ret.WriteString(gc.String())
	return ret.String(), nil
}
