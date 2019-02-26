package state

import (
	"fmt"
	abci "github.com/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
)

func calcValidatorsKey(height int64) []byte {
	return []byte(cmn.Fmt("validatorsKey:%v", height))
}

func calcConsensusParamsKey(height int64) []byte {
	return []byte(cmn.Fmt("consensusParamsKey:%v", height))
}

func calcABCIResponsesKey(height int64) []byte {
	return []byte(cmn.Fmt("abciResponsesKey:%v", height))
}

func LoadStateFromDBOrGenesisFile(stateDB dbm.DB, config *cfg.Config) (State, error) {
	state := LoadState(stateDB)
	if state.IsEmpty() {
		var err error
		state, err = MakeGenesisStateFromFile(config)
		if err != nil {
			return state, err
		}
		SaveState(stateDB, state)
	}

	return state, nil
}

func LoadStateFromDBOrGenesisDoc(stateDB dbm.DB, genesisDoc *types.GenesisDoc) (State, error) {
	state := LoadState(stateDB)
	if state.IsEmpty() {
		var err error
		state, err = MakeGenesisState(genesisDoc)
		if err != nil {
			return state, err
		}
		SaveState(stateDB, state)
	}

	return state, nil
}

func LoadState(db dbm.DB) State {
	return loadState(db, stateKey)
}

func loadState(db dbm.DB, key []byte) (state State) {
	buf := db.Get(key)
	if len(buf) == 0 {
		return state
	}

	err := cdc.UnmarshalBinaryBare(buf, &state)
	if err != nil {

		cmn.Exit(cmn.Fmt(`LoadState: Data has been corrupted or its spec has changed:
                %v\n`, err))
	}

	return state
}

func SaveState(db dbm.DB, s State) {
	saveState(db, s, stateKey)
}

func saveState(db dbm.DB, s State, key []byte) {
	nextHeight := s.LastBlockHeight + 1
	saveValidatorsInfo(db, nextHeight, s.LastHeightValidatorsChanged, s.Validators)
	saveConsensusParamsInfo(db, nextHeight, s.LastHeightConsensusParamsChanged, s.ConsensusParams)
	db.SetSync(stateKey, s.Bytes())
}

type ABCIResponses struct {
	DeliverTx	[]*abci.ResponseDeliverTx
	EndBlock	*abci.ResponseEndBlock
}

func NewABCIResponses(block *types.Block) *ABCIResponses {
	resDeliverTxs := make([]*abci.ResponseDeliverTx, block.NumTxs)
	if block.NumTxs == 0 {

		resDeliverTxs = nil
	}
	return &ABCIResponses{
		DeliverTx: resDeliverTxs,
	}
}

func (arz *ABCIResponses) Bytes() []byte {
	return cdc.MustMarshalBinaryBare(arz)
}

func (arz *ABCIResponses) ResultsHash() []byte {
	results := types.NewResults(arz.DeliverTx)
	return results.Hash()
}

func LoadABCITxResponses(db dbm.DB, tx cmn.HexBytes) (abci.ResponseDeliverTx, error) {
	var res = abci.ResponseDeliverTx{}
	buf := db.Get(tx)
	if len(buf) == 0 {
		return res, ErrNoABCITxResponseForTxHash{tx}
	}
	err := cdc.UnmarshalBinaryBare(buf, &res)
	if err != nil {

		cmn.Exit(cmn.Fmt(`LoadABCIResponses: Data has been corrupted or its spec has
                changed: %v\n`, err))
	}

	return res, nil
}

func LoadABCIResponses(db dbm.DB, height int64) (*ABCIResponses, error) {
	buf := db.Get(calcABCIResponsesKey(height))
	if len(buf) == 0 {
		return nil, ErrNoABCIResponsesForHeight{height}
	}

	abciResponses := new(ABCIResponses)
	err := cdc.UnmarshalBinaryBare(buf, abciResponses)
	if err != nil {

		cmn.Exit(cmn.Fmt(`LoadABCIResponses: Data has been corrupted or its spec has
                changed: %v\n`, err))
	}

	return abciResponses, nil
}

func saveABCIResponses(db dbm.DB, height int64, abciResponses *ABCIResponses) {
	db.SetSync(calcABCIResponsesKey(height), abciResponses.Bytes())
	for _, e := range abciResponses.DeliverTx {
		e.Height = height
		db.SetSync(e.TxHash, cdc.MustMarshalBinaryBare(e))
	}
}

type ValidatorsInfo struct {
	ValidatorSet		*types.ValidatorSet
	LastHeightChanged	int64
}

func (valInfo *ValidatorsInfo) Bytes() []byte {
	return cdc.MustMarshalBinaryBare(valInfo)
}

func LoadValidators(db dbm.DB, height int64) (*types.ValidatorSet, error) {
	valInfo := loadValidatorsInfo(db, height)
	if valInfo == nil {
		return nil, ErrNoValSetForHeight{height}
	}

	if valInfo.ValidatorSet == nil {
		valInfo = loadValidatorsInfo(db, valInfo.LastHeightChanged)
		if valInfo == nil {
			cmn.PanicSanity(fmt.Sprintf(`Couldn't find validators at height %d as
                        last changed from height %d`, valInfo.LastHeightChanged, height))
		}
	}

	return valInfo.ValidatorSet, nil
}

func loadValidatorsInfo(db dbm.DB, height int64) *ValidatorsInfo {
	buf := db.Get(calcValidatorsKey(height))
	if len(buf) == 0 {
		return nil
	}

	v := new(ValidatorsInfo)
	err := cdc.UnmarshalBinaryBare(buf, v)
	if err != nil {

		cmn.Exit(cmn.Fmt(`LoadValidators: Data has been corrupted or its spec has changed:
                %v\n`, err))
	}

	return v
}

func saveValidatorsInfo(db dbm.DB, nextHeight, changeHeight int64, valSet *types.ValidatorSet) {
	valInfo := &ValidatorsInfo{
		LastHeightChanged: changeHeight,
	}
	if changeHeight == nextHeight {
		valInfo.ValidatorSet = valSet
	}
	db.SetSync(calcValidatorsKey(nextHeight), valInfo.Bytes())
}

type ConsensusParamsInfo struct {
	ConsensusParams		types.ConsensusParams
	LastHeightChanged	int64
}

func (params ConsensusParamsInfo) Bytes() []byte {
	return cdc.MustMarshalBinaryBare(params)
}

func LoadConsensusParams(db dbm.DB, height int64) (types.ConsensusParams, error) {
	empty := types.ConsensusParams{}

	paramsInfo := loadConsensusParamsInfo(db, height)
	if paramsInfo == nil {
		return empty, ErrNoConsensusParamsForHeight{height}
	}

	if paramsInfo.ConsensusParams == empty {
		paramsInfo = loadConsensusParamsInfo(db, paramsInfo.LastHeightChanged)
		if paramsInfo == nil {
			cmn.PanicSanity(fmt.Sprintf(`Couldn't find consensus params at height %d as
                        last changed from height %d`, paramsInfo.LastHeightChanged, height))
		}
	}

	return paramsInfo.ConsensusParams, nil
}

func loadConsensusParamsInfo(db dbm.DB, height int64) *ConsensusParamsInfo {
	buf := db.Get(calcConsensusParamsKey(height))
	if len(buf) == 0 {
		return nil
	}

	paramsInfo := new(ConsensusParamsInfo)
	err := cdc.UnmarshalBinaryBare(buf, paramsInfo)
	if err != nil {

		cmn.Exit(cmn.Fmt(`LoadConsensusParams: Data has been corrupted or its spec has changed:
                %v\n`, err))
	}

	return paramsInfo
}

func saveConsensusParamsInfo(db dbm.DB, nextHeight, changeHeight int64, params types.ConsensusParams) {
	paramsInfo := &ConsensusParamsInfo{
		LastHeightChanged: changeHeight,
	}
	if changeHeight == nextHeight {
		paramsInfo.ConsensusParams = params
	}
	db.SetSync(calcConsensusParamsKey(nextHeight), paramsInfo.Bytes())
}
