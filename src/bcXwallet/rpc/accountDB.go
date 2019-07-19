package rpc

import (
	"bcXwallet/common"
	"blockchain/algorithm"
	"common/bcdb"
	"errors"
	"fmt"
	"path/filepath"
)

type DB struct {
	*bcdb.GILevelDB
}

const (
	countOfOnePage = 1000
)

var (
	db     DB
	dbName = "account"
)

func keyOfAccountNumber() []byte {
	return []byte("/bcbXWallet/accountNumber")
}

func keyOfWalletList(pageNumber uint64) []byte {
	return []byte(fmt.Sprintf("/bcbXWallet/walletList/page%d", pageNumber))
}

// Init DB
func InitDB() error {
	var err error

	dbPath := absolutePath(common.GetConfig().KeyStorePath)

	db.GILevelDB, err = bcdb.OpenDB(dbPath, "", "")

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

// IsExist - get true of account is exist or not
func (db *DB) IsExist(name string) (bool, error) {

	//获取账户信息
	acctBytes, err := db.Get([]byte(name))
	if err != nil {
		return false, err
	}

	if len(acctBytes) == 0 {
		return false, errors.New("account does not exist")
	}

	return true, nil
}

// Account - get account info with name and accessKey
func (db *DB) Account(name string, accessKey []byte) (*Account, error) {

	//获取账户信息
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

	//存储总的钱包数
	acctNumber++
	jsonCount, err := cdc.MarshalJSON(&acctNumber)
	if err != nil {
		return err
	}

	// batch set
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
