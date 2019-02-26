package core_types

import (
	"encoding/json"
	"strings"
	"time"

	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	cmn "github.com/tendermint/tmlibs/common"

	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/types"
)

type ResultBlockchainInfo struct {
	LastHeight	int64			`json:"last_height"`
	BlockMetas	[]*types.BlockMeta	`json:"block_metas"`
}

type ResultGenesis struct {
	Genesis *types.GenesisDoc `json:"genesis"`
}

type ResultBlock struct {
	BlockMeta	*types.BlockMeta	`json:"block_meta"`
	Block		*types.Block		`json:"block"`
	BlockSize	int			`json:"block_size"`
}

type ResultCommit struct {
	types.SignedHeader
	CanonicalCommit	bool	`json:"canonical"`
}

type ResultBlockResults struct {
	Height	int64			`json:"height"`
	Results	*state.ABCIResponses	`json:"results"`
}

func NewResultCommit(header *types.Header, commit *types.Commit,
	canonical bool) *ResultCommit {

	return &ResultCommit{
		SignedHeader: types.SignedHeader{
			Header:	header,
			Commit:	commit,
		},
		CanonicalCommit:	canonical,
	}
}

type SyncInfo struct {
	LatestBlockHash		cmn.HexBytes	`json:"latest_block_hash"`
	LatestAppHash		cmn.HexBytes	`json:"latest_app_hash"`
	LatestBlockHeight	int64		`json:"latest_block_height"`
	LatestBlockTime		time.Time	`json:"latest_block_time"`
	Syncing			bool		`json:"syncing"`
}

type ValidatorInfo struct {
	Address		string		`json:"address"`
	PubKey		crypto.PubKey	`json:"pub_key"`
	VotingPower	uint64		`json:"voting_power"`
	RewardAddr	string		`json:"reward_addr"`
	Name		string		`json:"name"`
}

type ResultStatus struct {
	NodeInfo	p2p.NodeInfo	`json:"node_info"`
	SyncInfo	SyncInfo	`json:"sync_info"`
	ValidatorInfo	ValidatorInfo	`json:"validator_info"`
}

func (s *ResultStatus) TxIndexEnabled() bool {
	if s == nil {
		return false
	}
	for _, s := range s.NodeInfo.Other {
		info := strings.Split(s, "=")
		if len(info) == 2 && info[0] == "tx_index" {
			return info[1] == "on"
		}
	}
	return false
}

type ResultNetInfo struct {
	Listening	bool		`json:"listening"`
	Listeners	[]string	`json:"listeners"`
	Peers		[]Peer		`json:"peers"`
}

type ResultDialSeeds struct {
	Log string `json:"log"`
}

type ResultDialPeers struct {
	Log string `json:"log"`
}

type Peer struct {
	p2p.NodeInfo		`json:"node_info"`
	IsOutbound		bool			`json:"is_outbound"`
	ConnectionStatus	p2p.ConnectionStatus	`json:"connection_status"`
}

type ResultValidators struct {
	BlockHeight	int64			`json:"block_height"`
	Validators	[]*types.Validator	`json:"validators"`
}

type ResultDumpConsensusState struct {
	RoundState	json.RawMessage		`json:"round_state"`
	PeerRoundStates	[]PeerRoundState	`json:"peer_round_states"`
}

type PeerRoundState struct {
	NodeAddress	string		`json:"node_address"`
	PeerRoundState	json.RawMessage	`json:"peer_round_state"`
}

type ResultBroadcastTx struct {
	Code	uint32		`json:"code"`
	Data	cmn.HexBytes	`json:"data"`
	Log	string		`json:"log"`

	Hash	cmn.HexBytes	`json:"hash"`
}

type ResultBroadcastTxCommit struct {
	CheckTx		abci.ResponseCheckTx	`json:"check_tx,omitempt"`
	DeliverTx	abci.ResponseDeliverTx	`json:"deliver_tx,omitempt"`
	Hash		cmn.HexBytes		`json:"hash,omitempt"`
	Height		int64			`json:"height,omitempt"`
}

type ResultTx struct {
	Hash		string			`json:"hash"`
	Height		int64			`json:"height"`
	Index		uint32			`json:"index"`
	DeliverResult	abci.ResponseDeliverTx	`json:"deliver_tx,omitempt"`
	CheckResult	abci.ResponseCheckTx	`json:"check_tx,omitempt"`
	Tx		types.Tx		`json:"tx"`
	Proof		types.TxProof		`json:"proof,omitempty"`
	StateCode	uint32			`json:"state_code"`
}

type ResultUnconfirmedTxs struct {
	N	int		`json:"n_txs"`
	Txs	[]types.Tx	`json:"txs"`
}

type ResultABCIInfo struct {
	Response abci.ResponseInfo `json:"response"`
}

type ResultABCIQuery struct {
	Response abci.ResponseQuery `json:"response"`
}

type (
	ResultUnsafeFlushMempool	struct{}
	ResultUnsafeProfile		struct{}
	ResultSubscribe			struct{}
	ResultUnsubscribe		struct{}
	ResultHealth			struct{}
)

type ResultEvent struct {
	Query	string			`json:"query"`
	Data	types.TMEventData	`json:"data"`
}
