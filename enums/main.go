// Support for declaring enums in Go
// and generating helpers (.go, .ts, .sql)
// All enums must be listed in enums.go
// with a comment for the label to display
package enums

import (
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

// EnumTable map type name to their values
type EnumTable map[string]EnumType

// AsLookupTable return a map with key TypeName.VarName
// and their resolved values.
// Strings are single quoted to be SQL compliant
func (t EnumTable) AsLookupTable() map[string]string {
	out := map[string]string{}
	for _, enum := range t {
		for _, value := range enum.Values {
			out[fmt.Sprintf("%s.%s", enum.Name, value.VarName)] = strings.ReplaceAll(value.Value, `"`, `'`)
		}
	}
	return out
}

// Lookup return the matching enum and it's related basic type,
// ok false if `ty` is not an enum.
func (t EnumTable) Lookup(ty *types.Named) (EnumType, *types.Basic, bool) {
	if enum, isEnum := t[ty.Obj().Name()]; isEnum {
		if basic, isBasic := ty.Underlying().(*types.Basic); isBasic {
			return enum, basic, true
		}
	}
	return EnumType{}, nil, false
}

type EnumType struct {
	Name   string
	Values []EnumValue
	IsInt  bool
}

// AsTuple returns a tuple of valid values
// compatible with SQL syntax.
// Ex: ('red', 'blue', 'green')
func (e EnumType) AsTuple() string {
	chunks := make([]string, len(e.Values))
	for i, val := range e.Values {
		chunks[i] = val.Value
	}
	out := fmt.Sprintf("(%s)", strings.Join(chunks, ", "))
	return strings.ReplaceAll(out, `"`, `'`)
}

// AsArray returns the code for a Go array
// containing all values.
func (e EnumType) AsArray() string {
	chunks := make([]string, len(e.Values))
	for i, val := range e.Values {
		chunks[i] = val.VarName
	}
	return fmt.Sprintf("[...]%s{%s}", e.Name, strings.Join(chunks, ", "))
}

// FetchEnums looks for an "enums.go" ending file in the package
// and all it's imports, restricted to the same "main" folder,
// which is <domain>/<org>/<main>, of `pa`.
// For example, if `pa` is github.com/gopher/lib/server,
// only the subpackages github.com/gopher/lib/... will be searched.
// Type with same local name will collide.
func FetchEnums(pa *packages.Package) (EnumTable, error) {
	chunks := strings.Split(pa.PkgPath, "/")
	var prefix string
	if len(chunks) >= 3 {
		prefix = strings.Join(chunks[:3], "/")
	}
	out := EnumTable{}
	err := fetchEnums(pa, out, prefix)
	return out, err
}
