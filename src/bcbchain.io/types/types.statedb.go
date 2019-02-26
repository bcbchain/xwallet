package types

import (
	"bcbchain.io/smc"
	"github.com/tendermint/go-crypto"
	"math/big"
)

type IssueToken struct {
	Address			smc.Address	`json:"address"`
	Owner			smc.Address	`json:"owner"`
	Version			string		`json:"version"`
	Name			string		`json:"name,omitempty"`
	Symbol			string		`json:"symbol"`
	TotalSupply		big.Int		`json:"totalSupply"`
	AddSupplyEnabled	bool		`json:"addSupplyEnabled"`
	BurnEnabled		bool		`json:"burnEnabled"`
	GasPrice		uint64		`json:"gasprice"`
}

type TokenBalance struct {
	Address	smc.Address	`json:"address"`
	Balance	big.Int		`json:"balance"`
}
type TokenBalances []TokenBalance

type Method struct {
	MethodId	string	`json:"methodId,omitempty"`
	Gas		uint64	`json:"gas,omitempty"`
	Prototype	string	`json:"prototype,omitempty"`
}

type Contract struct {
	Address		smc.Address	`json:"address,omitempty"`
	Owner		smc.Address	`json:"owner,omitempty"`
	Name		string		`json:"name,omitempty"`
	Version		string		`json:"version,omitempty"`
	CodeHash	string		`json:"codeHash,omitempty"`
	Methods		[]Method	`json:"methods,omitempty"`
	EffectHeight	uint64		`json:"effectHeight,omitempty"`
	LoseHeight	uint64		`json:"loseHeight,omitempty"`
}

type RewardStrategy struct {
	Strategy	[]Rewarder	`json:"rewardStrategy,omitempty"`
	EffectHeight	uint64		`json:"effectHeight,omitempty"`
}

type Rewarder struct {
	Name		string	`json:"name"`
	RewardPercent	string	`json:"rewardPercent"`
	Address		string	`json:"address"`
}

type AccountInfo struct {
	Nonce uint64
}

type TokenFee struct {
	MaxFee	uint64	`json:"maxFee"`
	MinFee	uint64	`json:"minFee"`
	Ratio	uint64	`json:"ratio"`
}

type AccountFee struct {
	Fee	uint64	`json:"fee"`
	Payer	string	`json:"payer"`
}

type UDCOrder struct {
	UDCState	string		`json:"udcstate,omitempty"`
	UDCHash		crypto.Hash	`json:"udchash,omitempty"`
	Nonce		uint64		`json:"nonce,omitempty"`
	ContractAddr	smc.Address	`json:"contractaddr,omitempty"`
	Owner		smc.Address	`json:"owner,omitempty"`
	Value		big.Int		`json:"value,omitempty"`
	MatureDate	string		`json:"maturedate,omitempty"`
}

type Validator struct {
	Name		string		`json:"name,omitempty"`
	NodePubKey	smc.PubKey	`json:"nodepubkey,omitempty"`
	NodeAddr	smc.Address	`json:"nodeaddr,omitempty"`
	RewardAddr	smc.Address	`json:"rewardaddr,omitempty"`
	Power		uint64		`json:"power,omitempty"`
}
