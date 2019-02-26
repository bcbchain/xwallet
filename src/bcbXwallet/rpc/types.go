package rpc

import (
	"bcbchain.io/keys"
)

const transferMethodID = "af0228bc"

type TransferParam struct {
	SmcAddress	keys.Address	`json:"smcAddress"`
	GasLimit	string		`json:"gasLimit"`
	Note		string		`json:"note"`
	To		keys.Address	`json:"to"`
	Value		string		`json:"value"`
}

type TransferOfflineParam struct {
	SmcAddress	keys.Address	`json:"smcAddress"`
	GasLimit	string		`json:"gasLimit"`
	Note		string		`json:"note"`
	Nonce		uint64		`json:"nonce"`
	To		keys.Address	`json:"to"`
	Value		string		`json:"value"`
}

type WalletCreateResult struct {
	AccessKey	string		`json:"accessKey"`
	WalletAddress	keys.Address	`json:"walletAddr"`
}

type WalletExportResult struct {
	PrivateKey	string		`json:"privateKey"`
	WalletAddress	keys.Address	`json:"walletAddr"`
}

type WalletImportResult struct {
	AccessKey	string		`json:"accessKey"`
	WalletAddress	keys.Address	`json:"walletAddr"`
}

type WalletListResult struct {
	Total		uint64		`json:"total"`
	WalletList	[]WalletItem	`json:"walletList"`
}

type WalletItem struct {
	Name		string		`json:"name"`
	WalletAddress	keys.Address	`json:"walletAddr"`
}

type TransferResult struct {
	Code	uint32	`json:"code"`
	Log	string	`json:"log"`
	Fee	uint64	`json:"fee"`
	TxHash	string	`json:"txHash"`
	Height	int64	`json:"height"`
}

type TransferOfflineResult struct {
	Tx string `json:"tx"`
}

type BlockHeightResult struct {
	LastBlock int64 `json:"lastBlock"`
}

type Message struct {
	SmcAddress	keys.Address	`json:"smcAddress"`
	SmcName		string		`json:"smcName"`
	Method		string		`json:"method"`
	To		string		`json:"to"`
	Value		string		`json:"value"`
}

type TxResult struct {
	TxHash		string		`json:"txHash"`
	TxTime		string		`json:"txTime"`
	Code		uint32		`json:"code"`
	Log		string		`json:"log"`
	BlockHash	string		`json:"blockHash"`
	BlockHeight	int64		`json:"blockHeight"`
	From		keys.Address	`json:"from"`
	Nonce		uint64		`json:"nonce"`
	GasLimit	uint64		`json:"gasLimit"`
	Fee		uint64		`json:"fee"`
	Note		string		`json:"note"`
	Messages	[]Message	`json:"messages"`
}

type BlockResult struct {
	BlockHeight	int64		`json:"blockHeight"`
	BlockHash	string		`json:"blockHash"`
	ParentHash	string		`json:"parentHash"`
	ChainID		string		`json:"chainID"`
	ValidatorHash	string		`json:"validatorHash"`
	ConsensusHash	string		`json:"consensusHash"`
	BlockTime	string		`json:"blockTime"`
	BlockSize	int		`json:"blockSize"`
	ProposerAddress	keys.Address	`json:"proposerAddress"`
	Txs		[]TxResult	`json:"txs"`
}

type BalanceResult struct {
	Balance string `json:"balance"`
}

type AllBalanceItemResult struct {
	TokenAddress	keys.Address	`json:"tokenAddress"`
	TokenName	string		`json:"tokenName"`
	Balance		string		`json:"balance"`
}

type NonceResult struct {
	Nonce uint64 `json:"nonce"`
}

type CommitTxResult struct {
	Code	uint32	`json:"code"`
	Log	string	`json:"log"`
	Fee	uint64	`json:"fee"`
	TxHash	string	`json:"txHash"`
	Height	int64	`json:"height"`
}

type VersionResult struct {
	Version string `json:"version"`
}

type MethodInfo struct {
	MethodID	uint32
	ParamData	[]byte
}
