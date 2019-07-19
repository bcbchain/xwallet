package rpc

import (
	"bcXwallet/common"
	"blockchain/algorithm"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/tx2"
	"blockchain/types"
	"common/bignumber_v1.0"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/go-crypto"
	"strings"
)

const (
	pattern = "^[a-zA-Z0-9_@.-]{1,40}$"
)

var cdc = amino.NewCodec()

func walletCreate(name, password string) (result *WalletCreateResult, err error) {
	logger := common.GetLogger()

	acct, accessKey, err := newAccount(name, password)
	if err != nil {
		logger.Info(err.Error())
		return
	}

	err = acct.Save(accessKey)
	if err != nil {
		return
	}

	result = new(WalletCreateResult)

	result.AccessKey = base58.Encode(accessKey)
	result.WalletAddress = acct.Address

	return
}

func walletExport(name, password, accessKey string, plainText bool) (result *WalletExportResult, err error) {

	accessKeyBytes := base58.Decode(accessKey)

	acct, err := db.Account(name, accessKeyBytes)
	if err != nil {
		fmt.Println("Failed to get account from database, please name or accessKey.")
		return
	}

	priKeyBytes, err := algorithm.DecryptWithPassword(acct.EncPrivateKey, []byte(password), accessKeyBytes)
	if err != nil {
		fmt.Println("Decrypt Password failed, please check password.")
		return
	}

	result = new(WalletExportResult)
	result.WalletAddress = acct.Address
	if plainText {
		result.PrivateKey = hex.EncodeToString(priKeyBytes)
	} else {
		result.PrivateKey = hex.EncodeToString(acct.EncPrivateKey)
	}

	return
}

func walletImport(name, privateKey, password, accessKey string, plainText bool) (result *WalletImportResult, err error) {

	result = new(WalletImportResult)

	isExist, _ := db.IsExist(name)
	if isExist {
		return nil, errors.New("The account of " + name + " is already exist!")
	}

	priKeyWithPWBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return result, errors.New(" The format of privateKey is wrong")
	}

	var priKeyBytes []byte
	var accessKeyBytes []byte
	if plainText {
		accessKeyBytes = crypto.CRandBytes(32)
		priKeyBytes = priKeyWithPWBytes
		accessKey = base58.Encode(accessKeyBytes)
	} else {
		if accessKey == "" {
			return nil, errors.New(" The accessKey can not be empty")
		}
		accessKeyBytes = base58.Decode(accessKey)

		priKeyBytes, err = algorithm.DecryptWithPassword(priKeyWithPWBytes, []byte(password), accessKeyBytes)
		if err != nil {
			return
		}
	}

	encPrivateKey := algorithm.EncryptWithPassword(priKeyBytes, []byte(password), accessKeyBytes)

	cfg := common.GetConfig()
	priKey := crypto.PrivKeyEd25519FromBytes(priKeyBytes)
	address := priKey.PubKey().Address(cfg.ChainID)

	acct := Account{
		Name:          name,
		Address:       address,
		EncPrivateKey: encPrivateKey,
		PrivateKey:    priKeyBytes,
	}

	err = acct.Save(accessKeyBytes)
	if err != nil {
		return
	}

	result.WalletAddress = address
	result.AccessKey = accessKey

	return
}

func walletList(pageNum uint64) (*WalletListResult, error) {
	wallet := new(WalletListResult)
	wallet.WalletList = make([]WalletItem, 0)

	var err error
	wallet.Total, err = db.AccountNumber()
	if err != nil {
		return nil, err
	}

	walletList, err := db.WalletList(pageNum)
	if err != nil {
		return nil, err
	}

	for _, walletItem := range walletList {
		item := WalletItem{}
		info := strings.Split(walletItem, "#")
		item.Name = info[0]
		item.WalletAddress = info[1]
		wallet.WalletList = append(wallet.WalletList, item)
	}

	return wallet, err
}

func transfer(name, accessKey string, gasLimit uint64, walletParams TransferParam) (result *TransferResult, err error) {

	config := common.GetConfig()
	result = new(TransferResult)

	accessKeyBytes := base58.Decode(accessKey)

	acct, err := db.Account(name, accessKeyBytes)
	if err != nil {
		return
	}

	// 获取nonce
	nonceResult, err := nonce(acct.Address)
	if err != nil {
		return
	}

	value := bignumber.NewNumberString(walletParams.Value)

	var txStr string

	if config.ChainVersion == "1" {
		txStr, err = PackAndSignTx(nonceResult.Nonce, gasLimit, walletParams.Note, walletParams.SmcAddress, walletParams.To, value.Bytes(), name, accessKey)
		if err != nil {
			return nil, err
		}
	} else if config.ChainVersion == "2" {
		var method uint32 = 0x44D8CA60
		v := bn.NewNumberStringBase(walletParams.Value, 10)
		V2Paramss := []interface{}{walletParams.To, v}
		prikey := "0x" + hex.EncodeToString(acct.PrivateKey)

		txStr = GenerateTx(walletParams.SmcAddress, method, V2Paramss, nonceResult.Nonce, int64(gasLimit), walletParams.Note, prikey)
	} else {
		return nil, errors.New("ChainVersion wrong, please check!")
	}

	commitResult, err := common.DoHttpRequestAndParse(config.NodeAddrSlice, txStr)
	if err != nil {
		return
	}

	if commitResult.CheckTx.Code != 200 {
		result.Log = commitResult.CheckTx.Log
		result.Code = commitResult.CheckTx.Code
	} else {
		result.Log = commitResult.DeliverTx.Log
		result.Code = commitResult.DeliverTx.Code
	}
	result.Fee = commitResult.DeliverTx.Fee
	result.Height = commitResult.Height
	result.TxHash = "0x" + hex.EncodeToString(commitResult.Hash)

	return
}

func walletTransferOffline(name, accessKey string, gasLimit uint64, walletParams TransferOfflineParam) (result *TransferOfflineResult, err error) {

	config := common.GetConfig()
	value := bignumber.NewNumberString(walletParams.Value)

	accessKeyBytes := base58.Decode(accessKey)

	acct, err := db.Account(name, accessKeyBytes)
	if err != nil {
		return
	}

	var txStr string

	if config.ChainVersion == "1" {
		txStr, err = PackAndSignTx(walletParams.Nonce, gasLimit, walletParams.Note, walletParams.SmcAddress, walletParams.To, value.Bytes(), name, accessKey)
		if err != nil {
			return nil, err
		}
	} else if config.ChainVersion == "2" {
		var method uint32 = 0x44D8CA60
		v := bn.NewNumberStringBase(walletParams.Value, 10)
		V2Paramss := []interface{}{walletParams.To, v}
		prikey := "0x" + hex.EncodeToString(acct.PrivateKey)

		txStr = GenerateTx(walletParams.SmcAddress, method, V2Paramss, walletParams.Nonce, int64(gasLimit), walletParams.Note, prikey)
	} else {
		return nil, errors.New("ChainVersion wrong, please check!")
	}

	result = new(TransferOfflineResult)
	result.Tx = txStr

	return
}

//GenerateTx generate tx with one contract method request
func GenerateTx(contract types.Address, method uint32, V2Paramss []interface{}, nonce uint64, gaslimit int64, note string, privKey string) string {
	items := tx2.WrapInvokeParams(V2Paramss...)
	message := types.Message{
		Contract: contract,
		MethodID: method,
		Items:    items,
	}
	payload := tx2.WrapPayload(nonce, gaslimit, note, message)
	return tx2.WrapTx(payload, privKey)
}
