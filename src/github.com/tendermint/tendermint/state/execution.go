package state

import (
	"encoding/binary"
	"fmt"
	"github.com/tendermint/tendermint/proxy"
	"os"
	"github.com/ebuchman/fail-test"
	"github.com/pkg/errors"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

var SyncTo int64 = 0

type BlockExecutor struct {
	db	dbm.DB

	proxyApp	proxy.AppConnConsensus

	eventBus	types.BlockEventPublisher

	mempool	types.Mempool
	evpool	types.EvidencePool

	logger	log.Logger
}

func NewBlockExecutor(db dbm.DB, logger log.Logger, proxyApp proxy.AppConnConsensus,
	mempool types.Mempool, evpool types.EvidencePool) *BlockExecutor {
	return &BlockExecutor{
		db:		db,
		proxyApp:	proxyApp,
		eventBus:	types.NopEventBus{},
		mempool:	mempool,
		evpool:		evpool,
		logger:		logger,
	}
}

func (blockExec *BlockExecutor) SetEventBus(eventBus types.BlockEventPublisher) {
	blockExec.eventBus = eventBus
}

func (blockExec *BlockExecutor) ValidateBlock(s State, block *types.Block) error {
	return validateBlock(blockExec.db, s, block)
}

func (blockExec *BlockExecutor) ApplyBlock(s State, blockID types.BlockID, block *types.Block) (State, error) {
	if SyncTo > 0 && block.Height >= SyncTo {
		fmt.Println("block height ===", block.Height, ", exiting...")
		os.Exit(0)
	}
	if err := blockExec.ValidateBlock(s, block); err != nil {
		return s, ErrInvalidBlock(err)
	}

	abciResponses, err := execBlockOnProxyApp(blockExec.logger, blockExec.proxyApp, block)
	if err != nil {
		return s, ErrProxyAppConn(err)
	}

	fail.Fail()

	saveABCIResponses(blockExec.db, block.Height, abciResponses)

	fail.Fail()

	s, err = updateState(s, blockID, block.Header, abciResponses)
	if err != nil {
		return s, fmt.Errorf("Commit failed for application: %v", err)
	}

	res, err := blockExec.Commit(block)
	if err != nil {
		return s, fmt.Errorf("Commit failed for application: %v", err)
	}

	fail.Fail()

	appState := abci.ByteToAppState(res.AppState)
	s.LastAppHash = appState.AppHash
	s.LastTxsHashList = make([][]byte, 0)
	for _, txhash := range appState.TxsHashList {
		s.LastTxsHashList = append(s.LastTxsHashList, txhash)
	}

	s.LastFee = appState.Fee
	s.LastAllocation = []abci.Allocation{}
	for _, r := range appState.Rewards {
		k := string(r.Key)
		v := binary.BigEndian.Uint64(r.Value)
		s.LastAllocation = append(s.LastAllocation, abci.Allocation{Addr: k, Fee: v})
	}

	blockExec.logger.Info("block apply", "height", block.Height, "got hash", s.LastAppHash, "s.LastTxsHashList", s.LastTxsHashList)

	SaveState(blockExec.db, s)

	fail.Fail()

	blockExec.evpool.Update(block)

	fireEvents(blockExec.logger, blockExec.eventBus, block, abciResponses)

	return s, nil
}

func (blockExec *BlockExecutor) Commit(block *types.Block) (*abci.ResponseCommit, error) {
	blockExec.mempool.Lock()
	defer blockExec.mempool.Unlock()

	err := blockExec.mempool.FlushAppConn()
	if err != nil {
		blockExec.logger.Error("Client error during mempool.FlushAppConn", "err", err)
		return nil, err
	}

	res, err := blockExec.proxyApp.CommitSync()
	if err != nil {
		blockExec.logger.Error("Client error during proxyAppConn.CommitSync", "err", err)
		return nil, err
	}

	blockExec.logger.Info("Committed state",
		"height", block.Height,
		"txs", block.NumTxs,
		"responseCommit", fmt.Sprintf("%v", res))

	if err := blockExec.mempool.Update(block.Height, block.Txs); err != nil {
		return nil, err
	}

	return res, nil
}

func execBlockOnProxyApp(logger log.Logger, proxyAppConn proxy.AppConnConsensus, block *types.Block) (*ABCIResponses, error) {
	var validTxs, invalidTxs = 0, 0

	txIndex := 0
	abciResponses := NewABCIResponses(block)

	proxyCb := func(req *abci.Request, res *abci.Response) {
		switch r := res.Value.(type) {
		case *abci.Response_DeliverTx:

			txRes := r.DeliverTx
			if txRes.Code == abci.CodeTypeOK {
				validTxs++
			} else {
				logger.Debug("Invalid tx", "code", txRes.Code, "log", txRes.Log)
				invalidTxs++
			}
			abciResponses.DeliverTx[txIndex] = txRes
			txIndex++
		}
	}
	proxyAppConn.SetResponseCallback(proxyCb)

	absentVals := make([]int32, 0)
	for valI, vote := range block.LastCommit.Precommits {
		if vote == nil {
			absentVals = append(absentVals, int32(valI))
		}
	}

	byzantineVals := make([]abci.Evidence, len(block.Evidence.Evidence))
	for i, ev := range block.Evidence.Evidence {
		byzantineVals[i] = abci.Evidence{
			PubKey:	ev.Address(),
			Height:	ev.Height(),
		}
	}

	res, err := proxyAppConn.BeginBlockSync(abci.RequestBeginBlock{
		Hash:			block.Hash(),
		Header:			types.TM2PB.Header(block.Header),
		AbsentValidators:	absentVals,
		ByzantineValidators:	byzantineVals,
	})
	if err != nil {
		logger.Error("Error in proxyAppConn.BeginBlock", "err", err)
		return nil, err
	}

	if res.Code != abci.CodeTypeOK {

		logger.Error("abci server  excute beginBlock failed", "log", res.Log)
		return nil, errors.New(res.Log)
	}

	for _, tx := range block.Txs {
		proxyAppConn.DeliverTxAsync(tx)
		if err := proxyAppConn.Error(); err != nil {
			return nil, err
		}
	}

	abciResponses.EndBlock, err = proxyAppConn.EndBlockSync(abci.RequestEndBlock{block.Height})
	if err != nil {
		logger.Error("Error in proxyAppConn.EndBlock", "err", err)
		return nil, err
	}

	logger.Info("Executed block", "height", block.Height, "validTxs", validTxs, "invalidTxs", invalidTxs)

	valUpdates := abciResponses.EndBlock.ValidatorUpdates
	if len(valUpdates) > 0 {
		logger.Info("Updates to validators", "updates", abci.ValidatorsString(valUpdates))
	}

	return abciResponses, nil
}

func updateValidators(currentSet *types.ValidatorSet, updates []abci.Validator) error {
	for _, v := range updates {
		pubkey, err := crypto.PubKeyFromBytes(v.PubKey)
		if err != nil {
			return err
		}

		address := pubkey.Address()
		power := v.Power

		if power < 0 {
			return fmt.Errorf("Power (%d) overflows int64", v.Power)
		}

		_, val := currentSet.GetByAddress(address)
		if val == nil {

			added := currentSet.Add(types.NewValidator(pubkey, power, v.RewardAddr, v.Name))
			if !added {
				return fmt.Errorf("Failed to add new validator %X with voting power %d", address, power)
			}
		} else if v.Power == 0 {

			_, removed := currentSet.Remove(address)
			if !removed {
				return fmt.Errorf("Failed to remove validator %X", address)
			}
		} else {

			val.VotingPower = power
			updated := currentSet.Update(types.NewValidator(pubkey, power, v.RewardAddr, v.Name))
			if !updated {
				return fmt.Errorf("Failed to update validator %X with voting power %d", address, power)
			}
		}
	}
	return nil
}

func updateState(s State, blockID types.BlockID, header *types.Header,
	abciResponses *ABCIResponses) (State, error) {

	prevValSet := s.Validators.Copy()
	nextValSet := prevValSet.Copy()

	lastHeightValsChanged := s.LastHeightValidatorsChanged
	if len(abciResponses.EndBlock.ValidatorUpdates) > 0 {
		err := updateValidators(nextValSet, abciResponses.EndBlock.ValidatorUpdates)
		if err != nil {
			return s, fmt.Errorf("Error changing validator set: %v", err)
		}

		lastHeightValsChanged = header.Height + 1
	}

	nextValSet.IncrementAccum(1)

	nextParams := s.ConsensusParams
	lastHeightParamsChanged := s.LastHeightConsensusParamsChanged
	if abciResponses.EndBlock.ConsensusParamUpdates != nil {

		nextParams = s.ConsensusParams.Update(abciResponses.EndBlock.ConsensusParamUpdates)
		err := nextParams.Validate()
		if err != nil {
			return s, fmt.Errorf("Error updating consensus params: %v", err)
		}

		lastHeightParamsChanged = header.Height + 1
	}

	return State{
		ChainID:				s.ChainID,
		LastBlockHeight:			header.Height,
		LastBlockTotalTx:			s.LastBlockTotalTx + header.NumTxs,
		LastBlockID:				blockID,
		LastBlockTime:				header.Time,
		Validators:				nextValSet,
		LastValidators:				s.Validators.Copy(),
		LastHeightValidatorsChanged:		lastHeightValsChanged,
		ConsensusParams:			nextParams,
		LastHeightConsensusParamsChanged:	lastHeightParamsChanged,
		LastResultsHash:			abciResponses.ResultsHash(),
		LastAppHash:				nil,
	}, nil
}

func fireEvents(logger log.Logger, eventBus types.BlockEventPublisher, block *types.Block, abciResponses *ABCIResponses) {

	txEventBuffer := types.NewTxEventBuffer(eventBus, int(block.NumTxs))
	for i, tx := range block.Data.Txs {
		txEventBuffer.PublishEventTx(types.EventDataTx{types.TxResult{
			Height:	block.Height,
			Index:	uint32(i),
			Tx:	tx,
			Result:	*(abciResponses.DeliverTx[i]),
		}})
	}

	eventBus.PublishEventNewBlock(types.EventDataNewBlock{block})
	eventBus.PublishEventNewBlockHeader(types.EventDataNewBlockHeader{block.Header})
	err := txEventBuffer.Flush()
	if err != nil {
		logger.Error("Failed to flush event buffer", "err", err)
	}
}

func ExecCommitBlock(appConnConsensus proxy.AppConnConsensus, block *types.Block, logger log.Logger) (*abci.ResponseCommit, error) {
	_, err := execBlockOnProxyApp(logger, appConnConsensus, block)
	if err != nil {
		logger.Error("Error executing block on proxy app", "height", block.Height, "err", err)
		return nil, err
	}

	res, err := appConnConsensus.CommitSync()
	if err != nil {
		logger.Error("Client error during proxyAppConn.CommitSync", "err", res)
		return nil, err
	}

	return res, nil
}
