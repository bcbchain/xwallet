package db

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cmn "github.com/tendermint/tmlibs/common"
)

func cleanupDBDir(dir, name string) {
	os.RemoveAll(filepath.Join(dir, name) + ".db")
}

func testBackendGetSetDelete(t *testing.T, backend DBBackendType) {

	dir, dirname := cmn.Tempdir(fmt.Sprintf("test_backend_%s_", backend))
	defer dir.Close()
	db := NewDB("testdb", backend, dirname)

	require.Nil(t, db.Get([]byte("")))

	require.Nil(t, db.Get(nil))

	key := []byte("abc")
	require.Nil(t, db.Get(key))

	db.Set(key, []byte(""))
	require.NotNil(t, db.Get(key))
	require.Empty(t, db.Get(key))

	db.Set(key, nil)
	require.NotNil(t, db.Get(key))
	require.Empty(t, db.Get(key))

	db.Delete(key)
	require.Nil(t, db.Get(key))
}

func TestBackendsGetSetDelete(t *testing.T) {
	for dbType := range backends {
		testBackendGetSetDelete(t, dbType)
	}
}

func withDB(t *testing.T, creator dbCreator, fn func(DB)) {
	name := cmn.Fmt("test_%x", cmn.RandStr(12))
	db, err := creator(name, "")
	defer cleanupDBDir("", name)
	assert.Nil(t, err)
	fn(db)
	db.Close()
}

func TestBackendsNilKeys(t *testing.T) {

	for dbType, creator := range backends {
		withDB(t, creator, func(db DB) {
			t.Run(fmt.Sprintf("Testing %s", dbType), func(t *testing.T) {

				expect := func(key, value []byte) {
					if len(key) == 0 {
						assert.Equal(t, db.Get(nil), db.Get([]byte("")))
						assert.Equal(t, db.Has(nil), db.Has([]byte("")))
					}
					assert.Equal(t, db.Get(key), value)
					assert.Equal(t, db.Has(key), value != nil)
				}

				expect(nil, nil)

				db.Set(nil, nil)
				expect(nil, []byte(""))

				db.Set(nil, []byte(""))
				expect(nil, []byte(""))

				db.Set(nil, []byte("abc"))
				expect(nil, []byte("abc"))
				db.Delete(nil)
				expect(nil, nil)

				db.Set(nil, []byte("abc"))
				expect(nil, []byte("abc"))
				db.Delete([]byte(""))
				expect(nil, nil)

				db.Set([]byte(""), []byte("abc"))
				expect(nil, []byte("abc"))
				db.Delete(nil)
				expect(nil, nil)

				db.Set([]byte(""), []byte("abc"))
				expect(nil, []byte("abc"))
				db.Delete([]byte(""))
				expect(nil, nil)

				db.SetSync(nil, []byte("abc"))
				expect(nil, []byte("abc"))
				db.DeleteSync(nil)
				expect(nil, nil)

				db.SetSync(nil, []byte("abc"))
				expect(nil, []byte("abc"))
				db.DeleteSync([]byte(""))
				expect(nil, nil)

				db.SetSync([]byte(""), []byte("abc"))
				expect(nil, []byte("abc"))
				db.DeleteSync(nil)
				expect(nil, nil)

				db.SetSync([]byte(""), []byte("abc"))
				expect(nil, []byte("abc"))
				db.DeleteSync([]byte(""))
				expect(nil, nil)
			})
		})
	}
}

func TestGoLevelDBBackend(t *testing.T) {
	name := cmn.Fmt("test_%x", cmn.RandStr(12))
	db := NewDB(name, GoLevelDBBackend, "")
	defer cleanupDBDir("", name)

	_, ok := db.(*GoLevelDB)
	assert.True(t, ok)
}
