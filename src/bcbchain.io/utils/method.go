package utils

import (
	"bytes"
	"encoding/binary"
	"bcbchain.io/algorithm"
	"bcbchain.io/keys"
	"bcbchain.io/rlp"
	"math/big"
)

type MethodInfo struct {
	MethodID	uint32
	ParamData	[]byte
}

type TransferParam struct {
	To	keys.Address
	Value	*big.Int
}

func GetTransferMethodAndParam(Address keys.Address, Value *big.Int) ([]byte, error) {

	var methodInfo MethodInfo

	transferParam := TransferParam{Address, Value}
	var err error
	methodInfo.ParamData, err = rlp.EncodeToBytes(transferParam)
	if err != nil {
		panic(err)
	}

	bytesBuffer := bytes.NewBuffer(algorithm.CalcMethodId("Transfer(smc.Address,big.Int)smc.Error"))
	binary.Read(bytesBuffer, binary.BigEndian, &methodInfo.MethodID)

	data, err := rlp.EncodeToBytes(methodInfo)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetMethodAndParam(Address keys.Address, Value *big.Int, MethondName string) ([]byte, error) {

	var methodInfo MethodInfo

	transferParam := TransferParam{Address, Value}
	var err error
	methodInfo.ParamData, err = rlp.EncodeToBytes(transferParam)
	if err != nil {
		panic(err)
	}

	bytesBuffer := bytes.NewBuffer(algorithm.CalcMethodId(MethondName))
	binary.Read(bytesBuffer, binary.BigEndian, &methodInfo.MethodID)

	data, err := rlp.EncodeToBytes(methodInfo)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetNewTokenMethodAndParam(Address keys.Address, Value *big.Int, MethondName string) ([]byte, error) {

	var methodInfo MethodInfo

	transferParam := TransferParam{Address, Value}
	var err error
	methodInfo.ParamData, err = rlp.EncodeToBytes(transferParam)
	if err != nil {
		panic(err)
	}

	bytesBuffer := bytes.NewBuffer(algorithm.CalcMethodId(MethondName))
	binary.Read(bytesBuffer, binary.BigEndian, &methodInfo.MethodID)

	data, err := rlp.EncodeToBytes(methodInfo)
	if err != nil {
		return nil, err
	}

	return data, nil
}
