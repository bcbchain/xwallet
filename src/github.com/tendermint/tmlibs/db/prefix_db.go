package db

import (
	"bytes"
	"fmt"
	"sync"
)

func IteratePrefix(db DB, prefix []byte) Iterator {
	var start, end []byte
	if len(prefix) == 0 {
		start = nil
		end = nil
	} else {
		start = cp(prefix)
		end = cpIncr(prefix)
	}
	return db.Iterator(start, end)
}

type prefixDB struct {
	mtx	sync.Mutex
	prefix	[]byte
	db	DB
}

func NewPrefixDB(db DB, prefix []byte) *prefixDB {
	return &prefixDB{
		prefix:	prefix,
		db:	db,
	}
}

func (pdb *prefixDB) Mutex() *sync.Mutex {
	return &(pdb.mtx)
}

func (pdb *prefixDB) Get(key []byte) []byte {
	pdb.mtx.Lock()
	defer pdb.mtx.Unlock()

	pkey := pdb.prefixed(key)
	value := pdb.db.Get(pkey)
	return value
}

func (pdb *prefixDB) Has(key []byte) bool {
	pdb.mtx.Lock()
	defer pdb.mtx.Unlock()

	return pdb.db.Has(pdb.prefixed(key))
}

func (pdb *prefixDB) Set(key []byte, value []byte) {
	pdb.mtx.Lock()
	defer pdb.mtx.Unlock()

	pkey := pdb.prefixed(key)
	pdb.db.Set(pkey, value)
}

func (pdb *prefixDB) SetSync(key []byte, value []byte) {
	pdb.mtx.Lock()
	defer pdb.mtx.Unlock()

	pdb.db.SetSync(pdb.prefixed(key), value)
}

func (pdb *prefixDB) Delete(key []byte) {
	pdb.mtx.Lock()
	defer pdb.mtx.Unlock()

	pdb.db.Delete(pdb.prefixed(key))
}

func (pdb *prefixDB) DeleteSync(key []byte) {
	pdb.mtx.Lock()
	defer pdb.mtx.Unlock()

	pdb.db.DeleteSync(pdb.prefixed(key))
}

func (pdb *prefixDB) Iterator(start, end []byte) Iterator {
	pdb.mtx.Lock()
	defer pdb.mtx.Unlock()

	var pstart, pend []byte
	pstart = append(cp(pdb.prefix), start...)
	if end == nil {
		pend = cpIncr(pdb.prefix)
	} else {
		pend = append(cp(pdb.prefix), end...)
	}
	return newPrefixIterator(
		pdb.prefix,
		start,
		end,
		pdb.db.Iterator(
			pstart,
			pend,
		),
	)
}

func (pdb *prefixDB) ReverseIterator(start, end []byte) Iterator {
	pdb.mtx.Lock()
	defer pdb.mtx.Unlock()

	var pstart, pend []byte
	if start == nil {

		pstart = cpIncr(pdb.prefix)
	} else {
		pstart = append(cp(pdb.prefix), start...)
	}
	if end == nil {

		pend = cpDecr(pdb.prefix)
	} else {
		pend = append(cp(pdb.prefix), end...)
	}
	ritr := pdb.db.ReverseIterator(pstart, pend)
	if start == nil {
		skipOne(ritr, cpIncr(pdb.prefix))
	}
	return newPrefixIterator(
		pdb.prefix,
		start,
		end,
		ritr,
	)
}

func (pdb *prefixDB) NewBatch() Batch {
	pdb.mtx.Lock()
	defer pdb.mtx.Unlock()

	return newPrefixBatch(pdb.prefix, pdb.db.NewBatch())
}

func (pdb *prefixDB) Close() {
	pdb.mtx.Lock()
	defer pdb.mtx.Unlock()

	pdb.db.Close()
}

func (pdb *prefixDB) Print() {
	fmt.Printf("prefix: %X\n", pdb.prefix)

	itr := pdb.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		key := itr.Key()
		value := itr.Value()
		fmt.Printf("[%X]:\t[%X]\n", key, value)
	}
}

func (pdb *prefixDB) Stats() map[string]string {
	stats := make(map[string]string)
	stats["prefixdb.prefix.string"] = string(pdb.prefix)
	stats["prefixdb.prefix.hex"] = fmt.Sprintf("%X", pdb.prefix)
	source := pdb.db.Stats()
	for key, value := range source {
		stats["prefixdb.source."+key] = value
	}
	return stats
}

func (pdb *prefixDB) prefixed(key []byte) []byte {
	return append(cp(pdb.prefix), key...)
}

type prefixBatch struct {
	prefix	[]byte
	source	Batch
}

func newPrefixBatch(prefix []byte, source Batch) prefixBatch {
	return prefixBatch{
		prefix:	prefix,
		source:	source,
	}
}

func (pb prefixBatch) Set(key, value []byte) {
	pkey := append(cp(pb.prefix), key...)
	pb.source.Set(pkey, value)
}

func (pb prefixBatch) Delete(key []byte) {
	pkey := append(cp(pb.prefix), key...)
	pb.source.Delete(pkey)
}

func (pb prefixBatch) Write() {
	pb.source.Write()
}

func (pb prefixBatch) WriteSync() {
	pb.source.WriteSync()
}

type prefixIterator struct {
	prefix	[]byte
	start	[]byte
	end	[]byte
	source	Iterator
	valid	bool
}

func newPrefixIterator(prefix, start, end []byte, source Iterator) prefixIterator {
	if !source.Valid() || !bytes.HasPrefix(source.Key(), prefix) {
		return prefixIterator{
			prefix:	prefix,
			start:	start,
			end:	end,
			source:	source,
			valid:	false,
		}
	} else {
		return prefixIterator{
			prefix:	prefix,
			start:	start,
			end:	end,
			source:	source,
			valid:	true,
		}
	}
}

func (itr prefixIterator) Domain() (start []byte, end []byte) {
	return itr.start, itr.end
}

func (itr prefixIterator) Valid() bool {
	return itr.valid && itr.source.Valid()
}

func (itr prefixIterator) Next() {
	if !itr.valid {
		panic("prefixIterator invalid, cannot call Next()")
	}
	itr.source.Next()
	if !itr.source.Valid() || !bytes.HasPrefix(itr.source.Key(), itr.prefix) {
		itr.source.Close()
		itr.valid = false
		return
	}
}

func (itr prefixIterator) Key() (key []byte) {
	if !itr.valid {
		panic("prefixIterator invalid, cannot call Key()")
	}
	return stripPrefix(itr.source.Key(), itr.prefix)
}

func (itr prefixIterator) Value() (value []byte) {
	if !itr.valid {
		panic("prefixIterator invalid, cannot call Value()")
	}
	return itr.source.Value()
}

func (itr prefixIterator) Close() {
	itr.source.Close()
}

func stripPrefix(key []byte, prefix []byte) (stripped []byte) {
	if len(key) < len(prefix) {
		panic("should not happen")
	}
	if !bytes.Equal(key[:len(prefix)], prefix) {
		panic("should not happne")
	}
	return key[len(prefix):]
}

func skipOne(itr Iterator, skipKey []byte) {
	if itr.Valid() {
		if bytes.Equal(itr.Key(), skipKey) {
			itr.Next()
		}
	}
}
