package raw

import (
	"github.com/Jeffail/gabs"
	"github.com/stackpulse/public-steps/common/env"
	"strings"
)

func Format(dataB []byte) (string, error) {
	ret := strings.Builder{}
	ret.WriteString(env.EndMarker())

	gc := gabs.New()
	_, _ = gc.Set(string(dataB), "output")

	ret.WriteString(gc.String())
	return ret.String(), nil
}
