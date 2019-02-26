package types

import (
	"bytes"
	"fmt"

	"github.com/tendermint/go-crypto"
)

type PrivValidator interface {
	GetAddress() crypto.Address
	GetPubKey() crypto.PubKey

	SignVote(chainID string, vote *Vote) error
	SignProposal(chainID string, proposal *Proposal) error
	SignHeartbeat(chainID string, heartbeat *Heartbeat) error
}

type PrivValidatorsByAddress []PrivValidator

func (pvs PrivValidatorsByAddress) Len() int {
	return len(pvs)
}

func (pvs PrivValidatorsByAddress) Less(i, j int) bool {
	return bytes.Compare([]byte(pvs[i].GetAddress()), []byte(pvs[j].GetAddress())) == -1
}

func (pvs PrivValidatorsByAddress) Swap(i, j int) {
	it := pvs[i]
	pvs[i] = pvs[j]
	pvs[j] = it
}

type MockPV struct {
	privKey crypto.PrivKey
}

func NewMockPV() *MockPV {
	return &MockPV{crypto.GenPrivKeyEd25519()}
}

func (pv *MockPV) GetAddress() crypto.Address {
	return pv.privKey.PubKey().Address()
}

func (pv *MockPV) GetPubKey() crypto.PubKey {
	return pv.privKey.PubKey()
}

func (pv *MockPV) SignVote(chainID string, vote *Vote) error {
	signBytes := vote.SignBytes(chainID)
	sig := pv.privKey.Sign(signBytes)
	vote.Signature = sig
	return nil
}

func (pv *MockPV) SignProposal(chainID string, proposal *Proposal) error {
	signBytes := proposal.SignBytes(chainID)
	sig := pv.privKey.Sign(signBytes)
	proposal.Signature = sig
	return nil
}

func (pv *MockPV) SignHeartbeat(chainID string, heartbeat *Heartbeat) error {
	sig := pv.privKey.Sign(heartbeat.SignBytes(chainID))
	heartbeat.Signature = sig
	return nil
}

func (pv *MockPV) String() string {
	return fmt.Sprintf("MockPV{%v}", pv.GetAddress())
}

func (pv *MockPV) DisableChecks() {

}
