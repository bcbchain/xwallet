package state

import (
	"bytes"
	"fmt"
	"github.com/tendermint/go-crypto"
	cfg "github.com/tendermint/tendermint/config"
	"time"

	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tendermint/types"
)

var (
	stateKey = []byte("stateKey")
)

type State struct {
	ChainID	string

	LastBlockHeight		int64
	LastBlockTotalTx	int64
	LastBlockID		types.BlockID
	LastBlockTime		time.Time
	LastFee			uint64
	LastAllocation		[]abci.Allocation

	Validators			*types.ValidatorSet
	LastValidators			*types.ValidatorSet
	LastHeightValidatorsChanged	int64

	ConsensusParams				types.ConsensusParams
	LastHeightConsensusParamsChanged	int64

	LastResultsHash	[]byte

	LastAppHash	[]byte

	LastTxsHashList	[][]byte
}

func (s State) Copy() State {
	return State{
		ChainID:	s.ChainID,

		LastBlockHeight:	s.LastBlockHeight,
		LastBlockTotalTx:	s.LastBlockTotalTx,
		LastBlockID:		s.LastBlockID,
		LastBlockTime:		s.LastBlockTime,
		LastFee:		s.LastFee,
		LastAllocation:		s.LastAllocation,

		Validators:			s.Validators.Copy(),
		LastValidators:			s.LastValidators.Copy(),
		LastHeightValidatorsChanged:	s.LastHeightValidatorsChanged,

		ConsensusParams:			s.ConsensusParams,
		LastHeightConsensusParamsChanged:	s.LastHeightConsensusParamsChanged,

		LastAppHash:		s.LastAppHash,
		LastTxsHashList:	s.LastTxsHashList,

		LastResultsHash:	s.LastResultsHash,
	}
}

func (s State) Equals(s2 State) bool {
	sbz, s2bz := s.Bytes(), s2.Bytes()
	return bytes.Equal(sbz, s2bz)
}

func (s State) Bytes() []byte {
	return cdc.MustMarshalBinaryBare(s)
}

func (s State) IsEmpty() bool {
	return s.Validators == nil
}

func (s State) GetValidators() (last *types.ValidatorSet, current *types.ValidatorSet) {
	return s.LastValidators, s.Validators
}

func (s State) MakeBlock(height int64, txs []types.Tx, commit *types.Commit, proposer crypto.Address, rewardAddr string, allocation []abci.Allocation) (*types.Block, *types.PartSet) {

	block := types.BCMakeBlock(height, txs, commit, s.LastTxsHashList, proposer, s.LastFee, rewardAddr, allocation)

	block.ChainID = s.ChainID
	block.TotalTxs = s.LastBlockTotalTx + block.NumTxs
	block.LastBlockID = s.LastBlockID
	block.ValidatorsHash = s.Validators.Hash()
	block.LastAppHash = s.LastAppHash

	block.ConsensusHash = s.ConsensusParams.Hash()
	block.LastResultsHash = s.LastResultsHash

	return block, block.MakePartSet(s.ConsensusParams.BlockGossip.BlockPartSizeBytes)
}

func MakeGenesisStateFromFile(config *cfg.Config) (State, error) {
	genDoc, err := MakeGenesisDocFromFile(config)
	if err != nil {
		return State{}, err
	}
	return MakeGenesisState(genDoc)
}

func MakeGenesisDocFromFile(config *cfg.Config) (*types.GenesisDoc, error) {
	return types.GenesisDocFromFile(config)
}

func MakeGenesisState(genDoc *types.GenesisDoc) (State, error) {
	err := genDoc.ValidateAndComplete()
	if err != nil {
		return State{}, fmt.Errorf("Error in genesis file: %v", err)
	}

	validators := make([]*types.Validator, len(genDoc.Validators))
	for i, val := range genDoc.Validators {
		pubKey := val.PubKey
		address := pubKey.Address()

		nodeName := val.Name

		power := val.Power
		if power < 0 {
			power = 0
		}

		validators[i] = &types.Validator{
			Address:	address,
			PubKey:		pubKey,
			VotingPower:	uint64(power),
			RewardAddr:	val.RewardAddr,
			Name:		nodeName,
		}
	}

	return State{

		ChainID:	genDoc.ChainID,

		LastBlockHeight:	0,
		LastBlockID:		types.BlockID{},
		LastBlockTime:		genDoc.GenesisTime,
		LastFee:		0,
		LastAllocation:		[]abci.Allocation{},

		Validators:			types.NewValidatorSet(validators),
		LastValidators:			types.NewValidatorSet(nil),
		LastHeightValidatorsChanged:	1,

		ConsensusParams:			*genDoc.ConsensusParams,
		LastHeightConsensusParamsChanged:	1,

		LastAppHash:	genDoc.AppHash,
	}, nil
}
