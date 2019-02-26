package dc_yuebao_v_1_0

import (
	atm "bcbchain.io/algorithm"
	"bcbchain.io/prototype"
	"bcbchain.io/rlp"
	"bcbchain.io/tx/tx"
	"bcbchain.io/utils"
	"math/big"
)

func YEBSetOwner(nonce, gasLimit, note, smcAddress, newOwner, privateKey string) string {
	items := make([]string, 0)
	items = append(items, newOwner)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.YEBSetOwner)),
		items,
		privateKey,
	)
}

func YEBPromiseEarningRate(nonce, gasLimit, note, smcAddress, earningRate, date, privateKey string) string {
	items := make([]string, 0)
	items = append(items, earningRate)
	items = append(items, date)

	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.YEBPromiseEarningRate)),
		items,
		privateKey,
	)
}

func YEBDepositsReceived(nonce, gasLimit, note, smcAddress, deposits, privateKey string) string {
	items := make([]string, 0)
	items = append(items, deposits)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.YEBDepositsReceived)),
		items,
		privateKey,
	)
}

func YEBDepositsWithdrawal(nonce, gasLimit, note, smcAddress, amount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, amount)

	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.YEBDepositsWithdrawal)),
		items,
		privateKey,
	)
}

func YEBDeposit(nonce, gasLimit, note, smcAddress, Deposit, privateKey string) string {
	items := make([]string, 0)
	items = append(items, Deposit)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.YEBDeposit)),
		items,
		privateKey,
	)
}

func YEBWithdraw(nonce, gasLimit, note, smcAddress, amount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, amount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.YEBWithdraw)),
		items,
		privateKey,
	)
}

func YEBWithdrawEarnings(nonce, gasLimit, note, smcAddress, amount, userAddress, privateKey string) string {
	items := make([]string, 0)
	items = append(items, amount)
	items = append(items, userAddress)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.YEBWithdrawEarnings)),
		items,
		privateKey,
	)
}

func YEBReinvest(nonce, gasLimit, note, smcAddress, amount, userAddress, privateKey string) string {
	items := make([]string, 0)
	items = append(items, amount)
	items = append(items, userAddress)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.YEBReinvest)),
		items,
		privateKey,
	)
}

func YEBSetOperators(nonce, gasLimit, note, smcAddress, address, privateKey string) string {
	items := make([]string, 0)
	items = append(items, address)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.YEBSetOperators)),
		items,
		privateKey,
	)
}

func DecodeYEBSetOwnerItems(data []byte) ([]string, string) {
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

func DecodeYEBPromiseEarningRateItems(data []byte) ([]string, string) {
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

func DecodeYEBDepositsReceivedItems(data []byte) ([]string, string) {
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

func DecodeYEBDepositsWithdrawalItems(data []byte) ([]string, string) {
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

func DecodeYEBDepositItems(data []byte) ([]string, string) {
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

func DecodeYEBWithdrawItems(data []byte) ([]string, string) {
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

func DecodeYEBWithdrawEarningsItems(data []byte) ([]string, string) {
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

	funds := new(big.Int).SetBytes(itemsBytes[0][:]).String()
	affAddr := string(itemsBytes[1][:])
	items = append(items, funds)
	items = append(items, affAddr)

	return items, errMsg
}

func DecodeYEBReinvestItems(data []byte) ([]string, string) {
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

	funds := new(big.Int).SetBytes(itemsBytes[0][:]).String()
	affAddr := string(itemsBytes[1][:])
	items = append(items, funds)
	items = append(items, affAddr)

	return items, errMsg
}
func DecodeYEBSetOperatorsItems(data []byte) ([]string, string) {
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
