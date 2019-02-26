package dummy

import (
	p2p "github.com/tendermint/tendermint/p2p"
	tmconn "github.com/tendermint/tendermint/p2p/conn"
	cmn "github.com/tendermint/tmlibs/common"
)

type peer struct {
	cmn.BaseService
	kv	map[string]interface{}
}

var _ p2p.Peer = (*peer)(nil)

func NewPeer() *peer {
	p := &peer{
		kv: make(map[string]interface{}),
	}
	p.BaseService = *cmn.NewBaseService(nil, "peer", p)
	return p
}

func (p *peer) ID() p2p.ID {
	return p2p.ID("dummy")
}

func (p *peer) IsOutbound() bool {
	return false
}

func (p *peer) IsPersistent() bool {
	return false
}

func (p *peer) NodeInfo() p2p.NodeInfo {
	return p2p.NodeInfo{}
}

func (p *peer) Status() tmconn.ConnectionStatus {
	return tmconn.ConnectionStatus{}
}

func (p *peer) Send(byte, []byte) bool {
	return true
}

func (p *peer) TrySend(byte, []byte) bool {
	return true
}

func (p *peer) Set(key string, value interface{}) {
	p.kv[key] = value
}

func (p *peer) Get(key string) interface{} {
	if value, ok := p.kv[key]; ok {
		return value
	}
	return nil
}
