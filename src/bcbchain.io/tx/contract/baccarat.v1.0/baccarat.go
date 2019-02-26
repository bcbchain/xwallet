package baccarat_v1_0

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

func BACSetOwner(nonce, gasLimit, note, smcAddress, newOwner, privateKey string) string {
	items := make([]string, 0)
	items = append(items, newOwner)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BACSetOwner)),
		items,
		privateKey,
	)
}

func BACSetSettings(nonce, gasLimit, note, smcAddress, setting, privateKey string) string {
	items := make([]string, 0)
	items = append(items, setting)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BACSetSettings)),
		items,
		privateKey,
	)
}

func BACSetRecvFeeInfo(nonce, gasLimit, note, smcAddress, recvFeeInfo, privateKey string) string {
	items := make([]string, 0)
	items = append(items, recvFeeInfo)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BACSetRecvFeeInfo)),
		items,
		privateKey,
	)
}

func BACSetSecretSigner(nonce, gasLimit, note, smcAddress, newSecretSigner, privateKey string) string {
	items := make([]string, 0)
	items = append(items, newSecretSigner)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BACSetSecretSigner)),
		items,
		privateKey,
	)
}

func BACWithdrawFunds(nonce, gasLimit, note, smcAddress, beneficiaryAddr, _withdrawAmount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, beneficiaryAddr)
	items = append(items, _withdrawAmount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BACWithdrawFunds)),
		items,
		privateKey,
	)
}

func BACPlaceBet(nonce, gasLimit, note, smcAddress, amount, betData, commitLastBlock, commit, signData, refAddress, privateKey string) string {
	items := make([]string, 0)
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
		utils.BytesToHex(atm.CalcMethodId(prototype.BACPlaceBet)),
		items,
		privateKey,
	)
}

func BACSettleBets(nonce, gasLimit, note, smcAddress, reveal, sellteCount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, reveal)
	items = append(items, sellteCount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BACSettleBets)),
		items,
		privateKey,
	)
}

func BACRefundBets(nonce, gasLimit, note, smcAddress, commit, refundCount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, commit)
	items = append(items, refundCount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BACRefundBets)),
		items,
		privateKey,
	)
}

func BACWithdrawWin(nonce, gasLimit, note, smcAddress, commit, privateKey string) string {
	items := make([]string, 0)
	items = append(items, commit)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.BACWithdrawWin)),
		items,
		privateKey,
	)
}

func DecodeBACSetOwnerItems(data []byte) ([]string, string) {
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

func DecodeBACSetSecretSignerItems(data []byte) ([]string, string) {
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

func DecodeBACSetSettingsItems(data []byte) ([]string, string) {
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

func DecodeBACSetRecvFeeInfoItems(data []byte) ([]string, string) {
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

func DecodeBACWithdrawFundsItems(data []byte) ([]string, string) {
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

	affAddr := string(itemsBytes[0][:])
	funds := new(big.Int).SetBytes(itemsBytes[1][:]).String()
	items = append(items, affAddr)
	items = append(items, funds)

	return items, errMsg
}

func DecodeBACPlaceBetItems(data []byte) ([]string, string) {
	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := ("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 6 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	amount := new(big.Int).SetBytes(itemsBytes[0][:]).String()
	betData := string(itemsBytes[1][:])
	commitLastBlock := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	commit := hex.EncodeToString(itemsBytes[3][:])
	signData := hex.EncodeToString(itemsBytes[4][:])
	address := hex.EncodeToString(itemsBytes[5][:])
	items = append(items, amount)
	items = append(items, betData)
	items = append(items, commitLastBlock)
	items = append(items, commit)
	items = append(items, signData)
	items = append(items, address)
	return items, errMsg
}

func DecodeBACSettleBetsItems(data []byte) ([]string, string) {
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

func DecodeBACRefundBetsItems(data []byte) ([]string, string) {
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

func DecodeBACWithdrawWinItems(data []byte) ([]string, string) {
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
