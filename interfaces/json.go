package interfaces

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/benoitkugler/structgen/loader"
)

// return the go code implementing JSON convertions
func (u Interface) json() string {
	var casesFrom, casesTo, kindDecls []string

	name := u.Name.Obj().Name()
	wrapperName := name + "Wrapper"

	for _, member := range u.Members {
		memberName := member.Obj().Name()
		kindValue := memberName
		kindVarName := memberName + name[0:2] + "Kind"
		kindDecls = append(kindDecls, fmt.Sprintf("%s = %q", kindVarName, kindValue))

		casesFrom = append(casesFrom, fmt.Sprintf(`case %q:
			var data %s
			err = json.Unmarshal(wr.Data, &data)
			out.Data = data
	`, kindValue, memberName))

		caseTo := fmt.Sprintf(`case %s:
			wr = wrapper{Kind: %q, Data: data}
		`, memberName, kindValue)
		casesTo = append(casesTo, caseTo)
	}

	codeKinds := fmt.Sprintf(`
	const (
		%s
	)
	`, strings.Join(kindDecls, "\n"))

	codeFrom := fmt.Sprintf(`func (out *%s) UnmarshalJSON(src []byte) error {
		var wr struct {
			Kind string
			Data json.RawMessage
		}
		err := json.Unmarshal(src, &wr)
		if err != nil {
			return err
		}
		switch wr.Kind {
			%s
		default:
			panic("exhaustive switch")
		}
		return err
	}
	`, wrapperName, strings.Join(casesFrom, ""))

	codeTo := fmt.Sprintf(`func (item %s) MarshalJSON() ([]byte, error) {
		type wrapper struct {
			Data interface{}
			Kind string
		}
		var wr wrapper
		switch data := item.Data.(type) {
		%s
		default:
			panic("exhaustive switch")
		}
		return json.Marshal(wr)
	}
	`, wrapperName, strings.Join(casesTo, ""))

	return fmt.Sprintf(`
	// %s may be used as replacements for %s 
	// when working with JSON
	type %s struct{
		Data %s
	}

	%s 

	%s

	%s
	`, wrapperName, name, wrapperName, name, codeFrom, codeTo, codeKinds)
}

type itfSlice struct {
	name string
	elem Interface
}

func (s itfSlice) Render() []loader.Declaration {
	out := s.elem.Render()
	out = append(out, loader.Declaration{
		Id:      s.name + "_json",
		Content: s.json(),
	})
	return out
}

func (s itfSlice) json() string {
	return fmt.Sprintf(`func (ct %s) MarshalJSON() ([]byte, error) {
		tmp := make([]%sWrapper, len(ct))
		for i, v := range ct {
			tmp[i].Data = v
		}
		return json.Marshal(tmp)
	}
	
	func (ct *%s) UnmarshalJSON(data []byte) error {
		var tmp []%sWrapper
		err := json.Unmarshal(data, &tmp)
		*ct = make(%s, len(tmp))
		for i, v := range tmp {
			(*ct)[i] = v.Data
		}
		return err
	}`, s.name, s.elem.Name.Obj().Name(), s.name, s.elem.Name.Obj().Name(), s.name)
}

type classField struct {
	goType *types.Var
	// may be empty if ignored
	type_ loader.Type
}

func (cf classField) requireWrapper() bool {
	_, isITF := cf.goType.Type().Underlying().(*types.Interface)
	return isITF
}

type class struct {
	goType *types.Named
	fields []classField
}

func (cl class) name() string { return cl.goType.Obj().Name() }

func myQualifier(pkg *types.Package) types.Qualifier {
	if pkg == nil {
		return nil
	}
	return func(other *types.Package) string {
		if pkg == other {
			return "" // same package; unqualified
		}
		return other.Name()
	}
}

func (cl class) Render() (out []loader.Declaration) {
	var requireWrapper bool
	// render the field decls
	for _, field := range cl.fields {
		if field.type_ != nil {
			out = append(out, field.type_.Render()...)
		}
		if field.requireWrapper() {
			requireWrapper = true
		}
	}

	// check for field requiring the wrapper
	if requireWrapper {
		var fieldsDef, fieldsAssignToWrapper, fieldsAssignFromWrapper []string
		for _, field := range cl.fields {
			fieldName := field.goType.Name()

			fieldType := types.TypeString(field.goType.Type(), myQualifier(cl.goType.Obj().Pkg()))

			fieldAssignToW := fmt.Sprintf("item.%s", fieldName)
			fieldAssignFromW := fmt.Sprintf("wr.%s", fieldName)
			if field.requireWrapper() { // wrapper field
				fieldType = fieldType + "Wrapper"
				fieldAssignToW = fmt.Sprintf("%s{item.%s}", fieldType, fieldName)
				fieldAssignFromW = fmt.Sprintf("wr.%s.Data", fieldName)
			}

			fieldsDef = append(fieldsDef, fieldName+" "+fieldType)
			fieldsAssignToWrapper = append(fieldsAssignToWrapper, fmt.Sprintf("%s: %s,", fieldName, fieldAssignToW))
			fieldsAssignFromWrapper = append(fieldsAssignFromWrapper, fmt.Sprintf("item.%s = %s", fieldName, fieldAssignFromW))
		}

		code := fmt.Sprintf(`
		func (item %s) MarshalJSON() ([]byte, error) {
			type wrapper struct {
				%s
			}
			wr := wrapper{
				%s
			}
			return json.Marshal(wr)
		}

		func (item *%s) UnmarshalJSON(src []byte) error {
			type wrapper struct {
				%s
			}
			var wr wrapper 
			err := json.Unmarshal(src, &wr)
			if err != nil {
				return err
			}
			%s
			return nil
		}
		`,
			cl.name(), strings.Join(fieldsDef, "\n"), strings.Join(fieldsAssignToWrapper, "\n"),
			cl.name(), strings.Join(fieldsDef, "\n"), strings.Join(fieldsAssignFromWrapper, "\n"),
		)
		out = append(out, loader.Declaration{
			Id:      cl.name() + "_json",
			Content: code,
		})
	}
	return out
}
