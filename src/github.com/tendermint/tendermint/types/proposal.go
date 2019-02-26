package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/tendermint/go-crypto"
)

var (
	ErrInvalidBlockPartSignature	= errors.New("Error invalid block part signature")
	ErrInvalidBlockPartHash		= errors.New("Error invalid block part hash")
)

type Proposal struct {
	Height			int64			`json:"height"`
	Round			int			`json:"round"`
	Timestamp		time.Time		`json:"timestamp"`
	BlockPartsHeader	PartSetHeader		`json:"block_parts_header"`
	POLRound		int			`json:"pol_round"`
	POLBlockID		BlockID			`json:"pol_block_id"`
	Signature		crypto.Signature	`json:"signature"`
}

func NewProposal(height int64, round int, blockPartsHeader PartSetHeader, polRound int, polBlockID BlockID) *Proposal {
	return &Proposal{
		Height:			height,
		Round:			round,
		Timestamp:		time.Now().UTC(),
		BlockPartsHeader:	blockPartsHeader,
		POLRound:		polRound,
		POLBlockID:		polBlockID,
	}
}

func (p *Proposal) String() string {
	return fmt.Sprintf("Proposal{%v/%v %v (%v,%v) %v @ %s}",
		p.Height, p.Round, p.BlockPartsHeader, p.POLRound,
		p.POLBlockID, p.Signature, CanonicalTime(p.Timestamp))
}

func (p *Proposal) SignBytes(chainID string) []byte {
	bz, err := cdc.MarshalJSON(CanonicalProposal(chainID, p))
	if err != nil {
		panic(err)
	}
	return bz
}
