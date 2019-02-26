package crypto

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"github.com/tendermint/ed25519"
	"github.com/tendermint/ed25519/extra25519"
	cmn "github.com/tendermint/tmlibs/common"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

type Address = string

type Hash = cmn.HexBytes

func PubKeyFromBytes(pubKeyBytes []byte) (pubKey PubKey, err error) {
	err = cdc.UnmarshalBinaryBare(pubKeyBytes, &pubKey)
	return
}

type PubKey interface {
	Address() Address
	Bytes() []byte
	VerifyBytes(msg []byte, sig Signature) bool
	Equals(PubKey) bool
}

var chainId []byte
var _ PubKey = PubKeyEd25519{}

type PubKeyEd25519 [32]byte

func SetChainId(cid string) {
	chainId = make([]byte, 0, 0)
	chainId = append(chainId, []byte(cid)...)
}

func GetChainId() string {
	if string(chainId) == "" {
		panic("crypto.SetChainId must be called first")
	}
	return string(chainId)
}

func PubKeyEd25519FromBytes(data []byte) PubKey {
	var pubkey PubKeyEd25519
	copy(pubkey[:], data)
	return pubkey
}

func (pubKey PubKeyEd25519) Address() Address {
	if string(chainId) == "" {
		panic("crypto.SetChainId must be called first")
	}

	hasherSHA3256 := sha3.New256()
	hasherSHA3256.Write(chainId[:])
	hasherSHA3256.Write(pubKey[:])
	sha := hasherSHA3256.Sum(nil)

	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha)
	rpd := hasherRIPEMD160.Sum(nil)

	hasher := ripemd160.New()
	hasher.Write(rpd)
	md := hasher.Sum(nil)

	addr := make([]byte, 0, 0)
	addr = append(addr, rpd...)
	addr = append(addr, md[:4]...)

	return string(chainId) + base58.Encode(addr)
}

func (pubKey PubKeyEd25519) Bytes() []byte {
	bz, err := cdc.MarshalBinaryBare(pubKey)
	if err != nil {
		panic(err)
	}
	return bz
}

func (pubKey PubKeyEd25519) VerifyBytes(msg []byte, sig_ Signature) bool {

	sig, ok := sig_.(SignatureEd25519)
	if !ok {
		return false
	}
	pubKeyBytes := [32]byte(pubKey)
	sigBytes := [64]byte(sig)
	return ed25519.Verify(&pubKeyBytes, msg, &sigBytes)
}

func (pubKey PubKeyEd25519) ToCurve25519() *[32]byte {
	keyCurve25519, pubKeyBytes := new([32]byte), [32]byte(pubKey)
	ok := extra25519.PublicKeyToCurve25519(keyCurve25519, &pubKeyBytes)
	if !ok {
		return nil
	}
	return keyCurve25519
}

func (pubKey PubKeyEd25519) String() string {
	return fmt.Sprintf("%X", pubKey[:])
}

func (pubKey PubKeyEd25519) Equals(other PubKey) bool {
	if otherEd, ok := other.(PubKeyEd25519); ok {
		return bytes.Equal(pubKey[:], otherEd[:])
	} else {
		return false
	}
}

var _ PubKey = PubKeySecp256k1{}

type PubKeySecp256k1 [33]byte

func (pubKey PubKeySecp256k1) Address() Address {
	hasherSHA256 := sha256.New()
	hasherSHA256.Write(pubKey[:])
	sha := hasherSHA256.Sum(nil)

	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha)
	return Address(hasherRIPEMD160.Sum(nil))
}

func (pubKey PubKeySecp256k1) Bytes() []byte {
	bz, err := cdc.MarshalBinaryBare(pubKey)
	if err != nil {
		panic(err)
	}
	return bz
}

func (pubKey PubKeySecp256k1) VerifyBytes(msg []byte, sig_ Signature) bool {

	sig, ok := sig_.(SignatureSecp256k1)
	if !ok {
		return false
	}

	pub__, err := secp256k1.ParsePubKey(pubKey[:], secp256k1.S256())
	if err != nil {
		return false
	}
	sig__, err := secp256k1.ParseDERSignature(sig[:], secp256k1.S256())
	if err != nil {
		return false
	}
	return sig__.Verify(Sha256(msg), pub__)
}

func (pubKey PubKeySecp256k1) String() string {
	return fmt.Sprintf("PubKeySecp256k1{%X}", pubKey[:])
}

func (pubKey PubKeySecp256k1) Equals(other PubKey) bool {
	if otherSecp, ok := other.(PubKeySecp256k1); ok {
		return bytes.Equal(pubKey[:], otherSecp[:])
	} else {
		return false
	}
}
