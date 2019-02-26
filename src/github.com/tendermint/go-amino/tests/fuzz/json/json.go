package fuzz_json

import (
	"github.com/tendermint/go-amino"
	"github.com/tendermint/go-amino/tests"
)

func Fuzz(data []byte) int {
	cdc := amino.NewCodec()
	cst := tests.ComplexSt{}
	err := cdc.UnmarshalJSON(data, &cst)
	if err != nil {
		return 0
	}
	return 1
}