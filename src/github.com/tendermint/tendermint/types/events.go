package types

import (
	"fmt"

	"github.com/tendermint/go-amino"
	tmpubsub "github.com/tendermint/tmlibs/pubsub"
	tmquery "github.com/tendermint/tmlibs/pubsub/query"
)

const (
	EventBond		= "Bond"
	EventCompleteProposal	= "CompleteProposal"
	EventDupeout		= "Dupeout"
	EventFork		= "Fork"
	EventLock		= "Lock"
	EventNewBlock		= "NewBlock"
	EventNewBlockHeader	= "NewBlockHeader"
	EventNewRound		= "NewRound"
	EventNewRoundStep	= "NewRoundStep"
	EventPolka		= "Polka"
	EventRebond		= "Rebond"
	EventRelock		= "Relock"
	EventTimeoutPropose	= "TimeoutPropose"
	EventTimeoutWait	= "TimeoutWait"
	EventTx			= "Tx"
	EventUnbond		= "Unbond"
	EventUnlock		= "Unlock"
	EventVote		= "Vote"
	EventProposalHeartbeat	= "ProposalHeartbeat"
)

type TMEventData interface {
	AssertIsTMEventData()
}

func (_ EventDataNewBlock) AssertIsTMEventData()		{}
func (_ EventDataNewBlockHeader) AssertIsTMEventData()		{}
func (_ EventDataTx) AssertIsTMEventData()			{}
func (_ EventDataRoundState) AssertIsTMEventData()		{}
func (_ EventDataVote) AssertIsTMEventData()			{}
func (_ EventDataProposalHeartbeat) AssertIsTMEventData()	{}
func (_ EventDataString) AssertIsTMEventData()			{}

func RegisterEventDatas(cdc *amino.Codec) {
	cdc.RegisterInterface((*TMEventData)(nil), nil)
	cdc.RegisterConcrete(EventDataNewBlock{}, "tendermint/event/NewBlock", nil)
	cdc.RegisterConcrete(EventDataNewBlockHeader{}, "tendermint/event/NewBlockHeader", nil)
	cdc.RegisterConcrete(EventDataTx{}, "tendermint/event/Tx", nil)
	cdc.RegisterConcrete(EventDataRoundState{}, "tendermint/event/RoundState", nil)
	cdc.RegisterConcrete(EventDataVote{}, "tendermint/event/Vote", nil)
	cdc.RegisterConcrete(EventDataProposalHeartbeat{}, "tendermint/event/ProposalHeartbeat", nil)
	cdc.RegisterConcrete(EventDataString(""), "tendermint/event/ProposalString", nil)
}

type EventDataNewBlock struct {
	Block *Block `json:"block"`
}

type EventDataNewBlockHeader struct {
	Header *Header `json:"header"`
}

type EventDataTx struct {
	TxResult
}

type EventDataProposalHeartbeat struct {
	Heartbeat *Heartbeat
}

type EventDataRoundState struct {
	Height	int64	`json:"height"`
	Round	int	`json:"round"`
	Step	string	`json:"step"`

	RoundState	interface{}	`json:"-"`
}

type EventDataVote struct {
	Vote *Vote
}

type EventDataString string

const (
	EventTypeKey	= "tm.event"

	TxHashKey	= "tx.hash"

	TxHeightKey	= "tx.height"
)

var (
	EventQueryBond			= QueryForEvent(EventBond)
	EventQueryUnbond		= QueryForEvent(EventUnbond)
	EventQueryRebond		= QueryForEvent(EventRebond)
	EventQueryDupeout		= QueryForEvent(EventDupeout)
	EventQueryFork			= QueryForEvent(EventFork)
	EventQueryNewBlock		= QueryForEvent(EventNewBlock)
	EventQueryNewBlockHeader	= QueryForEvent(EventNewBlockHeader)
	EventQueryNewRound		= QueryForEvent(EventNewRound)
	EventQueryNewRoundStep		= QueryForEvent(EventNewRoundStep)
	EventQueryTimeoutPropose	= QueryForEvent(EventTimeoutPropose)
	EventQueryCompleteProposal	= QueryForEvent(EventCompleteProposal)
	EventQueryPolka			= QueryForEvent(EventPolka)
	EventQueryUnlock		= QueryForEvent(EventUnlock)
	EventQueryLock			= QueryForEvent(EventLock)
	EventQueryRelock		= QueryForEvent(EventRelock)
	EventQueryTimeoutWait		= QueryForEvent(EventTimeoutWait)
	EventQueryVote			= QueryForEvent(EventVote)
	EventQueryProposalHeartbeat	= QueryForEvent(EventProposalHeartbeat)
	EventQueryTx			= QueryForEvent(EventTx)
)

func EventQueryTxFor(tx Tx) tmpubsub.Query {
	return tmquery.MustParse(fmt.Sprintf("%s='%s' AND %s='%X'", EventTypeKey, EventTx, TxHashKey, tx.Hash()))
}

func QueryForEvent(eventType string) tmpubsub.Query {
	return tmquery.MustParse(fmt.Sprintf("%s='%s'", EventTypeKey, eventType))
}

type BlockEventPublisher interface {
	PublishEventNewBlock(block EventDataNewBlock) error
	PublishEventNewBlockHeader(header EventDataNewBlockHeader) error
	PublishEventTx(EventDataTx) error
}

type TxEventPublisher interface {
	PublishEventTx(EventDataTx) error
}
