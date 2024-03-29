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
	case dartFloat: // use num to avoid issue with integers values
		return `double doubleFromJson(dynamic json) => (json as num).toDouble();
	
	double doubleToJson(double item) => item;
	
	`
	case dartString: // accept null
		return `String stringFromJson(dynamic json) => json == null ? "" : json as String;
		
		String stringToJson(String item) => item;
		
		`
	case dartBool, dartInt, dartAny:
		return fmt.Sprintf(`%s %sFromJson(dynamic json) => json as %s;
		
		%s %sToJson(%s item) => item;
		
		`, b, b.functionId(), b, b, b.functionId(), b)
	case dartTime:
		return `DateTime dateTimeFromJson(dynamic json) => DateTime.parse(json as String);

		dynamic dateTimeToJson(DateTime dt) => dt.toString();
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
	return fmt.Sprintf(`%s %sFromJson(dynamic json) => _%sExt.fromValue(json as %s);
	
	dynamic %sToJson(%s item) => item.toValue();
	
	`, en.name(), en.functionId(), en.name(), valueType,
		en.functionId(), en.name(),
	)
}

func (l list) json() string {
	// nil slices are jsonized as null, check for it then
	return fmt.Sprintf(`%s %sFromJson(dynamic json) {
		if (json == null) {
			return [];
		}
		return (json as List<dynamic>).map(%sFromJson).toList();
	}

	List<dynamic> %sToJson(%s item) {
		return item.map(%sToJson).toList();
	}
	`, l.name(), l.functionId(), l.element.functionId(),
		l.functionId(), l.name(), l.element.functionId(),
	)
}

func (d dict) json() string {
	keyFromJson := "k as " + d.key.name()
	if d.key.name() == "int" {
		keyFromJson = "int.parse(k)"
	}

	// nil dict are jsonized as null, check for it then
	return fmt.Sprintf(`%s %sFromJson(dynamic json) {
		if (json == null) {
			return {};
		}
		return (json as JSON).map((k,v) => MapEntry(%s, %sFromJson(v)));
	}
	
	Map<String, dynamic> %sToJson(%s item) {
		return item.map((k,v) => MapEntry(%sToJson(k).toString(), %sToJson(v)));
	}
	`, d.name(), d.functionId(), keyFromJson, d.element.functionId(),
		d.functionId(), d.name(), d.key.functionId(), d.element.functionId())
}

func (cl class) json() string {
	var fieldsFrom, fieldsTo []string
	for _, f := range cl.fields {
		fieldsFrom = append(fieldsFrom, fmt.Sprintf("%sFromJson(json['%s'])", f.type_.functionId(), f.name))
		fieldsTo = append(fieldsTo, fmt.Sprintf("%q : %sToJson(item.%s)", f.name, f.type_.functionId(), f.dartName()))
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
	
	`, cl.name_, cl.functionId(), cl.name_, strings.Join(fieldsFrom, ",\n"),
		cl.functionId(), cl.name_, strings.Join(fieldsTo, ",\n"),
	)
}

func (u union) json() string {
	var casesFrom, casesTo []string

	for i, member := range u.members {
		kindValue := member.tag
		casesFrom = append(casesFrom, fmt.Sprintf(`case %q:
			return %sFromJson(data);`, kindValue, member.type_.functionId()))

		caseTo := fmt.Sprintf(`if (item is %s) {
			return {'Kind': %q, 'Data': %sToJson(item)};
		}`, member.type_.name(), member.tag, member.type_.functionId())
		if i != 0 {
			caseTo = "else " + caseTo
		}
		casesTo = append(casesTo, caseTo)
	}

	codeFrom := fmt.Sprintf(`%s %sFromJson(dynamic json_) {
		final json = json_ as JSON;
		final kind = json['Kind'] as String;
		final data = json['Data'];
		switch (kind) {
			%s
		default:
			throw ("unexpected type");
		}
	}
	`, u.name_, u.functionId(), strings.Join(casesFrom, "\n"))

	codeTo := fmt.Sprintf(`JSON %sToJson(%s item) {
		%s else {
			throw ("unexpected type");
		}	
	}
	`, u.functionId(), u.name_, strings.Join(casesTo, ""))

	return codeFrom + "\n" + codeTo
}
