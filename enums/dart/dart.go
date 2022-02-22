package dart

import (
	"fmt"
	"strings"

	"github.com/benoitkugler/structgen/enums"
)

func EnumAsDart(e enums.Type) string {
	var names, values, labels []string
	for _, v := range e.Values {
		names = append(names, v.VarName)
		labels = append(labels, fmt.Sprintf("%q", v.Label))
		values = append(values, v.Value)
	}

	var fromValue string
	if e.IsInt { // want can just use Dart builtin enums
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
