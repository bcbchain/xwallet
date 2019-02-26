package keys_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbm "github.com/tendermint/tmlibs/db"

	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-crypto/keys"
	"github.com/tendermint/go-crypto/keys/words"
)

func TestKeyManagement(t *testing.T) {

	cstore := keys.New(
		dbm.NewMemDB(),
		words.MustLoadCodec("english"),
	)

	algo := keys.AlgoEd25519
	n1, n2, n3 := "personal", "business", "other"
	p1, p2 := "1234", "really-secure!@#$"

	l, err := cstore.List()
	require.Nil(t, err)
	assert.Empty(t, l)

	_, err = cstore.Get(n1)
	assert.NotNil(t, err)
	i, _, err := cstore.Create(n1, p1, algo)
	require.Equal(t, n1, i.Name)
	require.Nil(t, err)
	_, _, err = cstore.Create(n2, p2, algo)
	require.Nil(t, err)

	i2, err := cstore.Get(n2)
	assert.Nil(t, err)
	_, err = cstore.Get(n3)
	assert.NotNil(t, err)

	keyS, err := cstore.List()
	require.Nil(t, err)
	require.Equal(t, 2, len(keyS))

	assert.Equal(t, n2, keyS[0].Name)
	assert.Equal(t, n1, keyS[1].Name)
	assert.Equal(t, i2.PubKey, keyS[0].PubKey)

	err = cstore.Delete("bad name", "foo")
	require.NotNil(t, err)
	err = cstore.Delete(n1, p1)
	require.Nil(t, err)
	keyS, err = cstore.List()
	require.Nil(t, err)
	assert.Equal(t, 1, len(keyS))
	_, err = cstore.Get(n1)
	assert.NotNil(t, err)

}

func TestSignVerify(t *testing.T) {

	cstore := keys.New(
		dbm.NewMemDB(),
		words.MustLoadCodec("english"),
	)
	algo := keys.AlgoSecp256k1

	n1, n2 := "some dude", "a dudette"
	p1, p2 := "1234", "foobar"

	i1, _, err := cstore.Create(n1, p1, algo)
	require.Nil(t, err)

	i2, _, err := cstore.Create(n2, p2, algo)
	require.Nil(t, err)

	d1 := []byte("my first message")
	d2 := []byte("some other important info!")

	s11, pub1, err := cstore.Sign(n1, p1, d1)
	require.Nil(t, err)
	require.Equal(t, i1.PubKey, pub1)

	s12, pub1, err := cstore.Sign(n1, p1, d2)
	require.Nil(t, err)
	require.Equal(t, i1.PubKey, pub1)

	s21, pub2, err := cstore.Sign(n2, p2, d1)
	require.Nil(t, err)
	require.Equal(t, i2.PubKey, pub2)

	s22, pub2, err := cstore.Sign(n2, p2, d2)
	require.Nil(t, err)
	require.Equal(t, i2.PubKey, pub2)

	cases := []struct {
		key	crypto.PubKey
		data	[]byte
		sig	crypto.Signature
		valid	bool
	}{

		{i1.PubKey, d1, s11, true},

		{i1.PubKey, d2, s11, false},
		{i2.PubKey, d1, s11, false},
		{i1.PubKey, d1, s21, false},

		{i1.PubKey, d2, s12, true},
		{i2.PubKey, d1, s21, true},
		{i2.PubKey, d2, s22, true},
	}

	for i, tc := range cases {
		valid := tc.key.VerifyBytes(tc.data, tc.sig)
		assert.Equal(t, tc.valid, valid, "%d", i)
	}
}

func assertPassword(t *testing.T, cstore keys.Keybase, name, pass, badpass string) {
	err := cstore.Update(name, badpass, pass)
	assert.NotNil(t, err)
	err = cstore.Update(name, pass, pass)
	assert.Nil(t, err, "%+v", err)
}

func TestExportImport(t *testing.T) {

	db := dbm.NewMemDB()
	cstore := keys.New(
		db,
		words.MustLoadCodec("english"),
	)

	info, _, err := cstore.Create("john", "passphrase", keys.AlgoEd25519)
	assert.Nil(t, err)
	assert.Equal(t, info.Name, "john")
	addr := info.PubKey.Address()

	john, err := cstore.Get("john")
	assert.Nil(t, err)
	assert.Equal(t, john.Name, "john")
	assert.Equal(t, john.PubKey.Address(), addr)

	armor, err := cstore.Export("john")
	assert.Nil(t, err)

	err = cstore.Import("john2", armor)
	assert.Nil(t, err)

	john2, err := cstore.Get("john2")
	assert.Nil(t, err)

	assert.Equal(t, john.PubKey.Address(), addr)
	assert.Equal(t, john.Name, "john")
	assert.Equal(t, john, john2)
}

func TestAdvancedKeyManagement(t *testing.T) {

	cstore := keys.New(
		dbm.NewMemDB(),
		words.MustLoadCodec("english"),
	)

	algo := keys.AlgoSecp256k1
	n1, n2 := "old-name", "new name"
	p1, p2 := "1234", "foobar"

	_, _, err := cstore.Create(n1, p1, algo)
	require.Nil(t, err, "%+v", err)
	assertPassword(t, cstore, n1, p1, p2)

	err = cstore.Update(n1, "jkkgkg", p2)
	assert.NotNil(t, err)
	assertPassword(t, cstore, n1, p1, p2)

	err = cstore.Update(n1, p1, p2)
	assert.Nil(t, err)

	assertPassword(t, cstore, n1, p2, p1)

	_, err = cstore.Export(n1 + ".notreal")
	assert.NotNil(t, err)
	_, err = cstore.Export(" " + n1)
	assert.NotNil(t, err)
	_, err = cstore.Export(n1 + " ")
	assert.NotNil(t, err)
	_, err = cstore.Export("")
	assert.NotNil(t, err)
	exported, err := cstore.Export(n1)
	require.Nil(t, err, "%+v", err)

	err = cstore.Import(n2, exported)
	assert.Nil(t, err)

	err = cstore.Import(n2, exported)
	assert.NotNil(t, err)
}

func TestSeedPhrase(t *testing.T) {

	cstore := keys.New(
		dbm.NewMemDB(),
		words.MustLoadCodec("english"),
	)

	algo := keys.AlgoEd25519
	n1, n2 := "lost-key", "found-again"
	p1, p2 := "1234", "foobar"

	info, seed, err := cstore.Create(n1, p1, algo)
	require.Nil(t, err, "%+v", err)
	assert.Equal(t, n1, info.Name)
	assert.NotEmpty(t, seed)

	err = cstore.Delete(n1, p1)
	require.Nil(t, err, "%+v", err)
	_, err = cstore.Get(n1)
	require.NotNil(t, err)

	newInfo, err := cstore.Recover(n2, p2, seed)
	require.Nil(t, err, "%+v", err)
	assert.Equal(t, n2, newInfo.Name)
	assert.Equal(t, info.Address(), newInfo.Address())
	assert.Equal(t, info.PubKey, newInfo.PubKey)
}

func ExampleNew() {

	cstore := keys.New(
		dbm.NewMemDB(),
		words.MustLoadCodec("english"),
	)
	ed := keys.AlgoEd25519
	sec := keys.AlgoSecp256k1

	bob, _, err := cstore.Create("Bob", "friend", ed)
	if err != nil {

		fmt.Println(err)
	} else {

		fmt.Println(bob.Name)
	}
	cstore.Create("Alice", "secret", sec)
	cstore.Create("Carl", "mitm", ed)
	info, _ := cstore.List()
	for _, i := range info {
		fmt.Println(i.Name)
	}

	tx := []byte("deadbeef")
	sig, pub, err := cstore.Sign("Bob", "friend", tx)
	if err != nil {
		fmt.Println("don't accept real passphrase")
	}

	binfo, _ := cstore.Get("Bob")
	if !binfo.PubKey.Equals(bob.PubKey) {
		fmt.Println("Get and Create return different keys")
	}

	if pub.Equals(binfo.PubKey) {
		fmt.Println("signed by Bob")
	}
	if !pub.VerifyBytes(tx, sig) {
		fmt.Println("invalid signature")
	}

}
