package test

type DataType struct {
	DA int
	DB [4]bool
	DC []int
	DD MyEnumI
	DE MyEnumS
}

type Model struct {
	A int
	B string
	C []int
	D bool
	E map[string][]string
	F DataType
	G L
}

type MyEnumI int

type MyEnumS string

type L []DataType
