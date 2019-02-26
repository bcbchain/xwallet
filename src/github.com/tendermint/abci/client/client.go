package abcicli

import (
	"fmt"
	"sync"

	"github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
)

const (
	dialRetryIntervalSeconds	= 3
	echoRetryIntervalSeconds	= 1
)

type Client interface {
	cmn.Service

	SetResponseCallback(Callback)
	Error() error

	FlushAsync() *ReqRes
	EchoAsync(msg string) *ReqRes
	InfoAsync(types.RequestInfo) *ReqRes
	SetOptionAsync(types.RequestSetOption) *ReqRes
	DeliverTxAsync(tx []byte) *ReqRes
	CheckTxAsync(tx []byte) *ReqRes
	QueryAsync(types.RequestQuery) *ReqRes
	CommitAsync() *ReqRes
	InitChainAsync(types.RequestInitChain) *ReqRes
	BeginBlockAsync(types.RequestBeginBlock) *ReqRes
	EndBlockAsync(types.RequestEndBlock) *ReqRes

	FlushSync() error
	EchoSync(msg string) (*types.ResponseEcho, error)
	InfoSync(types.RequestInfo) (*types.ResponseInfo, error)
	SetOptionSync(types.RequestSetOption) (*types.ResponseSetOption, error)
	DeliverTxSync(tx []byte) (*types.ResponseDeliverTx, error)
	CheckTxSync(tx []byte) (*types.ResponseCheckTx, error)
	QuerySync(types.RequestQuery) (*types.ResponseQuery, error)
	CommitSync() (*types.ResponseCommit, error)
	InitChainSync(types.RequestInitChain) (*types.ResponseInitChain, error)
	BeginBlockSync(types.RequestBeginBlock) (*types.ResponseBeginBlock, error)
	EndBlockSync(types.RequestEndBlock) (*types.ResponseEndBlock, error)
}

func NewClient(addr, transport string, mustConnect bool) (client Client, err error) {
	switch transport {
	case "socket":
		client = NewSocketClient(addr, mustConnect)
	case "grpc":
		client = NewGRPCClient(addr, mustConnect)
	default:
		err = fmt.Errorf("Unknown abci transport %s", transport)
	}
	return
}

type Callback func(*types.Request, *types.Response)

type ReqRes struct {
	*types.Request
	*sync.WaitGroup
	*types.Response

	mtx	sync.Mutex
	done	bool
	cb	func(*types.Response)
}

func NewReqRes(req *types.Request) *ReqRes {
	return &ReqRes{
		Request:	req,
		WaitGroup:	waitGroup1(),
		Response:	nil,

		done:	false,
		cb:	nil,
	}
}

func (reqRes *ReqRes) SetCallback(cb func(res *types.Response)) {
	reqRes.mtx.Lock()

	if reqRes.done {
		reqRes.mtx.Unlock()
		cb(reqRes.Response)
		return
	}

	defer reqRes.mtx.Unlock()
	reqRes.cb = cb
}

func (reqRes *ReqRes) GetCallback() func(*types.Response) {
	reqRes.mtx.Lock()
	defer reqRes.mtx.Unlock()
	return reqRes.cb
}

func (reqRes *ReqRes) SetDone() {
	reqRes.mtx.Lock()
	reqRes.done = true
	reqRes.mtx.Unlock()
}

func waitGroup1() (wg *sync.WaitGroup) {
	wg = &sync.WaitGroup{}
	wg.Add(1)
	return
}
