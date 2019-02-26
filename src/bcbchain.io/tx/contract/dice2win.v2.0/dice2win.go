package dice2win2_0

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

func DWSetOwner(nonce, gasLimit, note, smcAddress, newOwner, privateKey string) string {
	items := make([]string, 0)
	items = append(items, newOwner)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWSetOwner2_0)),
		items,
		privateKey,
	)
}

func DWSetSecretSigner(nonce, gasLimit, note, smcAddress, newSecretSigner, privateKey string) string {
	items := make([]string, 0)
	items = append(items, newSecretSigner)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWSetSecretSigner2_0)),
		items,
		privateKey,
	)
}

func DWSetSettings(nonce, gasLimit, note, smcAddress, setting, privateKey string) string {
	items := make([]string, 0)
	items = append(items, setting)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWSetSettings2_0)),
		items,
		privateKey,
	)
}

func DWSetRecFeeInfo(nonce, gasLimit, note, smcAddress, recFeeInfo, privateKey string) string {
	items := make([]string, 0)
	items = append(items, recFeeInfo)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWSetRecvFeeInfo2_0)),
		items,
		privateKey,
	)
}

func DWWithdrawFunds(nonce, gasLimit, note, smcAddress, tokenName, beneficiaryAddr, withdrawAmount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, tokenName)
	items = append(items, beneficiaryAddr)
	items = append(items, withdrawAmount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWWithdrawFunds2_0)),
		items,
		privateKey,
	)
}

func DWPlaceBet(nonce, gasLimit, note, smcAddress, tokenName, amount, betMask, module, commitLastBlock, commit, signData, refAddress, privateKey string) string {
	items := make([]string, 0)
	items = append(items, tokenName)
	items = append(items, amount)
	items = append(items, betMask)
	items = append(items, module)
	items = append(items, commitLastBlock)
	items = append(items, commit)
	items = append(items, signData)
	items = append(items, refAddress)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWPlaceBet2_0)),
		items,
		privateKey,
	)
}

func DWSettleBets(nonce, gasLimit, note, smcAddress, reveal, sellteCount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, reveal)
	items = append(items, sellteCount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWSettleBets2_0)),
		items,
		privateKey,
	)
}

func DWRefundBets(nonce, gasLimit, note, smcAddress, commit, refundCount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, commit)
	items = append(items, refundCount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWRefundBets2_0)),
		items,
		privateKey,
	)
}

func DWWithdrawWin(nonce, gasLimit, note, smcAddress, commit, privateKey string) string {
	items := make([]string, 0)
	items = append(items, commit)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWWithdrawWin2_0)),
		items,
		privateKey,
	)
}

func DecodeDWSetOwnerItems(data []byte) ([]string, string) {
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

func DecodeDWSetSecretSignerItems(data []byte) ([]string, string) {
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

func DecodeDWSetSettingsItems(data []byte) ([]string, string) {
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

func DecodeDWSetRecFeeInfoItems(data []byte) ([]string, string) {
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

func DecodeDWWithdrawFundsItems(data []byte) ([]string, string) {
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

func DecodeDWPlaceBetItems(data []byte) ([]string, string) {
	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := ("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 8 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	tokenName := string(itemsBytes[0][:])
	amount := new(big.Int).SetBytes(itemsBytes[1][:]).String()
	betMask := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	module := strconv.FormatInt(common.Decode2Int64(itemsBytes[3][:]), 10)
	commitLastBlock := strconv.FormatInt(common.Decode2Int64(itemsBytes[4][:]), 10)
	commit := hex.EncodeToString(itemsBytes[5][:])
	signData := hex.EncodeToString(itemsBytes[6][:])
	refAddress := hex.EncodeToString(itemsBytes[7][:])

	items = append(items, tokenName)
	items = append(items, amount)
	items = append(items, betMask)
	items = append(items, module)
	items = append(items, commitLastBlock)
	items = append(items, commit)
	items = append(items, signData)
	items = append(items, refAddress)
	return items, errMsg
}

func DecodeDWSettleBetsItems(data []byte) ([]string, string) {
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

func DecodeDWRefundBetsItems(data []byte) ([]string, string) {
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

func DecodeDWWithdrawWinItems(data []byte) ([]string, string) {
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
