package ex

import (
	"github.com/tendermint/go-wire/data"
)

type Foo struct {
	FooInner "json:\"unwrap\""
}

var FooMapper = data.NewMapper(Foo{})

func (h Foo) MarshalJSON() ([]byte, error) {
	return FooMapper.ToJSON(h.FooInner)
}

func (h *Foo) UnmarshalJSON(data []byte) (err error) {
	parsed, err := FooMapper.FromJSON(data)
	if err == nil && parsed != nil {
		h.FooInner = parsed.(FooInner)
	}
	return err
}

func (h Foo) Unwrap() FooInner {
	hi := h.FooInner
	for wrap, ok := hi.(Foo); ok; wrap, ok = hi.(Foo) {
		hi = wrap.FooInner
	}
	return hi
}

func (h Foo) Empty() bool {
	return h.FooInner == nil
}

func init() {
	FooMapper.RegisterImplementation(Bling{}, "blng", 0x1)
}

func (hi Bling) Wrap() Foo {
	return Foo{hi}
}

func init() {
	FooMapper.RegisterImplementation(&Fuzz{}, "fzz", 0x2)
}

func (hi *Fuzz) Wrap() Foo {
	return Foo{hi}
}
