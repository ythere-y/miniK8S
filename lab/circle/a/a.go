package a

import (
	"minik8s/lab/circle/b"
	"minik8s/lab/circle/c"
)

type A struct {
	Pb *b.B
	Pc *c.C
}

func New(ic int) *A {
	a := &A{
		Pc: c.New(ic),
	}

	a.Pb = b.New(a)

	return a
}
func (a *A) GetC() *c.C {
	return a.Pc
}
