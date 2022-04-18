package jsonsql

import (
	"fmt"
	"strings"

	"github.com/benoitkugler/structgen/loader"
)

// SetupSQLCode should be included once before all functions definitions,
// in order to cleanup potential old definitions
const SetupSQLCode = `
CREATE OR REPLACE FUNCTION f_delfunc (OUT func_dropped int
)
AS $func$
DECLARE
    _sql text;
BEGIN
    SELECT
        count(*)::int,
        'DROP FUNCTION ' || string_agg(oid::regprocedure::text, '; DROP FUNCTION ')
    FROM
        pg_proc
    WHERE
        starts_with (proname, 'structgen_validate_json')
        AND pg_function_is_visible(oid) INTO func_dropped,
        _sql;
    -- only returned if trailing DROPs succeed
    IF func_dropped > 0 THEN
        -- only if function(s) found
        EXECUTE _sql;
    END IF;
END
$func$
LANGUAGE plpgsql;

SELECT
    f_delfunc ();

DROP FUNCTION f_delfunc;
`

const (
	vDynamic = `
	-- No validation : accept anything
	CREATE OR REPLACE FUNCTION %s (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		RETURN TRUE;
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;`

	vBasic = `
	CREATE OR REPLACE FUNCTION %s (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		RETURN jsonb_typeof(data) = '%s';
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;`

	vEnum = `
	CREATE OR REPLACE FUNCTION %s (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		RETURN jsonb_typeof(data) = '%s' AND %s IN %s;
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;`
)

func (b basic) Id() string { return string(b) }

func (b basic) Validations() []loader.Declaration {
	s := loader.Declaration{Id: FunctionName(b)}
	if b == Dynamic { // special case for Dynamic
		s.Content = fmt.Sprintf(vDynamic, FunctionName(b))
	} else {
		s.Content = fmt.Sprintf(vBasic, FunctionName(b), string(b))
	}
	return []loader.Declaration{s}
}

func (b enumValue) Id() string { return b.enumType.Name }

func (b enumValue) Validations() []loader.Declaration {
	s := loader.Declaration{Id: FunctionName(b)}
	typeCast := `data#>>'{}'`
	if b.enumType.IsInt {
		typeCast = "data::int"
	}
	s.Content = fmt.Sprintf(vEnum, FunctionName(b), string(b.basic), typeCast, b.enumType.AsTuple())
	return []loader.Declaration{s}
}

const vArray = `
	CREATE OR REPLACE FUNCTION %s (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		%s
		IF jsonb_typeof(data) != 'array' THEN RETURN FALSE; END IF;
		%s 
		RETURN (SELECT bool_and( %s(value) )  FROM jsonb_array_elements(data)) 
			%s;
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;`

func (b Array) Id() string {
	as := "array_"
	if b.length >= 0 {
		as += fmt.Sprintf("%d_", b.length)
	}
	return as + b.elem.Id()
}

func (b Array) Validations() []loader.Declaration {
	critereLength, acceptZeroLength := "", ""
	if b.length >= 0 {
		critereLength = fmt.Sprintf("AND jsonb_array_length(data) = %d", b.length)
	} else {
		acceptZeroLength = "IF jsonb_array_length(data) = 0 THEN RETURN TRUE; END IF;"
	}
	gardNull := ""
	if b.length == -1 { // accepts null
		gardNull = "IF jsonb_typeof(data) = 'null' THEN RETURN TRUE; END IF;"
	}

	out := b.elem.Validations() // recursion

	fn, elemFuncName := FunctionName(b), FunctionName(b.elem)
	content := fmt.Sprintf(vArray, fn, gardNull, acceptZeroLength, elemFuncName, critereLength)
	out = append(out, loader.Declaration{Id: fn, Content: content})

	return out
}

const vMap = `
	CREATE OR REPLACE FUNCTION %s (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		IF jsonb_typeof(data) = 'null' THEN -- accept null value coming from nil maps 
			RETURN TRUE;
		END IF;
		RETURN jsonb_typeof(data) = 'object'
			AND (SELECT bool_and( %s(value) ) FROM jsonb_each(data));
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;`

func (b Map) Id() string {
	return "map_" + b.elem.Id()
}

func (b Map) Validations() []loader.Declaration {
	out := b.elem.Validations() // recursion
	fn, elemFuncName := FunctionName(b), FunctionName(b.elem)
	content := fmt.Sprintf(vMap, fn, elemFuncName)

	out = append(out, loader.Declaration{Id: fn, Content: content})
	return out
}

const vStruct = `
	CREATE OR REPLACE FUNCTION %s (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		IF jsonb_typeof(data) != 'object' THEN 
			RETURN FALSE;
		END IF;
		RETURN (SELECT bool_and( 
			%s
		) FROM jsonb_each(data))  
		%s
		;
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;`

// to work around possible hash collision,
// we need to be able to check if an hash is already taken
var structIDsTable = map[uint32]Struct{}

func (b Struct) dump() []byte {
	var data []byte
	for _, f := range b.fields {
		data = append(data, f.key...)
		data = append(data, f.type_.Id()...)
	}
	return data
}

func (b Struct) Validations() (out []loader.Declaration) {
	var keys, checks []string
	for _, f := range b.fields {
		out = append(out, f.type_.Validations()...) // recursion
		keys = append(keys, fmt.Sprintf("'%s'", f.key))
		checks = append(checks, fmt.Sprintf("AND %s(data->'%s')", FunctionName(f.type_), f.key))
	}
	keyList := "key IN (" + strings.Join(keys, ", ") + ")"
	if len(keys) == 0 {
		keyList = "TRUE"
	}
	checkList := strings.Join(checks, "\n")
	fn := FunctionName(b)
	content := fmt.Sprintf(vStruct, fn, keyList, checkList)
	out = append(out, loader.Declaration{Id: fn, Content: content})

	return out
}
