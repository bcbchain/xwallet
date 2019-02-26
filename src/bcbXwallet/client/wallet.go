package client

import (
	rpc2 "bcbXwallet/rpc"
	"encoding/json"
	"fmt"
	"bcbchain.io/client"
	"strconv"
)

func WalletCreate(name, password, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc2.WalletCreateResult)
	_, err = rpc.Call("bcb_walletCreate", map[string]interface{}{"name": name, "password": password}, result)
	if err != nil {
		fmt.Printf("Cannot create wallet, name=%s, password=%s,\n error=%s \n", name, password, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}

func WalletExport(name, password, accessKey, url, plainText string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	bPlainText, err := strconv.ParseBool(plainText)
	if err != nil {
		return
	}

	result := new(rpc2.WalletExportResult)
	_, err = rpc.Call("bcb_walletExport", map[string]interface{}{"name": name, "password": password, "accessKey": accessKey, "plainText": bPlainText}, result)
	if err != nil {
		fmt.Printf("Cannot export wallet, name=%s, password=%s, accessKey=%s, plainText=%v,\n error=%s \n", name, password, accessKey, plainText, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}

func WalletImport(name, privateKey, password, accessKey, url, plainText string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	bPlainText, err := strconv.ParseBool(plainText)
	if err != nil {
		return
	}

	result := new(rpc2.WalletImportResult)
	_, err = rpc.Call("bcb_walletImport", map[string]interface{}{"name": name, "privateKey": privateKey, "password": password, "accessKey": accessKey, "plainText": bPlainText}, result)
	if err != nil {
		fmt.Printf("Cannot import wallet, name=%s, privateKey=%s, password=%s, accessKey=%s, plainText=%v,\n error=%s \n", name, privateKey, password, accessKey, plainText, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}

func WalletList(pageNum uint64, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc2.WalletListResult)
	_, err = rpc.Call("bcb_walletList", map[string]interface{}{"pageNum": pageNum}, result)
	if err != nil {
		fmt.Printf("Cannot list wallet, error=%s \n", err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}

func Transfer(name, accessKey, smcAddress, gasLimit, note, to, value, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	transferParam := rpc2.TransferParam{SmcAddress: smcAddress, GasLimit: gasLimit, Note: note, To: to, Value: value}

	result := new(rpc2.TransferResult)
	_, err = rpc.Call("bcb_transfer", map[string]interface{}{"name": name, "accessKey": accessKey, "walletParams": transferParam}, result)
	if err != nil {
		fmt.Printf("Cannot transfer, name=%s, accessKey=%s, walletParam=%v,\n error=%s \n", name, accessKey, transferParam, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}

func TransferOffline(name, accessKey, smcAddress, gasLimit, note, to, value, nonce, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	uNonce, err := strconv.ParseUint(nonce, 10, 64)
	if err != nil {
		return
	}
	transferParam := rpc2.TransferOfflineParam{SmcAddress: smcAddress, GasLimit: gasLimit, Note: note, Nonce: uNonce, To: to, Value: value}

	result := new(rpc2.TransferOfflineResult)
	_, err = rpc.Call("bcb_transferOffline", map[string]interface{}{"name": name, "accessKey": accessKey, "walletParams": transferParam}, result)
	if err != nil {
		fmt.Printf("Cannot transferOffline, name=%s, accessKey=%s, walletParam=%v,\n error=%s \n", name, accessKey, transferParam, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}
