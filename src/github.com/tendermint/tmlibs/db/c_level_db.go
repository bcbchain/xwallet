package db

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/jmhodges/levigo"
)

func init() {
	dbCreator := func(name string, dir string) (DB, error) {
		return NewCLevelDB(name, dir)
	}
	registerDBCreator(LevelDBBackend, dbCreator, true)
	registerDBCreator(CLevelDBBackend, dbCreator, false)
}

var _ DB = (*CLevelDB)(nil)

type CLevelDB struct {
	db	*levigo.DB
	ro	*levigo.ReadOptions
	wo	*levigo.WriteOptions
	woSync	*levigo.WriteOptions
}

func NewCLevelDB(name string, dir string) (*CLevelDB, error) {
	dbPath := filepath.Join(dir, name+".db")

	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(1 << 30))
	opts.SetCreateIfMissing(true)
	db, err := levigo.Open(dbPath, opts)
	if err != nil {
		return nil, err
	}
	ro := levigo.NewReadOptions()
	wo := levigo.NewWriteOptions()
	woSync := levigo.NewWriteOptions()
	woSync.SetSync(true)
	database := &CLevelDB{
		db:	db,
		ro:	ro,
		wo:	wo,
		woSync:	woSync,
	}
	return database, nil
}

func (db *CLevelDB) Get(key []byte) []byte {
	key = nonNilBytes(key)
	res, err := db.db.Get(db.ro, key)
	if err != nil {
		panic(err)
	}
	return res
}

func (db *CLevelDB) Has(key []byte) bool {
	return db.Get(key) != nil
}

func (db *CLevelDB) Set(key []byte, value []byte) {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	err := db.db.Put(db.wo, key, value)
	if err != nil {
		panic(err)
	}
}

func (db *CLevelDB) SetSync(key []byte, value []byte) {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	err := db.db.Put(db.woSync, key, value)
	if err != nil {
		panic(err)
	}
}

func (db *CLevelDB) Delete(key []byte) {
	key = nonNilBytes(key)
	err := db.db.Delete(db.wo, key)
	if err != nil {
		panic(err)
	}
}

func (db *CLevelDB) DeleteSync(key []byte) {
	key = nonNilBytes(key)
	err := db.db.Delete(db.woSync, key)
	if err != nil {
		panic(err)
	}
}

func (db *CLevelDB) DB() *levigo.DB {
	return db.db
}

func (db *CLevelDB) Close() {
	db.db.Close()
	db.ro.Close()
	db.wo.Close()
	db.woSync.Close()
}

func (db *CLevelDB) Print() {
	itr := db.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		key := itr.Key()
		value := itr.Value()
		fmt.Printf("[%X]:\t[%X]\n", key, value)
	}
}

func (db *CLevelDB) Stats() map[string]string {

	keys := []string{}

	stats := make(map[string]string)
	for _, key := range keys {
		str := db.db.PropertyValue(key)
		stats[key] = str
	}
	return stats
}

func (db *CLevelDB) NewBatch() Batch {
	batch := levigo.NewWriteBatch()
	return &cLevelDBBatch{db, batch}
}

type cLevelDBBatch struct {
	db	*CLevelDB
	batch	*levigo.WriteBatch
}

func (mBatch *cLevelDBBatch) Set(key, value []byte) {
	mBatch.batch.Put(key, value)
}

func (mBatch *cLevelDBBatch) Delete(key []byte) {
	mBatch.batch.Delete(key)
}

func (mBatch *cLevelDBBatch) Write() {
	err := mBatch.db.db.Write(mBatch.db.wo, mBatch.batch)
	if err != nil {
		panic(err)
	}
}

func (mBatch *cLevelDBBatch) WriteSync() {
	err := mBatch.db.db.Write(mBatch.db.woSync, mBatch.batch)
	if err != nil {
		panic(err)
	}
}

func (db *CLevelDB) Iterator(start, end []byte) Iterator {
	itr := db.db.NewIterator(db.ro)
	return newCLevelDBIterator(itr, start, end, false)
}

func (db *CLevelDB) ReverseIterator(start, end []byte) Iterator {
	panic("not implemented yet")
}

var _ Iterator = (*cLevelDBIterator)(nil)

type cLevelDBIterator struct {
	source		*levigo.Iterator
	start, end	[]byte
	isReverse	bool
	isInvalid	bool
}

func newCLevelDBIterator(source *levigo.Iterator, start, end []byte, isReverse bool) *cLevelDBIterator {
	if isReverse {
		panic("not implemented yet")
	}
	if start != nil {
		source.Seek(start)
	} else {
		source.SeekToFirst()
	}
	return &cLevelDBIterator{
		source:		source,
		start:		start,
		end:		end,
		isReverse:	isReverse,
		isInvalid:	false,
	}
}

func (itr cLevelDBIterator) Domain() ([]byte, []byte) {
	return itr.start, itr.end
}

func (itr cLevelDBIterator) Valid() bool {

	if itr.isInvalid {
		return false
	}

	itr.assertNoError()

	if !itr.source.Valid() {
		itr.isInvalid = true
		return false
	}

	var end = itr.end
	var key = itr.source.Key()
	if end != nil && bytes.Compare(end, key) <= 0 {
		itr.isInvalid = true
		return false
	}

	return true
}

func (itr cLevelDBIterator) Key() []byte {
	itr.assertNoError()
	itr.assertIsValid()
	return itr.source.Key()
}

func (itr cLevelDBIterator) Value() []byte {
	itr.assertNoError()
	itr.assertIsValid()
	return itr.source.Value()
}

func (itr cLevelDBIterator) Next() {
	itr.assertNoError()
	itr.assertIsValid()
	itr.source.Next()
}

func (itr cLevelDBIterator) Close() {
	itr.source.Close()
}

func (itr cLevelDBIterator) assertNoError() {
	if err := itr.source.GetError(); err != nil {
		panic(err)
	}
}

func (itr cLevelDBIterator) assertIsValid() {
	if !itr.Valid() {
		panic("cLevelDBIterator is invalid")
	}
}
