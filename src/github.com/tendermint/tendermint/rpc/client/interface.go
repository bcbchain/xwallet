package client

import (
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

type ABCIClient interface {
	ABCIInfo() (*ctypes.ResultABCIInfo, error)
	ABCIQuery(path string, data cmn.HexBytes) (*ctypes.ResultABCIQuery, error)
	ABCIQueryWithOptions(path string, data cmn.HexBytes,
		opts ABCIQueryOptions) (*ctypes.ResultABCIQuery, error)

	BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error)
	BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTx, error)
	BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error)
}

type SignClient interface {
	Block(height *int64) (*ctypes.ResultBlock, error)
	BlockResults(height *int64) (*ctypes.ResultBlockResults, error)
	Commit(height *int64) (*ctypes.ResultCommit, error)
	Validators(height *int64) (*ctypes.ResultValidators, error)
	Tx(hash []byte, prove bool) (*ctypes.ResultTx, error)
	TxSearch(query string, prove bool) ([]*ctypes.ResultTx, error)
}

type HistoryClient interface {
	Genesis() (*ctypes.ResultGenesis, error)
	BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error)
}

type StatusClient interface {
	Status() (*ctypes.ResultStatus, error)
	NumUnConfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error)
}

type Client interface {
	cmn.Service
	ABCIClient
	SignClient
	HistoryClient
	StatusClient
	EventsClient
}

type NetworkClient interface {
	NetInfo() (*ctypes.ResultNetInfo, error)
	DumpConsensusState() (*ctypes.ResultDumpConsensusState, error)
	Health() (*ctypes.ResultHealth, error)
}

type EventsClient interface {
	types.EventBusSubscriber
}
