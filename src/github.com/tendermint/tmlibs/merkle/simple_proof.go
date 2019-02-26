package merkle

import (
	"bytes"
	"fmt"
)

type SimpleProof struct {
	Aunts [][]byte `json:"aunts"`
}

func SimpleProofsFromHashers(items []Hasher) (rootHash []byte, proofs []*SimpleProof) {
	trails, rootSPN := trailsFromHashers(items)
	rootHash = rootSPN.Hash
	proofs = make([]*SimpleProof, len(items))
	for i, trail := range trails {
		proofs[i] = &SimpleProof{
			Aunts: trail.FlattenAunts(),
		}
	}
	return
}

func SimpleProofsFromMap(m map[string]Hasher) (rootHash []byte, proofs []*SimpleProof) {
	sm := NewSimpleMap()
	for k, v := range m {
		sm.Set(k, v)
	}
	sm.Sort()
	kvs := sm.kvs
	kvsH := make([]Hasher, 0, len(kvs))
	for _, kvp := range kvs {
		kvsH = append(kvsH, KVPair(kvp))
	}
	return SimpleProofsFromHashers(kvsH)
}

func (sp *SimpleProof) Verify(index int, total int, leafHash []byte, rootHash []byte) bool {
	computedHash := computeHashFromAunts(index, total, leafHash, sp.Aunts)
	return computedHash != nil && bytes.Equal(computedHash, rootHash)
}

func (sp *SimpleProof) String() string {
	return sp.StringIndented("")
}

func (sp *SimpleProof) StringIndented(indent string) string {
	return fmt.Sprintf(`SimpleProof{
%s  Aunts: %X
%s}`,
		indent, sp.Aunts,
		indent)
}

func computeHashFromAunts(index int, total int, leafHash []byte, innerHashes [][]byte) []byte {
	if index >= total || index < 0 || total <= 0 {
		return nil
	}
	switch total {
	case 0:
		panic("Cannot call computeHashFromAunts() with 0 total")
	case 1:
		if len(innerHashes) != 0 {
			return nil
		}
		return leafHash
	default:
		if len(innerHashes) == 0 {
			return nil
		}
		numLeft := (total + 1) / 2
		if index < numLeft {
			leftHash := computeHashFromAunts(index, numLeft, leafHash, innerHashes[:len(innerHashes)-1])
			if leftHash == nil {
				return nil
			}
			return SimpleHashFromTwoHashes(leftHash, innerHashes[len(innerHashes)-1])
		}
		rightHash := computeHashFromAunts(index-numLeft, total-numLeft, leafHash, innerHashes[:len(innerHashes)-1])
		if rightHash == nil {
			return nil
		}
		return SimpleHashFromTwoHashes(innerHashes[len(innerHashes)-1], rightHash)
	}
}

type SimpleProofNode struct {
	Hash	[]byte
	Parent	*SimpleProofNode
	Left	*SimpleProofNode
	Right	*SimpleProofNode
}

func (spn *SimpleProofNode) FlattenAunts() [][]byte {

	innerHashes := [][]byte{}
	for spn != nil {
		if spn.Left != nil {
			innerHashes = append(innerHashes, spn.Left.Hash)
		} else if spn.Right != nil {
			innerHashes = append(innerHashes, spn.Right.Hash)
		} else {
			break
		}
		spn = spn.Parent
	}
	return innerHashes
}

func trailsFromHashers(items []Hasher) (trails []*SimpleProofNode, root *SimpleProofNode) {

	switch len(items) {
	case 0:
		return nil, nil
	case 1:
		trail := &SimpleProofNode{items[0].Hash(), nil, nil, nil}
		return []*SimpleProofNode{trail}, trail
	default:
		lefts, leftRoot := trailsFromHashers(items[:(len(items)+1)/2])
		rights, rightRoot := trailsFromHashers(items[(len(items)+1)/2:])
		rootHash := SimpleHashFromTwoHashes(leftRoot.Hash, rightRoot.Hash)
		root := &SimpleProofNode{rootHash, nil, nil, nil}
		leftRoot.Parent = root
		leftRoot.Right = rightRoot
		rightRoot.Parent = root
		rightRoot.Left = leftRoot
		return append(lefts, rights...), root
	}
}
