package test

import (
	"reflect"
	"testing"
)

func TestMarshallJSON(t *testing.T) {
	var i union1 = member2{B: 4}
	b, err := union1MarshalJSON(i)
	if err != nil {
		t.Fatal(err)
	}
	i2, err := union1UnmarshalJSON(b)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(i, i2) {
		t.Fatal(i, i2)
	}
}
