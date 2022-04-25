package enums

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

// EnumValue decribe the value of one item in an enumeration type.
type EnumValue struct {
	TypeName string // the "parent" type
	VarName  string // the variable name of the item
	Value    string // the value, as Go Code, like 1 or "abc"
	Label    string // how to display the value
}

type enumLabels map[string]string // varName -> text

type enumsValue map[string]enumLabels // local type name -> datas

func (e enumsValue) add(typeName, varName, label string) {
	ts := e[typeName]
	if ts == nil {
		ts = enumLabels{}
	}
	ts[varName] = label
	e[typeName] = ts
}

// return the label values
// go/types doesn't include comments
func parse(enumFile *ast.File, f *token.FileSet) (enumsValue, error) {
	enums := enumsValue{}
	for _, decl := range enumFile.Decls {
		stm, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		var lastTypeName string
		for _, spec := range stm.Specs {
			s, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			varName := s.Names[0].String() // variable name to use in code

			if !token.IsExported(varName) { // ignore private variables
				continue
			}

			if s.Comment == nil {
				return nil, fmt.Errorf("value as comment expected at %s", f.Position(s.Pos()))
			}
			text := strings.TrimSpace(s.Comment.Text()) // label to display

			var typeName string
			if s.Type == nil {
				// for example in iotas : use last type name
				typeName = lastTypeName
			} else {
				typeName = s.Type.(*ast.Ident).String()
				lastTypeName = typeName
			}
			enums.add(typeName, varName, text)
		}
	}
	return enums, nil
}

// fetch constants value, as computed by Go
func aggregate(pa *packages.Package, enums enumsValue) EnumTable {
	out := EnumTable{}
	for _, name := range pa.Types.Scope().Names() {
		obj := pa.Types.Scope().Lookup(name)
		if decl, isConst := obj.(*types.Const); isConst {
			varName := decl.Name()
			if !decl.Exported() {
				continue
			}

			constVal := decl.Val().String()
			if named, isNamed := decl.Type().(*types.Named); isNamed {
				typeName := named.Obj().Name()
				isInt := false
				if basic, isBasic := named.Underlying().(*types.Basic); isBasic {
					switch basic.Kind() { // unsigned & int
					case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
						isInt = true
					}
				}
				// restriction to previsouly fetch enums
				if vals, isIn := enums[typeName]; isIn {
					dt := out[typeName]
					dt.Name, dt.IsInt = typeName, isInt // needed if first lookup
					dt.Values = append(dt.Values, EnumValue{
						TypeName: typeName,
						VarName:  varName,
						Value:    constVal,
						Label:    vals[varName],
					})
					out[typeName] = dt
				}
			}
		}
	}
	return out
}

func fetchEnums(pa *packages.Package, accu EnumTable, prefix string) error {
	for i, file := range pa.GoFiles {
		if strings.HasSuffix(file, "enums.go") {
			a := pa.Syntax[i]
			firstMap, err := parse(a, pa.Fset)
			if err != nil {
				return err
			}
			tmp := aggregate(pa, firstMap)
			for k, v := range tmp {
				accu[k] = v
			}
		}
	}
	for _, imp := range pa.Imports {
		ignore := prefix != "" && !strings.HasPrefix(imp.PkgPath, prefix)
		if ignore {
			continue
		}
		if err := fetchEnums(imp, accu, prefix); err != nil {
			return err
		}
	}
	return nil
}
