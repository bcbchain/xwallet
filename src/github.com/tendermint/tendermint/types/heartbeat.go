package types

import (
	"fmt"

	"github.com/tendermint/go-crypto"
	cmn "github.com/tendermint/tmlibs/common"
)

type Heartbeat struct {
	ValidatorAddress	crypto.Address		`json:"validator_address"`
	ValidatorIndex		int			`json:"validator_index"`
	Height			int64			`json:"height"`
	Round			int			`json:"round"`
	Sequence		int			`json:"sequence"`
	Signature		crypto.Signature	`json:"signature"`
}

func (heartbeat *Heartbeat) SignBytes(chainID string) []byte {
	bz, err := cdc.MarshalJSON(CanonicalHeartbeat(chainID, heartbeat))
	if err != nil {
		panic(err)
	}
	return bz
}

func (heartbeat *Heartbeat) Copy() *Heartbeat {
	if heartbeat == nil {
		return nil
	}
	heartbeatCopy := *heartbeat
	return &heartbeatCopy
}

func (heartbeat *Heartbeat) String() string {
	if heartbeat == nil {
		return "nil-heartbeat"
	}

	addr := heartbeat.ValidatorAddress[14:]
	return fmt.Sprintf("Heartbeat{%v:%X %v/%02d (%v) %v}",
		heartbeat.ValidatorIndex, cmn.Fingerprint([]byte(addr)),
		heartbeat.Height, heartbeat.Round, heartbeat.Sequence, heartbeat.Signature)
}
