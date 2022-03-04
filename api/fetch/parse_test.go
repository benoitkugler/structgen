package fetch

import (
	"fmt"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	pack, file, err := LoadSource("test/routes.go")
	if err != nil {
		t.Fatal(err)
	}
	ti := time.Now()
	apis := Parse(pack, file)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(time.Since(ti))
	fmt.Println(apis)
}
