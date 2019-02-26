package test

import (
	cmn "github.com/tendermint/tmlibs/common"
)

func MutateByteSlice(bytez []byte) []byte {

	if len(bytez) == 0 {
		panic("Cannot mutate an empty bytez")
	}

	mBytez := make([]byte, len(bytez))
	copy(mBytez, bytez)
	bytez = mBytez

	switch cmn.RandInt() % 2 {
	case 0:
		bytez[cmn.RandInt()%len(bytez)] += byte(cmn.RandInt()%255 + 1)
	case 1:
		pos := cmn.RandInt() % len(bytez)
		bytez = append(bytez[:pos], bytez[pos+1:]...)
	}
	return bytez
}
