package test

type s struct {
	V int
}

type b struct {
	D float64
}

type i interface {
	isI()
}

func (s) isI() {}
func (b) isI() {}

var (
	_ i = s{}
	_ i = b{}
)

type model struct {
	A     int
	Value i
}
