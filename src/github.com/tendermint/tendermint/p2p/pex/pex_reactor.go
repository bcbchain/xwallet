package pex

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/tendermint/go-amino"
	cmn "github.com/tendermint/tmlibs/common"

	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/p2p/conn"
)

type Peer = p2p.Peer

const (
	PexChannel	= byte(0x00)

	maxMsgSize	= 1048576

	defaultEnsurePeersPeriod	= 30 * time.Second
	defaultMinNumOutboundPeers	= 10

	defaultSeedDisconnectWaitPeriod	= 3 * time.Hour

	defaultCrawlPeerInterval	= 2 * time.Minute

	defaultCrawlPeersPeriod	= 30 * time.Second

	maxAttemptsToDial	= 16

	biasToSelectNewPeers	= 30
)

type PEXReactor struct {
	p2p.BaseReactor

	book			AddrBook
	config			*PEXReactorConfig
	ensurePeersPeriod	time.Duration

	requestsSent		*cmn.CMap
	lastReceivedRequests	*cmn.CMap

	attemptsToDial	sync.Map
}

type PEXReactorConfig struct {
	SeedMode	bool

	Seeds	[]string

	PrivatePeerIDs	[]string
}

type _attemptsToDial struct {
	number		int
	lastDialed	time.Time
}

func NewPEXReactor(b AddrBook, config *PEXReactorConfig) *PEXReactor {
	r := &PEXReactor{
		book:			b,
		config:			config,
		ensurePeersPeriod:	defaultEnsurePeersPeriod,
		requestsSent:		cmn.NewCMap(),
		lastReceivedRequests:	cmn.NewCMap(),
	}
	r.BaseReactor = *p2p.NewBaseReactor("PEXReactor", r)
	return r
}

func (r *PEXReactor) OnStart() error {
	if err := r.BaseReactor.OnStart(); err != nil {
		return err
	}
	err := r.book.Start()
	if err != nil && err != cmn.ErrAlreadyStarted {
		return err
	}

	if err := r.checkSeeds(); err != nil {
		return err
	}

	if r.config.SeedMode {
		go r.crawlPeersRoutine()
	} else {
		go r.ensurePeersRoutine()
	}
	return nil
}

func (r *PEXReactor) OnStop() {
	r.BaseReactor.OnStop()
	r.book.Stop()
}

func (r *PEXReactor) GetChannels() []*conn.ChannelDescriptor {
	return []*conn.ChannelDescriptor{
		{
			ID:			PexChannel,
			Priority:		1,
			SendQueueCapacity:	10,
		},
	}
}

func (r *PEXReactor) AddPeer(p Peer) {
	if p.IsOutbound() {

		if r.book.NeedMoreAddrs() {
			r.RequestAddrs(p)
		}
	} else {

		addr := p.NodeInfo().NetAddress()
		if !isAddrPrivate(addr, r.config.PrivatePeerIDs) {
			err := r.book.AddAddress(addr, addr)
			if err != nil {
				r.Logger.Error("Failed to add new address", "err", err)
			}
		}
	}
}

func (r *PEXReactor) RemovePeer(p Peer, reason interface{}) {
	id := string(p.ID())
	r.requestsSent.Delete(id)
	r.lastReceivedRequests.Delete(id)
}

func (r *PEXReactor) Receive(chID byte, src Peer, msgBytes []byte) {
	msg, err := DecodeMessage(msgBytes)
	if err != nil {
		r.Logger.Error("Error decoding message", "src", src, "chId", chID, "msg", msg, "err", err, "bytes", msgBytes)
		r.Switch.StopPeerForError(src, err)
		return
	}
	r.Logger.Debug("Received message", "src", src, "chId", chID, "msg", msg)

	switch msg := msg.(type) {
	case *pexRequestMessage:

		if err := r.receiveRequest(src); err != nil {
			r.Switch.StopPeerForError(src, err)
			return
		}

		if r.config.SeedMode {
			r.SendAddrs(src, r.book.GetSelectionWithBias(biasToSelectNewPeers))
			r.Switch.StopPeerGracefully(src)
		} else {
			r.SendAddrs(src, r.book.GetSelection())
		}

	case *pexAddrsMessage:

		if err := r.ReceiveAddrs(msg.Addrs, src); err != nil {
			r.Switch.StopPeerForError(src, err)
			return
		}
	default:
		r.Logger.Error(fmt.Sprintf("Unknown message type %v", reflect.TypeOf(msg)))
	}
}

func (r *PEXReactor) receiveRequest(src Peer) error {
	id := string(src.ID())
	v := r.lastReceivedRequests.Get(id)
	if v == nil {

		lastReceived := time.Time{}
		r.lastReceivedRequests.Set(id, lastReceived)
		return nil
	}

	lastReceived := v.(time.Time)
	if lastReceived.Equal(time.Time{}) {

		lastReceived = time.Now()
		r.lastReceivedRequests.Set(id, lastReceived)
		return nil
	}

	now := time.Now()
	if now.Sub(lastReceived) < r.ensurePeersPeriod/3 {
		return fmt.Errorf("Peer (%v) is sending too many PEX requests. Disconnecting", src.ID())
	}
	r.lastReceivedRequests.Set(id, now)
	return nil
}

func (r *PEXReactor) RequestAddrs(p Peer) {
	id := string(p.ID())
	if r.requestsSent.Has(id) {
		return
	}
	r.requestsSent.Set(id, struct{}{})
	p.Send(PexChannel, cdc.MustMarshalBinary(&pexRequestMessage{}))
}

func (r *PEXReactor) ReceiveAddrs(addrs []*p2p.NetAddress, src Peer) error {
	id := string(src.ID())

	if !r.requestsSent.Has(id) {
		return cmn.NewError("Received unsolicited pexAddrsMessage")
	}

	r.requestsSent.Delete(id)

	srcAddr := src.NodeInfo().NetAddress()
	for _, netAddr := range addrs {
		if netAddr != nil && !isAddrPrivate(netAddr, r.config.PrivatePeerIDs) {
			err := r.book.AddAddress(netAddr, srcAddr)
			if err != nil {
				r.Logger.Error("Failed to add new address", "err", err)
			}
		}
	}
	return nil
}

func (r *PEXReactor) SendAddrs(p Peer, netAddrs []*p2p.NetAddress) {
	p.Send(PexChannel, cdc.MustMarshalBinary(&pexAddrsMessage{Addrs: netAddrs}))
}

func (r *PEXReactor) SetEnsurePeersPeriod(d time.Duration) {
	r.ensurePeersPeriod = d
}

func (r *PEXReactor) ensurePeersRoutine() {
	var (
		seed	= cmn.NewRand()
		jitter	= seed.Int63n(r.ensurePeersPeriod.Nanoseconds())
	)

	if r.hasPotentialPeers() {
		time.Sleep(time.Duration(jitter))
	}

	r.ensurePeers()

	ticker := time.NewTicker(r.ensurePeersPeriod)
	for {
		select {
		case <-ticker.C:
			r.ensurePeers()
		case <-r.Quit():
			ticker.Stop()
			return
		}
	}
}

func (r *PEXReactor) ensurePeers() {
	var (
		out, in, dial	= r.Switch.NumPeers()
		numToDial	= defaultMinNumOutboundPeers - (out + dial)
	)
	r.Logger.Info(
		"Ensure peers",
		"numOutPeers", out,
		"numInPeers", in,
		"numDialing", dial,
		"numToDial", numToDial,
	)

	if numToDial <= 0 {
		return
	}

	newBias := cmn.MinInt(out, 8)*10 + 10

	toDial := make(map[p2p.ID]*p2p.NetAddress)

	maxAttempts := numToDial * 3

	for i := 0; i < maxAttempts && len(toDial) < numToDial; i++ {
		try := r.book.PickAddress(newBias)
		if try == nil {
			continue
		}
		if _, selected := toDial[try.ID]; selected {
			continue
		}
		if dialling := r.Switch.IsDialing(try.ID); dialling {
			continue
		}
		if connected := r.Switch.Peers().Has(try.ID); connected {
			continue
		}
		r.Logger.Info("Will dial address", "addr", try)
		toDial[try.ID] = try
	}

	for _, addr := range toDial {
		go r.dialPeer(addr)
	}

	if r.book.NeedMoreAddrs() {
		peers := r.Switch.Peers().List()
		peersCount := len(peers)
		if peersCount > 0 {
			peer := peers[cmn.RandInt()%peersCount]
			r.Logger.Info("We need more addresses. Sending pexRequest to random peer", "peer", peer)
			r.RequestAddrs(peer)
		}
	}

	if out+in+dial+len(toDial) == 0 {
		r.Logger.Info("No addresses to dial nor connected peers. Falling back to seeds")
		r.dialSeeds()
	}
}

func (r *PEXReactor) dialPeer(addr *p2p.NetAddress) {
	var attempts int
	var lastDialed time.Time
	if lAttempts, attempted := r.attemptsToDial.Load(addr.DialString()); attempted {
		attempts = lAttempts.(_attemptsToDial).number
		lastDialed = lAttempts.(_attemptsToDial).lastDialed
	}

	if attempts > maxAttemptsToDial {
		r.Logger.Error("Reached max attempts to dial", "addr", addr, "attempts", attempts)
		r.book.MarkBad(addr)
		return
	}

	if attempts > 0 {
		jitterSeconds := time.Duration(cmn.RandFloat64() * float64(time.Second))
		backoffDuration := jitterSeconds + ((1 << uint(attempts)) * time.Second)
		sinceLastDialed := time.Since(lastDialed)
		if sinceLastDialed < backoffDuration {
			r.Logger.Debug("Too early to dial", "addr", addr, "backoff_duration", backoffDuration, "last_dialed", lastDialed, "time_since", sinceLastDialed)
			return
		}
	}

	err := r.Switch.DialPeerWithAddress(addr, false)
	if err != nil {
		r.Logger.Error("Dialing failed", "addr", addr, "err", err, "attempts", attempts)

		if _, ok := err.(p2p.ErrSwitchAuthenticationFailure); ok {
			r.book.MarkBad(addr)
			r.attemptsToDial.Delete(addr.DialString())
		} else {
			r.book.MarkAttempt(addr)

			r.attemptsToDial.Store(addr.DialString(), _attemptsToDial{attempts + 1, time.Now()})
		}
	} else {

		r.attemptsToDial.Delete(addr.DialString())
	}
}

func (r *PEXReactor) checkSeeds() error {
	lSeeds := len(r.config.Seeds)
	if lSeeds == 0 {
		return nil
	}
	_, errs := p2p.NewNetAddressStrings(r.config.Seeds)
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PEXReactor) dialSeeds() {
	lSeeds := len(r.config.Seeds)
	if lSeeds == 0 {
		return
	}
	seedAddrs, _ := p2p.NewNetAddressStrings(r.config.Seeds)

	perm := cmn.RandPerm(lSeeds)

	for _, i := range perm {

		seedAddr := seedAddrs[i]
		err := r.Switch.DialPeerWithAddress(seedAddr, false)
		if err == nil {
			return
		}
		r.Switch.Logger.Error("Error dialing seed", "err", err, "seed", seedAddr)
	}
	r.Switch.Logger.Error("Couldn't connect to any seeds")
}

func (r *PEXReactor) AttemptsToDial(addr *p2p.NetAddress) int {
	lAttempts, attempted := r.attemptsToDial.Load(addr.DialString())
	if attempted {
		return lAttempts.(_attemptsToDial).number
	}
	return 0
}

func (r *PEXReactor) crawlPeersRoutine() {

	r.crawlPeers()

	ticker := time.NewTicker(defaultCrawlPeersPeriod)

	for {
		select {
		case <-ticker.C:
			r.attemptDisconnects()
			r.crawlPeers()
		case <-r.Quit():
			return
		}
	}
}

func (r *PEXReactor) hasPotentialPeers() bool {
	out, in, dial := r.Switch.NumPeers()

	return out+in+dial > 0 && len(r.book.ListOfKnownAddresses()) > 0
}

type crawlPeerInfo struct {
	Addr	*p2p.NetAddress

	LastAttempt	time.Time

	LastSuccess	time.Time
}

type oldestFirst []crawlPeerInfo

func (of oldestFirst) Len() int			{ return len(of) }
func (of oldestFirst) Swap(i, j int)		{ of[i], of[j] = of[j], of[i] }
func (of oldestFirst) Less(i, j int) bool	{ return of[i].LastAttempt.Before(of[j].LastAttempt) }

func (r *PEXReactor) getPeersToCrawl() []crawlPeerInfo {
	var of oldestFirst

	addrs := r.book.ListOfKnownAddresses()
	for _, addr := range addrs {
		if len(addr.ID()) == 0 {
			continue
		}

		of = append(of, crawlPeerInfo{
			Addr:		addr.Addr,
			LastAttempt:	addr.LastAttempt,
			LastSuccess:	addr.LastSuccess,
		})
	}
	sort.Sort(of)
	return of
}

func (r *PEXReactor) crawlPeers() {
	peerInfos := r.getPeersToCrawl()

	now := time.Now()

	for _, pi := range peerInfos {

		if now.Sub(pi.LastAttempt) < defaultCrawlPeerInterval {
			continue
		}

		err := r.Switch.DialPeerWithAddress(pi.Addr, false)
		if err != nil {
			r.book.MarkAttempt(pi.Addr)
			continue
		}

		peer := r.Switch.Peers().Get(pi.Addr.ID)
		r.RequestAddrs(peer)
	}
}

func (r *PEXReactor) attemptDisconnects() {
	for _, peer := range r.Switch.Peers().List() {
		if peer.Status().Duration < defaultSeedDisconnectWaitPeriod {
			continue
		}
		if peer.IsPersistent() {
			continue
		}
		r.Switch.StopPeerGracefully(peer)
	}
}

func isAddrPrivate(addr *p2p.NetAddress, privatePeerIDs []string) bool {
	for _, id := range privatePeerIDs {
		if string(addr.ID) == id {
			return true
		}
	}
	return false
}

type PexMessage interface{}

func RegisterPexMessage(cdc *amino.Codec) {
	cdc.RegisterInterface((*PexMessage)(nil), nil)
	cdc.RegisterConcrete(&pexRequestMessage{}, "tendermint/p2p/PexRequestMessage", nil)
	cdc.RegisterConcrete(&pexAddrsMessage{}, "tendermint/p2p/PexAddrsMessage", nil)
}

func DecodeMessage(bz []byte) (msg PexMessage, err error) {
	if len(bz) > maxMsgSize {
		return msg, fmt.Errorf("Msg exceeds max size (%d > %d)",
			len(bz), maxMsgSize)
	}
	err = cdc.UnmarshalBinary(bz, &msg)
	return
}

type pexRequestMessage struct {
}

func (m *pexRequestMessage) String() string {
	return "[pexRequest]"
}

type pexAddrsMessage struct {
	Addrs []*p2p.NetAddress
}

func (m *pexAddrsMessage) String() string {
	return fmt.Sprintf("[pexAddrs %v]", m.Addrs)
}
