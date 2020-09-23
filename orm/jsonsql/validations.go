package jsonsql

import (
	"fmt"
	"hash/adler32"
	"strings"
	"time"

	"github.com/benoitkugler/structgen/loader"
)

// should be included once before all functions definitions,
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

type sqlFunc struct {
	declId  string
	content string
}

func (s sqlFunc) Id() string     { return s.declId }
func (s sqlFunc) Render() string { return s.content }

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
)

func (b basic) Id() string { return string(b) }

func (b basic) AddValidation(l *loader.Declarations) {
	s := sqlFunc{declId: funcName(b)}
	if b == Dynamic { // special case for Dynamic
		s.content = fmt.Sprintf(vDynamic, funcName(b))
	} else {
		s.content = fmt.Sprintf(vBasic, funcName(b), string(b))
	}
	l.Add(s)
}

const vArray = `
	CREATE OR REPLACE FUNCTION %s (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		%s
		RETURN jsonb_typeof(data) = 'array'
			AND (SELECT bool_and( %s(value) )  FROM jsonb_array_elements(data)) 
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

func (b Array) AddValidation(l *loader.Declarations) {
	critereLength := ""
	if b.length >= 0 {
		critereLength = fmt.Sprintf("AND jsonb_array_length(data) = %d", b.length)
	}
	gardNull := ""
	if b.length == -1 { // accepts null
		gardNull = "IF jsonb_typeof(data) = 'null' THEN RETURN TRUE; END IF;"
	}
	b.elem.AddValidation(l) // recursion
	fn, elemFuncName := funcName(b), funcName(b.elem)
	content := fmt.Sprintf(vArray, fn, gardNull, elemFuncName, critereLength)
	l.Add(sqlFunc{declId: fn, content: content})
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

func (b Map) AddValidation(l *loader.Declarations) {
	b.elem.AddValidation(l) // recursion
	fn, elemFuncName := funcName(b), funcName(b.elem)
	content := fmt.Sprintf(vMap, fn, elemFuncName)
	l.Add(sqlFunc{declId: fn, content: content})
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

func (b Struct) Id() string {
	var data []byte
	for _, f := range b.fields {
		data = append(data, f.key...)
		data = append(data, f.type_.Id()...)
	}
	return fmt.Sprintf("struct_%d_%d", adler32.Checksum(data), time.Now().Nanosecond())
}

func (b Struct) AddValidation(l *loader.Declarations) {
	var keys, checks []string
	for _, f := range b.fields {
		f.type_.AddValidation(l) // recursion
		keys = append(keys, fmt.Sprintf("'%s'", f.key))
		checks = append(checks, fmt.Sprintf("AND %s(data->'%s')", funcName(f.type_), f.key))
	}
	keyList := "key IN (" + strings.Join(keys, ", ") + ")"
	if len(keys) == 0 {
		keyList = "TRUE"
	}
	checkList := strings.Join(checks, "\n")
	fn := funcName(b)
	content := fmt.Sprintf(vStruct, fn, keyList, checkList)
	l.Add(sqlFunc{declId: fn, content: content})
}