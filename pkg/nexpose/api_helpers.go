package nexpose

import (
	"net/url"
	"strings"
)

func encode(value string) string {

	value = url.QueryEscape(value)
	value = strings.Replace(value, ".", "%2E", -1)

	return value
}
