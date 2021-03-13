package base

import (
	"regexp"
	"strconv"
	"strings"
)

// Those characters are not considered as characters that needs escaping, so if the regex match it mean we DO need escaping
// (because the regex search for any characters NOT from those characters)
var notEscapedRe = regexp.MustCompile(`[^\w@%+=:,./-]`)

func quote(arg string) string {
	if notEscapedRe.MatchString(arg) {
		return strconv.Quote(arg)
	}
	return arg
}

func joinCommandLines(args []string) string {
	newArgs := make([]string, len(args))
	for i, arg := range args {
		newArgs[i] = quote(arg)
	}
	return strings.Join(newArgs, " ")
}
