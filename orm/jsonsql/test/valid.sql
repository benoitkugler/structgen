
	CREATE OR REPLACE FUNCTION structgen_validate_json_number (data jsonb)
		RETURNS boolean
		AS $$
	DECLARE
		is_valid boolean := jsonb_typeof(data) = 'number';
	BEGIN
		IF NOT is_valid THEN 
			RAISE WARNING '% is not a number', data;
		END IF;
		RETURN is_valid;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_boolean (data jsonb)
		RETURNS boolean
		AS $$
	DECLARE
		is_valid boolean := jsonb_typeof(data) = 'boolean';
	BEGIN
		IF NOT is_valid THEN 
			RAISE WARNING '% is not a boolean', data;
		END IF;
		RETURN is_valid;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_array_4_boolean (data jsonb)
		RETURNS boolean
		AS $$
	BEGIN
		
		IF jsonb_typeof(data) != 'array' THEN RETURN FALSE; END IF;
		 
		RETURN (SELECT bool_and( structgen_validate_json_boolean(value) )  FROM jsonb_array_elements(data)) 
			AND jsonb_array_length(data) = 4;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_array_number (data jsonb)
		RETURNS boolean
		AS $$
	BEGIN
		IF jsonb_typeof(data) = 'null' THEN RETURN TRUE; END IF;
		IF jsonb_typeof(data) != 'array' THEN RETURN FALSE; END IF;
		IF jsonb_array_length(data) = 0 THEN RETURN TRUE; END IF; 
		RETURN (SELECT bool_and( structgen_validate_json_number(value) )  FROM jsonb_array_elements(data)) 
			;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_MyEnumI (data jsonb)
		RETURNS boolean
		AS $$
	DECLARE
		is_valid boolean := jsonb_typeof(data) = 'number' AND data#>>'{}' IN (0, 1, 2);
	BEGIN
		IF NOT is_valid THEN 
			RAISE WARNING '% is not a MyEnumI', data;
		END IF;
		RETURN is_valid;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_MyEnumS (data jsonb)
		RETURNS boolean
		AS $$
	DECLARE
		is_valid boolean := jsonb_typeof(data) = 'string' AND data#>>'{}' IN ('456', '45');
	BEGIN
		IF NOT is_valid THEN 
			RAISE WARNING '% is not a MyEnumS', data;
		END IF;
		RETURN is_valid;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_tes_DataType (data jsonb)
		RETURNS boolean
		AS $$
	DECLARE 
		is_valid boolean;
	BEGIN
		IF jsonb_typeof(data) != 'object' THEN 
			RETURN FALSE;
		END IF;
		is_valid := (SELECT bool_and( 
			key IN ('DA', 'DB', 'DC', 'DD', 'DE')
		) FROM jsonb_each(data))  
		AND structgen_validate_json_number(data->'DA')
AND structgen_validate_json_array_4_boolean(data->'DB')
AND structgen_validate_json_array_number(data->'DC')
AND structgen_validate_json_MyEnumI(data->'DD')
AND structgen_validate_json_MyEnumS(data->'DE');

		RETURN is_valid;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_array_tes_DataType (data jsonb)
		RETURNS boolean
		AS $$
	BEGIN
		IF jsonb_typeof(data) = 'null' THEN RETURN TRUE; END IF;
		IF jsonb_typeof(data) != 'array' THEN RETURN FALSE; END IF;
		IF jsonb_array_length(data) = 0 THEN RETURN TRUE; END IF; 
		RETURN (SELECT bool_and( structgen_validate_json_tes_DataType(value) )  FROM jsonb_array_elements(data)) 
			;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_string (data jsonb)
		RETURNS boolean
		AS $$
	DECLARE
		is_valid boolean := jsonb_typeof(data) = 'string';
	BEGIN
		IF NOT is_valid THEN 
			RAISE WARNING '% is not a string', data;
		END IF;
		RETURN is_valid;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_array_string (data jsonb)
		RETURNS boolean
		AS $$
	BEGIN
		IF jsonb_typeof(data) = 'null' THEN RETURN TRUE; END IF;
		IF jsonb_typeof(data) != 'array' THEN RETURN FALSE; END IF;
		IF jsonb_array_length(data) = 0 THEN RETURN TRUE; END IF; 
		RETURN (SELECT bool_and( structgen_validate_json_string(value) )  FROM jsonb_array_elements(data)) 
			;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_map_array_string (data jsonb)
		RETURNS boolean
		AS $$
	BEGIN
		IF jsonb_typeof(data) = 'null' THEN -- accept null value coming from nil maps 
			RETURN TRUE;
		END IF;
		RETURN jsonb_typeof(data) = 'object'
			AND (SELECT bool_and( structgen_validate_json_array_string(value) ) FROM jsonb_each(data));
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_tes_Model (data jsonb)
		RETURNS boolean
		AS $$
	DECLARE 
		is_valid boolean;
	BEGIN
		IF jsonb_typeof(data) != 'object' THEN 
			RETURN FALSE;
		END IF;
		is_valid := (SELECT bool_and( 
			key IN ('A', 'B', 'C', 'D', 'E', 'F', 'G')
		) FROM jsonb_each(data))  
		AND structgen_validate_json_number(data->'A')
AND structgen_validate_json_string(data->'B')
AND structgen_validate_json_array_number(data->'C')
AND structgen_validate_json_boolean(data->'D')
AND structgen_validate_json_map_array_string(data->'E')
AND structgen_validate_json_tes_DataType(data->'F')
AND structgen_validate_json_array_tes_DataType(data->'G');

		RETURN is_valid;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_array_tes_Recursive (data jsonb)
		RETURNS boolean
		AS $$
	BEGIN
		IF jsonb_typeof(data) = 'null' THEN RETURN TRUE; END IF;
		IF jsonb_typeof(data) != 'array' THEN RETURN FALSE; END IF;
		IF jsonb_array_length(data) = 0 THEN RETURN TRUE; END IF; 
		RETURN (SELECT bool_and( structgen_validate_json_tes_Recursive(value) )  FROM jsonb_array_elements(data)) 
			;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_tes_Recursive (data jsonb)
		RETURNS boolean
		AS $$
	DECLARE 
		is_valid boolean;
	BEGIN
		IF jsonb_typeof(data) != 'object' THEN 
			RETURN FALSE;
		END IF;
		is_valid := (SELECT bool_and( 
			key IN ('B', 'A')
		) FROM jsonb_each(data))  
		AND structgen_validate_json_array_tes_Recursive(data->'B')
AND structgen_validate_json_number(data->'A');

		RETURN is_valid;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_tes_S (data jsonb)
		RETURNS boolean
		AS $$
	DECLARE 
		is_valid boolean;
	BEGIN
		IF jsonb_typeof(data) != 'object' THEN 
			RETURN FALSE;
		END IF;
		is_valid := (SELECT bool_and( 
			key IN ('A')
		) FROM jsonb_each(data))  
		AND structgen_validate_json_number(data->'A');

		RETURN is_valid;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_tes_myItf (data jsonb)
		RETURNS boolean
		AS $$
	BEGIN
		IF jsonb_typeof(data) != 'object' OR jsonb_typeof(data->'Kind') != 'number' OR jsonb_typeof(data->'Data') = 'null' THEN 
			RETURN FALSE;
		END IF;
		CASE 
			WHEN (data->'Kind')::int = 0 THEN 
 RETURN structgen_validate_json_array_tes_DataType(data->'Data');
WHEN (data->'Kind')::int = 1 THEN 
 RETURN structgen_validate_json_tes_S(data->'Data');
WHEN (data->'Kind')::int = 2 THEN 
 RETURN structgen_validate_json_array_number(data->'Data');
ELSE RETURN FALSE;
		END CASE;
	END;
	$$
	LANGUAGE 'plpgsql'
	IMMUTABLE;
