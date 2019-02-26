package merkle

import (
	"golang.org/x/crypto/ripemd160"
)

func SimpleHashFromTwoHashes(left []byte, right []byte) []byte {
	var hasher = ripemd160.New()
	err := encodeByteSlice(hasher, left)
	if err != nil {
		panic(err)
	}
	err = encodeByteSlice(hasher, right)
	if err != nil {
		panic(err)
	}
	return hasher.Sum(nil)
}

func SimpleHashFromHashes(hashes [][]byte) []byte {

	switch len(hashes) {
	case 0:
		return nil
	case 1:
		return hashes[0]
	default:
		left := SimpleHashFromHashes(hashes[:(len(hashes)+1)/2])
		right := SimpleHashFromHashes(hashes[(len(hashes)+1)/2:])
		return SimpleHashFromTwoHashes(left, right)
	}
}

func SimpleHashFromByteslices(bzs [][]byte) []byte {
	hashes := make([][]byte, len(bzs))
	for i, bz := range bzs {
		hashes[i] = SimpleHashFromBytes(bz)
	}
	return SimpleHashFromHashes(hashes)
}

func SimpleHashFromBytes(bz []byte) []byte {
	hasher := ripemd160.New()
	hasher.Write(bz)
	return hasher.Sum(nil)
}

func SimpleHashFromHashers(items []Hasher) []byte {
	hashes := make([][]byte, len(items))
	for i, item := range items {
		hash := item.Hash()
		hashes[i] = hash
	}
	return SimpleHashFromHashes(hashes)
}

func SimpleHashFromMap(m map[string]Hasher) []byte {
	sm := NewSimpleMap()
	for k, v := range m {
		sm.Set(k, v)
	}
	return sm.Hash()
}
