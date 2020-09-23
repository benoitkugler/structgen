package orm

import (
	"regexp"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func toLowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

var FnMap = template.FuncMap{
	"inc":     func(i int) int { return i + 1 },
	"snake":   toSnakeCase,
	"varname": toLowerFirst,
}
