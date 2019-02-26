package db

import (
	"fmt"
	"sync"

	cmn "github.com/tendermint/tmlibs/common"
)

func _fmt(f string, az ...interface{}) string {
	return fmt.Sprintf(f, az...)
}

type debugDB struct {
	label	string
	db	DB
}

func NewDebugDB(label string, db DB) debugDB {
	return debugDB{
		label:	label,
		db:	db,
	}
}

func (ddb debugDB) Mutex() *sync.Mutex	{ return nil }

func (ddb debugDB) Get(key []byte) (value []byte) {
	defer func() {
		fmt.Printf("%v.Get(%v) %v\n", ddb.label, cmn.Cyan(_fmt("%X", key)), cmn.Blue(_fmt("%X", value)))
	}()
	value = ddb.db.Get(key)
	return
}

func (ddb debugDB) Has(key []byte) (has bool) {
	defer func() {
		fmt.Printf("%v.Has(%v) %v\n", ddb.label, cmn.Cyan(_fmt("%X", key)), has)
	}()
	return ddb.db.Has(key)
}

func (ddb debugDB) Set(key []byte, value []byte) {
	fmt.Printf("%v.Set(%v, %v)\n", ddb.label, cmn.Cyan(_fmt("%X", key)), cmn.Yellow(_fmt("%X", value)))
	ddb.db.Set(key, value)
}

func (ddb debugDB) SetSync(key []byte, value []byte) {
	fmt.Printf("%v.SetSync(%v, %v)\n", ddb.label, cmn.Cyan(_fmt("%X", key)), cmn.Yellow(_fmt("%X", value)))
	ddb.db.SetSync(key, value)
}

func (ddb debugDB) SetNoLock(key []byte, value []byte) {
	fmt.Printf("%v.SetNoLock(%v, %v)\n", ddb.label, cmn.Cyan(_fmt("%X", key)), cmn.Yellow(_fmt("%X", value)))
	ddb.db.(atomicSetDeleter).SetNoLock(key, value)
}

func (ddb debugDB) SetNoLockSync(key []byte, value []byte) {
	fmt.Printf("%v.SetNoLockSync(%v, %v)\n", ddb.label, cmn.Cyan(_fmt("%X", key)), cmn.Yellow(_fmt("%X", value)))
	ddb.db.(atomicSetDeleter).SetNoLockSync(key, value)
}

func (ddb debugDB) Delete(key []byte) {
	fmt.Printf("%v.Delete(%v)\n", ddb.label, cmn.Red(_fmt("%X", key)))
	ddb.db.Delete(key)
}

func (ddb debugDB) DeleteSync(key []byte) {
	fmt.Printf("%v.DeleteSync(%v)\n", ddb.label, cmn.Red(_fmt("%X", key)))
	ddb.db.DeleteSync(key)
}

func (ddb debugDB) DeleteNoLock(key []byte) {
	fmt.Printf("%v.DeleteNoLock(%v)\n", ddb.label, cmn.Red(_fmt("%X", key)))
	ddb.db.(atomicSetDeleter).DeleteNoLock(key)
}

func (ddb debugDB) DeleteNoLockSync(key []byte) {
	fmt.Printf("%v.DeleteNoLockSync(%v)\n", ddb.label, cmn.Red(_fmt("%X", key)))
	ddb.db.(atomicSetDeleter).DeleteNoLockSync(key)
}

func (ddb debugDB) Iterator(start, end []byte) Iterator {
	fmt.Printf("%v.Iterator(%v, %v)\n", ddb.label, cmn.Cyan(_fmt("%X", start)), cmn.Blue(_fmt("%X", end)))
	return NewDebugIterator(ddb.label, ddb.db.Iterator(start, end))
}

func (ddb debugDB) ReverseIterator(start, end []byte) Iterator {
	fmt.Printf("%v.ReverseIterator(%v, %v)\n", ddb.label, cmn.Cyan(_fmt("%X", start)), cmn.Blue(_fmt("%X", end)))
	return NewDebugIterator(ddb.label, ddb.db.ReverseIterator(start, end))
}

func (ddb debugDB) NewBatch() Batch {
	fmt.Printf("%v.NewBatch()\n", ddb.label)
	return NewDebugBatch(ddb.label, ddb.db.NewBatch())
}

func (ddb debugDB) Close() {
	fmt.Printf("%v.Close()\n", ddb.label)
	ddb.db.Close()
}

func (ddb debugDB) Print() {
	ddb.db.Print()
}

func (ddb debugDB) Stats() map[string]string {
	return ddb.db.Stats()
}

type debugIterator struct {
	label	string
	itr	Iterator
}

func NewDebugIterator(label string, itr Iterator) debugIterator {
	return debugIterator{
		label:	label,
		itr:	itr,
	}
}

func (ditr debugIterator) Domain() (start []byte, end []byte) {
	defer func() {
		fmt.Printf("%v.itr.Domain() (%X,%X)\n", ditr.label, start, end)
	}()
	start, end = ditr.itr.Domain()
	return
}

func (ditr debugIterator) Valid() (ok bool) {
	defer func() {
		fmt.Printf("%v.itr.Valid() %v\n", ditr.label, ok)
	}()
	ok = ditr.itr.Valid()
	return
}

func (ditr debugIterator) Next() {
	fmt.Printf("%v.itr.Next()\n", ditr.label)
	ditr.itr.Next()
}

func (ditr debugIterator) Key() (key []byte) {
	fmt.Printf("%v.itr.Key() %v\n", ditr.label, cmn.Cyan(_fmt("%X", key)))
	key = ditr.itr.Key()
	return
}

func (ditr debugIterator) Value() (value []byte) {
	fmt.Printf("%v.itr.Value() %v\n", ditr.label, cmn.Blue(_fmt("%X", value)))
	value = ditr.itr.Value()
	return
}

func (ditr debugIterator) Close() {
	fmt.Printf("%v.itr.Close()\n", ditr.label)
	ditr.itr.Close()
}

type debugBatch struct {
	label	string
	bch	Batch
}

func NewDebugBatch(label string, bch Batch) debugBatch {
	return debugBatch{
		label:	label,
		bch:	bch,
	}
}

func (dbch debugBatch) Set(key, value []byte) {
	fmt.Printf("%v.batch.Set(%v, %v)\n", dbch.label, cmn.Cyan(_fmt("%X", key)), cmn.Yellow(_fmt("%X", value)))
	dbch.bch.Set(key, value)
}

func (dbch debugBatch) Delete(key []byte) {
	fmt.Printf("%v.batch.Delete(%v)\n", dbch.label, cmn.Red(_fmt("%X", key)))
	dbch.bch.Delete(key)
}

func (dbch debugBatch) Write() {
	fmt.Printf("%v.batch.Write()\n", dbch.label)
	dbch.bch.Write()
}

func (dbch debugBatch) WriteSync() {
	fmt.Printf("%v.batch.WriteSync()\n", dbch.label)
	dbch.bch.WriteSync()
}
