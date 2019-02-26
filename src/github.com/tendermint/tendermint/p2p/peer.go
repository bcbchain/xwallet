package p2p

import (
	"fmt"
	"net"
	"time"

	"github.com/tendermint/go-crypto"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/log"

	tmconn "github.com/tendermint/tendermint/p2p/conn"
)

type Peer interface {
	cmn.Service

	ID() ID
	IsOutbound() bool
	IsPersistent() bool
	NodeInfo() NodeInfo
	Status() tmconn.ConnectionStatus

	Send(byte, []byte) bool
	TrySend(byte, []byte) bool

	Set(string, interface{})
	Get(string) interface{}
}

type peerConn struct {
	outbound	bool
	persistent	bool
	config		*PeerConfig
	conn		net.Conn
}

func (pc peerConn) ID() ID {
	return PubKeyToID(pc.conn.(*tmconn.SecretConnection).RemotePubKey())
}

type peer struct {
	cmn.BaseService

	peerConn
	mconn	*tmconn.MConnection

	nodeInfo	NodeInfo
	channels	[]byte

	Data	*cmn.CMap
}

func newPeer(pc peerConn, nodeInfo NodeInfo,
	reactorsByCh map[byte]Reactor, chDescs []*tmconn.ChannelDescriptor,
	onPeerError func(Peer, interface{})) *peer {

	p := &peer{
		peerConn:	pc,
		nodeInfo:	nodeInfo,
		channels:	nodeInfo.Channels,
		Data:		cmn.NewCMap(),
	}
	p.mconn = createMConnection(pc.conn, p, reactorsByCh, chDescs, onPeerError, pc.config.MConfig)
	p.BaseService = *cmn.NewBaseService(nil, "Peer", p)
	return p
}

type PeerConfig struct {
	AuthEnc	bool	`mapstructure:"auth_enc"`

	HandshakeTimeout	time.Duration	`mapstructure:"handshake_timeout"`
	DialTimeout		time.Duration	`mapstructure:"dial_timeout"`

	MConfig	*tmconn.MConnConfig	`mapstructure:"connection"`

	Fuzz		bool		`mapstructure:"fuzz"`
	FuzzConfig	*FuzzConnConfig	`mapstructure:"fuzz_config"`
}

func DefaultPeerConfig() *PeerConfig {
	return &PeerConfig{
		AuthEnc:		true,
		HandshakeTimeout:	20,
		DialTimeout:		3,
		MConfig:		tmconn.DefaultMConnConfig(),
		Fuzz:			false,
		FuzzConfig:		DefaultFuzzConnConfig(),
	}
}

func newOutboundPeerConn(addr *NetAddress, config *PeerConfig, persistent bool, ourNodePrivKey crypto.PrivKey) (peerConn, error) {
	var pc peerConn

	conn, err := dial(addr, config)
	if err != nil {
		return pc, cmn.ErrorWrap(err, "Error creating peer")
	}

	pc, err = newPeerConn(conn, config, true, persistent, ourNodePrivKey)
	if err != nil {
		if err2 := conn.Close(); err2 != nil {
			return pc, cmn.ErrorWrap(err, err2.Error())
		}
		return pc, err
	}

	if config.AuthEnc && addr.ID != pc.ID() {
		if err2 := conn.Close(); err2 != nil {
			return pc, cmn.ErrorWrap(err, err2.Error())
		}
		return pc, ErrSwitchAuthenticationFailure{addr, pc.ID()}
	}
	return pc, nil
}

func newInboundPeerConn(conn net.Conn, config *PeerConfig, ourNodePrivKey crypto.PrivKey) (peerConn, error) {

	return newPeerConn(conn, config, false, false, ourNodePrivKey)
}

func newPeerConn(rawConn net.Conn,
	config *PeerConfig, outbound, persistent bool,
	ourNodePrivKey crypto.PrivKey) (pc peerConn, err error) {

	conn := rawConn

	if config.Fuzz {

		conn = FuzzConnAfterFromConfig(conn, 10*time.Second, config.FuzzConfig)
	}

	if config.AuthEnc {

		if err := conn.SetDeadline(time.Now().Add(config.HandshakeTimeout * time.Second)); err != nil {
			return pc, cmn.ErrorWrap(err, "Error setting deadline while encrypting connection")
		}

		conn, err = tmconn.MakeSecretConnection(conn, ourNodePrivKey)
		if err != nil {
			return pc, cmn.ErrorWrap(err, "Error creating peer")
		}
	}

	return peerConn{
		config:		config,
		outbound:	outbound,
		persistent:	persistent,
		conn:		conn,
	}, nil
}

func (p *peer) SetLogger(l log.Logger) {
	p.Logger = l
	p.mconn.SetLogger(l)
}

func (p *peer) OnStart() error {
	if err := p.BaseService.OnStart(); err != nil {
		return err
	}
	err := p.mconn.Start()
	return err
}

func (p *peer) OnStop() {
	p.BaseService.OnStop()
	p.mconn.Stop()
}

func (p *peer) ID() ID {
	return p.nodeInfo.ID
}

func (p *peer) IsOutbound() bool {
	return p.peerConn.outbound
}

func (p *peer) IsPersistent() bool {
	return p.peerConn.persistent
}

func (p *peer) NodeInfo() NodeInfo {
	return p.nodeInfo
}

func (p *peer) Status() tmconn.ConnectionStatus {
	return p.mconn.Status()
}

func (p *peer) Send(chID byte, msgBytes []byte) bool {
	if !p.IsRunning() {

		return false
	} else if !p.hasChannel(chID) {
		return false
	}
	return p.mconn.Send(chID, msgBytes)
}

func (p *peer) TrySend(chID byte, msgBytes []byte) bool {
	if !p.IsRunning() {
		return false
	} else if !p.hasChannel(chID) {
		return false
	}
	return p.mconn.TrySend(chID, msgBytes)
}

func (p *peer) Get(key string) interface{} {
	return p.Data.Get(key)
}

func (p *peer) Set(key string, data interface{}) {
	p.Data.Set(key, data)
}

func (p *peer) hasChannel(chID byte) bool {
	for _, ch := range p.channels {
		if ch == chID {
			return true
		}
	}

	p.Logger.Debug("Unknown channel for peer", "channel", chID, "channels", p.channels)
	return false
}

func (pc *peerConn) CloseConn() {
	pc.conn.Close()
}

func (pc *peerConn) HandshakeTimeout(ourNodeInfo NodeInfo, timeout time.Duration) (peerNodeInfo NodeInfo, err error) {

	if err := pc.conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return peerNodeInfo, cmn.ErrorWrap(err, "Error setting deadline")
	}

	var trs, _ = cmn.Parallel(
		func(_ int) (val interface{}, err error, abort bool) {
			_, err = cdc.MarshalBinaryWriter(pc.conn, ourNodeInfo)
			return
		},
		func(_ int) (val interface{}, err error, abort bool) {
			_, err = cdc.UnmarshalBinaryReader(pc.conn, &peerNodeInfo, int64(MaxNodeInfoSize()))
			return
		},
	)
	if err := trs.FirstError(); err != nil {
		return peerNodeInfo, cmn.ErrorWrap(err, "Error during handshake")
	}

	if err := pc.conn.SetDeadline(time.Time{}); err != nil {
		return peerNodeInfo, cmn.ErrorWrap(err, "Error removing deadline")
	}

	return peerNodeInfo, nil
}

func (p *peer) Addr() net.Addr {
	return p.peerConn.conn.RemoteAddr()
}

func (p *peer) CanSend(chID byte) bool {
	if !p.IsRunning() {
		return false
	}
	return p.mconn.CanSend(chID)
}

func (p *peer) String() string {
	if p.outbound {
		return fmt.Sprintf("Peer{%v %v out}", p.mconn, p.ID())
	}

	return fmt.Sprintf("Peer{%v %v in}", p.mconn, p.ID())
}

func dial(addr *NetAddress, config *PeerConfig) (net.Conn, error) {
	conn, err := addr.DialTimeout(config.DialTimeout * time.Second)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func createMConnection(conn net.Conn, p *peer, reactorsByCh map[byte]Reactor, chDescs []*tmconn.ChannelDescriptor,
	onPeerError func(Peer, interface{}), config *tmconn.MConnConfig) *tmconn.MConnection {

	onReceive := func(chID byte, msgBytes []byte) {
		reactor := reactorsByCh[chID]
		if reactor == nil {

			panic(cmn.Fmt("Unknown channel %X", chID))
		}
		reactor.Receive(chID, p, msgBytes)
	}

	onError := func(r interface{}) {
		onPeerError(p, r)
	}

	return tmconn.NewMConnectionWithConfig(conn, chDescs, onReceive, onError, config)
}
