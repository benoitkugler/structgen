package enums

import (
	"fmt"
	"testing"
	"time"

	"github.com/benoitkugler/structgen/loader"
)

func TestParse(t *testing.T) {
	// fn := "../../goACVE/server/directeurs/types.go"
	fn := "testenums.go"
	pa, err := loader.LoadSource(fn)
	if err != nil {
		t.Fatal(err)
	}
	ti := time.Now()
	l, err := FetchEnums(pa)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("recursively fetching enums :", time.Since(ti))
	if len(l) != 2 {
		t.Errorf("expected 2 enums, got %d", len(l))
	}
}
