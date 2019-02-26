package token_templet_v1_0

import (
	atm "bcbchain.io/algorithm"
	"bcbchain.io/prototype"
	"bcbchain.io/rlp"
	"bcbchain.io/tx/contract/common"
	"bcbchain.io/tx/tx"
	"bcbchain.io/utils"
	"math/big"
	"strconv"
	"unsafe"
)

var version string = ""

func TtSetGasPrice(nonce, gasLimit, note, currency, gasPrice, privateKey string) string {
	items := make([]string, 0)
	items = append(items, gasPrice)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		currency,
		utils.BytesToHex(atm.CalcMethodId(prototype.TtSetGasPrice)),
		items,
		privateKey)
}

func TtTransfer(nonce, gasLimit, note, tokenAddress, to, value, privateKey string) string {
	items := make([]string, 0)
	items = append(items, to)
	items = append(items, value)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		tokenAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.TtTransfer)),
		items,
		privateKey)
}

func TtBatchTransfer(nonce, gasLimit, note, tokenAddress, toList, value, privateKey string) string {
	items := make([]string, 0)
	items = append(items, toList)
	items = append(items, value)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		tokenAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.TtBatchTransfer)),
		items,
		privateKey)
}

func TtAddSupply(nonce, gasLimit, note, tokenAddress, value, privateKey string) string {
	items := make([]string, 0)
	items = append(items, value)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		tokenAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.TtAddSupply)),
		items,
		privateKey)
}

func TtBurn(nonce, gasLimit, note, tokenAddress, value, privateKey string) string {
	items := make([]string, 0)
	items = append(items, value)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		tokenAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.TtBurn)),
		items,
		privateKey)
}

func TtSetOwner(nonce, gasLimit, note, tokenAddress, newOwner, privateKey string) string {
	items := make([]string, 0)
	items = append(items, newOwner)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		tokenAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.TtSetOwner)),
		items,
		privateKey)
}

func DecodeTransferItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 2 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, string(itemsBytes[0][:]))
	items = append(items, new(big.Int).SetBytes(itemsBytes[1][:]).String())
	return items, errMsg
}

func DecodeBatchTransferItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 2 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, string(itemsBytes[0][:]))
	items = append(items, new(big.Int).SetBytes(itemsBytes[1][:]).String())
	return items, errMsg
}

func DecodeAddSupplyItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, new(big.Int).SetBytes(itemsBytes[0][:]).String())
	return items, errMsg
}

func DecodeBurnItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, new(big.Int).SetBytes(itemsBytes[0][:]).String())
	return items, errMsg
}

func DecodeSetGasPriceItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes[0]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed gasprice is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, strconv.FormatUint(common.Decode2Uint64(itemsBytes[0][:]), 10))
	return items, errMsg
}

func DecodeSetOwnerItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, string(itemsBytes[0][:]))
	return items, errMsg
}
