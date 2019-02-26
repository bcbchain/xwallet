package types

import (
	"bytes"
	"fmt"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tmlibs/merkle"
)

type ErrEvidenceInvalid struct {
	Evidence	Evidence
	ErrorValue	error
}

func NewEvidenceInvalidErr(ev Evidence, err error) *ErrEvidenceInvalid {
	return &ErrEvidenceInvalid{ev, err}
}

func (err *ErrEvidenceInvalid) Error() string {
	return fmt.Sprintf("Invalid evidence: %v. Evidence: %v", err.ErrorValue, err.Evidence)
}

type Evidence interface {
	Height() int64
	Address() string
	Index() int
	Hash() []byte
	Verify(chainID string) error
	Equal(Evidence) bool

	String() string
}

func RegisterEvidences(cdc *amino.Codec) {
	cdc.RegisterInterface((*Evidence)(nil), nil)
	cdc.RegisterConcrete(&DuplicateVoteEvidence{}, "tendermint/DuplicateVoteEvidence", nil)
}

type DuplicateVoteEvidence struct {
	PubKey	crypto.PubKey
	VoteA	*Vote
	VoteB	*Vote
}

func (dve *DuplicateVoteEvidence) String() string {
	return fmt.Sprintf("VoteA: %v; VoteB: %v", dve.VoteA, dve.VoteB)

}

func (dve *DuplicateVoteEvidence) Height() int64 {
	return dve.VoteA.Height
}

func (dve *DuplicateVoteEvidence) Address() string {
	return dve.PubKey.Address()
}

func (dve *DuplicateVoteEvidence) Index() int {
	return dve.VoteA.ValidatorIndex
}

func (dve *DuplicateVoteEvidence) Hash() []byte {
	return aminoHasher(dve).Hash()
}

func (dve *DuplicateVoteEvidence) Verify(chainID string) error {

	if dve.VoteA.Height != dve.VoteB.Height ||
		dve.VoteA.Round != dve.VoteB.Round ||
		dve.VoteA.Type != dve.VoteB.Type {
		return fmt.Errorf("DuplicateVoteEvidence Error: H/R/S does not match. Got %v and %v", dve.VoteA, dve.VoteB)
	}

	if dve.VoteA.ValidatorAddress == dve.VoteB.ValidatorAddress {
		return fmt.Errorf("DuplicateVoteEvidence Error: Validator addresses do not match. Got %X and %X", dve.VoteA.ValidatorAddress, dve.VoteB.ValidatorAddress)
	}

	if dve.VoteA.ValidatorIndex != dve.VoteB.ValidatorIndex {
		return fmt.Errorf("DuplicateVoteEvidence Error: Validator indices do not match. Got %d and %d", dve.VoteA.ValidatorIndex, dve.VoteB.ValidatorIndex)
	}

	if dve.VoteA.BlockID.Equals(dve.VoteB.BlockID) {
		return fmt.Errorf("DuplicateVoteEvidence Error: BlockIDs are the same (%v) - not a real duplicate vote", dve.VoteA.BlockID)
	}

	if !dve.PubKey.VerifyBytes(dve.VoteA.SignBytes(chainID), dve.VoteA.Signature) {
		return fmt.Errorf("DuplicateVoteEvidence Error verifying VoteA: %v", ErrVoteInvalidSignature)
	}
	if !dve.PubKey.VerifyBytes(dve.VoteB.SignBytes(chainID), dve.VoteB.Signature) {
		return fmt.Errorf("DuplicateVoteEvidence Error verifying VoteB: %v", ErrVoteInvalidSignature)
	}

	return nil
}

func (dve *DuplicateVoteEvidence) Equal(ev Evidence) bool {
	if _, ok := ev.(*DuplicateVoteEvidence); !ok {
		return false
	}

	dveHash := aminoHasher(dve).Hash()
	evHash := aminoHasher(ev).Hash()
	return bytes.Equal(dveHash, evHash)
}

type MockGoodEvidence struct {
	Height_		int64
	Address_	string
	Index_		int
}

func NewMockGoodEvidence(height int64, index int, address string) MockGoodEvidence {
	return MockGoodEvidence{height, address, index}
}

func (e MockGoodEvidence) Height() int64	{ return e.Height_ }
func (e MockGoodEvidence) Address() string	{ return e.Address_ }
func (e MockGoodEvidence) Index() int		{ return e.Index_ }
func (e MockGoodEvidence) Hash() []byte {
	return []byte(fmt.Sprintf("%d-%d", e.Height_, e.Index_))
}
func (e MockGoodEvidence) Verify(chainID string) error	{ return nil }
func (e MockGoodEvidence) Equal(ev Evidence) bool {
	e2 := ev.(MockGoodEvidence)
	return e.Height_ == e2.Height_ &&
		e.Address_ == e2.Address_ &&
		e.Index_ == e2.Index_
}
func (e MockGoodEvidence) String() string {
	return fmt.Sprintf("GoodEvidence: %d/%s/%d", e.Height_, e.Address_, e.Index_)
}

type MockBadEvidence struct {
	MockGoodEvidence
}

func (e MockBadEvidence) Verify(chainID string) error	{ return fmt.Errorf("MockBadEvidence") }
func (e MockBadEvidence) Equal(ev Evidence) bool {
	e2 := ev.(MockBadEvidence)
	return e.Height_ == e2.Height_ &&
		e.Address_ == e2.Address_ &&
		e.Index_ == e2.Index_
}
func (e MockBadEvidence) String() string {
	return fmt.Sprintf("BadEvidence: %d/%s/%d", e.Height_, e.Address_, e.Index_)
}

type EvidenceList []Evidence

func (evl EvidenceList) Hash() []byte {

	switch len(evl) {
	case 0:
		return nil
	case 1:
		return evl[0].Hash()
	default:
		left := EvidenceList(evl[:(len(evl)+1)/2]).Hash()
		right := EvidenceList(evl[(len(evl)+1)/2:]).Hash()
		return merkle.SimpleHashFromTwoHashes(left, right)
	}
}

func (evl EvidenceList) String() string {
	s := ""
	for _, e := range evl {
		s += fmt.Sprintf("%s\t\t", e)
	}
	return s
}

func (evl EvidenceList) Has(evidence Evidence) bool {
	for _, ev := range evl {
		if ev.Equal(evidence) {
			return true
		}
	}
	return false
}
