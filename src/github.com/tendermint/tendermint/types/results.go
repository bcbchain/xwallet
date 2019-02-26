package types

import (
	abci "github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/merkle"
)

type ABCIResult struct {
	Code	uint32		`json:"code"`
	Data	cmn.HexBytes	`json:"data"`
}

func (a ABCIResult) Hash() []byte {
	bz := aminoHash(a)
	return bz
}

type ABCIResults []ABCIResult

func NewResults(del []*abci.ResponseDeliverTx) ABCIResults {
	res := make(ABCIResults, len(del))
	for i, d := range del {
		res[i] = NewResultFromResponse(d)
	}
	return res
}

func NewResultFromResponse(response *abci.ResponseDeliverTx) ABCIResult {

	return ABCIResult{
		Code:	response.Code,
		Data:	[]byte(response.Data),
	}
}

func (a ABCIResults) Bytes() []byte {
	bz, err := cdc.MarshalBinary(a)
	if err != nil {
		panic(err)
	}
	return bz
}

func (a ABCIResults) Hash() []byte {
	return merkle.SimpleHashFromHashers(a.toHashers())
}

func (a ABCIResults) ProveResult(i int) merkle.SimpleProof {
	_, proofs := merkle.SimpleProofsFromHashers(a.toHashers())
	return *proofs[i]
}

func (a ABCIResults) toHashers() []merkle.Hasher {
	l := len(a)
	hashers := make([]merkle.Hasher, l)
	for i := 0; i < l; i++ {
		hashers[i] = a[i]
	}
	return hashers
}
