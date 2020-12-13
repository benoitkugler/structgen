// This package regroups some utility functions shared
// by the handlers
package utils

import (
	"go/types"
	"reflect"
	"regexp"
	"strings"
)

// since we can't acces the underlying name of named types,
// we check against this string, to detect time.Time
const timeString = "struct{wall uint64; ext int64; loc *time.Location}"

var reComment = regexp.MustCompile(`^// (\w+):(.+)`)

// GetFieldName returns the name of the field
// as it should appear.
// If first checks for `tagId`, then for 'json', and defaults to the Go name.
// A non exported field is ignored from UPDATE and INSERT statement (but not from Scans)
// If the tag contains "-", it's ignored : empty string is returned.
func GetFieldName(field *types.Var, fullTag, tagId string) (string, bool) {
	// use tags ts, defaults to json, defaults to go name
	sTag := reflect.StructTag(fullTag)
	tag := sTag.Get(tagId)
	if tag == "" {
		tag = sTag.Get("json")
	}
	tag = strings.Split(tag, ",")[0] // remove omitempty

	if tag == "-" { // ignored
		return "", false
	}
	if tag == "" {
		tag = field.Name() // go name
	}
	return tag, field.Exported()
}

// IsUnderlyingTime returns `true` is the underlying type
// of `typ` is time.Time
func IsUnderlyingTime(typ types.Type) bool {
	return typ.Underlying().String() == timeString
}

// IsSpecialComment returns a non empty tag if the comment
// has a special form // <tag>:<content>
func IsSpecialComment(comment string) (tag, content string) {
	match := reComment.FindStringSubmatch(comment)
	if len(match) > 0 {
		return match[1], match[2]
	}
	return "", ""
}
