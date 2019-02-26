package null

import (
	"errors"

	"github.com/tendermint/tendermint/state/txindex"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tmlibs/pubsub/query"
)

var _ txindex.TxIndexer = (*TxIndex)(nil)

type TxIndex struct{}

func (txi *TxIndex) Get(hash []byte) (*types.TxResult, error) {
	return nil, errors.New(`Indexing is disabled (set 'tx_index = "kv"' in config)`)
}

func (txi *TxIndex) AddBatch(batch *txindex.Batch) error {
	return nil
}

func (txi *TxIndex) Index(result *types.TxResult) error {
	return nil
}

func (txi *TxIndex) Search(q *query.Query) ([]*types.TxResult, error) {
	return []*types.TxResult{}, nil
}
