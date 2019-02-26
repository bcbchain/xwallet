package types

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"

	cmn "github.com/tendermint/tmlibs/common"
)

type P2PID string

type VoteSet struct {
	chainID	string
	height	int64
	round	int
	type_	byte
	valSet	*ValidatorSet

	mtx		sync.Mutex
	votesBitArray	*cmn.BitArray
	votes		[]*Vote
	sum		int64
	maj23		*BlockID
	votesByBlock	map[string]*blockVotes
	peerMaj23s	map[P2PID]BlockID
}

func NewVoteSet(chainID string, height int64, round int, type_ byte, valSet *ValidatorSet) *VoteSet {
	if height == 0 {
		cmn.PanicSanity("Cannot make VoteSet for height == 0, doesn't make sense.")
	}
	return &VoteSet{
		chainID:	chainID,
		height:		height,
		round:		round,
		type_:		type_,
		valSet:		valSet,
		votesBitArray:	cmn.NewBitArray(valSet.Size()),
		votes:		make([]*Vote, valSet.Size()),
		sum:		0,
		maj23:		nil,
		votesByBlock:	make(map[string]*blockVotes, valSet.Size()),
		peerMaj23s:	make(map[P2PID]BlockID),
	}
}

func (voteSet *VoteSet) ChainID() string {
	return voteSet.chainID
}

func (voteSet *VoteSet) Height() int64 {
	if voteSet == nil {
		return 0
	}
	return voteSet.height
}

func (voteSet *VoteSet) Round() int {
	if voteSet == nil {
		return -1
	}
	return voteSet.round
}

func (voteSet *VoteSet) Type() byte {
	if voteSet == nil {
		return 0x00
	}
	return voteSet.type_
}

func (voteSet *VoteSet) Size() int {
	if voteSet == nil {
		return 0
	}
	return voteSet.valSet.Size()
}

func (voteSet *VoteSet) AddVote(vote *Vote) (added bool, err error) {
	if voteSet == nil {
		cmn.PanicSanity("AddVote() on nil VoteSet")
	}
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()

	return voteSet.addVote(vote)
}

func (voteSet *VoteSet) addVote(vote *Vote) (added bool, err error) {
	if vote == nil {
		return false, ErrVoteNil
	}
	valIndex := vote.ValidatorIndex
	valAddr := vote.ValidatorAddress
	blockKey := vote.BlockID.Key()

	if valIndex < 0 {
		return false, errors.Wrap(ErrVoteInvalidValidatorIndex, "Index < 0")
	} else if len(valAddr) == 0 {
		return false, errors.Wrap(ErrVoteInvalidValidatorAddress, "Empty address")
	}

	if (vote.Height != voteSet.height) ||
		(vote.Round != voteSet.round) ||
		(vote.Type != voteSet.type_) {
		return false, errors.Wrapf(ErrVoteUnexpectedStep, "Got %d/%d/%d, expected %d/%d/%d",
			voteSet.height, voteSet.round, voteSet.type_,
			vote.Height, vote.Round, vote.Type)
	}

	lookupAddr, val := voteSet.valSet.GetByIndex(valIndex)
	if val == nil {
		return false, errors.Wrapf(ErrVoteInvalidValidatorIndex,
			"Cannot find validator %d in valSet of size %d", valIndex, voteSet.valSet.Size())
	}

	if valAddr != lookupAddr {
		return false, errors.Wrapf(ErrVoteInvalidValidatorAddress,
			"vote.ValidatorAddress (%X) does not match address (%X) for vote.ValidatorIndex (%d)\nEnsure the genesis file is correct across all validators.",
			valAddr, lookupAddr, valIndex)
	}

	if existing, ok := voteSet.getVote(valIndex, blockKey); ok {
		if existing.Signature.Equals(vote.Signature) {
			return false, nil
		}
		return false, errors.Wrapf(ErrVoteNonDeterministicSignature, "Existing vote: %v; New vote: %v", existing, vote)
	}

	if err := vote.Verify(voteSet.chainID, val.PubKey); err != nil {
		return false, errors.Wrapf(err, "Failed to verify vote with ChainID %s and PubKey %s", voteSet.chainID, val.PubKey)
	}

	added, conflicting := voteSet.addVerifiedVote(vote, blockKey, int64(val.VotingPower))
	if conflicting != nil {
		return added, NewConflictingVoteError(val, conflicting, vote)
	}
	if !added {
		cmn.PanicSanity("Expected to add non-conflicting vote")
	}
	return added, nil
}

func (voteSet *VoteSet) getVote(valIndex int, blockKey string) (vote *Vote, ok bool) {
	if existing := voteSet.votes[valIndex]; existing != nil && existing.BlockID.Key() == blockKey {
		return existing, true
	}
	if existing := voteSet.votesByBlock[blockKey].getByIndex(valIndex); existing != nil {
		return existing, true
	}
	return nil, false
}

func (voteSet *VoteSet) addVerifiedVote(vote *Vote, blockKey string, votingPower int64) (added bool, conflicting *Vote) {
	valIndex := vote.ValidatorIndex

	if existing := voteSet.votes[valIndex]; existing != nil {
		if existing.BlockID.Equals(vote.BlockID) {
			cmn.PanicSanity("addVerifiedVote does not expect duplicate votes")
		} else {
			conflicting = existing
		}

		if voteSet.maj23 != nil && voteSet.maj23.Key() == blockKey {
			voteSet.votes[valIndex] = vote
			voteSet.votesBitArray.SetIndex(valIndex, true)
		}

	} else {

		voteSet.votes[valIndex] = vote
		voteSet.votesBitArray.SetIndex(valIndex, true)
		voteSet.sum += votingPower
	}

	votesByBlock, ok := voteSet.votesByBlock[blockKey]
	if ok {
		if conflicting != nil && !votesByBlock.peerMaj23 {

			return false, conflicting
		}

	} else {

		if conflicting != nil {

			return false, conflicting
		}

		votesByBlock = newBlockVotes(false, voteSet.valSet.Size())
		voteSet.votesByBlock[blockKey] = votesByBlock

	}

	origSum := votesByBlock.sum
	quorum := voteSet.valSet.TotalVotingPower()*2/3 + 1

	votesByBlock.addVerifiedVote(vote, votingPower)

	if origSum < quorum && quorum <= votesByBlock.sum {

		if voteSet.maj23 == nil {
			maj23BlockID := vote.BlockID
			voteSet.maj23 = &maj23BlockID

			for i, vote := range votesByBlock.votes {
				if vote != nil {
					voteSet.votes[i] = vote
				}
			}
		}
	}

	return true, conflicting
}

func (voteSet *VoteSet) SetPeerMaj23(peerID P2PID, blockID BlockID) error {
	if voteSet == nil {
		cmn.PanicSanity("SetPeerMaj23() on nil VoteSet")
	}
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()

	blockKey := blockID.Key()

	if existing, ok := voteSet.peerMaj23s[peerID]; ok {
		if existing.Equals(blockID) {
			return nil
		}
		return fmt.Errorf("SetPeerMaj23: Received conflicting blockID from peer %v. Got %v, expected %v",
			peerID, blockID, existing)
	}
	voteSet.peerMaj23s[peerID] = blockID

	votesByBlock, ok := voteSet.votesByBlock[blockKey]
	if ok {
		if votesByBlock.peerMaj23 {
			return nil
		}
		votesByBlock.peerMaj23 = true

	} else {
		votesByBlock = newBlockVotes(true, voteSet.valSet.Size())
		voteSet.votesByBlock[blockKey] = votesByBlock

	}
	return nil
}

func (voteSet *VoteSet) BitArray() *cmn.BitArray {
	if voteSet == nil {
		return nil
	}
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()
	return voteSet.votesBitArray.Copy()
}

func (voteSet *VoteSet) BitArrayByBlockID(blockID BlockID) *cmn.BitArray {
	if voteSet == nil {
		return nil
	}
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()
	votesByBlock, ok := voteSet.votesByBlock[blockID.Key()]
	if ok {
		return votesByBlock.bitArray.Copy()
	}
	return nil
}

func (voteSet *VoteSet) GetByIndex(valIndex int) *Vote {
	if voteSet == nil {
		return nil
	}
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()
	return voteSet.votes[valIndex]
}

func (voteSet *VoteSet) GetByAddress(address string) *Vote {
	if voteSet == nil {
		return nil
	}
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()
	valIndex, val := voteSet.valSet.GetByAddress(address)
	if val == nil {
		cmn.PanicSanity("GetByAddress(address) returned nil")
	}
	return voteSet.votes[valIndex]
}

func (voteSet *VoteSet) HasTwoThirdsMajority() bool {
	if voteSet == nil {
		return false
	}
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()
	return voteSet.maj23 != nil
}

func (voteSet *VoteSet) IsCommit() bool {
	if voteSet == nil {
		return false
	}
	if voteSet.type_ != VoteTypePrecommit {
		return false
	}
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()
	return voteSet.maj23 != nil
}

func (voteSet *VoteSet) HasTwoThirdsAny() bool {
	if voteSet == nil {
		return false
	}
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()
	return voteSet.sum > voteSet.valSet.TotalVotingPower()*2/3
}

func (voteSet *VoteSet) HasAll() bool {
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()
	return voteSet.sum == voteSet.valSet.TotalVotingPower()
}

func (voteSet *VoteSet) TwoThirdsMajority() (blockID BlockID, ok bool) {
	if voteSet == nil {
		return BlockID{}, false
	}
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()
	if voteSet.maj23 != nil {
		return *voteSet.maj23, true
	}
	return BlockID{}, false
}

func (voteSet *VoteSet) String() string {
	if voteSet == nil {
		return "nil-VoteSet"
	}
	return voteSet.StringIndented("")
}

func (voteSet *VoteSet) StringIndented(indent string) string {
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()
	voteStrings := make([]string, len(voteSet.votes))
	for i, vote := range voteSet.votes {
		if vote == nil {
			voteStrings[i] = "nil-Vote"
		} else {
			voteStrings[i] = vote.String()
		}
	}
	return fmt.Sprintf(`VoteSet{
%s  H:%v R:%v T:%v
%s  %v
%s  %v
%s  %v
%s}`,
		indent, voteSet.height, voteSet.round, voteSet.type_,
		indent, strings.Join(voteStrings, "\n"+indent+"  "),
		indent, voteSet.votesBitArray,
		indent, voteSet.peerMaj23s,
		indent)
}

func (voteSet *VoteSet) MarshalJSON() ([]byte, error) {
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()
	voteStrings := make([]string, len(voteSet.votes))
	for i, vote := range voteSet.votes {
		if vote == nil {
			voteStrings[i] = "nil-Vote"
		} else {
			voteStrings[i] = vote.String()
		}
	}
	return cdc.MarshalJSON(struct {
		Votes		[]string		`json:"votes"`
		VotesBitArray	*cmn.BitArray		`json:"votes_bit_array"`
		PeerMaj23s	map[P2PID]BlockID	`json:"peer_maj_23s"`
	}{
		voteStrings, voteSet.votesBitArray, voteSet.peerMaj23s,
	})
}

func (voteSet *VoteSet) StringShort() string {
	if voteSet == nil {
		return "nil-VoteSet"
	}
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()
	return fmt.Sprintf(`VoteSet{H:%v R:%v T:%v +2/3:%v %v %v}`,
		voteSet.height, voteSet.round, voteSet.type_, voteSet.maj23, voteSet.votesBitArray, voteSet.peerMaj23s)
}

func (voteSet *VoteSet) MakeCommit() *Commit {
	if voteSet.type_ != VoteTypePrecommit {
		cmn.PanicSanity("Cannot MakeCommit() unless VoteSet.Type is VoteTypePrecommit")
	}
	voteSet.mtx.Lock()
	defer voteSet.mtx.Unlock()

	if voteSet.maj23 == nil {
		cmn.PanicSanity("Cannot MakeCommit() unless a blockhash has +2/3")
	}

	votesCopy := make([]*Vote, len(voteSet.votes))
	copy(votesCopy, voteSet.votes)
	return &Commit{
		BlockID:	*voteSet.maj23,
		Precommits:	votesCopy,
	}
}

type blockVotes struct {
	peerMaj23	bool
	bitArray	*cmn.BitArray
	votes		[]*Vote
	sum		int64
}

func newBlockVotes(peerMaj23 bool, numValidators int) *blockVotes {
	return &blockVotes{
		peerMaj23:	peerMaj23,
		bitArray:	cmn.NewBitArray(numValidators),
		votes:		make([]*Vote, numValidators),
		sum:		0,
	}
}

func (vs *blockVotes) addVerifiedVote(vote *Vote, votingPower int64) {
	valIndex := vote.ValidatorIndex
	if existing := vs.votes[valIndex]; existing == nil {
		vs.bitArray.SetIndex(valIndex, true)
		vs.votes[valIndex] = vote
		vs.sum += votingPower
	}
}

func (vs *blockVotes) getByIndex(index int) *Vote {
	if vs == nil {
		return nil
	}
	return vs.votes[index]
}

type VoteSetReader interface {
	Height() int64
	Round() int
	Type() byte
	Size() int
	BitArray() *cmn.BitArray
	GetByIndex(int) *Vote
	IsCommit() bool
}
