package p2p

import (
	"github.com/tendermint/tendermint/p2p/conn"
	cmn "github.com/tendermint/tmlibs/common"
)

type Reactor interface {
	cmn.Service

	SetSwitch(*Switch)

	GetChannels() []*conn.ChannelDescriptor

	AddPeer(peer Peer)

	RemovePeer(peer Peer, reason interface{})

	Receive(chID byte, peer Peer, msgBytes []byte)
}

type BaseReactor struct {
	cmn.BaseService
	Switch	*Switch
}

func NewBaseReactor(name string, impl Reactor) *BaseReactor {
	return &BaseReactor{
		BaseService:	*cmn.NewBaseService(nil, name, impl),
		Switch:		nil,
	}
}

func (br *BaseReactor) SetSwitch(sw *Switch) {
	br.Switch = sw
}
func (*BaseReactor) GetChannels() []*conn.ChannelDescriptor		{ return nil }
func (*BaseReactor) AddPeer(peer Peer)					{}
func (*BaseReactor) RemovePeer(peer Peer, reason interface{})		{}
func (*BaseReactor) Receive(chID byte, peer Peer, msgBytes []byte)	{}
