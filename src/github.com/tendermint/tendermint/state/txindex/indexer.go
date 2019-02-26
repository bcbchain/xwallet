package txindex

import (
	"errors"

	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tmlibs/pubsub/query"
)

type TxIndexer interface {
	AddBatch(b *Batch) error

	Index(result *types.TxResult) error

	Get(hash []byte) (*types.TxResult, error)

	Search(q *query.Query) ([]*types.TxResult, error)
}

type Batch struct {
	Ops []*types.TxResult
}

func NewBatch(n int) *Batch {
	return &Batch{
		Ops: make([]*types.TxResult, n),
	}
}

func (b *Batch) Add(result *types.TxResult) error {
	b.Ops[result.Index] = result
	return nil
}

func (b *Batch) Size() int {
	return len(b.Ops)
}

var ErrorEmptyHash = errors.New("Transaction hash cannot be empty")
