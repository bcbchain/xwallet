package bcdb

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"path/filepath"
	"strings"
	"net/http"
	"io"
	"github.com/tendermint/tmlibs/log"
)

type LevelDB struct {
	db *leveldb.DB
}

func nonNilBytes(bz []byte) []byte {
	if bz == nil {
		return []byte{}
	} else {
		return bz
	}
}

func OpenDB(name string, ip string, port string) (*LevelDB, error) {
	var dbPath string
	if strings.HasPrefix(name, "/") {
		dbPath = name + ".db"
	} else {
		home := os.Getenv("HOME")
		dbPath = filepath.Join(home, name+".db")
	}
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	database := &LevelDB{
		db: db,
	}
	return database, nil
}

func (db *LevelDB) Get(key []byte) ([]byte, error) {
	key = nonNilBytes(key)
	res, err := db.db.Get(key, nil)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return res, nil
}

func (db *LevelDB) Has(key []byte) bool {
	v, _ := db.Get(key)
	return v != nil
}

func (db *LevelDB) Set(key []byte, value []byte) error {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	return db.db.Put(key, value, nil)
}

func (db *LevelDB) SetSync(key []byte, value []byte) error {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	return db.db.Put(key, value, &opt.WriteOptions{Sync: true})
}

func (db *LevelDB) Delete(key []byte) error {
	key = nonNilBytes(key)
	return db.db.Delete(key, nil)
}

func (db *LevelDB) DeleteSync(key []byte) error {
	key = nonNilBytes(key)
	return db.db.Delete(key, &opt.WriteOptions{Sync: true})
}

func (db *LevelDB) Close() {
	db.db.Close()
}

func (db *LevelDB) Print() {
	str, _ := db.db.GetProperty("leveldb.stats")
	fmt.Printf("%v\n", str)

	iter := db.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("%s:%s\n", string(key), string(value))
	}
}

func (db *LevelDB) GetAllKey() []byte {

	var data []string

	iter := db.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		data = append(data, string(key))
	}
	if len(data) == 0 {
		return nil
	}
	keysBytes := strings.Join(data, ";")

	return []byte(keysBytes)
}

func (db *LevelDB) queryDB(w http.ResponseWriter, req *http.Request) {
	value, err := db.Get([]byte(req.RequestURI))
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, string(value))
	}
}

func (db *LevelDB) StartQueryDBServer(queryAddress string, logger log.Logger) {
	http.HandleFunc("/", db.queryDB)
	logger.Info("StartQueryDBServer", "address:", queryAddress)
	err := http.ListenAndServe(queryAddress, nil)
	if err != nil {
		logger.Error("ListenAndServe: ", "error", err)
	}
}

type LevelDBBatch struct {
	db	*LevelDB
	batch	*leveldb.Batch
}

func (db *LevelDB) NewBatch() *LevelDBBatch {
	batch := new(leveldb.Batch)
	return &LevelDBBatch{db, batch}
}

func (mBatch *LevelDBBatch) Set(key, value []byte) {
	mBatch.batch.Put(key, value)
}

func (mBatch *LevelDBBatch) Delete(key []byte) {
	mBatch.batch.Delete(key)
}

func (mBatch *LevelDBBatch) Commit() error {
	return mBatch.db.db.Write(mBatch.batch, nil)
}
