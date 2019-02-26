package wire

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"

	"golang.org/x/crypto/ripemd160"

	cmn "github.com/tendermint/tmlibs/common"
)

func MarshalBinary(o interface{}) ([]byte, error) {
	w, n, err := new(bytes.Buffer), new(int), new(error)
	WriteBinary(o, w, n, err)
	if *err != nil {
		return nil, *err
	}
	return w.Bytes(), nil
}

func UnmarshalBinary(bz []byte, ptr interface{}) error {
	r, n, err := bytes.NewBuffer(bz), new(int), new(error)
	ReadBinaryPtr(ptr, r, len(bz), n, err)
	return *err
}

func MarshalJSON(o interface{}) ([]byte, error) {
	w, n, err := new(bytes.Buffer), new(int), new(error)
	WriteJSON(o, w, n, err)
	if *err != nil {
		return nil, *err
	}
	return w.Bytes(), nil
}

func UnmarshalJSON(bz []byte, ptr interface{}) (err error) {
	ReadJSONPtr(ptr, bz, &err)
	return
}

func BinaryBytes(o interface{}) []byte {
	w, n, err := new(bytes.Buffer), new(int), new(error)
	WriteBinary(o, w, n, err)
	if *err != nil {
		cmn.PanicSanity(*err)
	}
	return w.Bytes()
}

func ReadBinaryBytes(d []byte, ptr interface{}) error {
	r, n, err := bytes.NewBuffer(d), new(int), new(error)
	ReadBinaryPtr(ptr, r, len(d), n, err)
	return *err
}

func JSONBytes(o interface{}) []byte {
	w, n, err := new(bytes.Buffer), new(int), new(error)
	WriteJSON(o, w, n, err)
	if *err != nil {
		cmn.PanicSanity(*err)
	}
	return w.Bytes()
}

func JSONBytesPretty(o interface{}) []byte {
	jsonBytes := JSONBytes(o)
	var object interface{}
	err := json.Unmarshal(jsonBytes, &object)
	if err != nil {
		cmn.PanicSanity(err)
	}
	jsonBytes, err = json.MarshalIndent(object, "", "\t")
	if err != nil {
		cmn.PanicSanity(err)
	}
	return jsonBytes
}

func ReadJSONBytes(d []byte, ptr interface{}) (err error) {
	ReadJSONPtr(ptr, d, &err)
	return
}

func BinaryEqual(a, b interface{}) bool {
	aBytes := BinaryBytes(a)
	bBytes := BinaryBytes(b)
	return bytes.Equal(aBytes, bBytes)
}

func BinaryCompare(a, b interface{}) int {
	aBytes := BinaryBytes(a)
	bBytes := BinaryBytes(b)
	return bytes.Compare(aBytes, bBytes)
}

func BinarySha256(o interface{}) []byte {
	hasher, n, err := sha256.New(), new(int), new(error)
	WriteBinary(o, hasher, n, err)
	if *err != nil {
		cmn.PanicSanity(*err)
	}
	return hasher.Sum(nil)
}

func BinaryRipemd160(o interface{}) []byte {
	hasher, n, err := ripemd160.New(), new(int), new(error)
	WriteBinary(o, hasher, n, err)
	if *err != nil {
		cmn.PanicSanity(*err)
	}
	return hasher.Sum(nil)
}
