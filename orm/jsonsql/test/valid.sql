CREATE OR REPLACE FUNCTION structgen_validate_json_number (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    RETURN jsonb_typeof(data) = 'number';
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_boolean (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    RETURN jsonb_typeof(data) = 'boolean';
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_4_boolean (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_boolean (value))
            FROM
                jsonb_array_elements(data))
        AND jsonb_array_length(data) = 4;
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_number (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_number (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_MyEnumI (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    RETURN jsonb_typeof(data) = 'number'
        AND data::int IN (0, 1, 2);
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_MyEnumS (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    RETURN jsonb_typeof(data) = 'string'
        AND data::text IN ('456', '45');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1883706743 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('DA', 'DB', 'DC', 'DD', 'DE'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'DA')
        AND structgen_validate_json_array_4_boolean (data -> 'DB')
        AND structgen_validate_json_array_number (data -> 'DC')
        AND structgen_validate_json_MyEnumI (data -> 'DD')
        AND structgen_validate_json_MyEnumS (data -> 'DE');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_string (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    RETURN jsonb_typeof(data) = 'string';
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_string (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_string (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_map_array_string (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        -- accept null value coming from nil maps
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'object'
        AND (
            SELECT
                bool_and(structgen_validate_json_array_string (value))
            FROM
                jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3716422242 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('A', 'B', 'C', 'D', 'E', 'F'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'A')
        AND structgen_validate_json_string (data -> 'B')
        AND structgen_validate_json_array_number (data -> 'C')
        AND structgen_validate_json_boolean (data -> 'D')
        AND structgen_validate_json_map_array_string (data -> 'E')
        AND structgen_validate_json_struct_1883706743 (data -> 'F');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

