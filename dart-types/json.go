package darttypes

import (
	"fmt"
	"strings"
)

// generate serialization and deserialization for JSON format
// for named type we generate helpers function, unless
// they are basic

func (n named) renderJSONconvertors(fromBody, toBody string) string {
	return fmt.Sprintf(`%s %sFromJson(dynamic json) {
		%s	
	}
	
	dynamic %sToJson(%s item) {
		%s	
	}
	`, n, n, fromBody,
		n, n, toBody,
	)
}

func (b basic) renderJSONconvertors() string {
	switch b {
	case dartString, dartBool, dartFloat, dartInt, dartAny:
		return fmt.Sprintf(`%s %sFromJson(dynamic json) {
			return json as %s;
		}

		dynamic %sToJson(%s item) {
			return item;
		}
		`, b, b, b, b, b)
	case dartTime:
		return `DateTime DateTimeFromJson(dynamic json) {
			return DateTime.parse(json as String);
		}

		dynamic DateTimeToJson(DateTime dt) {
			return dt.toString();
		}
		`
	default:
		panic("exhaustive switch")
	}
}

func (en enum) renderJSONconvertors() string {
	valueType := "String"
	if en.IsInt {
		valueType = "int"
	}
	return fmt.Sprintf(`%s %sFromJson(dynamic json) {
		return __%s.fromValue(json as %s);
	}
	
	dynamic %sToJson(%s item) {
		return item.toValue();
	}
	
	`, en.Name, en.Name, en.Name, valueType,
		en.Name, en.Name,
	)
}

func (named) fromJSONBody() string { return "" }
func (class) fromJSONBody() string { return "" }
func (enum) fromJSONBody() string  { return "" }
func (union) fromJSONBody() string { return "" } // TODO

func (named) toJSONBody() string { return "" }
func (class) toJSONBody() string { return "" }
func (enum) toJSONBody() string  { return "" }
func (union) toJSONBody() string { return "" } // TODO

func (b basic) fromJSONBody() string {
	return fmt.Sprintf("return %sFromJson(json);", b)
}

func (b basic) toJSONBody() string {
	return fmt.Sprintf("return %sToJson(item);", b)
}

func (l list) fromJSONBody() string {
	return fmt.Sprintf(`
		return (json as List<dynamic>).map(%sFromJson).toList();
	`, l.element.render())
}

func (l list) toJSONBody() string {
	return fmt.Sprintf(`
		return item.map(%sToJson).toList();
	`, l.element.render())
}

func (d dict) fromJSONBody() string {
	return fmt.Sprintf(`
		return json.map((k,v) => MapEntry(k as %s, %sFromJson(v)));
	`, d.key.render(), d.element.render())
}

func (d dict) toJSONBody() string {
	return fmt.Sprintf(`
		return item.map((k,v) => MapEntry(%sToJson(k), %sToJson(v)));
	`, d.key.render(), d.element.render())
}

func (cl class) renderJSONconvertors() string {
	var fieldsFrom, fieldsTo []string
	for _, f := range cl.fields {
		fieldsFrom = append(fieldsFrom, fmt.Sprintf("%sFromJson(json['%s'])", f.type_.render(), f.name))
		fieldsTo = append(fieldsTo, fmt.Sprintf("%q : %sToJson(item.%s)", f.name, f.type_.render(), f.name))
	}
	return fmt.Sprintf(`
	%s %sFromJson(JSON json) {
		return %s(
			%s
		);
	}
	
	JSON %sToJson(%s item) {
		return {
			%s
		};
	}
	
	`, cl.name, cl.name, cl.name, strings.Join(fieldsFrom, ",\n"),
		cl.name, cl.name, strings.Join(fieldsTo, ",\n"),
	)
}
