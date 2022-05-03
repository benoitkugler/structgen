package interfaces

import (
	"fmt"
	"strings"

	"github.com/benoitkugler/structgen/loader"
)

// return the go code implementing JSON convertions
func (u Interface) json() string {
	var casesFrom, casesTo, kinds []string

	name := u.Name.Obj().Name()
	wrapperName := name + "Wrapper"

	for i, member := range u.Members {
		memberName := member.Obj().Name()
		kinds = append(kinds, memberName+"Kind") // note that the index i is matching the iota used in kinds

		casesFrom = append(casesFrom, fmt.Sprintf(`case %d:
			var data %s
			err = json.Unmarshal(wr.Data, &data)
			out.Data = data
	`, i, memberName))

		caseTo := fmt.Sprintf(`case %s:
			wr = wrapper{Kind: %d, Data: data}
		`, memberName, i)
		casesTo = append(casesTo, caseTo)
	}

	codeKinds := fmt.Sprintf(`
	const (
		%s = iota 
		%s
	)
	`, kinds[0], strings.Join(kinds[1:], "\n"))

	codeFrom := fmt.Sprintf(`func (out *%s) UnmarshalJSON(src []byte) error {
		var wr struct {
			Data json.RawMessage
			Kind int
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
			Kind int
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

type class struct {
	fields []loader.Type
}

func (cl class) Render() (out []loader.Declaration) {
	for _, field := range cl.fields {
		out = append(out, field.Render()...)
	}
	return out
}
