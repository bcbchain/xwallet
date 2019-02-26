package data

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

var (
	Encoder		ByteEncoder	= hexEncoder{}
	HexEncoder			= hexEncoder{}
	B64Encoder			= base64Encoder{base64.URLEncoding}
	RawB64Encoder			= base64Encoder{base64.RawURLEncoding}
)

type Bytes []byte

func (b Bytes) Marshal() ([]byte, error) {
	return b, nil
}

func (b *Bytes) Unmarshal(data []byte) error {
	*b = data
	return nil
}

func (b Bytes) MarshalJSON() ([]byte, error) {
	return Encoder.Marshal(b)
}

func (b *Bytes) UnmarshalJSON(data []byte) error {
	ref := (*[]byte)(b)
	return Encoder.Unmarshal(ref, data)
}

func (b Bytes) Bytes() []byte {
	return b
}

func (b Bytes) String() string {
	raw, err := Encoder.Marshal(b)
	l := len(raw)
	if err != nil || l < 2 {
		return "Bytes<?????>"
	}
	return string(raw[1 : l-1])
}

type ByteEncoder interface {
	Marshal(bytes []byte) ([]byte, error)
	Unmarshal(dst *[]byte, src []byte) error
}

type hexEncoder struct{}

var _ ByteEncoder = hexEncoder{}

func (_ hexEncoder) Unmarshal(dst *[]byte, src []byte) (err error) {
	var s string
	err = json.Unmarshal(src, &s)
	if err != nil {
		return errors.Wrap(err, "parse string")
	}

	*dst, err = hex.DecodeString(s)
	return err
}

func (_ hexEncoder) Marshal(bytes []byte) ([]byte, error) {
	s := strings.ToUpper(hex.EncodeToString(bytes))
	return json.Marshal(s)
}

type base64Encoder struct {
	*base64.Encoding
}

var _ ByteEncoder = base64Encoder{}

func (e base64Encoder) Unmarshal(dst *[]byte, src []byte) (err error) {
	var s string
	err = json.Unmarshal(src, &s)
	if err != nil {
		return errors.Wrap(err, "parse string")
	}
	*dst, err = e.DecodeString(s)
	return err
}

func (e base64Encoder) Marshal(bytes []byte) ([]byte, error) {
	s := e.EncodeToString(bytes)
	return json.Marshal(s)
}
