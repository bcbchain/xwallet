package p2p

import (
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/p2p/conn"
	cmn "github.com/tendermint/tmlibs/common"
)

const (
	dialRandomizerIntervalMilliseconds	= 3000

	reconnectAttempts	= 20
	reconnectInterval	= 5 * time.Second

	reconnectBackOffAttempts	= 10
	reconnectBackOffBaseSeconds	= 3
)

type AddrBook interface {
	AddAddress(addr *NetAddress, src *NetAddress) error
	AddOurAddress(*NetAddress)
	OurAddress(*NetAddress) bool
	MarkGood(*NetAddress)
	RemoveAddress(*NetAddress)
	HasAddress(*NetAddress) bool
	Save()
}

type Switch struct {
	cmn.BaseService

	config		*cfg.P2PConfig
	peerConfig	*PeerConfig
	listeners	[]Listener
	reactors	map[string]Reactor
	chDescs		[]*conn.ChannelDescriptor
	reactorsByCh	map[byte]Reactor
	peers		*PeerSet
	dialing		*cmn.CMap
	nodeInfo	NodeInfo
	nodeKey		*NodeKey
	addrBook	AddrBook

	filterConnByAddr	func(net.Addr) error
	filterConnByID		func(ID) error

	rng	*cmn.Rand
}

func NewSwitch(config *cfg.P2PConfig) *Switch {
	sw := &Switch{
		config:		config,
		peerConfig:	DefaultPeerConfig(),
		reactors:	make(map[string]Reactor),
		chDescs:	make([]*conn.ChannelDescriptor, 0),
		reactorsByCh:	make(map[byte]Reactor),
		peers:		NewPeerSet(),
		dialing:	cmn.NewCMap(),
	}

	sw.rng = cmn.NewRand()

	sw.peerConfig.MConfig.FlushThrottle = time.Duration(config.FlushThrottleTimeout) * time.Millisecond
	sw.peerConfig.MConfig.SendRate = config.SendRate
	sw.peerConfig.MConfig.RecvRate = config.RecvRate
	sw.peerConfig.MConfig.MaxPacketMsgPayloadSize = config.MaxPacketMsgPayloadSize
	sw.peerConfig.AuthEnc = config.AuthEnc

	sw.BaseService = *cmn.NewBaseService(nil, "P2P Switch", sw)
	return sw
}

func (sw *Switch) AddReactor(name string, reactor Reactor) Reactor {

	reactorChannels := reactor.GetChannels()
	for _, chDesc := range reactorChannels {
		chID := chDesc.ID
		if sw.reactorsByCh[chID] != nil {
			cmn.PanicSanity(fmt.Sprintf("Channel %X has multiple reactors %v & %v", chID, sw.reactorsByCh[chID], reactor))
		}
		sw.chDescs = append(sw.chDescs, chDesc)
		sw.reactorsByCh[chID] = reactor
	}
	sw.reactors[name] = reactor
	reactor.SetSwitch(sw)
	return reactor
}

func (sw *Switch) Reactors() map[string]Reactor {
	return sw.reactors
}

func (sw *Switch) Reactor(name string) Reactor {
	return sw.reactors[name]
}

func (sw *Switch) AddListener(l Listener) {
	sw.listeners = append(sw.listeners, l)
}

func (sw *Switch) Listeners() []Listener {
	return sw.listeners
}

func (sw *Switch) IsListening() bool {
	return len(sw.listeners) > 0
}

func (sw *Switch) SetNodeInfo(nodeInfo NodeInfo) {
	sw.nodeInfo = nodeInfo
}

func (sw *Switch) NodeInfo() NodeInfo {
	return sw.nodeInfo
}

func (sw *Switch) SetNodeKey(nodeKey *NodeKey) {
	sw.nodeKey = nodeKey
}

func (sw *Switch) OnStart() error {

	for _, reactor := range sw.reactors {
		err := reactor.Start()
		if err != nil {
			return cmn.ErrorWrap(err, "failed to start %v", reactor)
		}
	}

	for _, listener := range sw.listeners {
		go sw.listenerRoutine(listener)
	}
	return nil
}

func (sw *Switch) OnStop() {

	for _, listener := range sw.listeners {
		listener.Stop()
	}
	sw.listeners = nil

	for _, peer := range sw.peers.List() {
		peer.Stop()
		sw.peers.Remove(peer)
	}

	sw.Logger.Debug("Switch: Stopping reactors")
	for _, reactor := range sw.reactors {
		reactor.Stop()
	}
}

func (sw *Switch) Broadcast(chID byte, msgBytes []byte) chan bool {
	successChan := make(chan bool, len(sw.peers.List()))
	sw.Logger.Debug("Broadcast", "channel", chID, "msgBytes", fmt.Sprintf("%X", msgBytes))
	var wg sync.WaitGroup
	for _, peer := range sw.peers.List() {
		wg.Add(1)
		go func(peer Peer) {
			defer wg.Done()
			success := peer.Send(chID, msgBytes)
			successChan <- success
		}(peer)
	}
	go func() {
		wg.Wait()
		close(successChan)
	}()
	return successChan
}

func (sw *Switch) NumPeers() (outbound, inbound, dialing int) {
	peers := sw.peers.List()
	for _, peer := range peers {
		if peer.IsOutbound() {
			outbound++
		} else {
			inbound++
		}
	}
	dialing = sw.dialing.Size()
	return
}

func (sw *Switch) Peers() IPeerSet {
	return sw.peers
}

func (sw *Switch) StopPeerForError(peer Peer, reason interface{}) {
	sw.Logger.Error("Stopping peer for error", "peer", peer, "err", reason)
	sw.stopAndRemovePeer(peer, reason)

	if peer.IsPersistent() {
		go sw.reconnectToPeer(peer)
	}
}

func (sw *Switch) StopPeerGracefully(peer Peer) {
	sw.Logger.Info("Stopping peer gracefully")
	sw.stopAndRemovePeer(peer, nil)
}

func (sw *Switch) stopAndRemovePeer(peer Peer, reason interface{}) {
	sw.peers.Remove(peer)
	peer.Stop()
	for _, reactor := range sw.reactors {
		reactor.RemovePeer(peer, reason)
	}
}

func (sw *Switch) reconnectToPeer(peer Peer) {

	netAddr := peer.NodeInfo().NetAddress()
	start := time.Now()
	sw.Logger.Info("Reconnecting to peer", "peer", peer)
	for i := 0; i < reconnectAttempts; i++ {
		if !sw.IsRunning() {
			return
		}

		err := sw.DialPeerWithAddress(netAddr, true)
		if err != nil {
			sw.Logger.Info("Error reconnecting to peer. Trying again", "tries", i, "err", err, "peer", peer)

			sw.randomSleep(reconnectInterval)
			continue
		} else {
			return
		}
	}

	sw.Logger.Error("Failed to reconnect to peer. Beginning exponential backoff",
		"peer", peer, "elapsed", time.Since(start))
	for i := 0; i < reconnectBackOffAttempts; i++ {
		if !sw.IsRunning() {
			return
		}

		sleepIntervalSeconds := math.Pow(reconnectBackOffBaseSeconds, float64(i))
		sw.randomSleep(time.Duration(sleepIntervalSeconds) * time.Second)
		err := sw.DialPeerWithAddress(netAddr, true)
		if err == nil {
			return
		}
		sw.Logger.Info("Error reconnecting to peer. Trying again", "tries", i, "err", err, "peer", peer)
	}
	sw.Logger.Error("Failed to reconnect to peer. Giving up", "peer", peer, "elapsed", time.Since(start))
}

func (sw *Switch) SetAddrBook(addrBook AddrBook) {
	sw.addrBook = addrBook
}

func (sw *Switch) MarkPeerAsGood(peer Peer) {
	if sw.addrBook != nil {
		sw.addrBook.MarkGood(peer.NodeInfo().NetAddress())
	}
}

func (sw *Switch) IsDialing(id ID) bool {
	return sw.dialing.Has(string(id))
}

func (sw *Switch) DialPeersAsync(addrBook AddrBook, peers []string, persistent bool) error {
	netAddrs, errs := NewNetAddressStrings(peers)

	for _, err := range errs {
		sw.Logger.Error("Error in peer's address", "err", err)
	}

	ourAddr := sw.nodeInfo.NetAddress()

	if addrBook != nil {

		for _, netAddr := range netAddrs {

			if !netAddr.Same(ourAddr) {
				if err := addrBook.AddAddress(netAddr, ourAddr); err != nil {
					sw.Logger.Error("Can't add peer's address to addrbook", "err", err)
				}
			}
		}

		addrBook.Save()
	}

	perm := sw.rng.Perm(len(netAddrs))
	for i := 0; i < len(perm); i++ {
		go func(i int) {
			j := perm[i]

			if netAddrs[j].Same(ourAddr) {
				return
			}

			sw.randomSleep(0)
			err := sw.DialPeerWithAddress(netAddrs[j], persistent)
			if err != nil {
				sw.Logger.Error("Error dialing peer", "err", err)
			}
		}(i)
	}
	return nil
}

func (sw *Switch) DialPeerWithAddress(addr *NetAddress, persistent bool) error {
	sw.dialing.Set(string(addr.ID), addr)
	defer sw.dialing.Delete(string(addr.ID))
	return sw.addOutboundPeerWithConfig(addr, sw.peerConfig, persistent)
}

func (sw *Switch) randomSleep(interval time.Duration) {
	r := time.Duration(sw.rng.Int63n(dialRandomizerIntervalMilliseconds)) * time.Millisecond
	time.Sleep(r + interval)
}

func (sw *Switch) FilterConnByAddr(addr net.Addr) error {
	if sw.filterConnByAddr != nil {
		return sw.filterConnByAddr(addr)
	}
	return nil
}

func (sw *Switch) FilterConnByID(id ID) error {
	if sw.filterConnByID != nil {
		return sw.filterConnByID(id)
	}
	return nil

}

func (sw *Switch) SetAddrFilter(f func(net.Addr) error) {
	sw.filterConnByAddr = f
}

func (sw *Switch) SetIDFilter(f func(ID) error) {
	sw.filterConnByID = f
}

func (sw *Switch) listenerRoutine(l Listener) {
	for {
		inConn, ok := <-l.Connections()
		if !ok {
			break
		}

		maxPeers := sw.config.MaxNumPeers
		if maxPeers <= sw.peers.Size() {
			sw.Logger.Info("Ignoring inbound connection: already have enough peers", "address", inConn.RemoteAddr().String(), "numPeers", sw.peers.Size(), "max", maxPeers)
			continue
		}

		err := sw.addInboundPeerWithConfig(inConn, sw.peerConfig)
		if err != nil {
			sw.Logger.Info("Ignoring inbound connection: error while adding peer", "address", inConn.RemoteAddr().String(), "err", err)
			continue
		}
	}

}

func (sw *Switch) addInboundPeerWithConfig(conn net.Conn, config *PeerConfig) error {
	peerConn, err := newInboundPeerConn(conn, config, sw.nodeKey.PrivKey)
	if err != nil {
		conn.Close()
		return err
	}
	if err = sw.addPeer(peerConn); err != nil {
		peerConn.CloseConn()
		return err
	}

	return nil
}

func (sw *Switch) addOutboundPeerWithConfig(addr *NetAddress, config *PeerConfig, persistent bool) error {
	sw.Logger.Info("Dialing peer", "address", addr)
	peerConn, err := newOutboundPeerConn(addr, config, persistent, sw.nodeKey.PrivKey)
	if err != nil {
		sw.Logger.Error("Failed to dial peer", "address", addr, "err", err)
		return err
	}

	if err := sw.addPeer(peerConn); err != nil {
		sw.Logger.Error("Failed to add peer", "address", addr, "err", err)
		peerConn.CloseConn()
		return err
	}
	return nil
}

func (sw *Switch) addPeer(pc peerConn) error {

	addr := pc.conn.RemoteAddr()
	if err := sw.FilterConnByAddr(addr); err != nil {
		return err
	}

	peerNodeInfo, err := pc.HandshakeTimeout(sw.nodeInfo, time.Duration(sw.peerConfig.HandshakeTimeout*time.Second))
	if err != nil {
		return err
	}

	peerID := peerNodeInfo.ID

	if pc.config.AuthEnc {
		connID := pc.ID()

		if peerID != connID {
			return fmt.Errorf("nodeInfo.ID() (%v) doesn't match conn.ID() (%v)",
				peerID, connID)
		}
	}

	if err := peerNodeInfo.Validate(); err != nil {
		return err
	}

	if sw.nodeKey.ID() == peerID {
		addr := peerNodeInfo.NetAddress()

		sw.addrBook.RemoveAddress(addr)

		sw.addrBook.AddOurAddress(addr)

		return ErrSwitchConnectToSelf
	}

	if sw.peers.Has(peerID) {
		return ErrSwitchDuplicatePeer
	}

	if err := sw.FilterConnByID(peerID); err != nil {
		return err
	}

	if err := sw.nodeInfo.CompatibleWith(peerNodeInfo); err != nil {
		return err
	}

	peer := newPeer(pc, peerNodeInfo, sw.reactorsByCh, sw.chDescs, sw.StopPeerForError)
	peer.SetLogger(sw.Logger.With("peer", addr))

	peer.Logger.Info("Successful handshake with peer", "peerNodeInfo", peerNodeInfo)

	if sw.IsRunning() {
		if err = sw.startInitPeer(peer); err != nil {
			return err
		}
	}

	if err := sw.peers.Add(peer); err != nil {
		return err
	}

	sw.Logger.Info("Added peer", "peer", peer)
	return nil
}

func (sw *Switch) startInitPeer(peer *peer) error {
	err := peer.Start()
	if err != nil {

		sw.Logger.Error("Error starting peer", "peer", peer, "err", err)
		return err
	}

	for _, reactor := range sw.reactors {
		reactor.AddPeer(peer)
	}

	return nil
}
