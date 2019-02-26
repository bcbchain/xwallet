package algorithm

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"math/big"
	"strconv"
	"strings"

	"bcbchain.io/smc"
	"github.com/btcsuite/btcutil/base58"
	"github.com/tendermint/go-crypto"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	var bytesBuffer *bytes.Buffer
	if len(b) < 8 {
		bytesBuffer = bytes.NewBuffer(make([]byte, 8-len(b)))
		bytesBuffer.Write(b)
	} else {
		bytesBuffer = bytes.NewBuffer(b)
	}

	var tmp int64
	err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	if err != nil {
		panic(err)
	}
	return int(tmp)
}

func CalcAddressFromCdcPubKey(pubKey []byte) crypto.Address {
	crptPubKey, err := crypto.PubKeyFromBytes(pubKey)
	if err != nil {
		panic(err)
	}
	return crptPubKey.Address()
}

func CheckAddress(chainID string, addr smc.Address) error {
	if !strings.HasPrefix(addr, chainID) {
		return errors.New("Address chainid is error!")
	}
	base58Addr := strings.Replace(addr, chainID, "", 1)
	addrData := base58.Decode(base58Addr)
	len := len(addrData)
	if len == 0 {
		return errors.New("Base58Addr parse error!")
	}

	hasher := ripemd160.New()
	hasher.Write(addrData[:len-4])
	md := hasher.Sum(nil)

	if bytes.Compare(md[:4], addrData[len-4:]) != 0 {
		return errors.New("Address checksum is error!")
	}

	return nil
}

func CalcContractAddress(chainID string, ownerAddr crypto.Address, contractName, version string) crypto.Address {
	hasherSHA3256 := sha3.New256()
	hasherSHA3256.Write([]byte(chainID))
	hasherSHA3256.Write([]byte(contractName))
	hasherSHA3256.Write([]byte(version))
	hasherSHA3256.Write([]byte(ownerAddr))
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

	return string(chainID) + base58.Encode(addr)
}

func CalcUdcHash(nonce uint64, token, owner crypto.Address, value big.Int, matureDate string) crypto.Hash {
	hasherSHA3256 := sha3.New256()
	hasherSHA3256.Write([]byte(strconv.FormatUint(nonce, 10)))
	hasherSHA3256.Write([]byte(token))
	hasherSHA3256.Write([]byte(owner))
	hasherSHA3256.Write(value.Bytes())
	hasherSHA3256.Write([]byte(matureDate))

	return hasherSHA3256.Sum(nil)
}

func CalcMethodId(protoType string) []byte {

	d := sha3.New256()
	d.Write([]byte(protoType))
	b := d.Sum(nil)
	return b[0:4]
}

func CalcCodeHash(code string) []byte {
	hasherSHA3256 := sha3.New256()
	hasherSHA3256.Write([]byte(code))
	return hasherSHA3256.Sum(nil)
}

func SHA3256(datas ...[]byte) []byte {

	hasherSHA3256 := sha3.New256()
	for _, data := range datas {
		hasherSHA3256.Write(data)
	}
	return hasherSHA3256.Sum(nil)
}

func GenSymmetrickeyFromPassword(password, keyword []byte) []byte {
	hasherSHA3256 := sha3.New256()
	hasherSHA3256.Write([]byte("7g$2HJJhh&&!^&!nNN8812MN31^%!@%*^&*&((&*152"))
	hasherSHA3256.Write(password[:])
	if keyword != nil {
		hasherSHA3256.Write(keyword[:])
	}
	sha := hasherSHA3256.Sum(nil)
	digest := md5.New()
	digest.Write(sha)
	return digest.Sum(nil)
}

func EncryptWithPassword(data, password, keyword []byte) []byte {
	if data == nil {
		return nil
	}
	key := GenSymmetrickeyFromPassword(password, keyword)
	enc, _ := aes.NewCipher(key)
	blockSize := enc.BlockSize()
	dat := make([]byte, len(data)+8)
	copy(dat, []byte{0x2e, 0x77, 0x61, 0x6c})
	copy(dat[4:], IntToBytes(len(data)))
	copy(dat[8:], data)
	if n := len(dat) % blockSize; n != 0 {
		m := blockSize - n
		for i := 0; i < m; i++ {
			dat = append(dat, 0)
		}
	}
	for i := 0; i < len(dat)/blockSize; i++ {
		enc.Encrypt(dat[i*blockSize:(i+1)*blockSize], dat[i*blockSize:(i+1)*blockSize])
	}
	return dat
}

func DecryptWithPassword(data, password, keyword []byte) ([]byte, error) {
	if data == nil || len(data) == 0 {
		return nil, errors.New("Cannot decrypt empty data")
	}

	key := GenSymmetrickeyFromPassword(password, keyword)
	dec, _ := aes.NewCipher(key)
	blockSize := dec.BlockSize()

	if len(data)%blockSize != 0 {
		return nil, errors.New("Decrypt data is not an integral multiple of a block")
	}

	dat := make([]byte, len(data))
	copy(dat, data)
	for i := 0; i < len(dat)/blockSize; i++ {
		dec.Decrypt(dat[i*blockSize:(i+1)*blockSize], dat[i*blockSize:(i+1)*blockSize])
	}
	if len(dat) < 8 {
		return nil, errors.New("Decrypt data failed!")
	}
	mac := make([]byte, 4)
	copy(mac, dat[:4])
	if bytes.Compare(mac, []byte{0x2e, 0x77, 0x61, 0x6c}) != 0 {
		return nil, errors.New("Decrypt data failed!")
	}
	size := BytesToInt(dat[4:8])
	return dat[8 : 8+size], nil
}
