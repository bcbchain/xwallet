package db

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/pkg/errors"
	cmn "github.com/tendermint/tmlibs/common"
)

const (
	keyPerm	= os.FileMode(0600)
	dirPerm	= os.FileMode(0700)
)

func init() {
	registerDBCreator(FSDBBackend, func(name string, dir string) (DB, error) {
		dbPath := filepath.Join(dir, name+".db")
		return NewFSDB(dbPath), nil
	}, false)
}

var _ DB = (*FSDB)(nil)

type FSDB struct {
	mtx	sync.Mutex
	dir	string
}

func NewFSDB(dir string) *FSDB {
	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		panic(errors.Wrap(err, "Creating FSDB dir "+dir))
	}
	database := &FSDB{
		dir: dir,
	}
	return database
}

func (db *FSDB) Get(key []byte) []byte {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	key = escapeKey(key)

	path := db.nameToPath(key)
	value, err := read(path)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		panic(errors.Wrapf(err, "Getting key %s (0x%X)", string(key), key))
	}
	return value
}

func (db *FSDB) Has(key []byte) bool {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	key = escapeKey(key)

	path := db.nameToPath(key)
	return cmn.FileExists(path)
}

func (db *FSDB) Set(key []byte, value []byte) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	db.SetNoLock(key, value)
}

func (db *FSDB) SetSync(key []byte, value []byte) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	db.SetNoLock(key, value)
}

func (db *FSDB) SetNoLock(key []byte, value []byte) {
	key = escapeKey(key)
	value = nonNilBytes(value)
	path := db.nameToPath(key)
	err := write(path, value)
	if err != nil {
		panic(errors.Wrapf(err, "Setting key %s (0x%X)", string(key), key))
	}
}

func (db *FSDB) Delete(key []byte) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	db.DeleteNoLock(key)
}

func (db *FSDB) DeleteSync(key []byte) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	db.DeleteNoLock(key)
}

func (db *FSDB) DeleteNoLock(key []byte) {
	key = escapeKey(key)
	path := db.nameToPath(key)
	err := remove(path)
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		panic(errors.Wrapf(err, "Removing key %s (0x%X)", string(key), key))
	}
}

func (db *FSDB) Close() {

}

func (db *FSDB) Print() {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	panic("FSDB.Print not yet implemented")
}

func (db *FSDB) Stats() map[string]string {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	panic("FSDB.Stats not yet implemented")
}

func (db *FSDB) NewBatch() Batch {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	panic("FSDB.NewBatch not yet implemented")
}

func (db *FSDB) Mutex() *sync.Mutex {
	return &(db.mtx)
}

func (db *FSDB) Iterator(start, end []byte) Iterator {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	keys, err := list(db.dir, start, end)
	if err != nil {
		panic(errors.Wrapf(err, "Listing keys in %s", db.dir))
	}
	sort.Strings(keys)
	return newMemDBIterator(db, keys, start, end)
}

func (db *FSDB) ReverseIterator(start, end []byte) Iterator {
	panic("not implemented yet")
}

func (db *FSDB) nameToPath(name []byte) string {
	n := url.PathEscape(string(name))
	return filepath.Join(db.dir, n)
}

func read(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func write(path string, d []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, keyPerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(d)
	if err != nil {
		return err
	}
	err = f.Sync()
	return err
}

func remove(path string) error {
	return os.Remove(path)
}

func list(dirPath string, start, end []byte) ([]string, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, name := range names {
		n, err := url.PathUnescape(name)
		if err != nil {
			return nil, fmt.Errorf("Failed to unescape %s while listing", name)
		}
		key := unescapeKey([]byte(n))
		if IsKeyInDomain(key, start, end, false) {
			keys = append(keys, string(key))
		}
	}
	return keys, nil
}

func escapeKey(key []byte) []byte {
	return []byte("k_" + string(key))
}
func unescapeKey(escKey []byte) []byte {
	if len(escKey) < 2 {
		panic(fmt.Sprintf("Invalid esc key: %x", escKey))
	}
	if string(escKey[:2]) != "k_" {
		panic(fmt.Sprintf("Invalid esc key: %x", escKey))
	}
	return escKey[2:]
}