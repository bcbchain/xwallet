package smc

import (
	"github.com/tendermint/tmlibs/common"
	"bcbchain.io/bcerrors"
	"golang.org/x/crypto/sha3"
	"math/big"
)

type Address = string

type ReceiptBytes []byte

type Hash = common.HexBytes

type Chromo = string

type Error = bcerrors.BCError

const HASH_LEN = 32

type PubKey = common.HexBytes

const PUBKEY_LEN = 32

const Max_Gas_Price = 1000000000

const Max_Name_Len = 40

func (r ReceiptBytes) String() string {
	return string(r)
}

type Receipt struct {
	Name		string	`json:"name"`
	ContractAddress	Address	`json:"contractAddress"`
	ReceiptBytes	[]byte	`json:"receiptBytes"`
	ReceiptHash	Hash	`json:"receiptHash"`
}

type ReceiptOfTransfer struct {
	Token	Address	`json:"token"`
	From	Address	`json:"from"`
	To	Address	`json:"to"`
	Value	big.Int	`json:"value"`
}

type ReceiptOfAddSupply struct {
	Token		Address	`json:"token"`
	Value		big.Int	`json:"value"`
	TotalSupply	big.Int	`json:"totalSupply"`
}

type ReceiptOfBurn struct {
	Token		Address	`json:"token"`
	Value		big.Int	`json:"value"`
	TotalSupply	big.Int	`json:"totalSupply"`
}

type ReceiptOfSetGasPrice struct {
	Token		Address	`json:"token"`
	GasPrice	uint64	`json:"gasPrice"`
}

type ReceiptOfSetOwner struct {
	ContractAddr	Address	`json:"contractAddr"`
	NewOwner	Address	`json:"newOwner"`
}

type ReceiptOfFee struct {
	Token	Address	`json:"token"`
	From	Address	`json:"from"`
	Value	uint64	`json:"value"`
}

type ReceiptOfNewToken struct {
	TokenAddress		Address	`json:"tokenAddress"`
	ContractAddress		Address	`json:"contractAddress"`
	AccountAddress		Address	`json:"accountAddress"`
	Owner			Address	`json:"owner"`
	Version			string	`json:"version"`
	Name			string	`json:"name,omitempty"`
	Symbol			string	`json:"symbol"`
	TotalSupply		big.Int	`json:"totalSupply"`
	AddSupplyEnabled	bool	`json:"addSupplyEnabled"`
	BurnEnabled		bool	`json:"burnEnabled"`
	GasPrice		uint64	`json:"gasprice"`
}

func CalcReceiptHash(name string, addr Address, receiptByte []byte) Hash {
	hasherSHA3256 := sha3.New256()
	hasherSHA3256.Write([]byte(name))
	hasherSHA3256.Write([]byte(addr))
	hasherSHA3256.Write(receiptByte)

	return hasherSHA3256.Sum(nil)
}
