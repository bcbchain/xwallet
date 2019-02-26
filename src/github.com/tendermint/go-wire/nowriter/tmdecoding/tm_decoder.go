package tmdecoding

import "time"

type TMDecoder interface {
	DecodeBool([]byte) (b bool, bytesRead int, err error)
	DecodeFloat32([]byte) (f float32, bytesRead int, err error)
	DecodeFloat64([]byte) (f float64, bytesRead int, err error)
	DecodeInt8([]byte) (i int8, bytesRead int, err error)
	DecodeInt16([]byte) (i int16, bytesRead int, err error)
	DecodeInt32([]byte) (i int32, bytesRead int, err error)
	DecodeInt64([]byte) (i int64, bytesRead int, err error)
	DecodeOctet([]byte) (b byte, bytesRead int, err error)
	DecodeOctets([]byte) (b []byte, bytesRead int, err error)
	DecodeTime([]byte) (t time.Time, bytesRead int, err error)
	DecodeUint8([]byte) (i uint8, bytesRead int, err error)
	DecodeUint16([]byte) (i uint16, bytesRead int, err error)
	DecodeUint16s([]byte) (iz []uint16, bytesRead int, err error)
	DecodeUint32([]byte) (i uint32, bytesRead int, err error)
	DecodeUint64([]byte) (i uint64, bytesRead int, err error)
	DecodeUvarint([]byte) (i uint, bytesRead int, err error)
	DecodeVarint([]byte) (i int, bytesRead int, err error)
}
