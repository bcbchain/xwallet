package types

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strings"

	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/merkle"
)

type ValidatorSet struct {
	Validators	[]*Validator	`json:"validators"`
	Proposer	*Validator	`json:"proposer"`

	totalVotingPower	int64
}

func NewValidatorSet(vals []*Validator) *ValidatorSet {
	validators := make([]*Validator, len(vals))
	for i, val := range vals {
		validators[i] = val.Copy()
	}
	sort.Sort(ValidatorsByAddress(validators))
	vs := &ValidatorSet{
		Validators: validators,
	}

	if vals != nil {
		vs.IncrementAccum(1)
	}

	return vs
}

func (valSet *ValidatorSet) IncrementAccum(times int) {

	validatorsHeap := cmn.NewHeap()
	for _, val := range valSet.Validators {

		val.Accum = safeAddClip(val.Accum, safeMulClip(int64(val.VotingPower), int64(times)))
		validatorsHeap.PushComparable(val, accumComparable{val})
	}

	for i := 0; i < times; i++ {
		mostest := validatorsHeap.Peek().(*Validator)
		if i == times-1 {
			valSet.Proposer = mostest
		}

		mostest.Accum = safeSubClip(mostest.Accum, valSet.TotalVotingPower())
		validatorsHeap.Update(mostest, accumComparable{mostest})
	}
}

func (valSet *ValidatorSet) Copy() *ValidatorSet {
	validators := make([]*Validator, len(valSet.Validators))
	for i, val := range valSet.Validators {

		validators[i] = val.Copy()
	}
	return &ValidatorSet{
		Validators:		validators,
		Proposer:		valSet.Proposer,
		totalVotingPower:	valSet.totalVotingPower,
	}
}

func (valSet *ValidatorSet) HasAddress(address string) bool {
	idx := sort.Search(len(valSet.Validators), func(i int) bool {
		return bytes.Compare([]byte(address), []byte(valSet.Validators[i].Address)) <= 0
	})
	return idx < len(valSet.Validators) && valSet.Validators[idx].Address == address
}

func (valSet *ValidatorSet) GetByAddress(address string) (index int, val *Validator) {
	idx := sort.Search(len(valSet.Validators), func(i int) bool {
		return bytes.Compare([]byte(address), []byte(valSet.Validators[i].Address)) <= 0
	})
	if idx < len(valSet.Validators) && valSet.Validators[idx].Address == address {
		return idx, valSet.Validators[idx].Copy()
	}
	return -1, nil
}

func (valSet *ValidatorSet) GetByIndex(index int) (address string, val *Validator) {
	if index < 0 || index >= len(valSet.Validators) {
		return "", nil
	}
	val = valSet.Validators[index]
	return val.Address, val.Copy()
}

func (valSet *ValidatorSet) Size() int {
	return len(valSet.Validators)
}

func (valSet *ValidatorSet) TotalVotingPower() int64 {
	if valSet.totalVotingPower == 0 {
		for _, val := range valSet.Validators {

			valSet.totalVotingPower = safeAddClip(valSet.totalVotingPower, int64(val.VotingPower))
		}
	}
	return valSet.totalVotingPower
}

func (valSet *ValidatorSet) GetProposer() (proposer *Validator) {
	if len(valSet.Validators) == 0 {
		return nil
	}
	if valSet.Proposer == nil {
		valSet.Proposer = valSet.findProposer()
	}
	return valSet.Proposer.Copy()
}

func (valSet *ValidatorSet) findProposer() *Validator {
	var proposer *Validator
	for _, val := range valSet.Validators {
		if proposer == nil || val.Address != proposer.Address {
			proposer = proposer.CompareAccum(val)
		}
	}
	return proposer
}

func (valSet *ValidatorSet) Hash() []byte {
	if len(valSet.Validators) == 0 {
		return nil
	}
	hashers := make([]merkle.Hasher, len(valSet.Validators))
	for i, val := range valSet.Validators {
		hashers[i] = val
	}
	return merkle.SimpleHashFromHashers(hashers)
}

func (valSet *ValidatorSet) Add(val *Validator) (added bool) {
	val = val.Copy()
	idx := sort.Search(len(valSet.Validators), func(i int) bool {
		return bytes.Compare([]byte(val.Address), []byte(valSet.Validators[i].Address)) <= 0
	})
	if idx >= len(valSet.Validators) {
		valSet.Validators = append(valSet.Validators, val)

		valSet.Proposer = nil
		valSet.totalVotingPower = 0
		return true
	} else if valSet.Validators[idx].Address == val.Address {
		return false
	} else {
		newValidators := make([]*Validator, len(valSet.Validators)+1)
		copy(newValidators[:idx], valSet.Validators[:idx])
		newValidators[idx] = val
		copy(newValidators[idx+1:], valSet.Validators[idx:])
		valSet.Validators = newValidators

		valSet.Proposer = nil
		valSet.totalVotingPower = 0
		return true
	}
}

func (valSet *ValidatorSet) Update(val *Validator) (updated bool) {
	index, sameVal := valSet.GetByAddress(val.Address)
	if sameVal == nil {
		return false
	}
	valSet.Validators[index] = val.Copy()

	valSet.Proposer = nil
	valSet.totalVotingPower = 0
	return true
}

func (valSet *ValidatorSet) Remove(address string) (val *Validator, removed bool) {
	idx := sort.Search(len(valSet.Validators), func(i int) bool {
		return bytes.Compare([]byte(address), []byte(valSet.Validators[i].Address)) <= 0
	})
	if idx >= len(valSet.Validators) || valSet.Validators[idx].Address != address {
		return nil, false
	}
	removedVal := valSet.Validators[idx]
	newValidators := valSet.Validators[:idx]
	if idx+1 < len(valSet.Validators) {
		newValidators = append(newValidators, valSet.Validators[idx+1:]...)
	}
	valSet.Validators = newValidators

	valSet.Proposer = nil
	valSet.totalVotingPower = 0
	return removedVal, true
}

func (valSet *ValidatorSet) Iterate(fn func(index int, val *Validator) bool) {
	for i, val := range valSet.Validators {
		stop := fn(i, val.Copy())
		if stop {
			break
		}
	}
}

func (valSet *ValidatorSet) VerifyCommit(chainID string, blockID BlockID, height int64, commit *Commit) error {
	if valSet.Size() != len(commit.Precommits) {
		return fmt.Errorf("Invalid commit -- wrong set size: %v vs %v", valSet.Size(), len(commit.Precommits))
	}
	if height != commit.Height() {
		return fmt.Errorf("Invalid commit -- wrong height: %v vs %v", height, commit.Height())
	}

	talliedVotingPower := int64(0)
	round := commit.Round()

	for idx, precommit := range commit.Precommits {

		if precommit == nil {
			continue
		}
		if precommit.Height != height {
			return fmt.Errorf("Invalid commit -- wrong height: %v vs %v", height, precommit.Height)
		}
		if precommit.Round != round {
			return fmt.Errorf("Invalid commit -- wrong round: %v vs %v", round, precommit.Round)
		}
		if precommit.Type != VoteTypePrecommit {
			return fmt.Errorf("Invalid commit -- not precommit @ index %v", idx)
		}
		_, val := valSet.GetByIndex(idx)

		precommitSignBytes := precommit.SignBytes(chainID)
		if !val.PubKey.VerifyBytes(precommitSignBytes, precommit.Signature) {
			return fmt.Errorf("Invalid commit -- invalid signature: %v", precommit)
		}
		if !blockID.Equals(precommit.BlockID) {
			continue
		}

		talliedVotingPower += int64(val.VotingPower)
	}

	if talliedVotingPower > valSet.TotalVotingPower()*2/3 {
		return nil
	}
	return fmt.Errorf("Invalid commit -- insufficient voting power: got %v, needed %v",
		talliedVotingPower, (valSet.TotalVotingPower()*2/3 + 1))
}

func (valSet *ValidatorSet) VerifyCommitAny(newSet *ValidatorSet, chainID string,
	blockID BlockID, height int64, commit *Commit) error {

	if newSet.Size() != len(commit.Precommits) {
		return cmn.NewError("Invalid commit -- wrong set size: %v vs %v", newSet.Size(), len(commit.Precommits))
	}
	if height != commit.Height() {
		return cmn.NewError("Invalid commit -- wrong height: %v vs %v", height, commit.Height())
	}

	oldVotingPower := int64(0)
	newVotingPower := int64(0)
	seen := map[int]bool{}
	round := commit.Round()

	for idx, precommit := range commit.Precommits {

		if precommit == nil {
			continue
		}
		if precommit.Height != height {

			return cmn.NewError("Blocks don't match - %d vs %d", round, precommit.Round)
		}
		if precommit.Round != round {
			return cmn.NewError("Invalid commit -- wrong round: %v vs %v", round, precommit.Round)
		}
		if precommit.Type != VoteTypePrecommit {
			return cmn.NewError("Invalid commit -- not precommit @ index %v", idx)
		}
		if !blockID.Equals(precommit.BlockID) {
			continue
		}

		vi, ov := valSet.GetByAddress(precommit.ValidatorAddress)
		if ov == nil || seen[vi] {
			continue
		}
		seen[vi] = true

		precommitSignBytes := precommit.SignBytes(chainID)
		if !ov.PubKey.VerifyBytes(precommitSignBytes, precommit.Signature) {
			return cmn.NewError("Invalid commit -- invalid signature: %v", precommit)
		}

		oldVotingPower += int64(ov.VotingPower)

		_, cv := newSet.GetByIndex(idx)
		if cv.PubKey.Equals(ov.PubKey) {

			newVotingPower += int64(cv.VotingPower)
		}
	}

	if oldVotingPower <= valSet.TotalVotingPower()*2/3 {
		return cmn.NewError("Invalid commit -- insufficient old voting power: got %v, needed %v",
			oldVotingPower, (valSet.TotalVotingPower()*2/3 + 1))
	} else if newVotingPower <= newSet.TotalVotingPower()*2/3 {
		return cmn.NewError("Invalid commit -- insufficient cur voting power: got %v, needed %v",
			newVotingPower, (newSet.TotalVotingPower()*2/3 + 1))
	}
	return nil
}

func (valSet *ValidatorSet) String() string {
	return valSet.StringIndented("")
}

func (valSet *ValidatorSet) StringIndented(indent string) string {
	if valSet == nil {
		return "nil-ValidatorSet"
	}
	valStrings := []string{}
	valSet.Iterate(func(index int, val *Validator) bool {
		valStrings = append(valStrings, val.String())
		return false
	})
	return fmt.Sprintf(`ValidatorSet{
%s  Proposer: %v
%s  Validators:
%s    %v
%s}`,
		indent, valSet.GetProposer().String(),
		indent,
		indent, strings.Join(valStrings, "\n"+indent+"  "),
		indent)

}

type ValidatorsByAddress []*Validator

func (vs ValidatorsByAddress) Len() int {
	return len(vs)
}

func (vs ValidatorsByAddress) Less(i, j int) bool {
	return bytes.Compare([]byte(vs[i].Address), []byte(vs[j].Address)) == -1
}

func (vs ValidatorsByAddress) Swap(i, j int) {
	it := vs[i]
	vs[i] = vs[j]
	vs[j] = it
}

type accumComparable struct {
	*Validator
}

func (ac accumComparable) Less(o interface{}) bool {
	other := o.(accumComparable).Validator
	larger := ac.CompareAccum(other)
	return bytes.Equal([]byte(larger.Address), []byte(ac.Address))
}

func RandValidatorSet(numValidators int, votingPower int64) (*ValidatorSet, []PrivValidator) {
	vals := make([]*Validator, numValidators)
	privValidators := make([]PrivValidator, numValidators)
	for i := 0; i < numValidators; i++ {
		val, privValidator := RandValidator(false, votingPower)
		vals[i] = val
		privValidators[i] = privValidator
	}
	valSet := NewValidatorSet(vals)
	sort.Sort(PrivValidatorsByAddress(privValidators))
	return valSet, privValidators
}

func safeMul(a int64, b int64) (int64, bool) {
	if a == 0 || b == 0 {
		return 0, false
	}
	if a == 1 {
		return b, false
	}
	if b == 1 {
		return a, false
	}
	if a == math.MinInt64 || b == math.MinInt64 {
		return -1, true
	}
	c := a * b
	return c, c/b != a
}

func safeAdd(a, b int64) (int64, bool) {
	if b > 0 && a > math.MaxInt64-b {
		return -1, true
	} else if b < 0 && a < math.MinInt64-b {
		return -1, true
	}
	return a + b, false
}

func safeSub(a, b int64) (int64, bool) {
	if b > 0 && a < math.MinInt64+b {
		return -1, true
	} else if b < 0 && a > math.MaxInt64+b {
		return -1, true
	}
	return a - b, false
}

func safeMulClip(a, b int64) int64 {
	c, overflow := safeMul(a, b)
	if overflow {
		if (a < 0 || b < 0) && !(a < 0 && b < 0) {
			return math.MinInt64
		}
		return math.MaxInt64
	}
	return c
}

func safeAddClip(a, b int64) int64 {
	c, overflow := safeAdd(a, b)
	if overflow {
		if b < 0 {
			return math.MinInt64
		}
		return math.MaxInt64
	}
	return c
}

func safeSubClip(a, b int64) int64 {
	c, overflow := safeSub(a, b)
	if overflow {
		if b > 0 {
			return math.MinInt64
		}
		return math.MaxInt64
	}
	return c
}
