package gcs

import (
	"fmt"
	"net/url"
	"strings"
)

func Parse(gcsUrl string) (bucket, object string, err error) {
	parsedUrl, err := url.Parse(gcsUrl)
	if err != nil {
		return "", "", fmt.Errorf("can't parse gcs url: %w", err)
	}

	if parsedUrl.Scheme != "gs" {
		return "", "", fmt.Errorf("not gcs scheme")
	}

	return parsedUrl.Host, strings.Trim(parsedUrl.Path, "/"), nil
}
