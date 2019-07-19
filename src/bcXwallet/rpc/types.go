package rpc

import (
	"blockchain/abciapp_v1.0/keys"
)

const transferMethodIDV1 = "af0228bc"
const transferMethodIDV2 = "44d8ca60"

// ----- param struct ----
type TransferParam struct {
	SmcAddress keys.Address `json:"smcAddress"`
	GasLimit   string       `json:"gasLimit"`
	Note       string       `json:"note"`
	To         keys.Address `json:"to"`
	Value      string       `json:"value"`
}

type TransferOfflineParam struct {
	SmcAddress keys.Address `json:"smcAddress"`
	GasLimit   string       `json:"gasLimit"`
	Note       string       `json:"note"`
	Nonce      uint64       `json:"nonce"`
	To         keys.Address `json:"to"`
	Value      string       `json:"value"`
}

// ----- result struct -----
// WalletCreateResult - create wallet result
type WalletCreateResult struct {
	AccessKey     string       `json:"accessKey"`
	WalletAddress keys.Address `json:"walletAddr"`
}

// WalletExportResult - export wallet result
type WalletExportResult struct {
	PrivateKey    string       `json:"privateKey"`
	WalletAddress keys.Address `json:"walletAddr"`
}

// WalletImportResult - import wallet result
type WalletImportResult struct {
	AccessKey     string       `json:"accessKey"`
	WalletAddress keys.Address `json:"walletAddr"`
}

// WalletListResult - list wallet
type WalletListResult struct {
	Total      uint64       `json:"total"`
	WalletList []WalletItem `json:"walletList"`
}

// WalletItemResult - wallet item
type WalletItem struct {
	Name          string       `json:"name"`
	WalletAddress keys.Address `json:"walletAddr"`
}

// TransferResult - transfer result
type TransferResult struct {
	Code   uint32 `json:"code"`
	Log    string `json:"log"`
	Fee    uint64 `json:"fee"`
	TxHash string `json:"txHash"`
	Height int64  `json:"height"`
}

// TransferResult - transfer result
type TransferOfflineResult struct {
	Tx string `json:"tx"`
}

// BlockHeightResult - block height result
type BlockHeightResult struct {
	LastBlock int64 `json:"lastBlock"`
}

// Message - message struct
type Message struct {
	SmcAddress keys.Address `json:"smcAddress"`
	SmcName    string       `json:"smcName"`
	Method     string       `json:"method"`
	To         string       `json:"to"`
	Value      string       `json:"value"`
}

// TxResult - transaction struct
type TxResult struct {
	TxHash      string       `json:"txHash"`
	TxTime      string       `json:"txTime"`
	Code        uint32       `json:"code"`
	Log         string       `json:"log"`
	BlockHash   string       `json:"blockHash"`
	BlockHeight int64        `json:"blockHeight"`
	From        keys.Address `json:"from"`
	Nonce       uint64       `json:"nonce"`
	GasLimit    uint64       `json:"gasLimit"`
	Fee         uint64       `json:"fee"`
	Note        string       `json:"note"`
	Messages    []Message    `json:"messages"`
}

// BlockResult - block struct
type BlockResult struct {
	BlockHeight     int64        `json:"blockHeight,omitempty"`
	BlockHash       string       `json:"blockHash,omitempty"`
	ParentHash      string       `json:"parentHash,omitempty"`
	ChainID         string       `json:"chainID,omitempty"`
	ValidatorHash   string       `json:"validatorHash,omitempty"`
	ConsensusHash   string       `json:"consensusHash,omitempty"`
	BlockTime       string       `json:"blockTime,omitempty"`
	BlockSize       int          `json:"blockSize,omitempty"`
	ProposerAddress keys.Address `json:"proposerAddress,omitempty"`
	Txs             []TxResult   `json:"txs,omitempty"`

	// simple result contain several blocks
	Result []SimpleBlockResult `json:"result,omitempty"`
}

// SimpleBlockResult simple block information contain height,hash and time
type SimpleBlockResult struct {
	BlockHeight int64  `json:"blockHeight"`
	BlockHash   string `json:"blockHash"`
	BlockTime   string `json:"blockTime"`
}

// BalanceResult - balance struct
type BalanceResult struct {
	Balance string `json:"balance"`
}

// AllBalanceItemResult - item of all balance struct
type AllBalanceItemResult struct {
	TokenAddress keys.Address `json:"tokenAddress"`
	TokenName    string       `json:"tokenName"`
	Balance      string       `json:"balance"`
}

// NonceResult - nonce struct
type NonceResult struct {
	Nonce uint64 `json:"nonce"`
}

// CommitTxResult - commit tx result
type CommitTxResult struct {
	Code   uint32 `json:"code"`
	Log    string `json:"log"`
	Fee    uint64 `json:"fee"`
	TxHash string `json:"txHash"`
	Height int64  `json:"height"`
}

// VersionResult - version struct
type VersionResult struct {
	Version string `json:"version"`
}

//-------------------------------------
// 定义交易数据结构

type MethodInfo struct {
	MethodID  uint32
	ParamData []byte
}
