package client

import (
	"context"

	"encoding/hex"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/rpc/core"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
	tmpubsub "github.com/tendermint/tmlibs/pubsub"
)

type Local struct {
	*types.EventBus
}

func NewLocal(node *nm.Node) *Local {
	node.ConfigureRPC()
	return &Local{
		EventBus: node.EventBus(),
	}
}

var (
	_	Client		= (*Local)(nil)
	_	NetworkClient	= Local{}
	_	EventsClient	= (*Local)(nil)
)

func (Local) Status() (*ctypes.ResultStatus, error) {
	return core.Status()
}

func (Local) NumUnConfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	return core.NumUnconfirmedTxs()
}

func (Local) ABCIInfo() (*ctypes.ResultABCIInfo, error) {
	return core.ABCIInfo()
}

func (c *Local) ABCIQuery(path string, data cmn.HexBytes) (*ctypes.ResultABCIQuery, error) {
	return c.ABCIQueryWithOptions(path, data, DefaultABCIQueryOptions)
}

func (Local) ABCIQueryWithOptions(path string, data cmn.HexBytes, opts ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	return core.ABCIQuery(path, data, opts.Height, opts.Trusted)
}

func (Local) BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	return core.BroadcastTxCommit(tx)
}

func (Local) BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return core.BroadcastTxAsync(tx)
}

func (Local) BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return core.BroadcastTxSync(tx)
}

func (Local) NetInfo() (*ctypes.ResultNetInfo, error) {
	return core.NetInfo()
}

func (Local) DumpConsensusState() (*ctypes.ResultDumpConsensusState, error) {
	return core.DumpConsensusState()
}

func (Local) Health() (*ctypes.ResultHealth, error) {
	return core.Health()
}

func (Local) DialSeeds(seeds []string) (*ctypes.ResultDialSeeds, error) {
	return core.UnsafeDialSeeds(seeds)
}

func (Local) DialPeers(peers []string, persistent bool) (*ctypes.ResultDialPeers, error) {
	return core.UnsafeDialPeers(peers, persistent)
}

func (Local) BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	return core.BlockchainInfo(minHeight, maxHeight)
}

func (Local) Genesis() (*ctypes.ResultGenesis, error) {
	return core.Genesis()
}

func (Local) Block(height *int64) (*ctypes.ResultBlock, error) {
	return core.Block(height)
}

func (Local) BlockResults(height *int64) (*ctypes.ResultBlockResults, error) {
	return core.BlockResults(height)
}

func (Local) Commit(height *int64) (*ctypes.ResultCommit, error) {
	return core.Commit(height)
}

func (Local) Validators(height *int64) (*ctypes.ResultValidators, error) {
	return core.Validators(height)
}

func (Local) Tx(hash []byte, prove bool) (*ctypes.ResultTx, error) {
	return core.Tx(hex.EncodeToString(hash), prove)
}

func (Local) TxSearch(query string, prove bool) ([]*ctypes.ResultTx, error) {
	return core.TxSearch(query, prove)
}

func (c *Local) Subscribe(ctx context.Context, subscriber string, query tmpubsub.Query, out chan<- interface{}) error {
	return c.EventBus.Subscribe(ctx, subscriber, query, out)
}

func (c *Local) Unsubscribe(ctx context.Context, subscriber string, query tmpubsub.Query) error {
	return c.EventBus.Unsubscribe(ctx, subscriber, query)
}

func (c *Local) UnsubscribeAll(ctx context.Context, subscriber string) error {
	return c.EventBus.UnsubscribeAll(ctx, subscriber)
}
