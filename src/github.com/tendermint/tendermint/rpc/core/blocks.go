package core

import (
	"fmt"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

func BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	if minHeight == 0 {
		minHeight = 1
	}

	if maxHeight == 0 {
		maxHeight = blockStore.Height()
	} else {
		maxHeight = cmn.MinInt64(blockStore.Height(), maxHeight)
	}

	const limit int64 = 20
	minHeight = cmn.MaxInt64(minHeight, maxHeight-limit)

	logger.Debug("BlockchainInfoHandler", "maxHeight", maxHeight, "minHeight", minHeight)

	if minHeight > maxHeight {
		return nil, fmt.Errorf("min height %d can't be greater than max height %d", minHeight, maxHeight)
	}

	blockMetas := []*types.BlockMeta{}
	for height := maxHeight; height >= minHeight; height-- {
		blockMeta := blockStore.LoadBlockMeta(height)
		blockMetas = append(blockMetas, blockMeta)
	}

	return &ctypes.ResultBlockchainInfo{blockStore.Height(), blockMetas}, nil
}

func Block(heightPtr *int64) (*ctypes.ResultBlock, error) {
	storeHeight := blockStore.Height()
	height, err := getHeight(storeHeight, heightPtr)
	if err != nil {
		return nil, err
	}

	blockMeta := blockStore.LoadBlockMeta(height)
	block := blockStore.LoadBlock(height)
	blockSize := calculateSize(block)
	return &ctypes.ResultBlock{BlockMeta: blockMeta, Block: block, BlockSize: blockSize}, nil
}

func calculateSize(block *types.Block) int {
	byteResult, err := cdc.MarshalBinary(block)
	if err != nil {
		return 0
	}
	return len(byteResult)
}

func Commit(heightPtr *int64) (*ctypes.ResultCommit, error) {
	storeHeight := blockStore.Height()
	height, err := getHeight(storeHeight, heightPtr)
	if err != nil {
		return nil, err
	}

	header := blockStore.LoadBlockMeta(height).Header

	if height == storeHeight {
		commit := blockStore.LoadSeenCommit(height)
		return ctypes.NewResultCommit(header, commit, false), nil
	}

	commit := blockStore.LoadBlockCommit(height)
	return ctypes.NewResultCommit(header, commit, true), nil
}

func BlockResults(heightPtr *int64) (*ctypes.ResultBlockResults, error) {
	storeHeight := blockStore.Height()
	height, err := getHeight(storeHeight, heightPtr)
	if err != nil {
		return nil, err
	}

	results, err := sm.LoadABCIResponses(stateDB, height)
	if err != nil {
		return nil, err
	}

	res := &ctypes.ResultBlockResults{
		Height:		height,
		Results:	results,
	}
	return res, nil
}

func getHeight(storeHeight int64, heightPtr *int64) (int64, error) {
	if heightPtr != nil {
		height := *heightPtr
		if height <= 0 {
			return 0, fmt.Errorf("Height must be greater than 0")
		}
		if height > storeHeight {
			return 0, fmt.Errorf("Height must be less than or equal to the current blockchain height")
		}
		return height, nil
	}
	return storeHeight, nil
}
