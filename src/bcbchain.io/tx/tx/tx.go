package tx

import (
	"bytes"
	"errors"
	"bcbchain.io/kms"
	"bcbchain.io/rlp"
	"bcbchain.io/types"
	"github.com/btcsuite/btcutil/base58"
	"github.com/tendermint/go-crypto"
	"strings"
)

func (tx *Transaction) TxGen(chainID, name, passphrase string) (string, error) {

	size, r, err := rlp.EncodeToReader(tx)
	if err != nil {
		return "", err
	}
	txBytes := make([]byte, size)
	r.Read(txBytes)

	sigInfo, err := kms.SignData(name, passphrase, txBytes)
	if err != nil {
		return "", err
	}

	size, r, err = rlp.EncodeToReader(sigInfo)
	if err != nil {
		return "", err
	}
	sigBytes := make([]byte, size)
	r.Read(sigBytes)

	txString := base58.Encode(txBytes)
	sigString := base58.Encode(sigBytes)

	MAC := string(chainID) + "<tx>"
	Version := "v1"
	SignerNumber := "<1>"

	return MAC + "." + Version + "." + txString + "." + SignerNumber + "." + sigString, nil
}

func (tx *Transaction) TxParse(chainID, txString string) (crypto.Address, error) {
	MAC := chainID + "<tx>"
	Version := "v1"
	SignerNumber := "<1>"
	strs := strings.Split(txString, ".")

	if strs[0] != MAC || strs[1] != Version || strs[3] != SignerNumber {
		return "", errors.New("tx data error")
	}

	txData := base58.Decode(strs[2])
	sigBytes := base58.Decode(strs[4])

	reader := bytes.NewReader(sigBytes)
	var siginfo types.Ed25519Sig
	err := rlp.Decode(reader, &siginfo)
	if err != nil {
		return "", err
	}

	if !siginfo.PubKey.VerifyBytes(txData, siginfo.SigValue) {
		return "", errors.New("verify sig fail")
	}

	reader = bytes.NewReader(txData)
	err = rlp.Decode(reader, tx)
	if err != nil {
		return "", err
	}

	crypto.SetChainId(chainID)
	return siginfo.PubKey.Address(), nil
}

func (qy *Query) QueryDataGen(chainID string, name, passphrase string) (string, error) {

	qyBytes, err := rlp.EncodeToBytes(qy)
	if err != nil {
		return "", err
	}

	sigInfo, err := kms.SignData(name, passphrase, qyBytes)
	if err != nil {
		return "", err
	}

	sigBytes, err := rlp.EncodeToBytes(sigInfo)
	if err != nil {
		return "", err
	}

	qyString := base58.Encode(qyBytes)
	sigString := base58.Encode(sigBytes)

	MAC := chainID + "<qy>"
	Version := "v1"
	SignerNumber := "<1>"

	return MAC + "." + Version + "." + qyString + "." + SignerNumber + "." + sigString, nil
}

func (qy *Query) QueryDataParse(chainID, txString string) (crypto.Address, error) {
	MAC := chainID + "<qy>"
	Version := "v1"
	SignerNumber := "<1>"
	strs := strings.Split(txString, ".")

	if strs[0] != MAC || strs[1] != Version || strs[3] != SignerNumber {
		return "", errors.New("tx data error")
	}

	txData := base58.Decode(strs[2])
	sigBytes := base58.Decode(strs[4])

	reader := bytes.NewReader(sigBytes)
	var siginfo types.Ed25519Sig
	err := rlp.Decode(reader, &siginfo)
	if err != nil {
		return "", err
	}

	if !siginfo.PubKey.VerifyBytes(txData, siginfo.SigValue) {
		return "", errors.New("verify sig fail")
	}
	crypto.SetChainId(chainID)
	siginfo.PubKey.Address()

	reader = bytes.NewReader(txData)
	err = rlp.Decode(reader, qy)
	if err != nil {
		return "", err
	}

	return siginfo.PubKey.Address(), nil
}
