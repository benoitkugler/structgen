package dart

import (
	"fmt"
	"testing"

	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
)

func TestDart(t *testing.T) {
	pa, err := loader.LoadSource("../testenums.go")
	if err != nil {
		t.Fatal(err)
	}
	l, err := enums.FetchEnums(pa)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range l {
		fmt.Println(EnumAsDart(e))
	}
}
