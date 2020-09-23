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

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1_155487247 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(TRUE)
        FROM
            jsonb_each(data));
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

CREATE OR REPLACE FUNCTION structgen_validate_json_number (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    RETURN jsonb_typeof(data) = 'number';
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1_155502241 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(TRUE)
        FROM
            jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3093437039_155507835 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_utilisateur', 'date_emission', 'tag'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_number (data -> 'id_utilisateur')
        AND structgen_validate_json_struct_1_155503749 (data -> 'date_emission')
        AND structgen_validate_json_string (data -> 'tag');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3550417834_155516581 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_commande', 'id_produit', 'quantite'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_commande')
        AND structgen_validate_json_number (data -> 'id_produit')
        AND structgen_validate_json_number (data -> 'quantite');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3550417834_155523031 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_commande', 'id_produit', 'quantite'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_commande')
        AND structgen_validate_json_number (data -> 'id_produit')
        AND structgen_validate_json_number (data -> 'quantite');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_3550417834_155524354 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_struct_3550417834_155524934 (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1_155530469 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(TRUE)
        FROM
            jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3094944368_155539125 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_utilisateur', 'date_emission', 'tag'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_number (data -> 'id_utilisateur')
        AND structgen_validate_json_struct_1_155535932 (data -> 'date_emission')
        AND structgen_validate_json_string (data -> 'tag');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_map_struct_3088194152_155540725 (data jsonb)
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
                bool_and(structgen_validate_json_struct_3083278945_155541546 (value))
            FROM
                jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2334919345_155547824 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('quantite', 'unite'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'quantite')
        AND structgen_validate_json_string (data -> 'unite');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

-- No validation : accept anything
CREATE OR REPLACE FUNCTION structgen_validate_json_ (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    RETURN TRUE;
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2860130181_155556992 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_utilisateur', 'id_ingredient', 'id_fournisseur', 'id_produit'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_utilisateur')
        AND structgen_validate_json_number (data -> 'id_ingredient')
        AND structgen_validate_json_number (data -> 'id_fournisseur')
        AND structgen_validate_json_number (data -> 'id_produit');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2860130181_155565540 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_utilisateur', 'id_ingredient', 'id_fournisseur', 'id_produit'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_utilisateur')
        AND structgen_validate_json_number (data -> 'id_ingredient')
        AND structgen_validate_json_number (data -> 'id_fournisseur')
        AND structgen_validate_json_number (data -> 'id_produit');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_2860130181_155566738 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_struct_2860130181_155567369 (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2689272702_155583188 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'nom', 'lieu'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_string (data -> 'nom')
        AND structgen_validate_json_string (data -> 'lieu');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2689272702_155591402 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'nom', 'lieu'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_string (data -> 'nom')
        AND structgen_validate_json_string (data -> 'lieu');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_map_struct_2689272702_155592485 (data jsonb)
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
                bool_and(structgen_validate_json_struct_2689272702_155593116 (value))
            FROM
                jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1373444784_155602359 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_sejour', 'nom', 'nb_personnes', 'couleur'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_number (data -> 'id_sejour')
        AND structgen_validate_json_string (data -> 'nom')
        AND structgen_validate_json_number (data -> 'nb_personnes')
        AND structgen_validate_json_string (data -> 'couleur');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1373444784_155611863 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_sejour', 'nom', 'nb_personnes', 'couleur'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_number (data -> 'id_sejour')
        AND structgen_validate_json_string (data -> 'nom')
        AND structgen_validate_json_number (data -> 'nb_personnes')
        AND structgen_validate_json_string (data -> 'couleur');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_map_struct_1373444784_155613001 (data jsonb)
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
                bool_and(structgen_validate_json_struct_1373444784_155613538 (value))
            FROM
                jsonb_each(data));
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

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1_155625301 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(TRUE)
        FROM
            jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2334919345_155631201 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('quantite', 'unite'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'quantite')
        AND structgen_validate_json_string (data -> 'unite');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3752929526_155634935 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'nom', 'unite', 'categorie', 'callories', 'conditionnement'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_string (data -> 'nom')
        AND structgen_validate_json_string (data -> 'unite')
        AND structgen_validate_json_string (data -> 'categorie')
        AND structgen_validate_json_struct_1_155626401 (data -> 'callories')
        AND structgen_validate_json_struct_2334919345_155632477 (data -> 'conditionnement');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2163611403_155643106 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_ingredient', 'id_produit', 'id_utilisateur'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_ingredient')
        AND structgen_validate_json_number (data -> 'id_produit')
        AND structgen_validate_json_number (data -> 'id_utilisateur');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2163611403_155652542 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_ingredient', 'id_produit', 'id_utilisateur'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_ingredient')
        AND structgen_validate_json_number (data -> 'id_produit')
        AND structgen_validate_json_number (data -> 'id_utilisateur');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_2163611403_155654002 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_struct_2163611403_155654539 (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1_155662275 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(TRUE)
        FROM
            jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2334919345_155666375 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('quantite', 'unite'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'quantite')
        AND structgen_validate_json_string (data -> 'unite');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3795069180_155669422 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'nom', 'unite', 'categorie', 'callories', 'conditionnement'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_string (data -> 'nom')
        AND structgen_validate_json_string (data -> 'unite')
        AND structgen_validate_json_string (data -> 'categorie')
        AND structgen_validate_json_struct_1_155663307 (data -> 'callories')
        AND structgen_validate_json_struct_2334919345_155667558 (data -> 'conditionnement');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_map_struct_3780061439_155671302 (data jsonb)
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
                bool_and(structgen_validate_json_struct_3807979777_155672589 (value))
            FROM
                jsonb_each(data));
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

CREATE OR REPLACE FUNCTION structgen_validate_json_array_7_boolean (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_boolean (value))
            FROM
                jsonb_array_elements(data))
        AND jsonb_array_length(data) = 7;
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3289518958_155682506 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_ingredient', 'quantite', 'cuisson'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_ingredient')
        AND structgen_validate_json_number (data -> 'quantite')
        AND structgen_validate_json_string (data -> 'cuisson');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3289518958_155687835 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_ingredient', 'quantite', 'cuisson'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_ingredient')
        AND structgen_validate_json_number (data -> 'quantite')
        AND structgen_validate_json_string (data -> 'cuisson');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_3289518958_155688961 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_struct_3289518958_155690486 (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_468069318_155702435 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_fournisseur', 'nom', 'jours_livraison', 'delai_commande', 'anticipation'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_number (data -> 'id_fournisseur')
        AND structgen_validate_json_string (data -> 'nom')
        AND structgen_validate_json_array_7_boolean (data -> 'jours_livraison')
        AND structgen_validate_json_number (data -> 'delai_commande')
        AND structgen_validate_json_number (data -> 'anticipation');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_468069318_155713345 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_fournisseur', 'nom', 'jours_livraison', 'delai_commande', 'anticipation'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_number (data -> 'id_fournisseur')
        AND structgen_validate_json_string (data -> 'nom')
        AND structgen_validate_json_array_7_boolean (data -> 'jours_livraison')
        AND structgen_validate_json_number (data -> 'delai_commande')
        AND structgen_validate_json_number (data -> 'anticipation');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_map_struct_468069318_155714791 (data jsonb)
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
                bool_and(structgen_validate_json_struct_468069318_155715607 (value))
            FROM
                jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1741621487_155722960 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('Int64', 'Valid'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'Int64')
        AND structgen_validate_json_boolean (data -> 'Valid');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_825169840_155727108 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_utilisateur', 'commentaire'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_struct_1741621487_155724329 (data -> 'id_utilisateur')
        AND structgen_validate_json_string (data -> 'commentaire');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3289518958_155734000 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_ingredient', 'quantite', 'cuisson'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_ingredient')
        AND structgen_validate_json_number (data -> 'quantite')
        AND structgen_validate_json_string (data -> 'cuisson');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_948900432_155736978 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_menu', 'LienIngredient'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_menu')
        AND structgen_validate_json_struct_3289518958_155735328 (data -> 'LienIngredient');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3289518958_155743943 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_ingredient', 'quantite', 'cuisson'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_ingredient')
        AND structgen_validate_json_number (data -> 'quantite')
        AND structgen_validate_json_string (data -> 'cuisson');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_949817944_155746696 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_menu', 'LienIngredient'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_menu')
        AND structgen_validate_json_struct_3289518958_155745296 (data -> 'LienIngredient');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_951390816_155748298 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_struct_950538839_155749265 (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3019705356_155753474 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_menu', 'id_recette'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_menu')
        AND structgen_validate_json_number (data -> 'id_recette');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3019705356_155757452 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_menu', 'id_recette'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_menu')
        AND structgen_validate_json_number (data -> 'id_recette');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_3019705356_155758532 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_struct_3019705356_155759018 (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1741621487_155765547 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('Int64', 'Valid'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'Int64')
        AND structgen_validate_json_boolean (data -> 'Valid');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_844961727_155769307 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_utilisateur', 'commentaire'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_struct_1741621487_155766752 (data -> 'id_utilisateur')
        AND structgen_validate_json_string (data -> 'commentaire');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_map_struct_820778924_155772903 (data jsonb)
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
                bool_and(structgen_validate_json_struct_826808241_155773783 (value))
            FROM
                jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2334919345_155786304 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('quantite', 'unite'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'quantite')
        AND structgen_validate_json_string (data -> 'unite');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2649371107_155795268 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_livraison', 'nom', 'conditionnement', 'prix', 'reference_fournisseur', 'colisage'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_number (data -> 'id_livraison')
        AND structgen_validate_json_string (data -> 'nom')
        AND structgen_validate_json_struct_2334919345_155787869 (data -> 'conditionnement')
        AND structgen_validate_json_number (data -> 'prix')
        AND structgen_validate_json_string (data -> 'reference_fournisseur')
        AND structgen_validate_json_number (data -> 'colisage');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2334919345_155807675 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('quantite', 'unite'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'quantite')
        AND structgen_validate_json_string (data -> 'unite');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2609131992_155813646 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_livraison', 'nom', 'conditionnement', 'prix', 'reference_fournisseur', 'colisage'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_number (data -> 'id_livraison')
        AND structgen_validate_json_string (data -> 'nom')
        AND structgen_validate_json_struct_2334919345_155809057 (data -> 'conditionnement')
        AND structgen_validate_json_number (data -> 'prix')
        AND structgen_validate_json_string (data -> 'reference_fournisseur')
        AND structgen_validate_json_number (data -> 'colisage');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_map_struct_2630496734_155815201 (data jsonb)
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
                bool_and(structgen_validate_json_struct_2647994851_155816224 (value))
            FROM
                jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1741621487_155822044 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('Int64', 'Valid'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'Int64')
        AND structgen_validate_json_boolean (data -> 'Valid');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_441326466_155830002 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_utilisateur', 'nom', 'mode_emploi'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_struct_1741621487_155823209 (data -> 'id_utilisateur')
        AND structgen_validate_json_string (data -> 'nom')
        AND structgen_validate_json_string (data -> 'mode_emploi');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3289518958_155838312 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_ingredient', 'quantite', 'cuisson'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_ingredient')
        AND structgen_validate_json_number (data -> 'quantite')
        AND structgen_validate_json_string (data -> 'cuisson');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2091062152_155841050 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_recette', 'LienIngredient'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_recette')
        AND structgen_validate_json_struct_3289518958_155839585 (data -> 'LienIngredient');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3289518958_155848113 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_ingredient', 'quantite', 'cuisson'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_ingredient')
        AND structgen_validate_json_number (data -> 'quantite')
        AND structgen_validate_json_string (data -> 'cuisson');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2090668934_155850709 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_recette', 'LienIngredient'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_recette')
        AND structgen_validate_json_struct_3289518958_155849295 (data -> 'LienIngredient');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_2090996616_155853532 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_struct_2090996614_155854438 (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1741621487_155864927 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('Int64', 'Valid'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'Int64')
        AND structgen_validate_json_boolean (data -> 'Valid');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_443685763_155869624 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_utilisateur', 'nom', 'mode_emploi'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_struct_1741621487_155866180 (data -> 'id_utilisateur')
        AND structgen_validate_json_string (data -> 'nom')
        AND structgen_validate_json_string (data -> 'mode_emploi');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_map_struct_446438277_155871087 (data jsonb)
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
                bool_and(structgen_validate_json_struct_449977223_155871992 (value))
            FROM
                jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1230185621_155883191 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_sejour', 'offset_personnes', 'jour_offset', 'horaire', 'anticipation'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_number (data -> 'id_sejour')
        AND structgen_validate_json_number (data -> 'offset_personnes')
        AND structgen_validate_json_number (data -> 'jour_offset')
        AND structgen_validate_json_ (data -> 'horaire')
        AND structgen_validate_json_number (data -> 'anticipation');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3025144856_155887142 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_repas', 'id_groupe'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_repas')
        AND structgen_validate_json_number (data -> 'id_groupe');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3025144856_155891612 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_repas', 'id_groupe'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_repas')
        AND structgen_validate_json_number (data -> 'id_groupe');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_3025144856_155892800 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_struct_3025144856_155894896 (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3289518958_155901627 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_ingredient', 'quantite', 'cuisson'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_ingredient')
        AND structgen_validate_json_number (data -> 'quantite')
        AND structgen_validate_json_string (data -> 'cuisson');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1316295337_155904339 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_repas', 'LienIngredient'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_repas')
        AND structgen_validate_json_struct_3289518958_155902910 (data -> 'LienIngredient');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3289518958_155910947 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_ingredient', 'quantite', 'cuisson'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_ingredient')
        AND structgen_validate_json_number (data -> 'quantite')
        AND structgen_validate_json_string (data -> 'cuisson');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1319310012_155915218 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_repas', 'LienIngredient'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_repas')
        AND structgen_validate_json_struct_3289518958_155913824 (data -> 'LienIngredient');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_1318654648_155916762 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_struct_1318982330_155917696 (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3222080626_155921812 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_repas', 'id_recette'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_repas')
        AND structgen_validate_json_number (data -> 'id_recette');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3222080626_155925747 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_repas', 'id_recette'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_repas')
        AND structgen_validate_json_number (data -> 'id_recette');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_3222080626_155926807 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_struct_3222080626_155927288 (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1230185621_155936345 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_sejour', 'offset_personnes', 'jour_offset', 'horaire', 'anticipation'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_number (data -> 'id_sejour')
        AND structgen_validate_json_number (data -> 'offset_personnes')
        AND structgen_validate_json_number (data -> 'jour_offset')
        AND structgen_validate_json_ (data -> 'horaire')
        AND structgen_validate_json_number (data -> 'anticipation');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_map_struct_1230185621_155937558 (data jsonb)
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
                bool_and(structgen_validate_json_struct_1230185621_155938151 (value))
            FROM
                jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1_155944741 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(TRUE)
        FROM
            jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1662851379_155948373 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_utilisateur', 'date_debut', 'nom'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_number (data -> 'id_utilisateur')
        AND structgen_validate_json_struct_1_155945967 (data -> 'date_debut')
        AND structgen_validate_json_string (data -> 'nom');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2281707320_155955851 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_utilisateur', 'id_sejour', 'id_fournisseur'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_utilisateur')
        AND structgen_validate_json_number (data -> 'id_sejour')
        AND structgen_validate_json_number (data -> 'id_fournisseur');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_2281707320_155961836 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_utilisateur', 'id_sejour', 'id_fournisseur'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_utilisateur')
        AND structgen_validate_json_number (data -> 'id_sejour')
        AND structgen_validate_json_number (data -> 'id_fournisseur');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_2281707320_155962884 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_struct_2281707320_155963415 (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1_155967794 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(TRUE)
        FROM
            jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1656101162_155971553 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'id_utilisateur', 'date_debut', 'nom'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_number (data -> 'id_utilisateur')
        AND structgen_validate_json_struct_1_155968814 (data -> 'date_debut')
        AND structgen_validate_json_string (data -> 'nom');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_map_struct_1651775780_155974409 (data jsonb)
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
                bool_and(structgen_validate_json_struct_1654724903_155975212 (value))
            FROM
                jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_map_boolean (data jsonb)
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
                bool_and(structgen_validate_json_boolean (value))
            FROM
                jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_4090762348_155988693 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'password', 'mail', 'prenom_nom'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_string (data -> 'password')
        AND structgen_validate_json_string (data -> 'mail')
        AND structgen_validate_json_string (data -> 'prenom_nom');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1499271403_155992905 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_utilisateur', 'id_fournisseur'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_utilisateur')
        AND structgen_validate_json_number (data -> 'id_fournisseur');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_1499271403_156006746 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id_utilisateur', 'id_fournisseur'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id_utilisateur')
        AND structgen_validate_json_number (data -> 'id_fournisseur');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_1499271403_156008021 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    RETURN jsonb_typeof(data) = 'array'
        AND (
            SELECT
                bool_and(structgen_validate_json_struct_1499271403_156008592 (value))
            FROM
                jsonb_array_elements(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_struct_4090762348_156015487 (data jsonb)
    RETURNS boolean
    AS $f$
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(KEY IN ('id', 'password', 'mail', 'prenom_nom'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_number (data -> 'id')
        AND structgen_validate_json_string (data -> 'password')
        AND structgen_validate_json_string (data -> 'mail')
        AND structgen_validate_json_string (data -> 'prenom_nom');
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION structgen_validate_json_map_struct_4090762348_156016687 (data jsonb)
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
                bool_and(structgen_validate_json_struct_4090762348_156017340 (value))
            FROM
                jsonb_each(data));
END;
$f$
LANGUAGE 'plpgsql'
IMMUTABLE;

