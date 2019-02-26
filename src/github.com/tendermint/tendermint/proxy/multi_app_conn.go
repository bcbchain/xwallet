package proxy

import (
	"github.com/pkg/errors"

	cmn "github.com/tendermint/tmlibs/common"
)

type AppConns interface {
	cmn.Service

	Mempool() AppConnMempool
	Consensus() AppConnConsensus
	Query() AppConnQuery
}

func NewAppConns(clientCreator ClientCreator, handshaker Handshaker) AppConns {
	return NewMultiAppConn(clientCreator, handshaker)
}

type Handshaker interface {
	Handshake(AppConns) error
}

type multiAppConn struct {
	cmn.BaseService

	handshaker	Handshaker

	mempoolConn	*appConnMempool
	consensusConn	*appConnConsensus
	queryConn	*appConnQuery

	clientCreator	ClientCreator
}

func NewMultiAppConn(clientCreator ClientCreator, handshaker Handshaker) *multiAppConn {
	multiAppConn := &multiAppConn{
		handshaker:	handshaker,
		clientCreator:	clientCreator,
	}
	multiAppConn.BaseService = *cmn.NewBaseService(nil, "multiAppConn", multiAppConn)
	return multiAppConn
}

func (app *multiAppConn) Mempool() AppConnMempool {
	return app.mempoolConn
}

func (app *multiAppConn) Consensus() AppConnConsensus {
	return app.consensusConn
}

func (app *multiAppConn) Query() AppConnQuery {
	return app.queryConn
}

func (app *multiAppConn) OnStart() error {

	querycli, err := app.clientCreator.NewABCIClient()
	if err != nil {
		return errors.Wrap(err, "Error creating ABCI client (query connection)")
	}
	querycli.SetLogger(app.Logger.With("module", "abci-client", "connection", "query"))
	if err := querycli.Start(); err != nil {
		return errors.Wrap(err, "Error starting ABCI client (query connection)")
	}
	app.queryConn = NewAppConnQuery(querycli)

	memcli, err := app.clientCreator.NewABCIClient()
	if err != nil {
		return errors.Wrap(err, "Error creating ABCI client (mempool connection)")
	}
	memcli.SetLogger(app.Logger.With("module", "abci-client", "connection", "mempool"))
	if err := memcli.Start(); err != nil {
		return errors.Wrap(err, "Error starting ABCI client (mempool connection)")
	}
	app.mempoolConn = NewAppConnMempool(memcli)

	concli, err := app.clientCreator.NewABCIClient()
	if err != nil {
		return errors.Wrap(err, "Error creating ABCI client (consensus connection)")
	}
	concli.SetLogger(app.Logger.With("module", "abci-client", "connection", "consensus"))
	if err := concli.Start(); err != nil {
		return errors.Wrap(err, "Error starting ABCI client (consensus connection)")
	}
	app.consensusConn = NewAppConnConsensus(concli)

	if app.handshaker != nil {
		return app.handshaker.Handshake(app)
	}

	return nil
}
