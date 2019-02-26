package token_byb_v1_0

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

func BYBInit(nonce, gasLimit, note, smcAddress, totalSupply, addSupplyEnabled, burnEnabled, privateKey string) string {
	items := make([]string, 0)
	items = append(items, totalSupply)
	items = append(items, addSupplyEnabled)
	items = append(items, burnEnabled)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BYBInit)),
		items,
		privateKey,
	)
}

func BYBSetOwner(nonce, gasLimit, note, smcAddress, toAddress, privateKey string) string {
	items := make([]string, 0)
	items = append(items, toAddress)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BYBSetOwner)),
		items,
		privateKey,
	)
}

func BYBSetGasPrice(nonce, gasLimit, note, smcAddress, gasPrice, privateKey string) string {
	items := make([]string, 0)
	items = append(items, gasPrice)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BYBSetGasPrice)),
		items,
		privateKey,
	)
}

func BYBAddsupply(nonce, gasLimit, note, smcAddress, value, privateKey string) string {
	items := make([]string, 0)
	items = append(items, value)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BYBAddSupply)),
		items,
		privateKey,
	)
}

func BYBBurn(nonce, gasLimit, note, smcAddress, value, privateKey string) string {
	items := make([]string, 0)
	items = append(items, value)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BYBBurn)),
		items,
		privateKey,
	)
}

func BYBNewBlackHole(nonce, gasLimit, note, smcAddress, blackHoleAddr, privateKey string) string {
	items := make([]string, 0)
	items = append(items, blackHoleAddr)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BYBNewBlackHole)),
		items,
		privateKey,
	)
}

func BYBNewStockHolder(nonce, gasLimit, note, smcAddress, stockHolderAddr, value, privateKey string) string {
	items := make([]string, 0)
	items = append(items, stockHolderAddr)
	items = append(items, value)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BYBNewStockHolder)),
		items,
		privateKey,
	)
}

func BYBDelStockHolder(nonce, gasLimit, note, smcAddress, stockHolderAddr, privateKey string) string {
	items := make([]string, 0)
	items = append(items, stockHolderAddr)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BYBDelStockHolder)),
		items,
		privateKey,
	)
}

func BYBChangeChromoOwnerShip(nonce, gasLimit, note, smcAddress, toAddress, chromo, privateKey string) string {
	items := make([]string, 0)
	items = append(items, chromo)
	items = append(items, toAddress)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BYBChangeChromoOwnerShip)),
		items,
		privateKey,
	)
}

func BYBTransfer(nonce, gasLimit, note, smcAddress, toAddress, value, privateKey string) string {
	items := make([]string, 0)
	items = append(items, toAddress)
	items = append(items, value)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BYBTransfer)),
		items,
		privateKey,
	)
}

func BYBTransferByChromo(nonce, gasLimit, note, smcAddress, toAddress, chromo, value, privateKey string) string {
	items := make([]string, 0)
	items = append(items, chromo)
	items = append(items, toAddress)
	items = append(items, value)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BYBTransferByChromo)),
		items,
		privateKey,
	)
}

func DecodeInitItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 3 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	value := new(big.Int).SetBytes(itemsBytes[0][:]).String()
	items = append(items, value)
	items = append(items, string(itemsBytes[1][:]))
	items = append(items, string(itemsBytes[2][:]))
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

	if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, string(itemsBytes[0][:]))
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
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed team is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, strconv.FormatUint(common.Decode2Uint64(itemsBytes[0]), 10))
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

	if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
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

	if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, new(big.Int).SetBytes(itemsBytes[0][:]).String())
	return items, errMsg
}

func DecodeNewBlackHoleItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, string(itemsBytes[0][:]))
	return items, errMsg
}

func DecodeNewStockHolderItems(data []byte) ([]string, string) {

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

func DecodeDelStockHolderItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, string(itemsBytes[0][:]))
	return items, errMsg
}

func DecodeChangeChromoOwnerShipItems(data []byte) ([]string, string) {

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
	items = append(items, string(itemsBytes[1][:]))
	return items, errMsg
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

func DecodeTransferByChromoItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 3 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, string(itemsBytes[0][:]))
	items = append(items, string(itemsBytes[1][:]))
	items = append(items, new(big.Int).SetBytes(itemsBytes[2][:]).String())
	return items, errMsg
}
