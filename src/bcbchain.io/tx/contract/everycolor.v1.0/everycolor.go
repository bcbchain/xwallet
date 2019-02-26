package everycolor_v1_0

import (
	"encoding/hex"
	"math/big"
	"strconv"

	atm "bcbchain.io/algorithm"
	"bcbchain.io/prototype"
	"bcbchain.io/rlp"
	"bcbchain.io/tx/contract/common"
	"bcbchain.io/tx/tx"
	"bcbchain.io/utils"
)

func ECSetOwner(nonce, gasLimit, note, smcAddress, newOwner, privateKey string) string {
	items := make([]string, 0)
	items = append(items, newOwner)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.ECSetOwner)),
		items,
		privateKey,
	)
}

func ECSetSecretSigner(nonce, gasLimit, note, smcAddress, newSecretSigner, privateKey string) string {
	items := make([]string, 0)
	items = append(items, newSecretSigner)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.ECSetSecretSigner)),
		items,
		privateKey,
	)
}

func ECSetSettings(nonce, gasLimit, note, smcAddress, setting, privateKey string) string {
	items := make([]string, 0)
	items = append(items, setting)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.ECSetSettings)),
		items,
		privateKey,
	)
}

func ECSetRecFeeInfo(nonce, gasLimit, note, smcAddress, recFeeInfo, privateKey string) string {
	items := make([]string, 0)
	items = append(items, recFeeInfo)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.ECSetRecvFeeInfo)),
		items,
		privateKey,
	)
}

func ECWithdrawFunds(nonce, gasLimit, note, smcAddress, tokenName, beneficiaryAddr, withdrawAmount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, tokenName)
	items = append(items, beneficiaryAddr)
	items = append(items, withdrawAmount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.ECWithdrawFunds)),
		items,
		privateKey,
	)
}

func ECPlaceBet(nonce, gasLimit, note, smcAddress, tokenName, amount, betData, commitLastBlock, commit, signData, refAddress, privateKey string) string {
	items := make([]string, 0)
	items = append(items, tokenName)
	items = append(items, amount)
	items = append(items, betData)
	items = append(items, commitLastBlock)
	items = append(items, commit)
	items = append(items, signData)
	items = append(items, refAddress)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.ECPlaceBet)),
		items,
		privateKey,
	)
}

func ECSettleBets(nonce, gasLimit, note, smcAddress, reveal, sellteCount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, reveal)
	items = append(items, sellteCount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.ECSettleBets)),
		items,
		privateKey,
	)
}

func ECRefundBets(nonce, gasLimit, note, smcAddress, commit, refundCount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, commit)
	items = append(items, refundCount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.ECRefundBets)),
		items,
		privateKey,
	)
}

func ECWithdrawWin(nonce, gasLimit, note, smcAddress, commit, privateKey string) string {
	items := make([]string, 0)
	items = append(items, commit)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.ECWithdrawWin)),
		items,
		privateKey,
	)
}

func DecodeECSetOwnerItems(data []byte) ([]string, string) {
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

func DecodeECSetSecretSignerItems(data []byte) ([]string, string) {
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
	items = append(items, hex.EncodeToString(itemsBytes[0][:]))
	return items, errMsg
}

func DecodeECSetSettingsItems(data []byte) ([]string, string) {
	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := ("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, string(itemsBytes[0][:]))
	return items, errMsg
}

func DecodeECSetRecFeeInfoItems(data []byte) ([]string, string) {
	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := ("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, string(itemsBytes[0][:]))
	return items, errMsg
}

func DecodeECWithdrawFundsItems(data []byte) ([]string, string) {
	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := ("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 3 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}
	tokenName := string(itemsBytes[0][:])
	affAddr := string(itemsBytes[1][:])
	funds := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	items = append(items, tokenName)
	items = append(items, affAddr)
	items = append(items, funds)

	return items, errMsg
}

func DecodeECPlaceBetItems(data []byte) ([]string, string) {
	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := ("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 7 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	tokenName := string(itemsBytes[0][:])
	amount := new(big.Int).SetBytes(itemsBytes[1][:]).String()
	betData := string(itemsBytes[2][:])
	commitLastBlock := strconv.FormatInt(common.Decode2Int64(itemsBytes[3][:]), 10)
	commit := hex.EncodeToString(itemsBytes[4][:])
	signData := hex.EncodeToString(itemsBytes[5][:])
	refAddress := hex.EncodeToString(itemsBytes[6][:])

	items = append(items, tokenName)
	items = append(items, amount)
	items = append(items, betData)
	items = append(items, commitLastBlock)
	items = append(items, commit)
	items = append(items, signData)
	items = append(items, refAddress)
	return items, errMsg
}

func DecodeECSettleBetsItems(data []byte) ([]string, string) {
	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := ("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 2 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}
	reveal := hex.EncodeToString(itemsBytes[0][:])
	sellteCount := strconv.FormatInt(common.Decode2Int64(itemsBytes[1][:]), 10)
	items = append(items, reveal)
	items = append(items, sellteCount)

	return items, errMsg
}

func DecodeECRefundBetsItems(data []byte) ([]string, string) {
	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := ("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	commit := hex.EncodeToString(itemsBytes[0][:])
	refundCount := strconv.FormatInt(common.Decode2Int64(itemsBytes[1][:]), 10)
	items = append(items, commit)
	items = append(items, refundCount)

	return items, errMsg
}

func DecodeECWithdrawWinItems(data []byte) ([]string, string) {
	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := ("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	commit := hex.EncodeToString(itemsBytes[0][:])
	items = append(items, commit)

	return items, errMsg
}
