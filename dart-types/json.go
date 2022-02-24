package darttypes

import (
	"fmt"
	"strings"
)

// generate serialization and deserialization for JSON format
// for named type we generate helpers function, unless
// they are basic

func (n named) json() string {
	// we directly call the underlying function to avoid boilerplate
	return ""
}

func (b basic) json() string {
	switch b {
	case dartString, dartBool, dartFloat, dartInt, dartAny:
		return fmt.Sprintf(`%s %sFromJson(dynamic json) => json as %s;
		
		dynamic %sToJson(%s item) => item;
		
		`, b, b, b, b, b)
	case dartTime:
		return `DateTime DateTimeFromJson(dynamic json) => DateTime.parse(json as String);

		dynamic DateTimeToJson(DateTime dt) => dt.toString();
		`
	default:
		panic("exhaustive switch")
	}
}

func (en enum) json() string {
	valueType := "String"
	if en.enum.IsInt {
		valueType = "int"
	}
	return fmt.Sprintf(`%s %sFromJson(dynamic json) => __%s.fromValue(json as %s);
	
	dynamic %sToJson(%s item) => item.toValue();
	
	`, en.enum.Name, en.enum.Name, en.enum.Name, valueType,
		en.enum.Name, en.enum.Name,
	)
}

func (l list) json() string {
	return fmt.Sprintf(`%s %sFromJson(dynamic json) {
		return (json as List<dynamic>).map(%sFromJson).toList();
	}

	dynamic %sToJson(%s item) {
		return item.map(%sToJson).toList();
	}
	`, l.name(), l.functionId(), l.element.functionId(),
		l.functionId(), l.name(), l.element.functionId(),
	)
}

func (d dict) json() string {
	return fmt.Sprintf(`%s %sFromJson(dynamic json) {
		return json.map((k,v) => MapEntry(k as %s, %sFromJson(v)));
	}
	
	dynamic %sToJson(%s item) {
		return item.map((k,v) => MapEntry(%sToJson(k), %sToJson(v)));
	}
	`, d.name(), d.functionId(), d.key.name(), d.element.functionId(),
		d.functionId(), d.name(), d.key.functionId(), d.element.functionId())
}

func (cl class) json() string {
	var fieldsFrom, fieldsTo []string
	for _, f := range cl.fields {
		fieldsFrom = append(fieldsFrom, fmt.Sprintf("%sFromJson(json['%s'])", f.type_.functionId(), f.name))
		fieldsTo = append(fieldsTo, fmt.Sprintf("%q : %sToJson(item.%s)", f.name, f.type_.functionId(), f.name))
	}
	return fmt.Sprintf(`
	%s %sFromJson(dynamic json_) {
		final json = (json_ as JSON);
		return %s(
			%s
		);
	}
	
	JSON %sToJson(%s item) {
		return {
			%s
		};
	}
	
	`, cl.name_, cl.name_, cl.name_, strings.Join(fieldsFrom, ",\n"),
		cl.name_, cl.name_, strings.Join(fieldsTo, ",\n"),
	)
}

func (u union) json() string {
	var casesFrom, casesTo []string

	for i, member := range u.members {
		casesFrom = append(casesFrom, fmt.Sprintf(`case  %d:
			return %sFromJson(data);`, i, member.functionId()))

		caseTo := fmt.Sprintf(`if (item is %s) {
			return {'Kind': %d, 'Data': %sToJson(item)};
		}`, member.name(), i, member.functionId())
		if i != 0 {
			caseTo = "else " + caseTo
		}
		casesTo = append(casesTo, caseTo)
	}

	codeFrom := fmt.Sprintf(`%s %sFromJson(dynamic json_) {
		final json = json_ as JSON;
		final kind = json['Kind'] as int;
		final data = json['Data'];
		switch (kind) {
			%s
		default:
			throw ("unexpected type");
		}
	}
	`, u.name_, u.name_, strings.Join(casesFrom, "\n"))

	codeTo := fmt.Sprintf(`JSON %sToJson(%s item) {
		%s else {
			throw ("unexpected type");
		}	
	}
	`, u.name_, u.name_, strings.Join(casesTo, ""))

	return codeFrom + "\n" + codeTo
}
