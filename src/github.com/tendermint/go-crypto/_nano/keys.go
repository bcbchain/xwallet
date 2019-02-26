package nano

import (
	"bytes"
	"encoding/hex"

	"github.com/pkg/errors"

	ledger "github.com/ethanfrey/ledger"

	crypto "github.com/tendermint/go-crypto"
	amino "github.com/tendermint/go-amino"
)

const (
	NameLedgerEd25519	= "ledger-ed25519"
	TypeLedgerEd25519	= 0x10

	Timeout	= 20
)

var device *ledger.Ledger

func getLedger() (*ledger.Ledger, error) {
	var err error
	if device == nil {
		device, err = ledger.FindLedger()
	}
	return device, err
}

func signLedger(device *ledger.Ledger, msg []byte) (pub crypto.PubKey, sig crypto.Signature, err error) {
	var resp []byte

	packets := generateSignRequests(msg)
	for _, pack := range packets {
		resp, err = device.Exchange(pack, Timeout)
		if err != nil {
			return pub, sig, err
		}
	}

	key, bsig, err := parseDigest(resp)
	if err != nil {
		return pub, sig, err
	}

	var b [32]byte
	copy(b[:], key)
	return PubKeyLedgerEd25519FromBytes(b), crypto.SignatureEd25519FromBytes(bsig), nil
}

type PrivKeyLedgerEd25519 struct {
	CachedPubKey crypto.PubKey
}

func NewPrivKeyLedgerEd25519() (crypto.PrivKey, error) {
	var pk PrivKeyLedgerEd25519

	_, err := pk.getPubKey()
	return pk.Wrap(), err
}

func (pk *PrivKeyLedgerEd25519) ValidateKey() error {

	pub, err := pk.forceGetPubKey()
	if err != nil {
		return err
	}

	if !pub.Equals(pk.CachedPubKey) {
		return errors.New("ledger doesn't match cached key")
	}
	return nil
}

func (pk *PrivKeyLedgerEd25519) AssertIsPrivKeyInner()	{}

func (pk *PrivKeyLedgerEd25519) Bytes() []byte {
	return amino.BinaryBytes(pk.Wrap())
}

func (pk *PrivKeyLedgerEd25519) Sign(msg []byte) crypto.Signature {

	dev, err := getLedger()
	if err != nil {
		panic(err)
	}

	pub, sig, err := signLedger(dev, msg)
	if err != nil {
		panic(err)
	}

	if pk.CachedPubKey.Empty() {
		pk.CachedPubKey = pub
	} else if !pk.CachedPubKey.Equals(pub) {
		panic("signed with a different key than stored")
	}
	return sig
}

func (pk *PrivKeyLedgerEd25519) PubKey() crypto.PubKey {
	key, err := pk.getPubKey()
	if err != nil {
		panic(err)
	}
	return key
}

func (pk *PrivKeyLedgerEd25519) getPubKey() (key crypto.PubKey, err error) {

	if pk.CachedPubKey.Empty() {
		pk.CachedPubKey, err = pk.forceGetPubKey()
	}
	return pk.CachedPubKey, err
}

func (pk *PrivKeyLedgerEd25519) forceGetPubKey() (key crypto.PubKey, err error) {
	dev, err := getLedger()
	if err != nil {
		return key, errors.New("Can't connect to ledger device")
	}
	key, _, err = signLedger(dev, []byte{0})
	if err != nil {
		return key, errors.New("Please open cosmos app on the ledger")
	}
	return key, err
}

func (pk *PrivKeyLedgerEd25519) Equals(other crypto.PrivKey) bool {
	if ledger, ok := other.Unwrap().(*PrivKeyLedgerEd25519); ok {
		return pk.CachedPubKey.Equals(ledger.CachedPubKey)
	}
	return false
}

type MockPrivKeyLedgerEd25519 struct {
	Msg	[]byte
	Pub	[KeyLength]byte
	Sig	[SigLength]byte
}

func NewMockKey(msg, pubkey, sig string) (pk MockPrivKeyLedgerEd25519) {
	var err error
	pk.Msg, err = hex.DecodeString(msg)
	if err != nil {
		panic(err)
	}

	bpk, err := hex.DecodeString(pubkey)
	if err != nil {
		panic(err)
	}
	bsig, err := hex.DecodeString(sig)
	if err != nil {
		panic(err)
	}

	copy(pk.Pub[:], bpk)
	copy(pk.Sig[:], bsig)
	return pk
}

var _ crypto.PrivKeyInner = MockPrivKeyLedgerEd25519{}

func (pk MockPrivKeyLedgerEd25519) AssertIsPrivKeyInner()	{}

func (pk MockPrivKeyLedgerEd25519) Bytes() []byte {
	return nil
}

func (pk MockPrivKeyLedgerEd25519) Sign(msg []byte) crypto.Signature {
	if !bytes.Equal(pk.Msg, msg) {
		panic("Mock key is for different msg")
	}
	return crypto.SignatureEd25519(pk.Sig).Wrap()
}

func (pk MockPrivKeyLedgerEd25519) PubKey() crypto.PubKey {
	return PubKeyLedgerEd25519FromBytes(pk.Pub)
}

func (pk MockPrivKeyLedgerEd25519) Equals(other crypto.PrivKey) bool {
	if mock, ok := other.Unwrap().(MockPrivKeyLedgerEd25519); ok {
		return bytes.Equal(mock.Pub[:], pk.Pub[:]) &&
			bytes.Equal(mock.Sig[:], pk.Sig[:]) &&
			bytes.Equal(mock.Msg, pk.Msg)
	}
	return false
}

type PubKeyLedgerEd25519 struct {
	crypto.PubKeyEd25519
}

func PubKeyLedgerEd25519FromBytes(key [32]byte) crypto.PubKey {
	return PubKeyLedgerEd25519{crypto.PubKeyEd25519(key)}.Wrap()
}

func (pk PubKeyLedgerEd25519) Bytes() []byte {
	return amino.BinaryBytes(pk.Wrap())
}

func (pk PubKeyLedgerEd25519) VerifyBytes(msg []byte, sig crypto.Signature) bool {
	hmsg := hashMsg(msg)
	return pk.PubKeyEd25519.VerifyBytes(hmsg, sig)
}

func (pk PubKeyLedgerEd25519) Equals(other crypto.PubKey) bool {
	if ledger, ok := other.Unwrap().(PubKeyLedgerEd25519); ok {
		return pk.PubKeyEd25519.Equals(ledger.PubKeyEd25519.Wrap())
	}
	return false
}

func init() {
	crypto.PrivKeyMapper.
		RegisterImplementation(&PrivKeyLedgerEd25519{}, NameLedgerEd25519, TypeLedgerEd25519).
		RegisterImplementation(MockPrivKeyLedgerEd25519{}, "mock-ledger", 0x11)

	crypto.PubKeyMapper.
		RegisterImplementation(PubKeyLedgerEd25519{}, NameLedgerEd25519, TypeLedgerEd25519)
}

func (pk *PrivKeyLedgerEd25519) Wrap() crypto.PrivKey {
	return crypto.PrivKey{PrivKeyInner: pk}
}

func (pk MockPrivKeyLedgerEd25519) Wrap() crypto.PrivKey {
	return crypto.PrivKey{PrivKeyInner: pk}
}

func (pk PubKeyLedgerEd25519) Wrap() crypto.PubKey {
	return crypto.PubKey{PubKeyInner: pk}
}
