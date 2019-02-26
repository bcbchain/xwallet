package conn

import (
	"bytes"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/ripemd160"

	"github.com/tendermint/go-crypto"
	cmn "github.com/tendermint/tmlibs/common"
)

const dataLenSize = 4
const dataMaxSize = 1024
const totalFrameSize = dataMaxSize + dataLenSize
const sealedFrameSize = totalFrameSize + secretbox.Overhead

type SecretConnection struct {
	conn		io.ReadWriteCloser
	recvBuffer	[]byte
	recvNonce	*[24]byte
	sendNonce	*[24]byte
	remPubKey	crypto.PubKey
	shrSecret	*[32]byte
}

func MakeSecretConnection(conn io.ReadWriteCloser, locPrivKey crypto.PrivKey) (*SecretConnection, error) {

	locPubKey := locPrivKey.PubKey()

	locEphPub, locEphPriv := genEphKeys()

	remEphPub, err := shareEphPubKey(conn, locEphPub)
	if err != nil {
		return nil, err
	}

	shrSecret := computeSharedSecret(remEphPub, locEphPriv)

	loEphPub, hiEphPub := sort32(locEphPub, remEphPub)

	locIsLeast := bytes.Equal(locEphPub[:], loEphPub[:])

	recvNonce, sendNonce := genNonces(loEphPub, hiEphPub, locIsLeast)

	challenge := genChallenge(loEphPub, hiEphPub)

	sc := &SecretConnection{
		conn:		conn,
		recvBuffer:	nil,
		recvNonce:	recvNonce,
		sendNonce:	sendNonce,
		shrSecret:	shrSecret,
	}

	locSignature := signChallenge(challenge, locPrivKey)

	authSigMsg, err := shareAuthSignature(sc, locPubKey, locSignature)
	if err != nil {
		return nil, err
	}
	remPubKey, remSignature := authSigMsg.Key, authSigMsg.Sig
	if !remPubKey.VerifyBytes(challenge[:], remSignature) {
		return nil, errors.New("Challenge verification failed")
	}

	sc.remPubKey = remPubKey
	return sc, nil
}

func (sc *SecretConnection) RemotePubKey() crypto.PubKey {
	return sc.remPubKey
}

func (sc *SecretConnection) Write(data []byte) (n int, err error) {
	for 0 < len(data) {
		var frame = make([]byte, totalFrameSize)
		var chunk []byte
		if dataMaxSize < len(data) {
			chunk = data[:dataMaxSize]
			data = data[dataMaxSize:]
		} else {
			chunk = data
			data = nil
		}
		chunkLength := len(chunk)
		binary.BigEndian.PutUint32(frame, uint32(chunkLength))
		copy(frame[dataLenSize:], chunk)

		var sealedFrame = make([]byte, sealedFrameSize)
		secretbox.Seal(sealedFrame[:0], frame, sc.sendNonce, sc.shrSecret)

		incr2Nonce(sc.sendNonce)

		_, err := sc.conn.Write(sealedFrame)
		if err != nil {
			return n, err
		}
		n += len(chunk)
	}
	return
}

func (sc *SecretConnection) Read(data []byte) (n int, err error) {
	if 0 < len(sc.recvBuffer) {
		n = copy(data, sc.recvBuffer)
		sc.recvBuffer = sc.recvBuffer[n:]
		return
	}

	sealedFrame := make([]byte, sealedFrameSize)
	_, err = io.ReadFull(sc.conn, sealedFrame)
	if err != nil {
		return
	}

	var frame = make([]byte, totalFrameSize)

	_, ok := secretbox.Open(frame[:0], sealedFrame, sc.recvNonce, sc.shrSecret)
	if !ok {
		return n, errors.New("Failed to decrypt SecretConnection")
	}
	incr2Nonce(sc.recvNonce)

	var chunkLength = binary.BigEndian.Uint32(frame)
	if chunkLength > dataMaxSize {
		return 0, errors.New("chunkLength is greater than dataMaxSize")
	}
	var chunk = frame[dataLenSize : dataLenSize+chunkLength]

	n = copy(data, chunk)
	sc.recvBuffer = chunk[n:]
	return
}

func (sc *SecretConnection) Close() error			{ return sc.conn.Close() }
func (sc *SecretConnection) LocalAddr() net.Addr		{ return sc.conn.(net.Conn).LocalAddr() }
func (sc *SecretConnection) RemoteAddr() net.Addr		{ return sc.conn.(net.Conn).RemoteAddr() }
func (sc *SecretConnection) SetDeadline(t time.Time) error	{ return sc.conn.(net.Conn).SetDeadline(t) }
func (sc *SecretConnection) SetReadDeadline(t time.Time) error {
	return sc.conn.(net.Conn).SetReadDeadline(t)
}
func (sc *SecretConnection) SetWriteDeadline(t time.Time) error {
	return sc.conn.(net.Conn).SetWriteDeadline(t)
}

func genEphKeys() (ephPub, ephPriv *[32]byte) {
	var err error
	ephPub, ephPriv, err = box.GenerateKey(crand.Reader)
	if err != nil {
		panic("Could not generate ephemeral keypairs")
	}
	return
}

func shareEphPubKey(conn io.ReadWriteCloser, locEphPub *[32]byte) (remEphPub *[32]byte, err error) {

	var trs, _ = cmn.Parallel(
		func(_ int) (val interface{}, err error, abort bool) {
			var _, err1 = cdc.MarshalBinaryWriter(conn, locEphPub)
			if err1 != nil {
				return nil, err1, true
			} else {
				return nil, nil, false
			}
		},
		func(_ int) (val interface{}, err error, abort bool) {
			var _remEphPub [32]byte
			var _, err2 = cdc.UnmarshalBinaryReader(conn, &_remEphPub, 1024*1024)
			if err2 != nil {
				return nil, err2, true
			} else {
				return _remEphPub, nil, false
			}
		},
	)

	if trs.FirstError() != nil {
		err = trs.FirstError()
		return
	}

	var _remEphPub = trs.FirstValue().([32]byte)
	return &_remEphPub, nil
}

func computeSharedSecret(remPubKey, locPrivKey *[32]byte) (shrSecret *[32]byte) {
	shrSecret = new([32]byte)
	box.Precompute(shrSecret, remPubKey, locPrivKey)
	return
}

func sort32(foo, bar *[32]byte) (lo, hi *[32]byte) {
	if bytes.Compare(foo[:], bar[:]) < 0 {
		lo = foo
		hi = bar
	} else {
		lo = bar
		hi = foo
	}
	return
}

func genNonces(loPubKey, hiPubKey *[32]byte, locIsLo bool) (recvNonce, sendNonce *[24]byte) {
	nonce1 := hash24(append(loPubKey[:], hiPubKey[:]...))
	nonce2 := new([24]byte)
	copy(nonce2[:], nonce1[:])
	nonce2[len(nonce2)-1] ^= 0x01
	if locIsLo {
		recvNonce = nonce1
		sendNonce = nonce2
	} else {
		recvNonce = nonce2
		sendNonce = nonce1
	}
	return
}

func genChallenge(loPubKey, hiPubKey *[32]byte) (challenge *[32]byte) {
	return hash32(append(loPubKey[:], hiPubKey[:]...))
}

func signChallenge(challenge *[32]byte, locPrivKey crypto.PrivKey) (signature crypto.Signature) {
	signature = locPrivKey.Sign(challenge[:])
	return
}

type authSigMessage struct {
	Key	crypto.PubKey
	Sig	crypto.Signature
}

func shareAuthSignature(sc *SecretConnection, pubKey crypto.PubKey, signature crypto.Signature) (recvMsg authSigMessage, err error) {

	var trs, _ = cmn.Parallel(
		func(_ int) (val interface{}, err error, abort bool) {
			var _, err1 = cdc.MarshalBinaryWriter(sc, authSigMessage{pubKey, signature})
			if err1 != nil {
				return nil, err1, true
			} else {
				return nil, nil, false
			}
		},
		func(_ int) (val interface{}, err error, abort bool) {
			var _recvMsg authSigMessage
			var _, err2 = cdc.UnmarshalBinaryReader(sc, &_recvMsg, 1024*1024)
			if err2 != nil {
				return nil, err2, true
			} else {
				return _recvMsg, nil, false
			}
		},
	)

	if trs.FirstError() != nil {
		err = trs.FirstError()
		return
	}

	var _recvMsg = trs.FirstValue().(authSigMessage)
	return _recvMsg, nil
}

func hash32(input []byte) (res *[32]byte) {
	hasher := sha256.New()
	hasher.Write(input)
	resSlice := hasher.Sum(nil)
	res = new([32]byte)
	copy(res[:], resSlice)
	return
}

func hash24(input []byte) (res *[24]byte) {
	hasher := ripemd160.New()
	hasher.Write(input)
	resSlice := hasher.Sum(nil)
	res = new([24]byte)
	copy(res[:], resSlice)
	return
}

func incr2Nonce(nonce *[24]byte) {
	incrNonce(nonce)
	incrNonce(nonce)
}

func incrNonce(nonce *[24]byte) {
	for i := 23; 0 <= i; i-- {
		nonce[i]++
		if nonce[i] != 0 {
			return
		}
	}
}
