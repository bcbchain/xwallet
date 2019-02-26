package kv

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/pubsub/query"

	"github.com/tendermint/tendermint/state/txindex"
	"github.com/tendermint/tendermint/types"
)

const (
	tagKeySeparator = "/"
)

var _ txindex.TxIndexer = (*TxIndex)(nil)

type TxIndex struct {
	store		dbm.DB
	tagsToIndex	[]string
	indexAllTags	bool
}

func NewTxIndex(store dbm.DB, options ...func(*TxIndex)) *TxIndex {
	txi := &TxIndex{store: store, tagsToIndex: make([]string, 0), indexAllTags: false}
	for _, o := range options {
		o(txi)
	}
	return txi
}

func IndexTags(tags []string) func(*TxIndex) {
	return func(txi *TxIndex) {
		txi.tagsToIndex = tags
	}
}

func IndexAllTags() func(*TxIndex) {
	return func(txi *TxIndex) {
		txi.indexAllTags = true
	}
}

func (txi *TxIndex) Get(hash []byte) (*types.TxResult, error) {
	if len(hash) == 0 {
		return nil, txindex.ErrorEmptyHash
	}

	rawBytes := txi.store.Get(hash)
	if rawBytes == nil {
		return nil, nil
	}

	txResult := new(types.TxResult)
	err := cdc.UnmarshalBinaryBare(rawBytes, &txResult)
	if err != nil {
		return nil, fmt.Errorf("Error reading TxResult: %v", err)
	}

	return txResult, nil
}

func (txi *TxIndex) AddBatch(b *txindex.Batch) error {
	storeBatch := txi.store.NewBatch()

	for _, result := range b.Ops {
		hash := result.Tx.Hash()

		for _, tag := range result.Result.Tags {
			if txi.indexAllTags || cmn.StringInSlice(string(tag.Key), txi.tagsToIndex) {
				storeBatch.Set(keyForTag(tag, result), hash)
			}
		}

		rawBytes, err := cdc.MarshalBinaryBare(result)
		if err != nil {
			return err
		}
		storeBatch.Set(hash, rawBytes)
	}

	storeBatch.Write()
	return nil
}

func (txi *TxIndex) Index(result *types.TxResult) error {
	b := txi.store.NewBatch()

	hash := result.Tx.Hash()

	for _, tag := range result.Result.Tags {
		if txi.indexAllTags || cmn.StringInSlice(string(tag.Key), txi.tagsToIndex) {
			b.Set(keyForTag(tag, result), hash)
		}
	}

	rawBytes, err := cdc.MarshalBinaryBare(result)
	if err != nil {
		return err
	}
	b.Set(hash, rawBytes)

	b.Write()
	return nil
}

func (txi *TxIndex) Search(q *query.Query) ([]*types.TxResult, error) {
	var hashes [][]byte
	var hashesInitialized bool

	conditions := q.Conditions()

	hash, err, ok := lookForHash(conditions)
	if err != nil {
		return nil, errors.Wrap(err, "error during searching for a hash in the query")
	} else if ok {
		res, err := txi.Get(hash)
		if res == nil {
			return []*types.TxResult{}, nil
		}
		return []*types.TxResult{res}, errors.Wrap(err, "error while retrieving the result")
	}

	skipIndexes := make([]int, 0)

	height, heightIndex := lookForHeight(conditions)
	if heightIndex >= 0 {
		skipIndexes = append(skipIndexes, heightIndex)
	}

	ranges, rangeIndexes := lookForRanges(conditions)
	if len(ranges) > 0 {
		skipIndexes = append(skipIndexes, rangeIndexes...)

		for _, r := range ranges {
			if !hashesInitialized {
				hashes = txi.matchRange(r, []byte(r.key))
				hashesInitialized = true
			} else {
				hashes = intersect(hashes, txi.matchRange(r, []byte(r.key)))
			}
		}
	}

	for i, c := range conditions {
		if cmn.IntInSlice(i, skipIndexes) {
			continue
		}

		if !hashesInitialized {
			hashes = txi.match(c, startKey(c, height))
			hashesInitialized = true
		} else {
			hashes = intersect(hashes, txi.match(c, startKey(c, height)))
		}
	}

	results := make([]*types.TxResult, len(hashes))
	i := 0
	for _, h := range hashes {
		results[i], err = txi.Get(h)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get Tx{%X}", h)
		}
		i++
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Height < results[j].Height
	})

	return results, nil
}

func lookForHash(conditions []query.Condition) (hash []byte, err error, ok bool) {
	for _, c := range conditions {
		if c.Tag == types.TxHashKey {
			decoded, err := hex.DecodeString(c.Operand.(string))
			return decoded, err, true
		}
	}
	return
}

func lookForHeight(conditions []query.Condition) (height int64, index int) {
	for i, c := range conditions {
		if c.Tag == types.TxHeightKey {
			return c.Operand.(int64), i
		}
	}
	return 0, -1
}

type queryRanges map[string]queryRange

type queryRange struct {
	key			string
	lowerBound		interface{}
	includeLowerBound	bool
	upperBound		interface{}
	includeUpperBound	bool
}

func (r queryRange) lowerBoundValue() interface{} {
	if r.lowerBound == nil {
		return nil
	}

	if r.includeLowerBound {
		return r.lowerBound
	} else {
		switch t := r.lowerBound.(type) {
		case int64:
			return t + 1
		case time.Time:
			return t.Unix() + 1
		default:
			panic("not implemented")
		}
	}
}

func (r queryRange) AnyBound() interface{} {
	if r.lowerBound != nil {
		return r.lowerBound
	} else {
		return r.upperBound
	}
}

func (r queryRange) upperBoundValue() interface{} {
	if r.upperBound == nil {
		return nil
	}

	if r.includeUpperBound {
		return r.upperBound
	} else {
		switch t := r.upperBound.(type) {
		case int64:
			return t - 1
		case time.Time:
			return t.Unix() - 1
		default:
			panic("not implemented")
		}
	}
}

func lookForRanges(conditions []query.Condition) (ranges queryRanges, indexes []int) {
	ranges = make(queryRanges)
	for i, c := range conditions {
		if isRangeOperation(c.Op) {
			r, ok := ranges[c.Tag]
			if !ok {
				r = queryRange{key: c.Tag}
			}
			switch c.Op {
			case query.OpGreater:
				r.lowerBound = c.Operand
			case query.OpGreaterEqual:
				r.includeLowerBound = true
				r.lowerBound = c.Operand
			case query.OpLess:
				r.upperBound = c.Operand
			case query.OpLessEqual:
				r.includeUpperBound = true
				r.upperBound = c.Operand
			}
			ranges[c.Tag] = r
			indexes = append(indexes, i)
		}
	}
	return ranges, indexes
}

func isRangeOperation(op query.Operator) bool {
	switch op {
	case query.OpGreater, query.OpGreaterEqual, query.OpLess, query.OpLessEqual:
		return true
	default:
		return false
	}
}

func (txi *TxIndex) match(c query.Condition, startKey []byte) (hashes [][]byte) {
	if c.Op == query.OpEqual {
		it := dbm.IteratePrefix(txi.store, startKey)
		defer it.Close()
		for ; it.Valid(); it.Next() {
			hashes = append(hashes, it.Value())
		}
	} else if c.Op == query.OpContains {

		it := txi.store.Iterator(nil, nil)
		defer it.Close()
		for ; it.Valid(); it.Next() {
			if !isTagKey(it.Key()) {
				continue
			}
			if strings.Contains(extractValueFromKey(it.Key()), c.Operand.(string)) {
				hashes = append(hashes, it.Value())
			}
		}
	} else {
		panic("other operators should be handled already")
	}
	return
}

func (txi *TxIndex) matchRange(r queryRange, prefix []byte) (hashes [][]byte) {

	hashesMap := make(map[string][]byte)

	lowerBound := r.lowerBoundValue()
	upperBound := r.upperBoundValue()

	it := dbm.IteratePrefix(txi.store, prefix)
	defer it.Close()
LOOP:
	for ; it.Valid(); it.Next() {
		if !isTagKey(it.Key()) {
			continue
		}
		switch r.AnyBound().(type) {
		case int64:
			v, err := strconv.ParseInt(extractValueFromKey(it.Key()), 10, 64)
			if err != nil {
				continue LOOP
			}
			include := true
			if lowerBound != nil && v < lowerBound.(int64) {
				include = false
			}
			if upperBound != nil && v > upperBound.(int64) {
				include = false
			}
			if include {
				hashesMap[fmt.Sprintf("%X", it.Value())] = it.Value()
			}

		}
	}
	hashes = make([][]byte, len(hashesMap))
	i := 0
	for _, h := range hashesMap {
		hashes[i] = h
		i++
	}
	return
}

func startKey(c query.Condition, height int64) []byte {
	var key string
	if height > 0 {
		key = fmt.Sprintf("%s/%v/%d", c.Tag, c.Operand, height)
	} else {
		key = fmt.Sprintf("%s/%v", c.Tag, c.Operand)
	}
	return []byte(key)
}

func isTagKey(key []byte) bool {
	return strings.Count(string(key), tagKeySeparator) == 3
}

func extractValueFromKey(key []byte) string {
	parts := strings.SplitN(string(key), tagKeySeparator, 3)
	return parts[1]
}

func keyForTag(tag cmn.KVPair, result *types.TxResult) []byte {
	return []byte(fmt.Sprintf("%s/%s/%d/%d", tag.Key, tag.Value, result.Height, result.Index))
}

func intersect(as, bs [][]byte) [][]byte {
	i := make([][]byte, 0, cmn.MinInt(len(as), len(bs)))
	for _, a := range as {
		for _, b := range bs {
			if bytes.Equal(a, b) {
				i = append(i, a)
			}
		}
	}
	return i
}
