package common

import (
	"bytes"
	"encoding/binary"
	"bcbchain.io/rlp"
)

func ByteSliceToInt64(b []byte) int64 {
	buf := bytes.NewBuffer(b)
	var v int64
	binary.Read(buf, binary.BigEndian, &v)
	return v
}

func Decode2Uint64(b []byte) uint64 {

	tx8 := make([]byte, 8)
	copy(tx8[len(tx8)-len(b):], b)

	return binary.BigEndian.Uint64(tx8[:])
}

func Decode2Int64(b []byte) int64 {

	return int64(Decode2Uint64(b))
}

func DecodeNoItems(data []byte) ([]string, string) {

	var errMsg string
	var itemsBytes = make([][]byte, 0)
	var items []string

	if err := rlp.DecodeBytes(data, &itemsBytes); err != nil {
		errMsg := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return items, errMsg
	}

	return items, errMsg
}
