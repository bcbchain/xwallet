package client

import (
	rpc3 "bcXwallet/rpc"
	"blockchain/abciapp_v1.0/keys"
	rpcclient "common/rpc/lib/client"
	"encoding/json"
	"fmt"
)

func BlockHeight(url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc3.BlockHeightResult)
	_, err = rpc.Call("bcb_blockHeight", map[string]interface{}{}, result)
	if err != nil {
		fmt.Printf("Cannot get block height, error=%s \n", err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))

	return
}

func Block(height int64, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc3.BlockResult)
	_, err = rpc.Call("bcb_block", map[string]interface{}{"height": height}, result)
	if err != nil {
		fmt.Printf("Cannot get block data, height=%d, error=%s \n", height, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))

	return
}

func Transaction(txHash, url string) (err error) {

	if txHash[:2] == "0x" {
		txHash = txHash[2:]
	}
	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc3.TxResult)
	_, err = rpc.Call("bcb_transaction", map[string]interface{}{"txHash": txHash}, result)
	if err != nil {
		fmt.Printf("Cannot get transaction, txHash=%s, error=%s \n", txHash, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))

	return
}

func Balance(address keys.Address, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc3.BalanceResult)
	_, err = rpc.Call("bcb_balance", map[string]interface{}{"address": address}, result)
	if err != nil {
		fmt.Printf("Cannot get balance, address=%s, error=%s \n", address, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))

	return
}

func BalanceOfToken(address, tokenAddress keys.Address, tokenName string, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc3.BalanceResult)
	_, err = rpc.Call("bcb_balanceOfToken", map[string]interface{}{"address": address, "tokenAddress": tokenAddress, "tokenName": tokenName}, result)
	if err != nil {
		fmt.Printf("Cannot get balance of token, address=%s, tokenAddress=%s, error=%s \n", address, tokenAddress, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))

	return
}

func AllBalance(address keys.Address, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new([]rpc3.AllBalanceItemResult)
	_, err = rpc.Call("bcb_allBalance", map[string]interface{}{"address": address}, result)
	if err != nil {
		fmt.Printf("Cannot all balance, address=%s, error=%s \n", address, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))

	return
}

func Nonce(address keys.Address, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc3.NonceResult)
	_, err = rpc.Call("bcb_nonce", map[string]interface{}{"address": address}, result)
	if err != nil {
		fmt.Printf("Cannot get nonce, address=%s, error=%s \n", address, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))

	return
}

func CommitTx(tx, url string) (err error) {

	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)

	result := new(rpc3.CommitTxResult)
	_, err = rpc.Call("bcb_commitTx", map[string]interface{}{"tx": tx}, result)
	if err != nil {
		fmt.Printf("Cannot commit transation, tx=%s, error=%s \n", tx, err.Error())
		return nil
	}

	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))

	return
}
