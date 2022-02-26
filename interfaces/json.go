package interfaces

import (
	"fmt"
	"strings"
)

// return the go code implementing JSON convertions
func (u Interface) json() string {
	var casesFrom, casesTo []string

	name := u.Name.Obj().Name()
	wrapperName := name + "Wrapper"

	for i, member := range u.Members {
		casesFrom = append(casesFrom, fmt.Sprintf(`case %d:
			var data %s
			err = json.Unmarshal(wr.Data, &data)
			out.data = data
	`, i, member.Obj().Name()))

		caseTo := fmt.Sprintf(`case %s:
			wr = wrapper{Kind: %d, Data: data}
		`, member.Obj().Name(), i)
		casesTo = append(casesTo, caseTo)
	}

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
		switch data := item.data.(type) {
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
		data %s
	}

	%s 

	%s
	`, wrapperName, name, wrapperName, name, codeFrom, codeTo)
}
