package amino

import (
	"bytes"
	"testing"
)

func TestReadByteSliceEquality(t *testing.T) {

	var encoded []byte
	var err error
	var cdc = NewCodec()

	var testBytes = []byte("ThisIsSomeTestArray")
	encoded, err = cdc.MarshalBinary(testBytes)
	if err != nil {
		t.Error(err.Error())
	}

	var testBytes2 []byte
	err = cdc.UnmarshalBinary(encoded, &testBytes2)
	if err != nil {
		t.Error(err.Error())
	}

	if !bytes.Equal(testBytes, testBytes2) {
		t.Error("Returned the wrong bytes")
	}

}
