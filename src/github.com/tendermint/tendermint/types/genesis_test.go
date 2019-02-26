package types

import (
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"
	"testing"
)

func TestGenesisBad(t *testing.T) {

	testCases := [][]byte{
		[]byte{},
		[]byte{1, 1, 1, 1, 1},
		[]byte(`{}`),
		[]byte(`{"chain_id":"mychain"}`),
		[]byte(`{"chain_id":"mychain","validators":[]}`),
		[]byte(`{"chain_id":"mychain","validators":[{}]}`),
		[]byte(`{"chain_id":"mychain","validators":null}`),
		[]byte(`{"chain_id":"mychain"}`),
		[]byte(`{"validators":[{"pub_key":{"type":"AC26791624DE60","value":"AT/+aaL1eB0477Mud9JMm8Sh8BIvOYlPGC9KkIUmFaE="},"power":10,"name":""}]}`),
	}

	for _, testCase := range testCases {
		_, err := GenesisDocFromJSON(testCase)
		assert.Error(t, err, "expected error for empty genDoc json")
	}
}

func TestGenesisGood(t *testing.T) {

	genDocBytes := []byte(`{"genesis_time":"0001-01-01T00:00:00Z","chain_id":"test-chain-QDKdJr","consensus_params":null,"validators":[{"pub_key":{"type":"AC26791624DE60","value":"AT/+aaL1eB0477Mud9JMm8Sh8BIvOYlPGC9KkIUmFaE="},"power":10,"name":""}],"app_hash":"","app_state":{"account_owner": "Bob"}}`)
	_, err := GenesisDocFromJSON(genDocBytes)
	assert.NoError(t, err, "expected no error for good genDoc json")

	baseGenDoc := &GenesisDoc{
		ChainID:	"abc",
		Validators:	[]GenesisValidator{{crypto.GenPrivKeyEd25519().PubKey(), 10, "myval"}},
	}
	genDocBytes, err = cdc.MarshalJSON(baseGenDoc)
	assert.NoError(t, err, "error marshalling genDoc")

	genDoc, err := GenesisDocFromJSON(genDocBytes)
	assert.NoError(t, err, "expected no error for valid genDoc json")
	assert.NotNil(t, genDoc.ConsensusParams, "expected consensus params to be filled in")

	genDocBytes, err = cdc.MarshalJSON(genDoc)
	assert.NoError(t, err, "error marshalling genDoc")
	genDoc, err = GenesisDocFromJSON(genDocBytes)
	assert.NoError(t, err, "expected no error for valid genDoc json")

	genDoc.ConsensusParams.BlockSize.MaxBytes = 0
	genDocBytes, err = cdc.MarshalJSON(genDoc)
	assert.NoError(t, err, "error marshalling genDoc")
	genDoc, err = GenesisDocFromJSON(genDocBytes)
	assert.Error(t, err, "expected error for genDoc json with block size of 0")
}
