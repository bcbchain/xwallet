package token_issue_v1_0

import (
	atm "bcbchain.io/algorithm"
	"bcbchain.io/prototype"
	"bcbchain.io/rlp"
	"bcbchain.io/tx/contract/common"
	"bcbchain.io/tx/tx"
	"bcbchain.io/utils"
	"github.com/tendermint/go-crypto"
	"math/big"
	"strconv"
	"unsafe"
)

var version string = ""

func TiNewToken(nonce, gasLimit, note, name, symbol, totalSupply, addSupplyEnabled, burnEnabled, initGasPrice, privateKey string) string {
	items := make([]string, 0)
	items = append(items, name)
	items = append(items, symbol)
	items = append(items, totalSupply)
	items = append(items, addSupplyEnabled)
	items = append(items, burnEnabled)
	items = append(items, initGasPrice)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		atm.CalcContractAddress(crypto.GetChainId(), tx.Owner, prototype.TokenIssue, version),
		utils.BytesToHex(atm.CalcMethodId(prototype.TiNewToken)),
		items,
		privateKey)
}

func DecodeNewTokenItems(data []byte) ([]string, string) {

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

	if len(itemsBytes[5]) > int(unsafe.Sizeof(uint64(0))) {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed gasprice is too short ") + "\",\"data\":\"\"}"
		return items, errMsg
	}

	name := string(itemsBytes[0][:])
	symbol := string(itemsBytes[1][:])
	totalSupply := new(big.Int).SetBytes(itemsBytes[2][:]).String()
	isAddSupply := string(itemsBytes[3][:])
	isBurn := string(itemsBytes[4][:])
	gasprice := strconv.FormatUint(common.Decode2Uint64(itemsBytes[5][:]), 10)

	items = append(items, name)
	items = append(items, symbol)
	items = append(items, totalSupply)
	items = append(items, isAddSupply)
	items = append(items, isBurn)
	items = append(items, gasprice)
	return items, errMsg
}
