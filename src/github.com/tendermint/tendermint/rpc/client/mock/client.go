package mock

import (
	"reflect"

	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/core"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

type Client struct {
	client.ABCIClient
	client.SignClient
	client.HistoryClient
	client.StatusClient
	client.EventsClient
	cmn.Service
}

var _ client.Client = Client{}

type Call struct {
	Name		string
	Args		interface{}
	Response	interface{}
	Error		error
}

func (c Call) GetResponse(args interface{}) (interface{}, error) {

	if c.Response == nil {
		if c.Error == nil {
			panic("Misconfigured call, you must set either Response or Error")
		}
		return nil, c.Error
	}

	if c.Error == nil {
		return c.Response, nil
	}

	if reflect.DeepEqual(args, c.Args) {
		return c.Response, nil
	}
	return nil, c.Error
}

func (c Client) Status() (*ctypes.ResultStatus, error) {
	return core.Status()
}

func (c Client) ABCIInfo() (*ctypes.ResultABCIInfo, error) {
	return core.ABCIInfo()
}

func (c Client) ABCIQuery(path string, data cmn.HexBytes) (*ctypes.ResultABCIQuery, error) {
	return c.ABCIQueryWithOptions(path, data, client.DefaultABCIQueryOptions)
}

func (c Client) ABCIQueryWithOptions(path string, data cmn.HexBytes, opts client.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	return core.ABCIQuery(path, data, opts.Height, opts.Trusted)
}

func (c Client) BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	return core.BroadcastTxCommit(tx)
}

func (c Client) BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return core.BroadcastTxAsync(tx)
}

func (c Client) BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return core.BroadcastTxSync(tx)
}

func (c Client) NetInfo() (*ctypes.ResultNetInfo, error) {
	return core.NetInfo()
}

func (c Client) DialSeeds(seeds []string) (*ctypes.ResultDialSeeds, error) {
	return core.UnsafeDialSeeds(seeds)
}

func (c Client) DialPeers(peers []string, persistent bool) (*ctypes.ResultDialPeers, error) {
	return core.UnsafeDialPeers(peers, persistent)
}

func (c Client) BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	return core.BlockchainInfo(minHeight, maxHeight)
}

func (c Client) Genesis() (*ctypes.ResultGenesis, error) {
	return core.Genesis()
}

func (c Client) Block(height *int64) (*ctypes.ResultBlock, error) {
	return core.Block(height)
}

func (c Client) Commit(height *int64) (*ctypes.ResultCommit, error) {
	return core.Commit(height)
}

func (c Client) Validators(height *int64) (*ctypes.ResultValidators, error) {
	return core.Validators(height)
}
