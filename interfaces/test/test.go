package test

type union1 interface {
	isUnion1()
}

func (member1) isUnion1() {}
func (member2) isUnion1() {}
func (member3) isUnion1() {}

type member1 int

type member2 struct {
	B int
}

type union2 interface {
	isUnion2()
}

func (member1) isUnion2() {}

type ITFSlice []union2
