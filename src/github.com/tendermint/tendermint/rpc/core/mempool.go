package core

import (
	"context"
	"fmt"
	"bcbchain.io/algorithm"
	"github.com/pkg/errors"
	abci "github.com/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
	"time"
)

func BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	txHash := cmn.HexBytes(algorithm.CalcCodeHash(string(tx)))

	mempool.GiTxCache(txHash, nil)

	go func() {
		err := mempool.CheckTx(tx, func(res *abci.Response) {
			mempool.GiTxCache(txHash, res.GetCheckTx())
		})
		if err != nil {
			fmt.Errorf("Error broadcasting transaction: %v", err)
		}
	}()

	return &ctypes.ResultBroadcastTx{Hash: txHash}, nil
}

func BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	txHash := cmn.HexBytes(algorithm.CalcCodeHash(string(tx)))

	resCh := make(chan *abci.Response, 1)
	err := mempool.CheckTx(tx, func(res *abci.Response) {
		resCh <- res
	})
	if err != nil {
		return nil, fmt.Errorf("Error broadcasting transaction: %v", err)
	}
	res := <-resCh
	r := res.GetCheckTx()

	return &ctypes.ResultBroadcastTx{
		Code:	r.Code,
		Data:	cmn.HexBytes(r.Data),
		Log:	r.Log,
		Hash:	txHash,
	}, nil
}

func BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	txHash := cmn.HexBytes(algorithm.CalcCodeHash(string(tx)))

	ctx, cancel := context.WithTimeout(context.Background(), subscribeTimeout)
	defer cancel()
	deliverTxResCh := make(chan interface{})
	q := types.EventQueryTxFor(tx)
	err := eventBus.Subscribe(ctx, "mempool", q, deliverTxResCh)
	if err != nil {
		err = errors.Wrap(err, "failed to subscribe to tx")
		logger.Error("Error on broadcastTxCommit", "err", err)
		return nil, fmt.Errorf("Error on broadcastTxCommit: %v", err)
	}
	defer eventBus.Unsubscribe(context.Background(), "mempool", q)

	checkTxResCh := make(chan *abci.Response, 1)
	err = mempool.CheckTx(tx, func(res *abci.Response) {
		checkTxResCh <- res
	})
	if err != nil {
		logger.Error("Error on broadcastTxCommit", "err", err)
		return nil, fmt.Errorf("Error on broadcastTxCommit: %v", err)
	}
	checkTxRes := <-checkTxResCh
	checkTxR := checkTxRes.GetCheckTx()

	if checkTxR.Code != abci.CodeTypeOK {

		return &ctypes.ResultBroadcastTxCommit{
			CheckTx:	*checkTxR,
			DeliverTx:	abci.ResponseDeliverTx{},
			Hash:		txHash,
		}, nil
	}

	timer := time.NewTimer(60 * 2 * time.Second)
	select {
	case deliverTxResMsg := <-deliverTxResCh:
		deliverTxRes := deliverTxResMsg.(types.EventDataTx)

		deliverTxR := deliverTxRes.Result
		logger.Trace("DeliverTx passed ", "tx", cmn.HexBytes(tx), "response", deliverTxR)
		logger.Trace("DeliverTx byte passed ", "tx", []byte(tx), "txHash", tx.Hash())
		return &ctypes.ResultBroadcastTxCommit{
			CheckTx:	*checkTxR,
			DeliverTx:	deliverTxR,
			Hash:		txHash,
			Height:		deliverTxRes.Height,
		}, nil
	case <-timer.C:
		logger.Error("failed to include tx")
		return &ctypes.ResultBroadcastTxCommit{
			CheckTx:	*checkTxR,
			DeliverTx:	abci.ResponseDeliverTx{},
			Hash:		txHash,
		}, fmt.Errorf("Timed out waiting for transaction to be included in a block")
	}
}

func UnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	txs := mempool.Reap(-1)
	return &ctypes.ResultUnconfirmedTxs{len(txs), txs}, nil
}

func NumUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	return &ctypes.ResultUnconfirmedTxs{N: mempool.Size()}, nil
}
