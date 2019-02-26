package rpc

import (
	"encoding/binary"
	"errors"
	atm "bcbchain.io/algorithm"
	"bcbchain.io/keys"
	"bcbchain.io/prototype"
	"bcbchain.io/rlp"
	"bcbchain.io/types"
	"github.com/btcsuite/btcutil/base58"
	"github.com/tendermint/go-crypto"
)

type BcbXTransaction struct {
	Nonce		uint64
	GasLimit	uint64
	Note		string
	To		keys.Address
	Data		[]byte
}

func PackAndSignTx(nonce, gasLimit uint64, note, tokenAddress, toAddress string, value []byte, name, accessKey string) (string, error) {

	var mi MethodInfo
	var err error

	methodId := atm.CalcMethodId(prototype.TtTransfer)

	mi.MethodID = binary.BigEndian.Uint32(methodId)

	var itemsBytes = make([][]byte, 0)

	itemsBytes = append(itemsBytes, []byte(toAddress))
	itemsBytes = append(itemsBytes, value)

	mi.ParamData, err = rlp.EncodeToBytes(itemsBytes)
	if err != nil {
		return "", err
	}

	data, err := rlp.EncodeToBytes(mi)
	if err != nil {
		return "", err
	}

	tx1 := NewTransaction(nonce, gasLimit, note, tokenAddress, data)
	return tx1.TxGen(name, accessKey)
}

func NewTransaction(nonce uint64, gasLimit uint64, note string, to keys.Address, data []byte) BcbXTransaction {
	tx := BcbXTransaction{
		Nonce:		nonce,
		GasLimit:	gasLimit,
		Note:		note,
		To:		to,
		Data:		data,
	}
	return tx
}

func (tx *BcbXTransaction) TxGen(name, accessKey string) (string, error) {

	size, r, err := rlp.EncodeToReader(tx)
	if err != nil {
		return "", err
	}
	txBytes := make([]byte, size)
	_, _ = r.Read(txBytes)

	sigInfo, err := SignData(name, accessKey, txBytes)
	if err != nil {
		return "", err
	}

	size, r, err = rlp.EncodeToReader(sigInfo)
	if err != nil {
		return "", err
	}
	sigBytes := make([]byte, size)
	_, _ = r.Read(sigBytes)

	txString := base58.Encode(txBytes)
	sigString := base58.Encode(sigBytes)

	MAC := string(crypto.GetChainId()) + "<tx>"
	Version := "v1"
	SignerNumber := "<1>"

	return MAC + "." + Version + "." + txString + "." + SignerNumber + "." + sigString, nil
}

func SignData(name, accessKey string, data []byte) (*types.Ed25519Sig, error) {
	if name == "" || accessKey == "" {
		return nil, errors.New("user name and accessKey cannot to te empty")
	}

	if name != "" && len(name) > 40 {
		return nil, errors.New("user name length only can be 1-40")
	}
	if len(data) <= 0 {
		return nil, errors.New("user data which wants be signed length needs more than 0")
	}

	accessKeyBytes := base58.Decode(accessKey)

	acct, err2 := db.Account(name, accessKeyBytes)
	if acct == nil {
		return nil, err2
	}

	priKey := crypto.PrivKeyEd25519FromBytes(acct.PrivateKey)
	pubKey := priKey.PubKey()

	sigInfo := types.Ed25519Sig{
		SigType:	"ed25519",
		PubKey:		pubKey.(crypto.PubKeyEd25519),
		SigValue:	priKey.Sign(data).(crypto.SignatureEd25519),
	}

	return &sigInfo, nil
}
