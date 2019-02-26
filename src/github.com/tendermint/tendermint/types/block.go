package types

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/tendermint/tendermint/softforks"

	"github.com/tendermint/abci/types"
	"strings"
	"sync"
	"time"

	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/merkle"
	"golang.org/x/crypto/ripemd160"
)

type Block struct {
	mtx		sync.Mutex
	*Header		`json:"header"`
	*Data		`json:"data"`
	Evidence	EvidenceData	`json:"evidence"`
	LastCommit	*Commit		`json:"last_commit"`
}

func MakeBlock(height int64, txs []Tx, commit *Commit) *Block {
	block := &Block{
		Header: &Header{
			Height:	height,
			Time:	time.Now(),
			NumTxs:	int64(len(txs)),
		},
		LastCommit:	commit,
		Data: &Data{
			Txs: txs,
		},
	}
	block.fillHeader()
	return block
}

func BCMakeBlock(height int64, txs []Tx, commit *Commit, txHashList [][]byte,
	proposer string, lastFee uint64, rewardAddr string, lastAllocation []types.Allocation) *Block {

	block := &Block{
		Header: &Header{
			Height:			height,
			Time:			time.Now(),
			NumTxs:			int64(len(txs)),
			LastFee:		lastFee,
			LastAllocation:		lastAllocation,
			ProposerAddress:	proposer,
			RewardAddress:		rewardAddr,
		},
		LastCommit:	commit,
		Data: &Data{
			Txs:			txs,
			LastTxsHashList:	txHashList,
		},
	}
	softforks.Init()
	if softforks.IsForkForV1023233(height) {
		r := make([]byte, 32)
		_, e := rand.Read(r)
		if e != nil {
			panic(e)
		}
		block.Header.RandomOfBlock = r
	}
	block.fillHeader()
	return block

}

func (b *Block) AddEvidence(evidence []Evidence) {
	b.Evidence.Evidence = append(b.Evidence.Evidence, evidence...)
}

func (b *Block) ValidateBasic() error {
	if b == nil {
		return errors.New("Nil blocks are invalid")
	}
	b.mtx.Lock()
	defer b.mtx.Unlock()

	newTxs := int64(len(b.Data.Txs))
	if b.NumTxs != newTxs {
		return fmt.Errorf("Wrong Block.Header.NumTxs. Expected %v, got %v", newTxs, b.NumTxs)
	}
	if !bytes.Equal(b.LastCommitHash, b.LastCommit.Hash()) {
		return fmt.Errorf("Wrong Block.Header.LastCommitHash.  Expected %v, got %v", b.LastCommitHash, b.LastCommit.Hash())
	}
	if b.Header.Height != 1 {
		if err := b.LastCommit.ValidateBasic(); err != nil {
			return err
		}
	}
	if !bytes.Equal(b.DataHash, b.Data.Hash()) {
		return fmt.Errorf("Wrong Block.Header.DataHash.  Expected %v, got %v", b.DataHash, b.Data.Hash())
	}
	if !bytes.Equal(b.EvidenceHash, b.Evidence.Hash()) {
		return errors.New(cmn.Fmt("Wrong Block.Header.EvidenceHash.  Expected %v, got %v", b.EvidenceHash, b.Evidence.Hash()))
	}
	return nil
}

func (b *Block) fillHeader() {
	if b.LastCommitHash == nil {
		b.LastCommitHash = b.LastCommit.Hash()
	}
	if b.DataHash == nil {
		b.DataHash = b.Data.Hash()
	}
	if b.EvidenceHash == nil {
		b.EvidenceHash = b.Evidence.Hash()
	}
}

func (b *Block) Hash() cmn.HexBytes {
	if b == nil {
		return nil
	}
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if b == nil || b.Header == nil || b.Data == nil || b.LastCommit == nil {
		return nil
	}
	b.fillHeader()
	return b.Header.Hash()
}

func (b *Block) MakePartSet(partSize int) *PartSet {
	if b == nil {
		return nil
	}
	b.mtx.Lock()
	defer b.mtx.Unlock()

	bz, err := cdc.MarshalBinary(b)
	if err != nil {
		panic(err)
	}
	return NewPartSetFromData(bz, partSize)
}

func (b *Block) HashesTo(hash []byte) bool {
	if len(hash) == 0 {
		return false
	}
	if b == nil {
		return false
	}
	return bytes.Equal(b.Hash(), hash)
}

func (b *Block) String() string {
	return b.StringIndented("")
}

func (b *Block) StringIndented(indent string) string {
	if b == nil {
		return "nil-Block"
	}
	return fmt.Sprintf(`Block{
%s  %v
%s  %v
%s  %v
%s  %v
%s}#%v`,
		indent, b.Header.StringIndented(indent+"  "),
		indent, b.Data.StringIndented(indent+"  "),
		indent, b.Evidence.StringIndented(indent+"  "),
		indent, b.LastCommit.StringIndented(indent+"  "),
		indent, b.Hash())
}

func (b *Block) StringShort() string {
	if b == nil {
		return "nil-Block"
	}
	return fmt.Sprintf("Block#%v", b.Hash())
}

type Allocation []types.Allocation

type Header struct {
	ChainID	string		`json:"chain_id"`
	Height	int64		`json:"height"`
	Time	time.Time	`json:"time"`
	NumTxs	int64		`json:"num_txs"`

	LastBlockID	BlockID	`json:"last_block_id"`
	TotalTxs	int64	`json:"total_txs"`

	LastCommitHash	cmn.HexBytes	`json:"last_commit_hash"`
	DataHash	cmn.HexBytes	`json:"data_hash"`

	ValidatorsHash	cmn.HexBytes	`json:"validators_hash"`
	ConsensusHash	cmn.HexBytes	`json:"consensus_hash"`
	LastAppHash	cmn.HexBytes	`json:"last_app_hash"`
	LastResultsHash	cmn.HexBytes	`json:"last_results_hash"`

	EvidenceHash	cmn.HexBytes	`json:"evidence_hash"`

	LastFee		uint64		`json:"last_fee"`
	LastAllocation	Allocation	`json:"last_allocation"`
	ProposerAddress	string		`json:"proposer_address"`
	RewardAddress	string		`json:"reward_address"`

	RandomOfBlock	cmn.HexBytes	`json:"random_of_block,omitempty"`
}

func (h *Header) Hash() cmn.HexBytes {
	if h == nil || len(h.ValidatorsHash) == 0 {
		return nil
	}
	softforks.Init()
	if softforks.IsForkForV1023233(h.Height) {
		return merkle.SimpleHashFromMap(map[string]merkle.Hasher{
			"ChainID":		aminoHasher(h.ChainID),
			"Height":		aminoHasher(h.Height),
			"Time":			aminoHasher(h.Time),
			"NumTxs":		aminoHasher(h.NumTxs),
			"TotalTxs":		aminoHasher(h.TotalTxs),
			"LastBlockID":		aminoHasher(h.LastBlockID),
			"LastCommit":		aminoHasher(h.LastCommitHash),
			"Data":			aminoHasher(h.DataHash),
			"Validators":		aminoHasher(h.ValidatorsHash),
			"LastApp":		aminoHasher(h.LastAppHash),
			"Consensus":		aminoHasher(h.ConsensusHash),
			"Results":		aminoHasher(h.LastResultsHash),
			"Evidence":		aminoHasher(h.EvidenceHash),
			"LastFee":		aminoHasher(h.LastFee),
			"LastAllocation":	aminoHasher(h.LastAllocation),
			"Proposer":		aminoHasher(h.ProposerAddress),
			"RewardAddr":		aminoHasher(h.RewardAddress),
			"RandomOfBlock":	aminoHasher(h.RandomOfBlock),
		})
	} else {
		return merkle.SimpleHashFromMap(map[string]merkle.Hasher{
			"ChainID":		aminoHasher(h.ChainID),
			"Height":		aminoHasher(h.Height),
			"Time":			aminoHasher(h.Time),
			"NumTxs":		aminoHasher(h.NumTxs),
			"TotalTxs":		aminoHasher(h.TotalTxs),
			"LastBlockID":		aminoHasher(h.LastBlockID),
			"LastCommit":		aminoHasher(h.LastCommitHash),
			"Data":			aminoHasher(h.DataHash),
			"Validators":		aminoHasher(h.ValidatorsHash),
			"LastApp":		aminoHasher(h.LastAppHash),
			"Consensus":		aminoHasher(h.ConsensusHash),
			"Results":		aminoHasher(h.LastResultsHash),
			"Evidence":		aminoHasher(h.EvidenceHash),
			"LastFee":		aminoHasher(h.LastFee),
			"LastAllocation":	aminoHasher(h.LastAllocation),
			"Proposer":		aminoHasher(h.ProposerAddress),
			"RewardAddr":		aminoHasher(h.RewardAddress),
		})
	}
}

func (as *Allocation) StringIndented(indent string) string {
	if as == nil {
		return "[]"
	}
	res := "["
	for i, v := range *as {
		res += fmt.Sprintf(`{Addr:%s, Fee:%d}`, v.Addr, v.Fee)
		if i != len(*as)-1 {
			res += "," + indent
		}
	}
	return res + "]"
}

func (h *Header) StringIndented(indent string) string {
	if h == nil {
		return "nil-Header"
	}
	return fmt.Sprintf(`Header{
%s  ChainID:        %v
%s  Height:         %v
%s  Time:           %v
%s  NumTxs:         %v
%s  TotalTxs:       %v
%s  LastBlockID:    %v
%s  LastCommit:     %v
%s  Data:           %v
%s  Validators:     %v
%s  LastApp:        %v
%s  Consensus:      %v
%s  Results:        %v
%s  Evidence:       %v
%s  LastFee:        %v
%s  LastAllocation: %v
%s  Proposer:       %v
%s  RewardAddr:     %v
%s  RandomOfBlock:  %v
%s}#%v`,
		indent, h.ChainID,
		indent, h.Height,
		indent, h.Time,
		indent, h.NumTxs,
		indent, h.TotalTxs,
		indent, h.LastBlockID,
		indent, h.LastCommitHash,
		indent, h.DataHash,
		indent, h.ValidatorsHash,
		indent, h.LastAppHash,
		indent, h.ConsensusHash,
		indent, h.LastResultsHash,
		indent, h.EvidenceHash,
		indent, h.LastFee,
		indent, h.LastAllocation.StringIndented(" "),
		indent, h.ProposerAddress,
		indent, h.RewardAddress,
		indent, h.RandomOfBlock,
		indent, h.Hash())
}

type Commit struct {
	BlockID		BlockID	`json:"block_id"`
	Precommits	[]*Vote	`json:"precommits"`

	firstPrecommit	*Vote
	hash		cmn.HexBytes
	bitArray	*cmn.BitArray
}

func (commit *Commit) FirstPrecommit() *Vote {
	if len(commit.Precommits) == 0 {
		return nil
	}
	if commit.firstPrecommit != nil {
		return commit.firstPrecommit
	}
	for _, precommit := range commit.Precommits {
		if precommit != nil {
			commit.firstPrecommit = precommit
			return precommit
		}
	}
	return &Vote{
		Type: VoteTypePrecommit,
	}
}

func (commit *Commit) Height() int64 {
	if len(commit.Precommits) == 0 {
		return 0
	}
	return commit.FirstPrecommit().Height
}

func (commit *Commit) Round() int {
	if len(commit.Precommits) == 0 {
		return 0
	}
	return commit.FirstPrecommit().Round
}

func (commit *Commit) Type() byte {
	return VoteTypePrecommit
}

func (commit *Commit) Size() int {
	if commit == nil {
		return 0
	}
	return len(commit.Precommits)
}

func (commit *Commit) BitArray() *cmn.BitArray {
	if commit.bitArray == nil {
		commit.bitArray = cmn.NewBitArray(len(commit.Precommits))
		for i, precommit := range commit.Precommits {

			commit.bitArray.SetIndex(i, precommit != nil)
		}
	}
	return commit.bitArray
}

func (commit *Commit) GetByIndex(index int) *Vote {
	return commit.Precommits[index]
}

func (commit *Commit) IsCommit() bool {
	return len(commit.Precommits) != 0
}

func (commit *Commit) ValidateBasic() error {
	if commit.BlockID.IsZero() {
		return errors.New("Commit cannot be for nil block")
	}
	if len(commit.Precommits) == 0 {
		return errors.New("No precommits in commit")
	}
	height, round := commit.Height(), commit.Round()

	for _, precommit := range commit.Precommits {

		if precommit == nil {
			continue
		}

		if precommit.Type != VoteTypePrecommit {
			return fmt.Errorf("Invalid commit vote. Expected precommit, got %v",
				precommit.Type)
		}

		if precommit.Height != height {
			return fmt.Errorf("Invalid commit precommit height. Expected %v, got %v",
				height, precommit.Height)
		}

		if precommit.Round != round {
			return fmt.Errorf("Invalid commit precommit round. Expected %v, got %v",
				round, precommit.Round)
		}
	}
	return nil
}

func (commit *Commit) Hash() cmn.HexBytes {
	if commit.hash == nil {
		bs := make([]merkle.Hasher, len(commit.Precommits))
		for i, precommit := range commit.Precommits {
			bs[i] = aminoHasher(precommit)
		}
		commit.hash = merkle.SimpleHashFromHashers(bs)
	}
	return commit.hash
}

func (commit *Commit) StringIndented(indent string) string {
	if commit == nil {
		return "nil-Commit"
	}
	precommitStrings := make([]string, len(commit.Precommits))
	for i, precommit := range commit.Precommits {
		precommitStrings[i] = precommit.String()
	}
	return fmt.Sprintf(`Commit{
%s  BlockID:    %v
%s  Precommits: %v
%s}#%v`,
		indent, commit.BlockID,
		indent, strings.Join(precommitStrings, "\n"+indent+"  "),
		indent, commit.hash)
}

type SignedHeader struct {
	Header	*Header	`json:"header"`
	Commit	*Commit	`json:"commit"`
}

type Data struct {
	Txs	Txs	`json:"txs"`

	hash	cmn.HexBytes

	LastTxsHashList	HashList	`json:"lastTxsHashList"`
}

func (data *Data) Hash() cmn.HexBytes {
	if data == nil {
		return (Txs{}).Hash()
	}
	if data.hash == nil {
		data.hash = data.Txs.Hash()
	}
	return data.hash
}

func (data *Data) StringIndented(indent string) string {
	if data == nil {
		return "nil-Data"
	}
	txStrings := make([]string, cmn.MinInt(len(data.Txs), 21))
	for i, tx := range data.Txs {
		if i == 20 {
			txStrings[i] = fmt.Sprintf("... (%v total)", len(data.Txs))
			break
		}
		txStrings[i] = fmt.Sprintf("Tx:%v", tx)
	}
	return fmt.Sprintf(`Data{
%s  %v
%s}#%v`,
		indent, strings.Join(txStrings, "\n"+indent+"  "),
		indent, data.hash)
}

type EvidenceData struct {
	Evidence	EvidenceList	`json:"evidence"`

	hash	cmn.HexBytes
}

func (data *EvidenceData) Hash() cmn.HexBytes {
	if data.hash == nil {
		data.hash = data.Evidence.Hash()
	}
	return data.hash
}

func (data *EvidenceData) StringIndented(indent string) string {
	if data == nil {
		return "nil-Evidence"
	}
	evStrings := make([]string, cmn.MinInt(len(data.Evidence), 21))
	for i, ev := range data.Evidence {
		if i == 20 {
			evStrings[i] = fmt.Sprintf("... (%v total)", len(data.Evidence))
			break
		}
		evStrings[i] = fmt.Sprintf("Evidence:%v", ev)
	}
	return fmt.Sprintf(`Data{
%s  %v
%s}#%v`,
		indent, strings.Join(evStrings, "\n"+indent+"  "),
		indent, data.hash)
	return ""
}

type BlockID struct {
	Hash		cmn.HexBytes	`json:"hash"`
	PartsHeader	PartSetHeader	`json:"parts"`
}

func (blockID BlockID) IsZero() bool {
	return len(blockID.Hash) == 0 && blockID.PartsHeader.IsZero()
}

func (blockID BlockID) Equals(other BlockID) bool {
	return bytes.Equal(blockID.Hash, other.Hash) &&
		blockID.PartsHeader.Equals(other.PartsHeader)
}

func (blockID BlockID) Key() string {
	bz, err := cdc.MarshalBinaryBare(blockID.PartsHeader)
	if err != nil {
		panic(err)
	}
	return string(blockID.Hash) + string(bz)
}

func (blockID BlockID) String() string {
	return fmt.Sprintf(`%v:%v`, blockID.Hash, blockID.PartsHeader)
}

type hasher struct {
	item interface{}
}

func (h hasher) Hash() []byte {
	hasher := ripemd160.New()
	if h.item != nil && !cmn.IsTypedNil(h.item) && !cmn.IsEmpty(h.item) {
		bz, err := cdc.MarshalBinaryBare(h.item)
		if err != nil {
			panic(err)
		}
		_, err = hasher.Write(bz)
		if err != nil {
			panic(err)
		}
	}
	return hasher.Sum(nil)

}

func aminoHash(item interface{}) []byte {
	h := hasher{item}
	return h.Hash()
}

func aminoHasher(item interface{}) merkle.Hasher {
	return hasher{item}
}
