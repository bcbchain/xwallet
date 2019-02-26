package amino_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-amino"
)

type SimpleStruct struct {
	String	string
	Bytes	[]byte
	Time	time.Time
}

func newSimpleStruct() SimpleStruct {
	s := SimpleStruct{
		String:	"hello",
		Bytes:	[]byte("goodbye"),
		Time:	time.Now().UTC().Truncate(time.Millisecond),
	}
	return s
}

func TestMarshalUnmarshalBinaryPointer0(t *testing.T) {

	var s = newSimpleStruct()
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinary(s)
	assert.Nil(t, err)

	var s2 SimpleStruct
	err = cdc.UnmarshalBinary(b, &s2)
	assert.Nil(t, err)
	assert.Equal(t, s, s2)

}

func TestMarshalUnmarshalBinaryPointer1(t *testing.T) {

	var s = newSimpleStruct()
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinary(&s)
	assert.Nil(t, err)

	var s2 SimpleStruct
	err = cdc.UnmarshalBinary(b, &s2)
	assert.Nil(t, err)
	assert.Equal(t, s, s2)

}

func TestMarshalUnmarshalBinaryPointer2(t *testing.T) {

	var s = newSimpleStruct()
	var ptr = &s
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinary(&ptr)
	assert.Nil(t, err)

	var s2 SimpleStruct
	err = cdc.UnmarshalBinary(b, &s2)
	assert.Nil(t, err)
	assert.Equal(t, s, s2)

}

func TestMarshalUnmarshalBinaryPointer3(t *testing.T) {

	var s = newSimpleStruct()
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinary(s)
	assert.Nil(t, err)

	var s2 *SimpleStruct
	err = cdc.UnmarshalBinary(b, &s2)
	assert.Nil(t, err)
	assert.Equal(t, s, *s2)
}

func TestMarshalUnmarshalBinaryPointer4(t *testing.T) {

	var s = newSimpleStruct()
	var ptr = &s
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinary(&ptr)
	assert.Nil(t, err)

	var s2 *SimpleStruct
	err = cdc.UnmarshalBinary(b, &s2)
	assert.Nil(t, err)
	assert.Equal(t, s, *s2)

}