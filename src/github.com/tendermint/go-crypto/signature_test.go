package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/ed25519"
	amino "github.com/tendermint/go-amino"
)

func TestSignAndValidateEd25519(t *testing.T) {

	privKey := GenPrivKeyEd25519()
	pubKey := privKey.PubKey()

	msg := CRandBytes(128)
	sig := privKey.Sign(msg)

	assert.True(t, pubKey.VerifyBytes(msg, sig))

	sigEd := sig.(SignatureEd25519)
	sigEd[7] ^= byte(0x01)
	sig = sigEd

	assert.False(t, pubKey.VerifyBytes(msg, sig))
}

func TestSignAndValidateSecp256k1(t *testing.T) {
	privKey := GenPrivKeySecp256k1()
	pubKey := privKey.PubKey()

	msg := CRandBytes(128)
	sig := privKey.Sign(msg)

	assert.True(t, pubKey.VerifyBytes(msg, sig))

	sigEd := sig.(SignatureSecp256k1)
	sigEd[3] ^= byte(0x01)
	sig = sigEd

	assert.False(t, pubKey.VerifyBytes(msg, sig))
}

func TestSignatureEncodings(t *testing.T) {
	cases := []struct {
		privKey		PrivKey
		sigSize		int
		sigPrefix	amino.PrefixBytes
	}{
		{
			privKey:	GenPrivKeyEd25519(),
			sigSize:	ed25519.SignatureSize,
			sigPrefix:	[4]byte{0x3d, 0xa1, 0xdb, 0x2a},
		},
		{
			privKey:	GenPrivKeySecp256k1(),
			sigSize:	0,
			sigPrefix:	[4]byte{0x16, 0xe1, 0xfe, 0xea},
		},
	}

	for _, tc := range cases {

		pubKey := tc.privKey.PubKey()

		msg := CRandBytes(128)
		sig := tc.privKey.Sign(msg)

		bin, err := cdc.MarshalBinaryBare(sig)
		require.Nil(t, err, "%+v", err)
		if tc.sigSize != 0 {

			assert.Equal(t, tc.sigSize+amino.PrefixBytesLen+1, len(bin))
		}
		assert.Equal(t, tc.sigPrefix[:], bin[0:amino.PrefixBytesLen])

		sig2 := Signature(nil)
		err = cdc.UnmarshalBinaryBare(bin, &sig2)
		require.Nil(t, err, "%+v", err)
		assert.EqualValues(t, sig, sig2)
		assert.True(t, pubKey.VerifyBytes(msg, sig2))

	}
}
