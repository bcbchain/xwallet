package core

import (
	cm "github.com/tendermint/tendermint/consensus"
	"github.com/tendermint/tendermint/p2p"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/types"
)

func Validators(heightPtr *int64) (*ctypes.ResultValidators, error) {
	storeHeight := blockStore.Height()
	height, err := getHeight(storeHeight, heightPtr)
	if err != nil {
		return nil, err
	}

	validators, err := sm.LoadValidators(stateDB, height)
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultValidators{height, validators.Validators}, nil
}

func DumpConsensusState() (*ctypes.ResultDumpConsensusState, error) {
	peers := p2pSwitch.Peers().List()
	peerRoundStates := make([]ctypes.PeerRoundState, 0)
	for _, peer := range peers {
		if isPrivate(peer) {
			continue
		}
		peerState := peer.Get(types.PeerStateKey).(*cm.PeerState)
		peerRoundState, err := peerState.GetRoundStateJSON()
		if err != nil {
			return nil, err
		}
		prs := ctypes.PeerRoundState{
			NodeAddress:	p2p.IDAddressString(peer.ID(), peer.NodeInfo().ListenAddr),
			PeerRoundState:	peerRoundState,
		}
		peerRoundStates = append(peerRoundStates, prs)
	}
	roundState, err := consensusState.GetRoundStateJSON()
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultDumpConsensusState{RoundState: roundState, PeerRoundStates: peerRoundStates}, nil
}
