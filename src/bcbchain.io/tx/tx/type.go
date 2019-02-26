package tx

import (
	"bcbchain.io/keys"
	"math/big"
)

type MethodInfo struct {
	MethodID	uint32
	ParamData	[]byte
}

type BigNumber = big.Int

type Method struct {
	MethodID	uint32
	Prototype	string
}

type Transaction struct {
	Nonce		uint64
	GasLimit	uint64
	Note		string
	To		keys.Address
	Data		[]byte
}

const MAX_SIZE_NOTE = 256

type Query struct {
	QueryKey string
}

func NewTransaction(nonce uint64, gaslimit uint64, note string, to keys.Address, data []byte) Transaction {
	tx := Transaction{
		Nonce:		nonce,
		GasLimit:	gaslimit,
		Note:		note,
		To:		to,
		Data:		data,
	}
	return tx
}
