package transferagency_v1_0

import (
	"math/big"

	atm "bcbchain.io/algorithm"
	"bcbchain.io/prototype"
	"bcbchain.io/rlp"
	"bcbchain.io/tx/tx"
	"bcbchain.io/utils"
)

func TACSetManager(nonce, gasLimit, note, smcAddress, address, privateKey string) string {
	items := make([]string, 0)
	items = append(items, address)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.TACSetManager)),
		items,
		privateKey,
	)
}

func TACSetTokenFee(nonce, gasLimit, note, smcAddress, strTokenFee, privateKey string) string {
	items := make([]string, 0)
	items = append(items, strTokenFee)

	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.TACSetTokenFee)),
		items,
		privateKey,
	)
}

func TACTrasfer(nonce, gasLimit, note, smcAddress, tokenName, targetAddress, amount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, tokenName)
	items = append(items, targetAddress)
	items = append(items, amount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.TACTransfer)),
		items,
		privateKey,
	)
}

func TACWithdrawFunds(nonce, gasLimit, note, smcAddress, tokenName, amount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, tokenName)
	items = append(items, amount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.TACWithdrawFunds)),
		items,
		privateKey,
	)
}

func DecodeTACSetManagerItems(data []byte) ([]string, string) {
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

func DecodeTACSetTokenFeeItems(data []byte) ([]string, string) {
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

func DecodeTACTransferItems(data []byte) ([]string, string) {
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

	funds := new(big.Int).SetBytes(itemsBytes[0][:]).String()
	items = append(items, funds)

	return items, errMsg
}

func DecodeTACWithdrawFundsItems(data []byte) ([]string, string) {
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

	funds := new(big.Int).SetBytes(itemsBytes[0][:]).String()
	items = append(items, funds)

	return items, errMsg
}
