package state

import (
	"testing"

	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

func TestValidateBlock(t *testing.T) {
	state := state()

	blockExec := NewBlockExecutor(dbm.NewMemDB(), log.TestingLogger(), nil, nil, nil)

	block := makeBlock(state, 1)
	err := blockExec.ValidateBlock(state, block)
	require.NoError(t, err)

	block = makeBlock(state, 1)
	block.ChainID = "not-the-real-one"
	err = blockExec.ValidateBlock(state, block)
	require.Error(t, err)

	block = makeBlock(state, 1)
	block.Height += 10
	err = blockExec.ValidateBlock(state, block)
	require.Error(t, err)

	block = makeBlock(state, 1)
	block.TotalTxs += 10
	err = blockExec.ValidateBlock(state, block)
	require.Error(t, err)

	block = makeBlock(state, 1)
	block.LastBlockID.PartsHeader.Total += 10
	err = blockExec.ValidateBlock(state, block)
	require.Error(t, err)

	block = makeBlock(state, 1)
	block.AppHash = []byte("wrong app hash")
	err = blockExec.ValidateBlock(state, block)
	require.Error(t, err)

	block = makeBlock(state, 1)
	block.ConsensusHash = []byte("wrong consensus hash")
	err = blockExec.ValidateBlock(state, block)
	require.Error(t, err)

	block = makeBlock(state, 1)
	block.LastResultsHash = []byte("wrong results hash")
	err = blockExec.ValidateBlock(state, block)
	require.Error(t, err)

	block = makeBlock(state, 1)
	block.ValidatorsHash = []byte("wrong validators hash")
	err = blockExec.ValidateBlock(state, block)
	require.Error(t, err)
}
