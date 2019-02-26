package rpc

import (
	"errors"
	"bcbchain.io/algorithm"
	"bcbchain.io/keys"
	"github.com/tendermint/go-crypto"
)

type Account struct {
	EncPrivateKey	[]byte		`json:"encPrivateKey"`
	PrivateKey	[]byte		`json:"privateKey"`
	Name		string		`json:"name"`
	Address		keys.Address	`json:"address"`
	Hash		[]byte		`json:"hash"`
}

func newAccount(name, password string) (*Account, []byte, error) {
	isExist, _ := db.IsExist(name)
	if isExist {
		return nil, nil, errors.New("The account of " + name + " is already exist!")
	}

	priKey := crypto.GenPrivKeyEd25519()
	priKeyByte := priKey[:]

	accessKey := crypto.CRandBytes(32)
	priKeyWithPWBytes := algorithm.EncryptWithPassword(priKeyByte, []byte(password), accessKey)

	acct := Account{
		Name:		name,
		Address:	priKey.PubKey().Address(),
		EncPrivateKey:	priKeyWithPWBytes,
		PrivateKey:	priKeyByte,
	}

	return &acct, accessKey, nil
}

func (acct *Account) Save(accessKey []byte) error {
	return db.SetAccount(acct, accessKey)
}
