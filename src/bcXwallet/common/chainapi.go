package common

import (
	types2 "blockchain/abciapp_v1.0/types"
	rpcclient "common/rpc/lib/client"
	"encoding/json"
	"errors"
	"strings"
)

//网络请求和结果解析
func DoHttpRequestAndParseEx(nodeAddrSlice []string, methodName string, params map[string]interface{}, result interface{}) (err error) {

	for i, nodeAddr := range nodeAddrSlice {
		rpc := rpcclient.NewJSONRPCClientEx(nodeAddr, "", true)
		_, err = rpc.Call(methodName, params, result)
		if err == nil {
			break
		} else {
			if i == len(nodeAddrSlice)-1 {
				splitErr := strings.Split(err.Error(), ":")
				return errors.New(strings.Trim(splitErr[len(splitErr)-1], " "))
			}
		}
	}

	return
}

//网络请求和结果解析
func DoHttpRequestAndParse(nodeAddrSlice []string, txStr string) (*types2.ResultBroadcastTxCommit, error) {

	result := new(types2.ResultBroadcastTxCommit)

	for i, nodeAddr := range nodeAddrSlice {
		rpc := rpcclient.NewJSONRPCClientEx(nodeAddr, "", true)
		_, err := rpc.Call("broadcast_tx_commit", map[string]interface{}{"tx": []byte(txStr)}, result)
		if err == nil {
			break
		} else {
			if i == len(nodeAddrSlice)-1 {
				splitErr := strings.Split(err.Error(), ":")
				return nil, errors.New(strings.Trim(splitErr[len(splitErr)-1], " "))
			}
		}
	}

	return result, nil
}

func DoHttpQueryAndParse(nodeAddrSlice []string, key string, data interface{}) (err error) {

	value, err := DoHttpQuery(nodeAddrSlice, key)
	if err != nil {
		return
	}

	err = json.Unmarshal(value, data)

	return
}

func DoHttpQuery(nodeAddrSlice []string, key string) (value []byte, err error) {

	result := new(types2.ResultABCIQuery)
	for i, nodeAddr := range nodeAddrSlice {
		rpc := rpcclient.NewJSONRPCClientEx(nodeAddr, "", true)
		_, err = rpc.Call("abci_query", map[string]interface{}{"path": key}, result)
		if err == nil {
			break
		} else {
			if i == len(nodeAddrSlice)-1 {
				splitErr := strings.Split(err.Error(), ":")
				return nil, errors.New(strings.Trim(splitErr[len(splitErr)-1], " "))
			}
		}
	}
	value = result.Response.Value

	return
}
