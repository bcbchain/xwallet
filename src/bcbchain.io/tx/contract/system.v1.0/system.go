package system_v1_0

import (
	"encoding/binary"
	"encoding/hex"
	atm "bcbchain.io/algorithm"
	"bcbchain.io/prototype"
	"bcbchain.io/rlp"

	"bcbchain.io/tx/tx"
	"bcbchain.io/utils"
	"github.com/tendermint/go-crypto"
	"math/big"
	"strconv"
	"unsafe"
)

var version string = ""

func SysNewValidator(nonce, gasLimit, note, name, pubkey, rewardaddr, power, privateKey string) string {
	items := make([]string, 0)
	items = append(items, name)

	items = append(items, pubkey)
	items = append(items, rewardaddr)
	items = append(items, power)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		atm.CalcContractAddress(crypto.GetChainId(), tx.Owner, prototype.System, version),
		utils.BytesToHex(atm.CalcMethodId(prototype.SysNewValidator)),
		items,
		privateKey)
}

func SysSetValidatorPower(nonce, gasLimit, note, pubkey, power, privateKey string) string {
	items := make([]string, 0)
	items = append(items, pubkey)
	items = append(items, power)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		atm.CalcContractAddress(crypto.GetChainId(), tx.Owner, prototype.System, version),
		utils.BytesToHex(atm.CalcMethodId(prototype.SysSetPower)),
		items,
		privateKey)
}
func SysSetValidatorRewardAddr(nonce, gasLimit, note, pubkey, reward, privateKey string) string {
	items := make([]string, 0)
	items = append(items, pubkey)
	items = append(items, reward)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		atm.CalcContractAddress(crypto.GetChainId(), tx.Owner, prototype.System, version),
		utils.BytesToHex(atm.CalcMethodId(prototype.SysSetRewardAddr)),
		items,
		privateKey)
}

func SysSetPower(nonce, gasLimit, note, pubkey, power, privateKey string) string {
	items := make([]string, 0)
	items = append(items, pubkey)
	items = append(items, power)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		atm.CalcContractAddress(crypto.GetChainId(), tx.Owner, prototype.System, version),
		utils.BytesToHex(atm.CalcMethodId(prototype.SysSetPower)),
		items,
		privateKey)
}

func SysSetRewardAddr(nonce, gasLimit, note, pubkey, rewardaddr, privateKey string) string {
	items := make([]string, 0)
	items = append(items, pubkey)
	items = append(items, rewardaddr)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		atm.CalcContractAddress(crypto.GetChainId(), tx.Owner, prototype.System, version),
		utils.BytesToHex(atm.CalcMethodId(prototype.SysSetRewardAddr)),
		items,
		privateKey)
}

func SysSetRewardStrategy(nonce, gasLimit, note, strategy, effectHeight, privateKey string) string {
	items := make([]string, 0)
	items = append(items, strategy)
	items = append(items, effectHeight)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		atm.CalcContractAddress(crypto.GetChainId(), tx.Owner, prototype.System, version),
		utils.BytesToHex(atm.CalcMethodId(prototype.SysSetRewardStrategy)),
		items,
		privateKey)
}

func SysDeployInternalContract(nonce, gasLimit, note, smcAddress, name, version, contractPrototype, gas, codeHash, effectHeight, privateKey string) string {
	items := make([]string, 0)
	items = append(items, name)
	items = append(items, version)
	items = append(items, contractPrototype)
	items = append(items, gas)
	items = append(items, codeHash)
	items = append(items, effectHeight)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SysDeployInternalContract)),
		items,
		privateKey,
	)
}

func SysForbidInternalContract(nonce, gasLimit, note, smcAddress, contractAddr, effectHeight, privateKey string) string {
	items := make([]string, 0)
	items = append(items, contractAddr)
	items = append(items, effectHeight)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SysForbidInternalContract)),
		items,
		privateKey,
	)
}

func DecodeNewValidatorItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes) < 4 {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed param is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	if len(itemsBytes[3]) != int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed gasprice is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, string(itemsBytes[0][:]))
	items = append(items, hex.EncodeToString(itemsBytes[1][:]))
	items = append(items, string(itemsBytes[2][:]))
	items = append(items, strconv.FormatUint(binary.BigEndian.Uint64(itemsBytes[3][:]), 10))
	return items, errMsg
}

func DecodeSetPowerItems(data []byte) ([]string, string) {

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

	if len(itemsBytes[1]) != int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed gasprice is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, hex.EncodeToString(itemsBytes[0][:]))
	items = append(items, strconv.FormatUint(binary.BigEndian.Uint64(itemsBytes[1][:]), 10))
	return items, errMsg
}

func DecodeSetRewardAddrItems(data []byte) ([]string, string) {

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

	items = append(items, "0x"+hex.EncodeToString(itemsBytes[0][:]))
	items = append(items, string(itemsBytes[1][:]))
	return items, errMsg
}

func DecodeForbidInternalContractItems(data []byte) ([]string, string) {

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

	if len(itemsBytes[1]) != int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed gasprice is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	items = append(items, string(itemsBytes[0][:]))
	items = append(items, strconv.FormatUint(binary.BigEndian.Uint64(itemsBytes[1][:]), 10))
	return items, errMsg
}

func DecodeDeployInternalContractItems(data []byte) ([]string, string) {

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

	if len(itemsBytes[5]) != int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed gasprice is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	name := string(itemsBytes[0][:])
	version := string(itemsBytes[1][:])
	protoTypes := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	gas := string(itemsBytes[3][:])
	codeHash := string(itemsBytes[4][:])
	effectHeight := strconv.FormatUint(binary.BigEndian.Uint64(itemsBytes[5][:]), 10)

	items = append(items, name)
	items = append(items, version)
	items = append(items, protoTypes)
	items = append(items, gas)
	items = append(items, codeHash)
	items = append(items, effectHeight)
	return items, errMsg
}

func DecodeSetRewardStrategyItems(data []byte) ([]string, string) {

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
	items = append(items, strconv.FormatUint(binary.BigEndian.Uint64(itemsBytes[1][:]), 10))
	return items, errMsg
}
