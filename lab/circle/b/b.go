package b

import "minik8s/lab/circle/c"

type B struct {
	Pa a
}
type a interface {
	GetC() *c.C
}

func New(a a) *B {
	return &B{
		Pa: a,
	}
}

func (b *B) DisplayC() {
	b.Pa.GetC().Show()
}
