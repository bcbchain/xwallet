package crypto_gods_v1_0

import (
	atm "bcbchain.io/algorithm"
	"bcbchain.io/prototype"
	"bcbchain.io/rlp"
	"bcbchain.io/tx/contract/common"
	tx1 "bcbchain.io/tx/tx"
	"bcbchain.io/utils"
	"math/big"
	"strconv"
	"unsafe"
)

func CgsInit(nonce, gasLimit, note, smcAddress, privateKey string) string {
	items := make([]string, 0)
	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsInit)),
		items,
		privateKey,
	)
}

func CgsActivate(nonce, gasLimit, note, smcAddress, privateKey string) string {
	items := make([]string, 0)
	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsActivate)),
		items,
		privateKey,
	)
}

func CgsSetOwner(nonce, gasLimit, note, smcAddress, newOwner, privateKey string) string {
	items := make([]string, 0)
	items = append(items, newOwner)
	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsActivate)),
		items,
		privateKey,
	)
}

func CgsRegisterNameXid(nonce, gasLimit, note, smcAddress, name, affCode, bcb, privateKey string) string {
	items := make([]string, 0)
	items = append(items, name)
	items = append(items, affCode)
	items = append(items, bcb)
	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsRegisterNameXid)),
		items,
		privateKey,
	)
}

func CgsRegisterNameXaddr(nonce, gasLimit, note, smcAddress, name, affCode, bcb, privateKey string) string {
	items := make([]string, 0)
	items = append(items, name)
	items = append(items, affCode)
	items = append(items, bcb)
	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsRegisterNameXaddr)),
		items,
		privateKey,
	)
}

func CgsRegisterNameXname(nonce, gasLimit, note, smcAddress, name, affCode, bcb, privateKey string) string {
	items := make([]string, 0)
	items = append(items, name)
	items = append(items, affCode)
	items = append(items, bcb)
	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsRegisterNameXname)),
		items,
		privateKey,
	)
}

func CgsBuy(nonce, gasLimit, note, smcAddress, team, bcb, privateKey string) string {
	items := make([]string, 0)
	items = append(items, team)
	items = append(items, bcb)
	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsBuy)),
		items,
		privateKey,
	)
}

func CgsBuyXid(nonce, gasLimit, note, smcAddress, team, bcb, affCode, privateKey string) string {
	items := make([]string, 0)
	items = append(items, affCode)
	items = append(items, team)
	items = append(items, bcb)

	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsBuyXid)),
		items,
		privateKey,
	)
}

func CgsBuyXAddress(nonce, gasLimit, note, smcAddress, team, bcb, affAddr, privateKey string) string {
	items := make([]string, 0)
	items = append(items, affAddr)
	items = append(items, team)
	items = append(items, bcb)

	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsBuyXaddr)),
		items,
		privateKey,
	)
}

func CgsBuyXName(nonce, gasLimit, note, smcAddress, team, bcb, affName, privateKey string) string {
	items := make([]string, 0)
	items = append(items, affName)
	items = append(items, team)
	items = append(items, bcb)

	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsBuyXname)),
		items,
		privateKey,
	)
}

func CgsReloadXid(nonce, gasLimit, note, smcAddress, team, bcb, affCode, privateKey string) string {
	items := make([]string, 0)
	items = append(items, affCode)
	items = append(items, team)
	items = append(items, bcb)
	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsReloadXid)),
		items,
		privateKey,
	)
}

func CgsReloadXaddr(nonce, gasLimit, note, smcAddress, team, bcb, affCode, privateKey string) string {
	items := make([]string, 0)
	items = append(items, affCode)
	items = append(items, team)
	items = append(items, bcb)
	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsReloadXaddr)),
		items,
		privateKey,
	)
}

func CgsReloadXname(nonce, gasLimit, note, smcAddress, team, bcb, affCode, privateKey string) string {
	items := make([]string, 0)
	items = append(items, affCode)
	items = append(items, team)
	items = append(items, bcb)
	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsReloadXname)),
		items,
		privateKey,
	)
}

func CgsWithDraw(nonce, gasLimit, note, smcAddress, privateKey string) string {
	items := make([]string, 0)
	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsWithDraw)),
		items,
		privateKey,
	)
}

func CgsSetSettings(nonce, gasLimit, note, smcAddress, settings, privateKey string) string {
	items := make([]string, 0)
	items = append(items, settings)
	return tx1.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.CgsSetSetting)),
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

	return items, errMsg
}

func DecodeActivateItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	return items, errMsg
}

func DecodeBuyItems(data []byte) ([]string, string) {

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

	if len(itemsBytes[0]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed team is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	team := strconv.FormatInt(common.Decode2Int64(itemsBytes[0][:]), 10)
	value := new(big.Int).SetBytes(itemsBytes[1][:]).String()
	items = append(items, team)
	items = append(items, value)

	return items, errMsg
}

func DecodeBuyXidItems(data []byte) ([]string, string) {

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

	if len(itemsBytes[0]) != int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed affiliate code is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes[1]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed team is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	affCode := strconv.FormatInt(common.Decode2Int64(itemsBytes[0][:]), 10)
	team := strconv.FormatInt(common.Decode2Int64(itemsBytes[1][:]), 10)
	value := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	items = append(items, affCode)
	items = append(items, team)
	items = append(items, value)

	return items, errMsg
}

func DecodeBuyXaddrItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes) < 3 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes[1]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed team is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	affAddr := string(itemsBytes[0][:])
	team := strconv.FormatInt(common.Decode2Int64(itemsBytes[1][:]), 10)
	value := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	items = append(items, affAddr)
	items = append(items, team)
	items = append(items, value)

	return items, errMsg
}

func DecodeBuyXnameItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes) < 3 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes[1]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed team is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	affName := string(itemsBytes[0][:])
	team := strconv.FormatInt(common.Decode2Int64(itemsBytes[1][:]), 10)
	value := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	items = append(items, affName)
	items = append(items, team)
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
	items = append(items, string(itemsBytes[0][:]))

	return items, errMsg
}

func DecodeSetSettingItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}
	items = append(items, string(itemsBytes[0][:]))

	return items, errMsg
}

func DecodeGetCurrentRoundInfoItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	return items, errMsg
}

func DecodeGetBuyPriceItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	return items, errMsg
}

func DecodeGetPriceItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}
	value := new(big.Int).SetBytes(itemsBytes[0][:]).String()
	items = append(items, value)

	return items, errMsg
}

func DecodeGetKeysItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}
	value := new(big.Int).SetBytes(itemsBytes[0][:]).String()
	items = append(items, value)

	return items, errMsg
}

func DecodeGetPlayerInfoByAddressItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}
	items = append(items, string(itemsBytes[0][:]))

	return items, errMsg
}

func DecodeGetPlayerVaultItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes) < 1 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}
	items = append(items, strconv.FormatInt(common.Decode2Int64(itemsBytes[1][:]), 10))

	return items, errMsg
}

func DecodeGetTimeLeftItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	return items, errMsg
}

func DecodeWithDrawItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	return items, errMsg
}

func DecodeReloadXidItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes) < 3 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes[0]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed affiliate code is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes[1]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed team is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	affCode := strconv.FormatInt(common.Decode2Int64(itemsBytes[0][:]), 10)
	team := strconv.FormatInt(common.Decode2Int64(itemsBytes[1][:]), 10)
	value := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	items = append(items, affCode)
	items = append(items, team)
	items = append(items, value)

	return items, errMsg
}

func DecodeReloadXaddrItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes) < 3 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes[1]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed team is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	affAddr := string(itemsBytes[0][:])
	team := strconv.FormatInt(common.Decode2Int64(itemsBytes[1][:]), 10)
	value := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	items = append(items, affAddr)
	items = append(items, team)
	items = append(items, value)

	return items, errMsg
}

func DecodeReloadXnameItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes) < 3 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes[1]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed team is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	affName := string(itemsBytes[0][:])
	team := strconv.FormatInt(common.Decode2Int64(itemsBytes[1][:]), 10)
	value := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	items = append(items, affName)
	items = append(items, team)
	items = append(items, value)

	return items, errMsg
}

func DecodeRegisterNameXidItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes) < 3 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes[1]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed affiliate code is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	name := string(itemsBytes[0][:])
	affCode := strconv.FormatInt(common.Decode2Int64(itemsBytes[1][:]), 10)
	value := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	items = append(items, name)
	items = append(items, affCode)
	items = append(items, value)

	return items, errMsg
}

func DecodeRegisterNameXaddrItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	} else if len(itemsBytes) < 3 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	name := string(itemsBytes[0][:])
	affAddr := string(itemsBytes[1][:])
	value := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	items = append(items, name)
	items = append(items, affAddr)
	items = append(items, value)

	return items, errMsg
}

func DecodeRegisterNameXnameItems(data []byte) ([]string, string) {

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

	name := string(itemsBytes[0][:])
	affName := string(itemsBytes[1][:])
	value := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	items = append(items, name)
	items = append(items, affName)
	items = append(items, value)

	return items, errMsg
}
