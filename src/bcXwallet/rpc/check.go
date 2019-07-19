package rpc

import (
	"blockchain/abciapp_v1.0/smc"
	"bytes"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func checkName(name string) error {
	valid, err := regexp.Match(pattern, []byte(name))
	if err != nil {
		return errors.New("Regular expression error=" + err.Error())
	}
	if !valid {
		return errors.New(`Name contains by [letters, numbers, "_", "@", "." and "-"] and length must be [1-40] `)
	}

	return nil
}

func checkPrivateKey(privateKey string, plainText bool) error {
	switch plainText {
	case true:
		if len(privateKey) != 128 && len(privateKey) != 64 {
			return errors.New("The length of privateKey is wrong ")
		}
	case false:
		if len(privateKey) != 160 {
			return errors.New("The length of privateKey is wrong ")
		}
	}

	return nil
}

func checkPassword(s string) (flag bool) {
	ascOther := ` !"#$%&'()*+,-/:;<=>?[]\^{|}~@_.` + "`"
	count := 0
	number := false
	upper := false
	lower := false
	special := false
	other := true
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
			count++
		case unicode.IsUpper(c):
			upper = true
			count++
		case unicode.IsLower(c):
			lower = true
			count++
		case strings.Contains(ascOther, string(c)):
			special = true
			count++
		default:
			other = false
		}
	}

	flag = number && upper && lower && special && other && 8 <= count && count <= 20

	return
}

// nolint
func checkAddress(chainID string, addr smc.Address) error {
	if !strings.HasPrefix(addr, chainID) {
		return errors.New("Address chainID is error! ")
	}
	base58Addr := strings.Replace(addr, chainID, "", 1)
	addrData := base58.Decode(base58Addr)
	dataLen := len(addrData)
	if dataLen < 4 {
		return errors.New("Base58Addr parse error! ")
	}

	hasher := ripemd160.New()
	hasher.Write(addrData[:dataLen-4])
	md := hasher.Sum(nil)

	if bytes.Compare(md[:4], addrData[dataLen-4:]) != 0 {
		return errors.New("Address checksum is error! ")
	}

	return nil
}

func requireUint64(valueStr string) (uint64, error) {
	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}
