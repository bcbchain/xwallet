package privval

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

const (
	stepNone	int8	= 0
	stepPropose	int8	= 1
	stepPrevote	int8	= 2
	stepPrecommit	int8	= 3
)

func voteToStep(vote *types.Vote) int8 {
	switch vote.Type {
	case types.VoteTypePrevote:
		return stepPrevote
	case types.VoteTypePrecommit:
		return stepPrecommit
	default:
		cmn.PanicSanity("Unknown vote type")
		return 0
	}
}

type FilePV struct {
	Address		crypto.Address		`json:"address"`
	PubKey		crypto.PubKey		`json:"pub_key"`
	LastHeight	int64			`json:"last_height"`
	LastRound	int			`json:"last_round"`
	LastStep	int8			`json:"last_step"`
	LastSignature	crypto.Signature	`json:"last_signature,omitempty"`
	LastSignBytes	cmn.HexBytes		`json:"last_signbytes,omitempty"`
	PrivKey		crypto.PrivKey		`json:"priv_key"`

	filePath	string
	mtx		sync.Mutex
}

func (pv *FilePV) GetAddress() crypto.Address {
	return pv.Address
}

func (pv *FilePV) GetPubKey() crypto.PubKey {
	return pv.PubKey
}

func GenFilePV(filePath string) *FilePV {
	privKey := crypto.GenPrivKeyEd25519()
	return &FilePV{
		Address:	privKey.PubKey().Address(),
		PubKey:		privKey.PubKey(),
		PrivKey:	privKey,
		LastStep:	stepNone,
		filePath:	filePath,
	}
}

func LoadFilePV(filePath string) *FilePV {
	pvJSONBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		cmn.Exit(err.Error())
	}
	pv := &FilePV{}
	err = cdc.UnmarshalJSON(pvJSONBytes, &pv)
	if err != nil {
		cmn.Exit(cmn.Fmt("Error reading PrivValidator from %v: %v\n", filePath, err))
	}

	pv.filePath = filePath
	return pv
}

func LoadOrGenFilePV(filePath string) *FilePV {
	var pv *FilePV
	if cmn.FileExists(filePath) {
		pv = LoadFilePV(filePath)
	} else {
		pv = GenFilePV(filePath)
		pv.Save()
	}
	return pv
}

func (pv *FilePV) Save() {
	pv.mtx.Lock()
	defer pv.mtx.Unlock()
	pv.save()
}

func (pv *FilePV) save() {
	outFile := pv.filePath
	if outFile == "" {
		panic("Cannot save PrivValidator: filePath not set")
	}
	jsonBytes, err := cdc.MarshalJSONIndent(pv, "", "  ")
	if err != nil {
		panic(err)
	}
	err = cmn.WriteFileAtomic(outFile, jsonBytes, 0600)
	if err != nil {
		panic(err)
	}
}

func (pv *FilePV) Reset() {
	var sig crypto.Signature
	pv.LastHeight = 0
	pv.LastRound = 0
	pv.LastStep = 0
	pv.LastSignature = sig
	pv.LastSignBytes = nil
	pv.Save()
}

func (pv *FilePV) SignVote(chainID string, vote *types.Vote) error {
	pv.mtx.Lock()
	defer pv.mtx.Unlock()
	if err := pv.signVote(chainID, vote); err != nil {
		return errors.New(cmn.Fmt("Error signing vote: %v", err))
	}
	return nil
}

func (pv *FilePV) SignProposal(chainID string, proposal *types.Proposal) error {
	pv.mtx.Lock()
	defer pv.mtx.Unlock()
	if err := pv.signProposal(chainID, proposal); err != nil {
		return fmt.Errorf("Error signing proposal: %v", err)
	}
	return nil
}

func (pv *FilePV) checkHRS(height int64, round int, step int8) (bool, error) {
	if pv.LastHeight > height {
		return false, errors.New("Height regression")
	}

	if pv.LastHeight == height {
		if pv.LastRound > round {
			return false, errors.New("Round regression")
		}

		if pv.LastRound == round {
			if pv.LastStep > step {
				return false, errors.New("Step regression")
			} else if pv.LastStep == step {
				if pv.LastSignBytes != nil {
					if pv.LastSignature == nil {
						panic("pv: LastSignature is nil but LastSignBytes is not!")
					}
					return true, nil
				}
				return false, errors.New("No LastSignature found")
			}
		}
	}
	return false, nil
}

func (pv *FilePV) signVote(chainID string, vote *types.Vote) error {
	height, round, step := vote.Height, vote.Round, voteToStep(vote)
	signBytes := vote.SignBytes(chainID)

	sameHRS, err := pv.checkHRS(height, round, step)
	if err != nil {
		return err
	}

	if sameHRS {
		if bytes.Equal(signBytes, pv.LastSignBytes) {
			vote.Signature = pv.LastSignature
		} else if timestamp, ok := checkVotesOnlyDifferByTimestamp(pv.LastSignBytes, signBytes); ok {
			vote.Timestamp = timestamp
			vote.Signature = pv.LastSignature
		} else {
			err = fmt.Errorf("Conflicting data")
		}
		return err
	}

	sig := pv.PrivKey.Sign(signBytes)
	pv.saveSigned(height, round, step, signBytes, sig)
	vote.Signature = sig
	return nil
}

func (pv *FilePV) signProposal(chainID string, proposal *types.Proposal) error {
	height, round, step := proposal.Height, proposal.Round, stepPropose
	signBytes := proposal.SignBytes(chainID)

	sameHRS, err := pv.checkHRS(height, round, step)
	if err != nil {
		return err
	}

	if sameHRS {
		if bytes.Equal(signBytes, pv.LastSignBytes) {
			proposal.Signature = pv.LastSignature
		} else if timestamp, ok := checkProposalsOnlyDifferByTimestamp(pv.LastSignBytes, signBytes); ok {
			proposal.Timestamp = timestamp
			proposal.Signature = pv.LastSignature
		} else {
			err = fmt.Errorf("Conflicting data")
		}
		return err
	}

	sig := pv.PrivKey.Sign(signBytes)
	pv.saveSigned(height, round, step, signBytes, sig)
	proposal.Signature = sig
	return nil
}

func (pv *FilePV) saveSigned(height int64, round int, step int8,
	signBytes []byte, sig crypto.Signature) {

	pv.LastHeight = height
	pv.LastRound = round
	pv.LastStep = step
	pv.LastSignature = sig
	pv.LastSignBytes = signBytes
	pv.save()
}

func (pv *FilePV) SignHeartbeat(chainID string, heartbeat *types.Heartbeat) error {
	pv.mtx.Lock()
	defer pv.mtx.Unlock()
	heartbeat.Signature = pv.PrivKey.Sign(heartbeat.SignBytes(chainID))
	return nil
}

func (pv *FilePV) String() string {
	return fmt.Sprintf("PrivValidator{%v LH:%v, LR:%v, LS:%v}", pv.GetAddress(), pv.LastHeight, pv.LastRound, pv.LastStep)
}

func checkVotesOnlyDifferByTimestamp(lastSignBytes, newSignBytes []byte) (time.Time, bool) {
	var lastVote, newVote types.CanonicalJSONVote
	if err := cdc.UnmarshalJSON(lastSignBytes, &lastVote); err != nil {
		panic(fmt.Sprintf("LastSignBytes cannot be unmarshalled into vote: %v", err))
	}
	if err := cdc.UnmarshalJSON(newSignBytes, &newVote); err != nil {
		panic(fmt.Sprintf("signBytes cannot be unmarshalled into vote: %v", err))
	}

	lastTime, err := time.Parse(types.TimeFormat, lastVote.Timestamp)
	if err != nil {
		panic(err)
	}

	now := types.CanonicalTime(time.Now())
	lastVote.Timestamp = now
	newVote.Timestamp = now
	lastVoteBytes, _ := cdc.MarshalJSON(lastVote)
	newVoteBytes, _ := cdc.MarshalJSON(newVote)

	return lastTime, bytes.Equal(newVoteBytes, lastVoteBytes)
}

func checkProposalsOnlyDifferByTimestamp(lastSignBytes, newSignBytes []byte) (time.Time, bool) {
	var lastProposal, newProposal types.CanonicalJSONProposal
	if err := cdc.UnmarshalJSON(lastSignBytes, &lastProposal); err != nil {
		panic(fmt.Sprintf("LastSignBytes cannot be unmarshalled into proposal: %v", err))
	}
	if err := cdc.UnmarshalJSON(newSignBytes, &newProposal); err != nil {
		panic(fmt.Sprintf("signBytes cannot be unmarshalled into proposal: %v", err))
	}

	lastTime, err := time.Parse(types.TimeFormat, lastProposal.Timestamp)
	if err != nil {
		panic(err)
	}

	now := types.CanonicalTime(time.Now())
	lastProposal.Timestamp = now
	newProposal.Timestamp = now
	lastProposalBytes, _ := cdc.MarshalJSON(lastProposal)
	newProposalBytes, _ := cdc.MarshalJSON(newProposal)

	return lastTime, bytes.Equal(newProposalBytes, lastProposalBytes)
}
