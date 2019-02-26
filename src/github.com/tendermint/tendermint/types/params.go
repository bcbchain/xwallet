package types

import (
	abci "github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/merkle"
)

const (
	MaxBlockSizeBytes = 104857600
)

type ConsensusParams struct {
	BlockSize	`json:"block_size_params"`
	TxSize		`json:"tx_size_params"`
	BlockGossip	`json:"block_gossip_params"`
	EvidenceParams	`json:"evidence_params"`
}

type BlockSize struct {
	MaxBytes	int	`json:"max_bytes"`
	MaxTxs		int	`json:"max_txs"`
	MaxGas		int64	`json:"max_gas"`
}

type TxSize struct {
	MaxBytes	int	`json:"max_bytes"`
	MaxGas		int64	`json:"max_gas"`
}

type BlockGossip struct {
	BlockPartSizeBytes int `json:"block_part_size_bytes"`
}

type EvidenceParams struct {
	MaxAge int64 `json:"max_age"`
}

func DefaultConsensusParams() *ConsensusParams {
	return &ConsensusParams{
		DefaultBlockSize(),
		DefaultTxSize(),
		DefaultBlockGossip(),
		DefaultEvidenceParams(),
	}
}

func DefaultBlockSize() BlockSize {
	return BlockSize{
		MaxBytes:	22020096,
		MaxTxs:		100000,
		MaxGas:		-1,
	}
}

func DefaultTxSize() TxSize {
	return TxSize{
		MaxBytes:	10240,
		MaxGas:		-1,
	}
}

func DefaultBlockGossip() BlockGossip {
	return BlockGossip{
		BlockPartSizeBytes: 65536,
	}
}

func DefaultEvidenceParams() EvidenceParams {
	return EvidenceParams{
		MaxAge: 100000,
	}
}

func (params *ConsensusParams) Validate() error {

	if params.BlockSize.MaxBytes <= 0 {
		return cmn.NewError("BlockSize.MaxBytes must be greater than 0. Got %d", params.BlockSize.MaxBytes)
	}
	if params.BlockGossip.BlockPartSizeBytes <= 0 {
		return cmn.NewError("BlockGossip.BlockPartSizeBytes must be greater than 0. Got %d", params.BlockGossip.BlockPartSizeBytes)
	}

	if params.BlockSize.MaxBytes > MaxBlockSizeBytes {
		return cmn.NewError("BlockSize.MaxBytes is too big. %d > %d",
			params.BlockSize.MaxBytes, MaxBlockSizeBytes)
	}
	return nil
}

func (params *ConsensusParams) Hash() []byte {
	return merkle.SimpleHashFromMap(map[string]merkle.Hasher{
		"block_gossip_part_size_bytes":	aminoHasher(params.BlockGossip.BlockPartSizeBytes),
		"block_size_max_bytes":		aminoHasher(params.BlockSize.MaxBytes),
		"block_size_max_gas":		aminoHasher(params.BlockSize.MaxGas),
		"block_size_max_txs":		aminoHasher(params.BlockSize.MaxTxs),
		"tx_size_max_bytes":		aminoHasher(params.TxSize.MaxBytes),
		"tx_size_max_gas":		aminoHasher(params.TxSize.MaxGas),
	})
}

func (params ConsensusParams) Update(params2 *abci.ConsensusParams) ConsensusParams {
	res := params

	if params2 == nil {
		return res
	}

	if params2.BlockSize != nil {
		if params2.BlockSize.MaxBytes > 0 {
			res.BlockSize.MaxBytes = int(params2.BlockSize.MaxBytes)
		}
		if params2.BlockSize.MaxTxs > 0 {
			res.BlockSize.MaxTxs = int(params2.BlockSize.MaxTxs)
		}
		if params2.BlockSize.MaxGas > 0 {
			res.BlockSize.MaxGas = params2.BlockSize.MaxGas
		}
	}
	if params2.TxSize != nil {
		if params2.TxSize.MaxBytes > 0 {
			res.TxSize.MaxBytes = int(params2.TxSize.MaxBytes)
		}
		if params2.TxSize.MaxGas > 0 {
			res.TxSize.MaxGas = params2.TxSize.MaxGas
		}
	}
	if params2.BlockGossip != nil {
		if params2.BlockGossip.BlockPartSizeBytes > 0 {
			res.BlockGossip.BlockPartSizeBytes = int(params2.BlockGossip.BlockPartSizeBytes)
		}
	}
	return res
}
