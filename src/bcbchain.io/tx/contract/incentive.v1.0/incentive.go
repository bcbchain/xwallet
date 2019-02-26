package incentive_v_1_0

import (
	"math/big"

	atm "bcbchain.io/algorithm"
	"bcbchain.io/prototype"
	"bcbchain.io/rlp"
	"bcbchain.io/tx/tx"
	"bcbchain.io/utils"
)

func ICTSetManager(nonce, gasLimit, note, smcAddress, address, privateKey string) string {
	items := make([]string, 0)
	items = append(items, address)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.ICTSetManager)),
		items,
		privateKey,
	)
}

func ICTSetPayManager(nonce, gasLimit, note, smcAddress, address, privateKey string) string {
	items := make([]string, 0)
	items = append(items, address)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.ICTSetPayManager)),
		items,
		privateKey,
	)
}

func ICTSetEmployeeBonuses(nonce, gasLimit, note, smcAddress, strEmployeeBonuses, privateKey string) string {
	items := make([]string, 0)
	items = append(items, strEmployeeBonuses)

	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.ICTSetEmployeeBonuses)),
		items,
		privateKey,
	)
}

func ICTPayBonuses(nonce, gasLimit, note, smcAddress, totalAmount, privateKey string) string {
	items := make([]string, 0)
	items = append(items, totalAmount)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.ICTPayBonuses)),
		items,
		privateKey,
	)
}

func DecodeICTSetManagerItems(data []byte) ([]string, string) {
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

func DecodeICTSetPayManagerItems(data []byte) ([]string, string) {
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

func DecodeICTSetEmployeeBonusesItems(data []byte) ([]string, string) {
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

func DecodeICTPayBonusesItems(data []byte) ([]string, string) {
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
