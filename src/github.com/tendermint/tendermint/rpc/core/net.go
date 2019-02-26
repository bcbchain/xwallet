package core

import (
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/p2p"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func NetInfo() (*ctypes.ResultNetInfo, error) {
	listening := p2pSwitch.IsListening()
	var listeners []string
	for _, listener := range p2pSwitch.Listeners() {
		listeners = append(listeners, listener.String())
	}
	var peers []ctypes.Peer
	for _, peer := range p2pSwitch.Peers().List() {
		if isPrivate(peer) {
			continue
		}
		peers = append(peers, ctypes.Peer{
			NodeInfo:		peer.NodeInfo(),
			IsOutbound:		peer.IsOutbound(),
			ConnectionStatus:	peer.Status(),
		})
	}
	return &ctypes.ResultNetInfo{
		Listening:	listening,
		Listeners:	listeners,
		Peers:		peers,
	}, nil
}

func isPrivate(peer p2p.Peer) bool {
	for _, id := range privatePeerIDs {
		if string(peer.NodeInfo().NetAddress().ID) == id {
			return true
		}
	}
	return false
}

func UnsafeDialSeeds(seeds []string) (*ctypes.ResultDialSeeds, error) {
	if len(seeds) == 0 {
		return &ctypes.ResultDialSeeds{}, errors.New("No seeds provided")
	}

	logger.Info("DialSeeds", "addrBook", addrBook, "seeds", seeds)
	err := p2pSwitch.DialPeersAsync(addrBook, seeds, false)
	if err != nil {
		return &ctypes.ResultDialSeeds{}, err
	}
	return &ctypes.ResultDialSeeds{Log: "Dialing seeds in progress. See /net_info for details"}, nil
}

func UnsafeDialPeers(peers []string, persistent bool) (*ctypes.ResultDialPeers, error) {
	if len(peers) == 0 {
		return &ctypes.ResultDialPeers{}, errors.New("No peers provided")
	}

	logger.Info("DialPeers", "addrBook", addrBook, "peers", peers, "persistent", persistent)
	err := p2pSwitch.DialPeersAsync(addrBook, peers, persistent)
	if err != nil {
		return &ctypes.ResultDialPeers{}, err
	}
	return &ctypes.ResultDialPeers{Log: "Dialing peers in progress. See /net_info for details"}, nil
}

func Genesis() (*ctypes.ResultGenesis, error) {
	return &ctypes.ResultGenesis{Genesis: genDoc}, nil
}
