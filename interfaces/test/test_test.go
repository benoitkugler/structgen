package test

import (
	"reflect"
	"testing"
)

func TestMarshallJSON(t *testing.T) {
	var i union1 = member2{B: 4}
	b, err := union1Wrapper{data: i}.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var tmp union1Wrapper
	err = tmp.UnmarshalJSON(b)
	if err != nil {
		t.Fatal(err)
	}
	i2 := tmp.data
	if !reflect.DeepEqual(i, i2) {
		t.Fatal(i, i2)
	}
}
