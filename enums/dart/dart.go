package dart

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/benoitkugler/structgen/enums"
)

func EnumAsDart(e enums.Type) string {
	if e.IsInt {
		// we have to sort by values, which must be ints
		sort.Slice(e.Values, func(i, j int) bool {
			vi, err := strconv.Atoi(e.Values[i].Value)
			if err != nil {
				panic(err)
			}
			vj, err := strconv.Atoi(e.Values[j].Value)
			if err != nil {
				panic(err)
			}
			return vi < vj
		})
	}

	var names, values, labels []string
	for _, v := range e.Values {
		names = append(names, v.VarName)
		labels = append(labels, fmt.Sprintf("%q", v.Label))
		values = append(values, v.Value)
	}

	var fromValue string
	if e.IsInt { // we can just use Dart builtin enums

		fromValue = fmt.Sprintf(`static %s fromValue(int i) {
			return i as %s;
		}
		
		int toValue() {
			return this.index;
		}
		`, e.Name, e.Name)
	} else { // add lookup array
		fromValue = fmt.Sprintf(`
		static const _values = [
			%s
		];
		static %s fromValue(String s) {
			return _values.indexOf(s) as %s;
		}

		String toValue() {
			return _values[this.index];
		}
		`, strings.Join(values, ", "), e.Name, e.Name)
	}

	return fmt.Sprintf(`enum  %s {
		%s
	}
	
	extension __%s on %s {
		static const _labels = [
			%s
		];
		
		String label() { return _labels[this.index]; }

		%s
	}
	`, e.Name, strings.Join(names, ", "), e.Name, e.Name, strings.Join(labels, ", "), fromValue)
}
