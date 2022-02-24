package interfaces

import (
	"fmt"
	"strings"
)

// return the go code implementing JSON convertions
func (u Interface) json() string {
	var casesFrom, casesTo []string

	for i, member := range u.Members {
		casesFrom = append(casesFrom, fmt.Sprintf(`case %d:
			var out %s
			err = json.Unmarshal(wr.Data, &out)
			return out, err
	`, i, member.Obj().Name()))

		caseTo := fmt.Sprintf(`case %s:
			wr = wrapper{Kind: %d, Data: item}
		`, member.Obj().Name(), i)
		casesTo = append(casesTo, caseTo)
	}

	codeFrom := fmt.Sprintf(`func %sUnmarshalJSON(src []byte) (%s, error) {
		var wr struct {
			Data json.RawMessage
			Kind int
		}
		err := json.Unmarshal(src, &wr)
		if err != nil {
			return nil, err
		}
		switch wr.Kind {
			%s
		default:
			panic("exhaustive switch")
		}
	}
	`, u.Name.Obj().Name(), u.Name.Obj().Name(), strings.Join(casesFrom, ""))

	codeTo := fmt.Sprintf(`func %sMarshalJSON(item %s) ([]byte, error) {
		type wrapper struct {
			Data interface{}
			Kind int
		}
		var wr wrapper
		switch item.(type) {
		%s
		default:
			panic("exhaustive switch")
		}
		return json.Marshal(wr)
	}
	`, u.Name.Obj().Name(), u.Name.Obj().Name(), strings.Join(casesTo, ""))

	return codeFrom + "\n" + codeTo
}
