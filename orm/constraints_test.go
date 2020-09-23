package orm

import (
	"fmt"
	"testing"

	"github.com/benoitkugler/structgen/loader"
)

func TestEnums(t *testing.T) {
	s := "smdmsmdl #Enum1.Value1_1 #Enum2.Value1)"
	_, err := parseComment(s, map[string]string{
		"Enum1.Value1_1": "7879",
		"Enum2.Value1":   `'string value'`,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseUnique(t *testing.T) {
	fmt.Println(IsUniqueConstraint(loader.Comment{Content: "ADD UNIQUE(id_camp, id_personne,  id_groupe)"})) // ""
	fmt.Println(IsUniqueConstraint(loader.Comment{Content: "ADD UNIQUE(id_camp)"}))                          // "id_camp"
	fmt.Println(IsUniqueConstraint(loader.Comment{Content: "ADD UNIQUE(id_camp )"}))                         // "id_camp"
}
