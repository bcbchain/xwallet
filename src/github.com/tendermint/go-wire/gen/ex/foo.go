package ex

import (
	"fmt"

	"github.com/tendermint/tmlibs/common"
)

type FooInner interface {
	Bar() int
}

type Bling struct{}

func (b Bling) Bar() int {
	return common.RandInt()
}

type Fuzz struct{}

func (f *Fuzz) Bar() int {
	fmt.Println("hello")
	return 42
}
