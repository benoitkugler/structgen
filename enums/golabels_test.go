package enums

import (
	"fmt"
	"testing"

	"github.com/benoitkugler/structgen/loader"
)

func TestGo(t *testing.T) {
	pa, err := loader.LoadSource("testenums.go")
	if err != nil {
		t.Fatal(err)
	}
	l, err := FetchEnums(pa)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range l {
		fmt.Println(e.labels())
	}
}
