package token_basic_v1_0

import (
	atm "bcbchain.io/algorithm"
	"bcbchain.io/prototype"
	"bcbchain.io/rlp"
	"bcbchain.io/tx/contract/common"

	"bcbchain.io/tx/tx"
	"bcbchain.io/utils"
	"github.com/tendermint/go-crypto"
	"strconv"
	"unsafe"
)

var version string = ""

func TbTransfer(nonce, gasLimit, note, to, value, privateKey string) string {
	items := make([]string, 0)
	items = append(items, to)
	items = append(items, value)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		atm.CalcContractAddress(crypto.GetChainId(), tx.Owner, prototype.TokenBasic, version),
		utils.BytesToHex(atm.CalcMethodId(prototype.TbTransfer)),
		items,
		privateKey)
}

func TbSetGasBasePrice(nonce, gasLimit, note, gasPrice, privateKey string) string {
	items := make([]string, 0)
	items = append(items, gasPrice)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		atm.CalcContractAddress(crypto.GetChainId(), tx.Owner, prototype.TokenBasic, version),
		utils.BytesToHex(atm.CalcMethodId(prototype.TbSetGasBasePrice)),
		items,
		privateKey)
}

func TbSetGasPrice(nonce, gasLimit, note, gasPrice, privateKey string) string {
	items := make([]string, 0)
	items = append(items, gasPrice)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		atm.CalcContractAddress(crypto.GetChainId(), tx.Owner, prototype.TokenBasic, version),
		utils.BytesToHex(atm.CalcMethodId(prototype.TbSetGasPrice)),
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
	return items, errMsg
}

func DecodeSetGasBasePriceItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes[0]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed gasbaseprice is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, strconv.FormatUint(common.Decode2Uint64(itemsBytes[0][:]), 10))
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
