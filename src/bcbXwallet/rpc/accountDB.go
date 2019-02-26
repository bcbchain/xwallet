package rpc

import (
	"bcbXwallet/common"
	"errors"
	"fmt"
	"bcbchain.io/algorithm"
	"bcbchain.io/bcdb"
	"path/filepath"
)

type DB struct {
	*bcdb.LevelDB
}

const (
	countOfOnePage = 1000
)

var (
	db	DB
	dbName	= "account"
)

func keyOfAccountNumber() []byte {
	return []byte("/bcbXWallet/accountNumber")
}

func keyOfWalletList(pageNumber uint64) []byte {
	return []byte(fmt.Sprintf("/bcbXWallet/walletList/page%d", pageNumber))
}

func InitDB() error {
	var err error

	dbPath := absolutePath(common.GetConfig().KeyStorePath)

	db.LevelDB, err = bcdb.OpenDB(dbPath, "", "")

	return err
}

func absolutePath(path string) string {
	if filepath.IsAbs(path) {
		path = filepath.Join(path, dbName)
	} else {
		dir, err := common.CurrentDirectory()
		if err != nil {
			panic(err)
		}
		path = filepath.Join(dir, path, dbName)

	}

	return path
}

func (db *DB) IsExist(name string) (bool, error) {

	acctBytes, err := db.Get([]byte(name))
	if err != nil {
		return false, err
	}

	if len(acctBytes) == 0 {
		return false, errors.New("account does not exist")
	}

	return true, nil
}

func (db *DB) Account(name string, accessKey []byte) (*Account, error) {

	acctBytes, err := db.Get([]byte(name))
	if err != nil {
		return nil, err
	}
	if len(acctBytes) == 0 {
		return nil, errors.New("Account does not exist ")
	}

	jsonBytes, err := algorithm.DecryptWithPassword(acctBytes, nil, accessKey)
	if err != nil {
		return nil, fmt.Errorf("The accessKey is wrong ")
	}

	var acct = Account{}
	err = cdc.UnmarshalJSON(jsonBytes, &acct)
	if err != nil {
		return nil, errors.New("UnmarshalJSON error : " + err.Error())
	}

	return &acct, nil
}

func (db *DB) SetAccount(acct *Account, accessKey []byte) error {

	jsonBytes, err := cdc.MarshalJSON(acct)
	if err != nil {
		return err
	}
	walBytes := algorithm.EncryptWithPassword(jsonBytes, nil, accessKey)

	acctNumber, err := db.AccountNumber()
	if err != nil {
		panic(err)
	}

	pageNumber := acctNumber/countOfOnePage + 1

	walletList, err := db.WalletList(pageNumber)
	if err != nil {
		return err
	}
	walletList = append(walletList, acct.Name+"#"+acct.Address)
	jsonList, err := cdc.MarshalJSON(&walletList)
	if err != nil {
		return err
	}

	acctNumber++
	jsonCount, err := cdc.MarshalJSON(&acctNumber)
	if err != nil {
		return err
	}

	dbBatch := db.NewBatch()
	dbBatch.Set([]byte(acct.Name), walBytes)
	dbBatch.Set([]byte(keyOfWalletList(pageNumber)), jsonList)
	dbBatch.Set([]byte(keyOfAccountNumber()), jsonCount)

	return dbBatch.Commit()
}

func (db *DB) WalletList(pageNumber uint64) ([]string, error) {

	bytes, err := db.Get(keyOfWalletList(pageNumber))
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, nil
	}

	list := make([]string, 0)
	err = cdc.UnmarshalJSON(bytes, &list)

	return list, err
}

func (db *DB) AccountNumber() (uint64, error) {

	bytes, err := db.Get(keyOfAccountNumber())
	if err != nil {
		return 0, err
	}

	if len(bytes) == 0 {
		return 0, nil
	}

	number := uint64(0)
	err = cdc.UnmarshalJSON(bytes, &number)

	return number, err
}
