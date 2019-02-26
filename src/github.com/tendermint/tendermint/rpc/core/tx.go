package core

import (
	"bytes"
	"encoding/hex"
	"fmt"
	abci "github.com/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/state/txindex/null"
	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
	tmquery "github.com/tendermint/tmlibs/pubsub/query"
)

func Tx(hash string, prove bool) (*ctypes.ResultTx, error) {

	var stateCode uint32
	var height int64
	var index uint32
	check := true

	deTx, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	dResult, err := sm.LoadABCITxResponses(stateDB, cmn.HexBytes(deTx))
	if err == nil {
		height = dResult.Height
	}

	results, err := sm.LoadABCIResponses(stateDB, height)
	if err != nil {
		return nil, err
	}

	for i, item := range results.DeliverTx {
		if bytes.Compare(item.TxHash, deTx) == 0 {
			index = uint32(i)
			break
		}
	}

	var checkResult abci.ResponseCheckTx
	checkRes, errCheck := mempool.GiTxSearch(hash)
	if errCheck == nil {
		if checkRes != nil {
			checkResult = *checkRes
		}
	} else {
		check = false
	}

	if dResult.Height == 0 {
		if checkRes == nil {
			if !check {
				return nil, errCheck
			} else {
				stateCode = 1
				checkResult = abci.ResponseCheckTx{}
			}
		} else if checkRes.Code == 2018 {
			checkResult = abci.ResponseCheckTx{}
			stateCode = 1
		} else {
			stateCode = 2
		}
	} else {
		if checkRes == nil {
			if !check {
				stateCode = 3
				checkResult = abci.ResponseCheckTx{}
			} else {
				stateCode = 5
			}
		} else {
			stateCode = 4
		}

	}

	return &ctypes.ResultTx{
		Hash:		string(hash),
		Height:		height,
		Index:		index,
		DeliverResult:	dResult,
		CheckResult:	checkResult,
		StateCode:	stateCode,
	}, nil
}

func TxSearch(query string, prove bool) ([]*ctypes.ResultTx, error) {

	if _, ok := txIndexer.(*null.TxIndex); ok {
		return nil, fmt.Errorf("Transaction indexing is disabled")
	}

	q, err := tmquery.New(query)
	if err != nil {
		return nil, err
	}

	results, err := txIndexer.Search(q)
	if err != nil {
		return nil, err
	}

	apiResults := make([]*ctypes.ResultTx, len(results))
	var proof types.TxProof
	for i, r := range results {
		height := r.Height
		index := r.Index

		if prove {
			block := blockStore.LoadBlock(height)
			proof = block.Data.Txs.Proof(int(index))
		}

		apiResults[i] = &ctypes.ResultTx{
			Hash:		string(r.Tx.Hash()),
			Height:		height,
			Index:		index,
			DeliverResult:	r.Result,
			Tx:		r.Tx,
			Proof:		proof,
		}
	}

	return apiResults, nil
}
