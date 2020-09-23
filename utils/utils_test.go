package utils

import (
	"fmt"
	"testing"
)

func TestComment(t *testing.T) {
	fmt.Println(IsSpecialComment("// dlskl:mùsldsllùmd"))
	fmt.Println(IsSpecialComment("// lmsksmkd"))
}
