package merkle

import (
	cmn "github.com/tendermint/tmlibs/common"
	"golang.org/x/crypto/ripemd160"
)

type SimpleMap struct {
	kvs	cmn.KVPairs
	sorted	bool
}

func NewSimpleMap() *SimpleMap {
	return &SimpleMap{
		kvs:	nil,
		sorted:	false,
	}
}

func (sm *SimpleMap) Set(key string, value Hasher) {
	sm.sorted = false

	khash := SimpleHashFromBytes([]byte(key))

	vhash := value.Hash()

	sm.kvs = append(sm.kvs, cmn.KVPair{
		Key:	khash,
		Value:	vhash,
	})
}

func (sm *SimpleMap) Hash() []byte {
	sm.Sort()
	return hashKVPairs(sm.kvs)
}

func (sm *SimpleMap) Sort() {
	if sm.sorted {
		return
	}
	sm.kvs.Sort()
	sm.sorted = true
}

func (sm *SimpleMap) KVPairs() cmn.KVPairs {
	sm.Sort()
	kvs := make(cmn.KVPairs, len(sm.kvs))
	copy(kvs, sm.kvs)
	return kvs
}

type KVPair cmn.KVPair

func (kv KVPair) Hash() []byte {
	hasher := ripemd160.New()
	err := encodeByteSlice(hasher, kv.Key)
	if err != nil {
		panic(err)
	}
	err = encodeByteSlice(hasher, kv.Value)
	if err != nil {
		panic(err)
	}
	return hasher.Sum(nil)
}

func hashKVPairs(kvs cmn.KVPairs) []byte {
	kvsH := make([]Hasher, 0, len(kvs))
	for _, kvp := range kvs {
		kvsH = append(kvsH, KVPair(kvp))
	}
	return SimpleHashFromHashers(kvsH)
}
