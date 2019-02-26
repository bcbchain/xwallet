package amino

import (
	"encoding/binary"
	"io"
	"math"
	"time"
)

func EncodeInt8(w io.Writer, i int8) (err error) {
	return EncodeVarint(w, int64(i))
}

func EncodeInt16(w io.Writer, i int16) (err error) {
	return EncodeVarint(w, int64(i))
}

func EncodeInt32(w io.Writer, i int32) (err error) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(i))
	_, err = w.Write(buf[:])
	return
}

func EncodeInt64(w io.Writer, i int64) (err error) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(i))
	_, err = w.Write(buf[:])
	return err
}

func EncodeVarint(w io.Writer, i int64) (err error) {
	var buf [10]byte
	n := binary.PutVarint(buf[:], i)
	_, err = w.Write(buf[0:n])
	return
}

func VarintSize(i int64) int {
	var buf [10]byte
	n := binary.PutVarint(buf[:], i)
	return n
}

func EncodeByte(w io.Writer, b byte) (err error) {
	return EncodeUvarint(w, uint64(b))
}

func EncodeUint8(w io.Writer, u uint8) (err error) {
	return EncodeUvarint(w, uint64(u))
}

func EncodeUint16(w io.Writer, u uint16) (err error) {
	return EncodeUvarint(w, uint64(u))
}

func EncodeUint32(w io.Writer, u uint32) (err error) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], u)
	_, err = w.Write(buf[:])
	return
}

func EncodeUint64(w io.Writer, u uint64) (err error) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], u)
	_, err = w.Write(buf[:])
	return
}

func EncodeUvarint(w io.Writer, u uint64) (err error) {
	var buf [10]byte
	n := binary.PutUvarint(buf[:], u)
	_, err = w.Write(buf[0:n])
	return
}

func UvarintSize(u uint64) int {
	var buf [10]byte
	n := binary.PutUvarint(buf[:], u)
	return n
}

func EncodeBool(w io.Writer, b bool) (err error) {
	if b {
		err = EncodeUint8(w, 1)
	} else {
		err = EncodeUint8(w, 0)
	}
	return
}

func EncodeFloat32(w io.Writer, f float32) (err error) {
	return EncodeUint32(w, math.Float32bits(f))
}

func EncodeFloat64(w io.Writer, f float64) (err error) {
	return EncodeUint64(w, math.Float64bits(f))
}

func EncodeTime(w io.Writer, t time.Time) (err error) {
	var s = t.Unix()
	var ns = int32(t.Nanosecond())

	err = encodeFieldNumberAndTyp3(w, 1, Typ3_8Byte)
	if err != nil {
		return
	}
	err = EncodeInt64(w, s)
	if err != nil {
		return
	}

	err = encodeFieldNumberAndTyp3(w, 2, Typ3_4Byte)
	if err != nil {
		return
	}
	err = EncodeInt32(w, ns)
	if err != nil {
		return
	}

	err = EncodeByte(w, byte(0x04))
	return
}

func EncodeByteSlice(w io.Writer, bz []byte) (err error) {
	err = EncodeUvarint(w, uint64(len(bz)))
	if err != nil {
		return
	}
	_, err = w.Write(bz)
	return
}

func ByteSliceSize(bz []byte) int {
	return UvarintSize(uint64(len(bz))) + len(bz)
}

func EncodeString(w io.Writer, s string) (err error) {
	return EncodeByteSlice(w, []byte(s))
}
