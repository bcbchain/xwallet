package types

import (
	"bytes"
	"errors"
	"fmt"

	abci "github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/merkle"
)

type Tx []byte

func (tx Tx) Hash() []byte {
	return aminoHasher(tx).Hash()
}

func (tx Tx) String() string {
	return fmt.Sprintf("Tx{%X}", []byte(tx))
}

type Txs []Tx

type HashList [][]byte

func (txs Txs) Hash() []byte {

	switch len(txs) {
	case 0:
		return nil
	case 1:
		return txs[0].Hash()
	default:
		left := Txs(txs[:(len(txs)+1)/2]).Hash()
		right := Txs(txs[(len(txs)+1)/2:]).Hash()
		return merkle.SimpleHashFromTwoHashes(left, right)
	}
}

func (txs Txs) Index(tx Tx) int {
	for i := range txs {
		if bytes.Equal(txs[i], tx) {
			return i
		}
	}
	return -1
}

func (txs Txs) IndexByHash(hash []byte) int {
	for i := range txs {
		if bytes.Equal(txs[i].Hash(), hash) {
			return i
		}
	}
	return -1
}

func (txs Txs) Proof(i int) TxProof {
	l := len(txs)
	hashers := make([]merkle.Hasher, l)
	for i := 0; i < l; i++ {
		hashers[i] = txs[i]
	}
	root, proofs := merkle.SimpleProofsFromHashers(hashers)

	return TxProof{
		Index:		i,
		Total:		l,
		RootHash:	root,
		Data:		txs[i],
		Proof:		*proofs[i],
	}
}

type TxProof struct {
	Index, Total	int
	RootHash	cmn.HexBytes
	Data		Tx
	Proof		merkle.SimpleProof
}

func (tp TxProof) LeafHash() []byte {
	return tp.Data.Hash()
}

func (tp TxProof) Validate(dataHash []byte) error {
	if !bytes.Equal(dataHash, tp.RootHash) {
		return errors.New("Proof matches different data hash")
	}
	if tp.Index < 0 {
		return errors.New("Proof index cannot be negative")
	}
	if tp.Total <= 0 {
		return errors.New("Proof total must be positive")
	}
	valid := tp.Proof.Verify(tp.Index, tp.Total, tp.LeafHash(), tp.RootHash)
	if !valid {
		return errors.New("Proof is not internally consistent")
	}
	return nil
}

type TxResult struct {
	Height	int64			`json:"height"`
	Index	uint32			`json:"index"`
	Tx	Tx			`json:"tx"`
	Result	abci.ResponseDeliverTx	`json:"result"`
}
