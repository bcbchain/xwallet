package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/hex"
	"io"
	"sync"

	. "github.com/tendermint/tmlibs/common"
)

var gRandInfo *randInfo

func init() {
	gRandInfo = &randInfo{}
	gRandInfo.MixEntropy(randBytes(32))
}

func MixEntropy(seedBytes []byte) {
	gRandInfo.MixEntropy(seedBytes)
}

func randBytes(numBytes int) []byte {
	b := make([]byte, numBytes)
	_, err := crand.Read(b)
	if err != nil {
		PanicCrisis(err)
	}
	return b
}

func CRandBytes(numBytes int) []byte {
	b := make([]byte, numBytes)
	_, err := gRandInfo.Read(b)
	if err != nil {
		PanicCrisis(err)
	}
	return b
}

func CRandHex(numDigits int) string {
	return hex.EncodeToString(CRandBytes(numDigits / 2))
}

func CReader() io.Reader {
	return gRandInfo
}

type randInfo struct {
	mtx		sync.Mutex
	seedBytes	[32]byte
	cipherAES256	cipher.Block
	streamAES256	cipher.Stream
	reader		io.Reader
}

func (ri *randInfo) MixEntropy(seedBytes []byte) {
	ri.mtx.Lock()
	defer ri.mtx.Unlock()

	hashBytes := Sha256(seedBytes)
	hashBytes32 := [32]byte{}
	copy(hashBytes32[:], hashBytes)
	ri.seedBytes = xorBytes32(ri.seedBytes, hashBytes32)

	var err error
	ri.cipherAES256, err = aes.NewCipher(ri.seedBytes[:])
	if err != nil {
		PanicSanity("Error creating AES256 cipher: " + err.Error())
	}

	ri.streamAES256 = cipher.NewCTR(ri.cipherAES256, randBytes(aes.BlockSize))

	ri.reader = &cipher.StreamReader{S: ri.streamAES256, R: crand.Reader}
}

func (ri *randInfo) Read(b []byte) (n int, err error) {
	ri.mtx.Lock()
	defer ri.mtx.Unlock()
	return ri.reader.Read(b)
}

func xorBytes32(bytesA [32]byte, bytesB [32]byte) (res [32]byte) {
	for i, b := range bytesA {
		res[i] = b ^ bytesB[i]
	}
	return res
}
