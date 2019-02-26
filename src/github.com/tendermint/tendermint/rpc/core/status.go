package core

import (
	"time"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

func Status() (*ctypes.ResultStatus, error) {
	latestHeight := blockStore.Height()
	var (
		latestBlockMeta		*types.BlockMeta
		latestBlockHash		cmn.HexBytes
		latestAppHash		cmn.HexBytes
		latestBlockTimeNano	int64
	)
	if latestHeight != 0 {
		latestBlockMeta = blockStore.LoadBlockMeta(latestHeight)
		latestBlockHash = latestBlockMeta.BlockID.Hash
		latestAppHash = latestBlockMeta.Header.LastAppHash
		latestBlockTimeNano = latestBlockMeta.Header.Time.UnixNano()
	}

	latestBlockTime := time.Unix(0, latestBlockTimeNano)

	var votingPower uint64
	var rewardAddr string
	var nodeName string
	if val := validatorAtHeight(latestHeight); val != nil {
		votingPower = val.VotingPower
		rewardAddr = val.RewardAddr
		nodeName = val.Name
	}

	result := &ctypes.ResultStatus{
		NodeInfo:	p2pSwitch.NodeInfo(),
		SyncInfo: ctypes.SyncInfo{
			LatestBlockHash:	latestBlockHash,
			LatestAppHash:		latestAppHash,
			LatestBlockHeight:	latestHeight,
			LatestBlockTime:	latestBlockTime,
			Syncing:		consensusReactor.FastSync(),
		},
		ValidatorInfo: ctypes.ValidatorInfo{
			Address:	pubKey.Address(),
			PubKey:		pubKey,
			VotingPower:	votingPower,
			RewardAddr:	rewardAddr,
			Name:		nodeName,
		},
	}

	return result, nil
}

func validatorAtHeight(h int64) *types.Validator {
	lastBlockHeight, vals := consensusState.GetValidators()

	privValAddress := pubKey.Address()

	if lastBlockHeight == h {
		for _, val := range vals {
			if val.Address == privValAddress {
				return val
			}
		}
	}

	if lastBlockHeight > h {
		vals, err := sm.LoadValidators(stateDB, h)
		if err != nil {

			return nil
		}
		_, val := vals.GetByAddress(privValAddress)
		return val
	}

	return nil
}
