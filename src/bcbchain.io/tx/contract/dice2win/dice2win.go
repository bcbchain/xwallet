package dice2win

import (
	"encoding/hex"
	atm "bcbchain.io/algorithm"
	"bcbchain.io/prototype"
	"bcbchain.io/rlp"
	"bcbchain.io/tx/contract/common"
	"bcbchain.io/tx/tx"
	"bcbchain.io/utils"
	"math/big"
	"strconv"
)

func DWSetOwner(nonce, gasLimit, note, smcAddress, newOwner, privateKey string) string {
	items := make([]string, 0)
	items = append(items, newOwner)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWSetOwner)),
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
		utils.BytesToHex(atm.CalcMethodId(prototype.DWSetSecretSigner)),
		items,
		privateKey,
	)
}

func DWWithdrawFunds(nonce, gasLimit, note, smcAddress, beneficiaryAddr, _withdrawAmount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, beneficiaryAddr)
	items = append(items, _withdrawAmount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWWithdrawFunds)),
		items,
		privateKey,
	)
}

func DWPlaceBet(nonce, gasLimit, note, smcAddress, amount, betMask, modulo, commitLastBlock, commit, signData, refAddress, privateKey string) string {
	items := make([]string, 0)
	items = append(items, amount)
	items = append(items, betMask)
	items = append(items, modulo)
	items = append(items, commitLastBlock)
	items = append(items, commit)
	items = append(items, signData)
	items = append(items, refAddress)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWPlaceBet)),
		items,
		privateKey,
	)
}

func DWSettleBet(nonce, gasLimit, note, smcAddress, reveal, cleanCommit, privateKey string) string {
	items := make([]string, 0)
	items = append(items, reveal)
	items = append(items, cleanCommit)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWSettleBet)),
		items,
		privateKey,
	)
}

func DWRefundBet(nonce, gasLimit, note, smcAddress, commit, privateKey string) string {
	items := make([]string, 0)
	items = append(items, commit)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWRefundBet)),
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
		utils.BytesToHex(atm.CalcMethodId(prototype.DWSetSettings)),
		items,
		privateKey,
	)
}

func DWClearStorage(nonce, gasLimit, note, smcAddress, cleanCommits, privateKey string) string {
	items := make([]string, 0)
	items = append(items, cleanCommits)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.DWClearStorage)),
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
		utils.BytesToHex(atm.CalcMethodId(prototype.DWSetRecFeeInfo)),
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

func DecodeDWWithdrawFundsItems(data []byte) ([]string, string) {
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

func DecodeDWPlaceBetItems(data []byte) ([]string, string) {
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

	amount := new(big.Int).SetBytes(itemsBytes[0][:]).String()
	betMask := new(big.Int).SetBytes(itemsBytes[1][:]).String()

	moudlo := strconv.FormatInt(common.Decode2Int64(itemsBytes[2][:]), 10)
	commitLastBlock := new(big.Int).SetBytes(itemsBytes[3][:]).String()
	commit := hex.EncodeToString(itemsBytes[4][:])
	signData := hex.EncodeToString(itemsBytes[5][:])
	address := hex.EncodeToString(itemsBytes[6][:])
	items = append(items, amount)
	items = append(items, betMask)
	items = append(items, moudlo)
	items = append(items, commitLastBlock)
	items = append(items, commit)
	items = append(items, signData)
	items = append(items, address)
	return items, errMsg
}

func DecodeDWSettleBetItems(data []byte) ([]string, string) {
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
	cleanCommit := hex.EncodeToString(itemsBytes[1][:])
	items = append(items, reveal)
	items = append(items, cleanCommit)

	return items, errMsg
}

func DecodeDWRefundBetItems(data []byte) ([]string, string) {
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

func DecodeDWClearStorageItems(data []byte) ([]string, string) {
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

	commit := string(itemsBytes[0][:])
	items = append(items, commit)

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
