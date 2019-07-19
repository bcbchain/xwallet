package rpc

import (
	"bcXwallet/common"
	"blockchain/abciapp_v1.0/keys"
	"blockchain/algorithm"
	"bufio"
	"github.com/bgentry/speakeasy"
	"github.com/pkg/errors"
	"github.com/tendermint/go-crypto"
)

// MinPassLength is the minimum acceptable password length
const (
	MinPassLength = 8
	MaxPassLength = 20
)

// ----- account struct -----
type Account struct {
	EncPrivateKey []byte       `json:"encPrivateKey"`
	PrivateKey    []byte       `json:"privateKey"`
	Name          string       `json:"name"`
	Address       keys.Address `json:"address"`
	Hash          []byte       `json:"hash"`
}

func newAccount(name, password string) (*Account, []byte, error) {
	cfg := common.GetConfig()
	isExist, _ := db.IsExist(name)
	if isExist {
		return nil, nil, errors.New("The account of " + name + " is already exist!")
	}

	priKey := crypto.GenPrivKeyEd25519()
	priKeyByte := priKey[:]

	accessKey := crypto.CRandBytes(32)
	priKeyWithPWBytes := algorithm.EncryptWithPassword(priKeyByte, []byte(password), accessKey)

	acct := Account{
		Name:          name,
		Address:       priKey.PubKey().Address(cfg.ChainID),
		EncPrivateKey: priKeyWithPWBytes,
		PrivateKey:    priKeyByte,
	}

	return &acct, accessKey, nil
}

func (acct *Account) Save(accessKey []byte) error {
	return db.SetAccount(acct, accessKey)
}

// GetPassword will prompt for a password one-time (to sign a tx)
// It enforces the password length
func getPassword(prompt string, buf *bufio.Reader) (pass string, err error) {
	pass, err = speakeasy.Ask(prompt)
	if err != nil {
		return "", err
	}
	if len(pass) < MinPassLength {
		return "", errors.Errorf("Password must be at least %d characters", MinPassLength)
	}

	if len(pass) > MaxPassLength {
		return "", errors.Errorf("Password must be at most %d characters", MaxPassLength)
	}

	return pass, nil
}
