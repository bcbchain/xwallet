package types

import (
	abci "github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
)

type Mempool interface {
	Lock()
	Unlock()

	Size() int
	CheckTx(Tx, func(*abci.Response)) error
	Reap(int) Txs
	Update(height int64, txs Txs) error
	Flush()
	FlushAppConn() error

	TxsAvailable() <-chan int64
	EnableTxsAvailable()
	GiTxSearch(string) (*abci.ResponseCheckTx, error)
	GiTxCache(cmn.HexBytes, interface{})
}

type MockMempool struct {
}

func (m MockMempool) Lock()							{}
func (m MockMempool) Unlock()							{}
func (m MockMempool) Size() int							{ return 0 }
func (m MockMempool) CheckTx(tx Tx, cb func(*abci.Response)) error		{ return nil }
func (m MockMempool) Reap(n int) Txs						{ return Txs{} }
func (m MockMempool) Update(height int64, txs Txs) error			{ return nil }
func (m MockMempool) Flush()							{}
func (m MockMempool) FlushAppConn() error					{ return nil }
func (m MockMempool) TxsAvailable() <-chan int64				{ return make(chan int64) }
func (m MockMempool) EnableTxsAvailable()					{}
func (m MockMempool) GiTxSearch(tx string) (*abci.ResponseCheckTx, error)	{ return nil, nil }
func (m MockMempool) GiTxCache(tx cmn.HexBytes, a interface{})			{}

type BlockStoreRPC interface {
	Height() int64

	LoadBlockMeta(height int64) *BlockMeta
	LoadBlock(height int64) *Block
	LoadBlockPart(height int64, index int) *Part

	LoadBlockCommit(height int64) *Commit
	LoadSeenCommit(height int64) *Commit
}

type BlockStore interface {
	BlockStoreRPC
	SaveBlock(block *Block, blockParts *PartSet, seenCommit *Commit)
}

type EvidencePool interface {
	PendingEvidence() []Evidence
	AddEvidence(Evidence) error
	Update(*Block)
}

type MockEvidencePool struct {
}

func (m MockEvidencePool) PendingEvidence() []Evidence	{ return nil }
func (m MockEvidencePool) AddEvidence(Evidence) error	{ return nil }
func (m MockEvidencePool) Update(*Block)		{}
