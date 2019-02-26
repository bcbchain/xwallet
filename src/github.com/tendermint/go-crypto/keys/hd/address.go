package hd

import (
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"math/big"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

func ComputeBTCAddress(pubKeyHex string, chainCodeHex string, path string, index int32) string {
	pubKeyBytes := DerivePublicKeyForPath(
		HexDecode(pubKeyHex),
		HexDecode(chainCodeHex),
		fmt.Sprintf("%v/%v", path, index),
	)
	return BTCAddrFromPubKeyBytes(pubKeyBytes)
}

func ComputePrivateKey(mprivHex string, chainHex string, path string, index int32) string {
	privKeyBytes := DerivePrivateKeyForPath(
		HexDecode(mprivHex),
		HexDecode(chainHex),
		fmt.Sprintf("%v/%v", path, index),
	)
	return HexEncode(privKeyBytes)
}

func ComputeBTCAddressForPrivKey(privKey string) string {
	pubKeyBytes := PubKeyBytesFromPrivKeyBytes(HexDecode(privKey), true)
	return BTCAddrFromPubKeyBytes(pubKeyBytes)
}

func SignBTCMessage(privKey string, message string, compress bool) string {
	prefixBytes := []byte("Bitcoin Signed Message:\n")
	messageBytes := []byte(message)
	bytes := []byte{}
	bytes = append(bytes, byte(len(prefixBytes)))
	bytes = append(bytes, prefixBytes...)
	bytes = append(bytes, byte(len(messageBytes)))
	bytes = append(bytes, messageBytes...)
	privKeyBytes := HexDecode(privKey)
	x, y := btcec.S256().ScalarBaseMult(privKeyBytes)
	ecdsaPubKey := ecdsa.PublicKey{
		Curve:	btcec.S256(),
		X:	x,
		Y:	y,
	}
	ecdsaPrivKey := &btcec.PrivateKey{
		PublicKey:	ecdsaPubKey,
		D:		new(big.Int).SetBytes(privKeyBytes),
	}
	sigbytes, err := btcec.SignCompact(btcec.S256(), ecdsaPrivKey, CalcHash256(bytes), compress)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(sigbytes)
}

func ComputeMastersFromSeed(seed string) (string, string, string) {
	key, data := []byte("Bitcoin seed"), []byte(seed)
	secret, chain := I64(key, data)
	pubKeyBytes := PubKeyBytesFromPrivKeyBytes(secret, true)
	return HexEncode(pubKeyBytes), HexEncode(secret), HexEncode(chain)
}

func ComputeWIF(privKey string, compress bool) string {
	return WIFFromPrivKeyBytes(HexDecode(privKey), compress)
}

func ComputeBTCTxId(rawTxHex string) string {
	return HexEncode(ReverseBytes(CalcHash256(HexDecode(rawTxHex))))
}

func DerivePrivateKeyForPath(privKeyBytes []byte, chainCode []byte, path string) []byte {
	data := privKeyBytes
	parts := strings.Split(path, "/")
	for _, part := range parts {
		prime := part[len(part)-1:] == "'"

		if prime {
			part = part[:len(part)-1]
		}
		i, err := strconv.Atoi(part)
		if err != nil {
			panic(err)
		}
		if i < 0 {
			panic(errors.New("index too large."))
		}
		data, chainCode = DerivePrivateKey(data, chainCode, uint32(i), prime)

	}
	return data
}

func DerivePublicKeyForPath(pubKeyBytes []byte, chainCode []byte, path string) []byte {
	data := pubKeyBytes
	parts := strings.Split(path, "/")
	for _, part := range parts {
		prime := part[len(part)-1:] == "'"
		if prime {
			panic(errors.New("cannot do a prime derivation from public key"))
		}
		i, err := strconv.Atoi(part)
		if err != nil {
			panic(err)
		}
		if i < 0 {
			panic(errors.New("index too large."))
		}
		data, chainCode = DerivePublicKey(data, chainCode, uint32(i))

	}
	return data
}

func DerivePrivateKey(privKeyBytes []byte, chainCode []byte, index uint32, prime bool) ([]byte, []byte) {
	var data []byte
	if prime {
		index = index | 0x80000000
		data = append([]byte{byte(0)}, privKeyBytes...)
	} else {
		public := PubKeyBytesFromPrivKeyBytes(privKeyBytes, true)
		data = public
	}
	data = append(data, uint32ToBytes(index)...)
	data2, chainCode2 := I64(chainCode, data)
	x := addScalars(privKeyBytes, data2)
	return x, chainCode2
}

func DerivePublicKey(pubKeyBytes []byte, chainCode []byte, index uint32) ([]byte, []byte) {
	data := []byte{}
	data = append(data, pubKeyBytes...)
	data = append(data, uint32ToBytes(index)...)
	data2, chainCode2 := I64(chainCode, data)
	data2p := PubKeyBytesFromPrivKeyBytes(data2, true)
	return addPoints(pubKeyBytes, data2p), chainCode2
}

func addPoints(a []byte, b []byte) []byte {
	ap, err := btcec.ParsePubKey(a, btcec.S256())
	if err != nil {
		panic(err)
	}
	bp, err := btcec.ParsePubKey(b, btcec.S256())
	if err != nil {
		panic(err)
	}
	sumX, sumY := btcec.S256().Add(ap.X, ap.Y, bp.X, bp.Y)
	sum := &btcec.PublicKey{
		Curve:	btcec.S256(),
		X:	sumX,
		Y:	sumY,
	}
	return sum.SerializeCompressed()
}

func addScalars(a []byte, b []byte) []byte {
	aInt := new(big.Int).SetBytes(a)
	bInt := new(big.Int).SetBytes(b)
	sInt := new(big.Int).Add(aInt, bInt)
	x := sInt.Mod(sInt, btcec.S256().N).Bytes()
	x2 := [32]byte{}
	copy(x2[32-len(x):], x)
	return x2[:]
}

func uint32ToBytes(i uint32) []byte {
	b := [4]byte{}
	binary.BigEndian.PutUint32(b[:], i)
	return b[:]
}

func HexEncode(b []byte) string {
	return hex.EncodeToString(b)
}

func HexDecode(str string) []byte {
	b, _ := hex.DecodeString(str)
	return b
}

func I64(key []byte, data []byte) ([]byte, []byte) {
	mac := hmac.New(sha512.New, key)
	mac.Write(data)
	I := mac.Sum(nil)
	return I[:32], I[32:]
}

const (
	btcPrefixPubKeyHash	= byte(0x00)
	btcPrefixPrivKey	= byte(0x80)
)

func BTCAddrFromPubKeyBytes(pubKeyBytes []byte) string {
	versionPrefix := btcPrefixPubKeyHash
	h160 := CalcHash160(pubKeyBytes)
	h160 = append([]byte{versionPrefix}, h160...)
	checksum := CalcHash256(h160)
	b := append(h160, checksum[:4]...)
	return base58.Encode(b)
}

func BTCAddrBytesFromPubKeyBytes(pubKeyBytes []byte) (addrBytes []byte, checksum []byte) {
	versionPrefix := btcPrefixPubKeyHash
	h160 := CalcHash160(pubKeyBytes)
	_h160 := append([]byte{versionPrefix}, h160...)
	checksum = CalcHash256(_h160)[:4]
	return h160, checksum
}

func WIFFromPrivKeyBytes(privKeyBytes []byte, compress bool) string {
	versionPrefix := btcPrefixPrivKey
	bytes := append([]byte{versionPrefix}, privKeyBytes...)
	if compress {
		bytes = append(bytes, byte(1))
	}
	checksum := CalcHash256(bytes)
	bytes = append(bytes, checksum[:4]...)
	return base58.Encode(bytes)
}

func PubKeyBytesFromPrivKeyBytes(privKeyBytes []byte, compress bool) (pubKeyBytes []byte) {
	x, y := btcec.S256().ScalarBaseMult(privKeyBytes)
	pub := &btcec.PublicKey{
		Curve:	btcec.S256(),
		X:	x,
		Y:	y,
	}

	if compress {
		return pub.SerializeCompressed()
	}
	return pub.SerializeUncompressed()
}

func CalcHash(data []byte, hasher hash.Hash) []byte {
	hasher.Write(data)
	return hasher.Sum(nil)
}

func CalcHash160(data []byte) []byte {
	return CalcHash(CalcHash(data, sha256.New()), ripemd160.New())
}

func CalcHash256(data []byte) []byte {
	return CalcHash(CalcHash(data, sha256.New()), sha256.New())
}

func CalcSha512(data []byte) []byte {
	return CalcHash(data, sha512.New())
}

func ReverseBytes(buf []byte) []byte {
	var res []byte
	if len(buf) == 0 {
		return res
	}

	blen := len(buf)
	res = make([]byte, blen)
	mid := blen / 2
	for left := 0; left <= mid; left++ {
		right := blen - left - 1
		res[left] = buf[right]
		res[right] = buf[left]
	}
	return res
}
