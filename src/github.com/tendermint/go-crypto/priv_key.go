package crypto

import (
	"crypto/subtle"

	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/tendermint/ed25519"
	"github.com/tendermint/ed25519/extra25519"
	. "github.com/tendermint/tmlibs/common"
)

func PrivKeyFromBytes(privKeyBytes []byte) (privKey PrivKey, err error) {
	err = cdc.UnmarshalBinaryBare(privKeyBytes, &privKey)
	return
}

type PrivKey interface {
	Bytes() []byte
	Sign(msg []byte) Signature
	PubKey() PubKey
	Equals(PrivKey) bool
}

var _ PrivKey = PrivKeyEd25519{}

type PrivKeyEd25519 [64]byte

func PrivKeyEd25519FromBytes(data []byte) PrivKey {
	var privkey PrivKeyEd25519
	copy(privkey[:], data)
	return privkey
}

func (privKey PrivKeyEd25519) Bytes() []byte {
	bz, err := cdc.MarshalBinaryBare(privKey)
	if err != nil {
		panic(err)
	}
	return bz
}

func (privKey PrivKeyEd25519) Sign(msg []byte) Signature {
	privKeyBytes := [64]byte(privKey)
	signatureBytes := ed25519.Sign(&privKeyBytes, msg)
	return SignatureEd25519(*signatureBytes)
}

func (privKey PrivKeyEd25519) PubKey() PubKey {
	privKeyBytes := [64]byte(privKey)
	pubBytes := *ed25519.MakePublicKey(&privKeyBytes)
	return PubKeyEd25519(pubBytes)
}

func (privKey PrivKeyEd25519) Equals(other PrivKey) bool {
	if otherEd, ok := other.(PrivKeyEd25519); ok {
		return subtle.ConstantTimeCompare(privKey[:], otherEd[:]) == 1
	} else {
		return false
	}
}

func (privKey PrivKeyEd25519) ToCurve25519() *[32]byte {
	keyCurve25519 := new([32]byte)
	privKeyBytes := [64]byte(privKey)
	extra25519.PrivateKeyToCurve25519(keyCurve25519, &privKeyBytes)
	return keyCurve25519
}

func (privKey PrivKeyEd25519) Generate(index int) PrivKeyEd25519 {
	bz, err := cdc.MarshalBinaryBare(struct {
		PrivKey	[64]byte
		Index	int
	}{privKey, index})
	if err != nil {
		panic(err)
	}
	newBytes := Sha256(bz)
	var newKey [64]byte
	copy(newKey[:], newBytes)
	return PrivKeyEd25519(newKey)
}

func GenPrivKeyEd25519() PrivKeyEd25519 {
	privKeyBytes := new([64]byte)
	copy(privKeyBytes[:32], CRandBytes(32))
	ed25519.MakePublicKey(privKeyBytes)
	return PrivKeyEd25519(*privKeyBytes)
}

func GenPrivKeyEd25519FromSecret(secret []byte) PrivKeyEd25519 {
	privKey32 := Sha256(secret)
	privKeyBytes := new([64]byte)
	copy(privKeyBytes[:32], privKey32)
	ed25519.MakePublicKey(privKeyBytes)
	return PrivKeyEd25519(*privKeyBytes)
}

var _ PrivKey = PrivKeySecp256k1{}

type PrivKeySecp256k1 [32]byte

func (privKey PrivKeySecp256k1) Bytes() []byte {
	bz, err := cdc.MarshalBinaryBare(privKey)
	if err != nil {
		panic(err)
	}
	return bz
}

func (privKey PrivKeySecp256k1) Sign(msg []byte) Signature {
	priv__, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey[:])
	sig__, err := priv__.Sign(Sha256(msg))
	if err != nil {
		PanicSanity(err)
	}
	return SignatureSecp256k1(sig__.Serialize())
}

func (privKey PrivKeySecp256k1) PubKey() PubKey {
	_, pub__ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey[:])
	var pub PubKeySecp256k1
	copy(pub[:], pub__.SerializeCompressed())
	return pub
}

func (privKey PrivKeySecp256k1) Equals(other PrivKey) bool {
	if otherSecp, ok := other.(PrivKeySecp256k1); ok {
		return subtle.ConstantTimeCompare(privKey[:], otherSecp[:]) == 1
	} else {
		return false
	}
}

func GenPrivKeySecp256k1() PrivKeySecp256k1 {
	privKeyBytes := [32]byte{}
	copy(privKeyBytes[:], CRandBytes(32))
	priv, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKeyBytes[:])
	copy(privKeyBytes[:], priv.Serialize())
	return PrivKeySecp256k1(privKeyBytes)
}

func GenPrivKeySecp256k1FromSecret(secret []byte) PrivKeySecp256k1 {
	privKey32 := Sha256(secret)
	priv, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey32)
	privKeyBytes := [32]byte{}
	copy(privKeyBytes[:], priv.Serialize())
	return PrivKeySecp256k1(privKeyBytes)
}
