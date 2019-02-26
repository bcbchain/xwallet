package common

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"bcbchain.io/fs"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

func ToHex(val uint64) string {
	valBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(valBytes, val)
	return string("0x") + hex.EncodeToString(valBytes)
}

func BytesToHex(valBytes []byte) string {
	return string("0x") + hex.EncodeToString(valBytes)
}

func FloatToHex(val string) string {
	resultInt := FloatStrToBigInt(val)
	return BytesToHex(resultInt.Bytes())
}

func ParseHexString(hexStr string, fieldName string, lenConstraint int) error {
	if len(hexStr)%2 != 0 {
		return errors.New(fieldName + " must be hex string with even length")
	}
	hexBytes, _ := hex.DecodeString(hexStr)
	if lenConstraint > 0 && len(hexBytes) != lenConstraint {
		return errors.New(fieldName + " must be " + strconv.Itoa(lenConstraint*2) + " hex-chars")
	}
	return nil
}

func UintToHex(val uint64) string {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	return string("0x") + hex.EncodeToString(buf)
}

func UintToBigInt(val uint64) big.Int {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	return *new(big.Int).SetBytes(buf)
}

func IntToHex(val int64) string {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(val))
	return string("0x") + hex.EncodeToString(buf)
}

func IntToByte(val int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(val))
	return buf
}

func JudgeFloatStr(val string) {
	pattern := `^\d+(\.\d{0,9})?$`
	valid, err := regexp.Match(pattern, []byte(val))
	if err != nil {
		logger.Info("Regular expression error")
		panic(" Regular expression error")
	}
	if !valid {
		logger.Info(`The money is illegal. It can only be float and >= 0.000000001`)
		panic(`The money is illegal. It can only be float and >= 0.000000001`)
	}

}

func FloatStrToBigInt(val string) big.Int {
	JudgeFloatStr(val)
	var valStr []string
	valInt := big.NewInt(0)
	valFloat := big.NewInt(0)
	if strings.Contains(val, ".") {
		valStr = strings.Split(val, ".")
		var valFloatStr string
		length := len(valStr[1])
		if length > 9 {
			valFloatStr = (valStr[1])[0:9]
			length = 9
		} else {
			valFloatStr = valStr[1]
		}
		intNum, _ := strconv.ParseInt(valFloatStr, 10, 64)
		valFloat = big.NewInt(intNum)
		var mulTemp int64 = 1
		for i := 0; i < (9 - length); i++ {
			mulTemp *= 10
		}
		valFloat.Mul(valFloat, big.NewInt(mulTemp))
	} else {
		valStr = append(valStr, val)
	}

	intStr := valStr[0]
	for i := 0; i < len(intStr); i++ {
		intNum, _ := strconv.ParseInt(string(intStr[i]), 10, 64)
		valInt.Mul(valInt, big.NewInt(10))
		valInt.Add(valInt, big.NewInt(intNum))
	}
	valInt.Mul(valInt, big.NewInt(1E9))
	return *valInt.Add(valInt, valFloat)
}

func UintStringToHex(str string) string {
	gasLimitUint, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return ""
	}
	gasLimitStr := UintToHex(gasLimitUint)
	return gasLimitStr
}

func WriteToFile(path, fileName string, data []byte) (string, error) {
	if isHavePath, _ := fs.PathExists(path); !isHavePath {
		_, err := fs.MakeDir(path)
		if err != nil {
			return "", err
		}
	}

	err := ioutil.WriteFile(path+fileName, data, 0600)
	if err != nil {
		fmt.Println("Generating file failure")
		return "", errors.New(" Generating file failure")
	}
	return path + fileName, nil
}
