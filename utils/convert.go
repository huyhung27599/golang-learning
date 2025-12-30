package utils

import (
	"regexp"
	"strings"
)

func camelToSnake(camel string) string {
	return strings.ToLower(regexp.MustCompile("([A-Z])").ReplaceAllString(camel, "_$1"))
}