package words

import (
	"testing"

	asrt "github.com/stretchr/testify/assert"

	cmn "github.com/tendermint/tmlibs/common"
)

var codecs = []ECC{
	NewIBMCRC16(),
	NewSCSICRC16(),
	NewCCITTCRC16(),
	NewIEEECRC32(),
	NewCastagnoliCRC32(),
	NewKoopmanCRC32(),
	NewISOCRC64(),
	NewECMACRC64(),
}

func TestECCPasses(t *testing.T) {
	assert := asrt.New(t)

	checks := append(codecs, NoECC{})

	for _, check := range checks {
		for i := 0; i < 2000; i++ {
			numBytes := cmn.RandInt()%60 + 1
			data := cmn.RandBytes(numBytes)

			checked := check.AddECC(data)
			res, err := check.CheckECC(checked)
			if assert.Nil(err, "%#v: %+v", check, err) {
				assert.Equal(data, res, "%v", check)
			}
		}
	}
}

func TestECCFails(t *testing.T) {
	assert := asrt.New(t)

	checks := codecs
	attempts := 2000

	for _, check := range checks {
		failed := 0
		for i := 0; i < attempts; i++ {
			numBytes := cmn.RandInt()%60 + 1
			data := cmn.RandBytes(numBytes)
			_, err := check.CheckECC(data)
			if err != nil {
				failed += 1
			}
		}

		assert.InDelta(attempts, failed, 1, "%v", check)
	}
}
