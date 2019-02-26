package p2p

import (
	"sync"
)

type IPeerSet interface {
	Has(key ID) bool
	Get(key ID) Peer
	List() []Peer
	Size() int
}

type PeerSet struct {
	mtx	sync.Mutex
	lookup	map[ID]*peerSetItem
	list	[]Peer
}

type peerSetItem struct {
	peer	Peer
	index	int
}

func NewPeerSet() *PeerSet {
	return &PeerSet{
		lookup:	make(map[ID]*peerSetItem),
		list:	make([]Peer, 0, 256),
	}
}

func (ps *PeerSet) Add(peer Peer) error {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	if ps.lookup[peer.ID()] != nil {
		return ErrSwitchDuplicatePeer
	}

	index := len(ps.list)

	ps.list = append(ps.list, peer)
	ps.lookup[peer.ID()] = &peerSetItem{peer, index}
	return nil
}

func (ps *PeerSet) Has(peerKey ID) bool {
	ps.mtx.Lock()
	_, ok := ps.lookup[peerKey]
	ps.mtx.Unlock()
	return ok
}

func (ps *PeerSet) Get(peerKey ID) Peer {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	item, ok := ps.lookup[peerKey]
	if ok {
		return item.peer
	}
	return nil
}

func (ps *PeerSet) Remove(peer Peer) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	item := ps.lookup[peer.ID()]
	if item == nil {
		return
	}

	index := item.index

	newList := make([]Peer, len(ps.list)-1)
	copy(newList, ps.list)

	if index == len(ps.list)-1 {
		ps.list = newList
		delete(ps.lookup, peer.ID())
		return
	}

	lastPeer := ps.list[len(ps.list)-1]
	lastPeerKey := lastPeer.ID()
	lastPeerItem := ps.lookup[lastPeerKey]
	newList[index] = lastPeer
	lastPeerItem.index = index
	ps.list = newList
	delete(ps.lookup, peer.ID())
}

func (ps *PeerSet) Size() int {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	return len(ps.list)
}

func (ps *PeerSet) List() []Peer {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	return ps.list
}
