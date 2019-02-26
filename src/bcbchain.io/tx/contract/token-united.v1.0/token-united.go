package token_united_v1_0

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

func UtNewToken(nonce, gasLimit, note, smcAddress, name, symbol, privateKey, totalSupply,
	addSupplyEnabled, burnEnabled, gasPrice string) string {
	items := make([]string, 0)
	items = append(items, name)
	items = append(items, symbol)
	items = append(items, totalSupply)
	items = append(items, addSupplyEnabled)
	items = append(items, burnEnabled)
	items = append(items, gasPrice)

	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.UtNewToken)),
		items,
		privateKey,
	)
}

func UtTransfer(nonce, gasLimit, note, smcAddress, to, privateKey, value string) string {
	items := make([]string, 0)
	items = append(items, to)
	items = append(items, value)

	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.UtTransfer)),
		items,
		privateKey,
	)
}

func UtSetOwner(nonce, gasLimit, note, smcAddress, newOwner, privateKey string) string {
	items := make([]string, 0)
	items = append(items, newOwner)

	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.UtSetOwner)),
		items,
		privateKey,
	)
}

func UtSetGasPrice(nonce, gasLimit, note, smcAddress, privateKey, gasPrice string) string {
	items := make([]string, 0)
	items = append(items, gasPrice)

	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.UtSetGasPrice)),
		items,
		privateKey,
	)
}

func UtSetFee(nonce, gasLimit, note, smcAddress, privateKey, ratio, maxFee, minFee string) string {
	items := make([]string, 0)
	items = append(items, ratio)
	items = append(items, maxFee)
	items = append(items, minFee)

	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.UtSetFee)),
		items,
		privateKey,
	)
}

func UtSetGasPayer(nonce, gasLimit, note, smcAddress, payer, privateKey string) string {
	items := make([]string, 0)
	items = append(items, payer)

	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.UtSetGasPayer)),
		items,
		privateKey,
	)
}

func UtAddSupply(nonce, gasLimit, note, smcAddress, privateKey, value string) string {
	items := make([]string, 0)
	items = append(items, value)

	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.UtAddSupply)),
		items,
		privateKey,
	)
}

func UtBurn(nonce, gasLimit, note, smcAddress, privateKey, value string) string {
	items := make([]string, 0)
	items = append(items, value)

	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.UtBurn)),
		items,
		privateKey,
	)
}

func UtWithdraw(nonce, gasLimit, note, smcAddress, privateKey string) string {
	items := make([]string, 0)

	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.UtWithdraw)),
		items,
		privateKey,
	)
}

func DecodeNewTokenItems(data []byte) ([]string, string) {
	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 6 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes[5]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed gas price invalid ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	value := new(big.Int).SetBytes(itemsBytes[2][:]).String()

	items = append(items, string(itemsBytes[0][:]))
	items = append(items, string(itemsBytes[1][:]))
	items = append(items, value)
	items = append(items, string(itemsBytes[3][:]))
	items = append(items, string(itemsBytes[4][:]))
	items = append(items, strconv.FormatUint(common.Decode2Uint64(itemsBytes[5]), 10))
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

	value := new(big.Int).SetBytes(itemsBytes[1][:]).String()

	items = append(items, string(itemsBytes[0][:]))
	items = append(items, value)
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

	if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes[0]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed gas price invalid ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, strconv.FormatUint(common.Decode2Uint64(itemsBytes[0]), 10))
	return items, errMsg
}

func DecodeSetGasPayerItems(data []byte) ([]string, string) {
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
func DecodeSetFeeItems(data []byte) ([]string, string) {
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

	if len(itemsBytes[0]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed max fee invalid ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes[1]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed min fee invalid ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes[2]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed percent invalid ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, strconv.FormatUint(common.Decode2Uint64(itemsBytes[0]), 10))
	items = append(items, strconv.FormatUint(common.Decode2Uint64(itemsBytes[1]), 10))
	items = append(items, strconv.FormatUint(common.Decode2Uint64(itemsBytes[2]), 10))

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

	value := new(big.Int).SetBytes(itemsBytes[0][:]).String()

	items = append(items, value)
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

	value := new(big.Int).SetBytes(itemsBytes[0][:]).String()

	items = append(items, value)
	return items, errMsg
}
