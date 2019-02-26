package client

import (
	rpc2 "bcbXwallet/rpc"
	"encoding/json"
	"fmt"
	"bcbchain.io/client"
	"bcbchain.io/keys"
)

func BlockHeight(url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc2.BlockHeightResult)
	_, err = rpc.Call("bcb_blockHeight", map[string]interface{}{}, result)
	if err != nil {
		fmt.Printf("Cannot get block height, error=%s \n", err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}

func Block(height int64, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc2.BlockResult)
	_, err = rpc.Call("bcb_block", map[string]interface{}{"height": height}, result)
	if err != nil {
		fmt.Printf("Cannot get block data, height=%d, error=%s \n", height, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}

func Transaction(txHash, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc2.TxResult)
	_, err = rpc.Call("bcb_transaction", map[string]interface{}{"txHash": txHash}, result)
	if err != nil {
		fmt.Printf("Cannot get transaction, txHash=%s, error=%s \n", txHash, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}

func Balance(address keys.Address, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc2.BalanceResult)
	_, err = rpc.Call("bcb_balance", map[string]interface{}{"address": address}, result)
	if err != nil {
		fmt.Printf("Cannot get balance, address=%s, error=%s \n", address, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}

func BalanceOfToken(address, tokenAddress keys.Address, tokenName string, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc2.BalanceResult)
	_, err = rpc.Call("bcb_balanceOfToken", map[string]interface{}{"address": address, "tokenAddress": tokenAddress, "tokenName": tokenName}, result)
	if err != nil {
		fmt.Printf("Cannot get balance of token, address=%s, tokenAddress=%s, error=%s \n", address, tokenAddress, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}

func AllBalance(address keys.Address, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new([]rpc2.AllBalanceItemResult)
	_, err = rpc.Call("bcb_allBalance", map[string]interface{}{"address": address}, result)
	if err != nil {
		fmt.Printf("Cannot all balance, address=%s, error=%s \n", address, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}

func Nonce(address keys.Address, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc2.NonceResult)
	_, err = rpc.Call("bcb_nonce", map[string]interface{}{"address": address}, result)
	if err != nil {
		fmt.Printf("Cannot get nonce, address=%s, error=%s \n", address, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}

func CommitTx(tx, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc2.CommitTxResult)
	_, err = rpc.Call("bcb_commitTx", map[string]interface{}{"tx": tx}, result)
	if err != nil {
		fmt.Printf("Cannot commit transation, tx=%s, error=%s \n", tx, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(jsIndent))

	return
}
