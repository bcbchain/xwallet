package tx

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"bcbchain.io/rlp"
	"bcbchain.io/tx/tx"
	"bcbchain.io/utils"
	"github.com/tendermint/go-crypto"
	"io/ioutil"
	"strconv"
	"strings"
)

func InitWrapper(genesisFile string) error {
	genesis := &tx.GenesisParameter{}
	jsonBytes, err := ioutil.ReadFile(genesisFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBytes, genesis)
	if err != nil {
		return err
	}
	tx.Owner = genesis.AppState.Token.OwnerAddress
	crypto.SetChainId(genesis.ChainId)

	return nil
}

func PackAndSignTx(nonce, gasLimit, note, to, methodId string, items []string, privateKey string) string {

	nonceInt, err := utils.ParseHexUint64(nonce, "nonce")
	if err != nil {
		return err.Error()
	}

	gasLimitInt, err := utils.ParseHexUint64(gasLimit, "gasLimit")
	if err != nil {
		return err.Error()
	}

	toAddress, err := utils.ParseRawAddress(to, "to")
	if err != nil {
		return err.Error()
	}

	var mi tx.MethodInfo

	_, err = utils.ParseHexUint32(methodId, "methodId")
	if err != nil {
		return err.Error()
	}
	dataBytes, _ := hex.DecodeString(string([]byte(methodId[2:])))
	mi.MethodID = binary.BigEndian.Uint32(dataBytes)

	var itemsBytes = make([]([]byte), 0)
	for i, item := range items {
		var itemBytes []byte
		if strings.HasPrefix(item, "0x") {

			if strings.Contains(item, ",") {
				addrs := strings.Split(item, ",")
				var addrStr string
				for _, value := range addrs {
					if strings.HasPrefix(value, "0x") {
						addrStr += strings.TrimPrefix(value, "0x")
					}
				}
				itemBytes, err = hex.DecodeString(addrStr)
			} else {
				itemBytes, err = utils.ParseHexString(item, string("item[")+strconv.Itoa(i)+"]", 0)
			}

			if err != nil {
				return err.Error()
			}

		} else {
			itemBytes = []byte(item)
		}
		itemsBytes = append(itemsBytes, itemBytes)
	}
	mi.ParamData, err = rlp.EncodeToBytes(itemsBytes)
	if err != nil {
		return err.Error()
	}

	data, err := rlp.EncodeToBytes(mi)
	if err != nil {
		return err.Error()
	}

	ss := strings.Split(privateKey, ":")
	if len(ss) != 2 {
		panic("privateKey format error")
	}
	tx1 := tx.NewTransaction(nonceInt, gasLimitInt, note, toAddress, data)
	txStr, err := tx1.TxGen(crypto.GetChainId(), ss[0], ss[1])

	if err != nil {
		errInfo := string("{\"code\":-2, \"message\":\"tx.Transaction.TxGen failed(") + err.Error() + ")\",\"data\":\"\"}"
		return errInfo
	}
	return txStr
}
