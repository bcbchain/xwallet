package rpc

import (
	"blockchain/abciapp_v1.0/keys"
	"blockchain/abciapp_v1.0/prototype"
	atm "blockchain/algorithm"
	"blockchain/smcsdk/sdk/rlp"
	"blockchain/types"
	"encoding/binary"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"github.com/tendermint/go-crypto"
)

type BcbXTransaction struct {
	Nonce    uint64       // 交易发起者发起交易的计数值，从1开始，必须单调增长，增长步长为1。
	GasLimit uint64       // 交易发起者愿意为执行此次交易支付的GAS数量的最大值。
	Note     string       // UTF-8编码的备注信息，要求小于256个字符。
	To       keys.Address // 合约地址
	Data     []byte       // 调用智能合约所需要的参数，RLP编码格式。
}

func PackAndSignTx(nonce, gasLimit uint64, note, tokenAddress, toAddress string, value []byte, name, accessKey string) (string, error) {

	var mi MethodInfo
	var err error

	methodId := atm.CalcMethodId(prototype.TtTransfer)
	//dataBytes, _ := hex.DecodeString(string(methodId))
	mi.MethodID = binary.BigEndian.Uint32(methodId)

	var itemsBytes = make([][]byte, 0)

	itemsBytes = append(itemsBytes, []byte(toAddress))
	itemsBytes = append(itemsBytes, value)

	mi.ParamData, err = rlp.EncodeToBytes(itemsBytes)
	if err != nil {
		return "", err
	}

	data, err := rlp.EncodeToBytes(mi)
	if err != nil {
		return "", err
	}

	tx1 := NewTransaction(nonce, gasLimit, note, tokenAddress, data)
	return tx1.TxGen(name, accessKey)
}

func NewTransaction(nonce uint64, gasLimit uint64, note string, to keys.Address, data []byte) BcbXTransaction {
	tx := BcbXTransaction{
		Nonce:    nonce,
		GasLimit: gasLimit,
		Note:     note,
		To:       to,
		Data:     data,
	}
	return tx
}

// 定义生成交易的接口函数，其中tx.Data已经按RLP进行编码
//返回构造好的交易数据，MAC.Version.Payload.<1>.Signature，Payload和Signature格式是RLP编码后的HexString
func (tx *BcbXTransaction) TxGen(name, accessKey string) (string, error) {
	//RLP编码tx
	size, r, err := rlp.EncodeToReader(tx)
	if err != nil {
		return "", err
	}
	txBytes := make([]byte, size)
	_, _ = r.Read(txBytes)

	sigInfo, err := SignData(name, accessKey, txBytes)
	if err != nil {
		return "", err
	}

	//RLP编码签名信息
	size, r, err = rlp.EncodeToReader(sigInfo)
	if err != nil {
		return "", err
	}
	sigBytes := make([]byte, size)
	_, _ = r.Read(sigBytes) //转换为字节流

	txString := base58.Encode(txBytes)
	sigString := base58.Encode(sigBytes)

	MAC := string(crypto.GetChainId()) + "<tx>"
	Version := "v1"
	SignerNumber := "<1>"

	return MAC + "." + Version + "." + txString + "." + SignerNumber + "." + sigString, nil
}

func SignData(name, accessKey string, data []byte) (*types.Ed25519Sig, error) {
	if name == "" || accessKey == "" {
		return nil, errors.New("user name and accessKey cannot to te empty")
	}

	if name != "" && len(name) > 40 {
		return nil, errors.New("user name length only can be 1-40")
	}
	if len(data) <= 0 {
		return nil, errors.New("user data which wants be signed length needs more than 0")
	}

	accessKeyBytes := base58.Decode(accessKey)

	acct, err2 := db.Account(name, accessKeyBytes)
	if acct == nil {
		return nil, err2
	}

	priKey := crypto.PrivKeyEd25519FromBytes(acct.PrivateKey)
	pubKey := priKey.PubKey()

	sigInfo := types.Ed25519Sig{
		SigType:  "ed25519",
		PubKey:   pubKey.(crypto.PubKeyEd25519),
		SigValue: priKey.Sign(data).(crypto.SignatureEd25519),
	}

	return &sigInfo, nil
}
