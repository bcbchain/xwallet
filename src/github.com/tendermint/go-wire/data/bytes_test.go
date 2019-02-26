package data_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	data "github.com/tendermint/go-wire/data"
)

func TestMarshal(t *testing.T) {
	assert := assert.New(t)

	b := []byte("hello world")
	dataB := data.Bytes(b)
	b2, err := dataB.Marshal()
	assert.Nil(err)
	assert.Equal(b, b2)

	var dataB2 data.Bytes
	err = (&dataB2).Unmarshal(b)
	assert.Nil(err)
	assert.Equal(dataB, dataB2)
}

func TestEncoders(t *testing.T) {
	assert := assert.New(t)

	hex := data.HexEncoder
	b64 := data.B64Encoder
	rb64 := data.RawB64Encoder
	cases := []struct {
		encoder		data.ByteEncoder
		input, expected	[]byte
	}{

		{hex, []byte(`"1A2B3C4D"`), []byte{0x1a, 0x2b, 0x3c, 0x4d}},
		{hex, []byte(`"DE14"`), []byte{0xde, 0x14}},

		{hex, []byte(`0123`), nil},
		{hex, []byte(`"dewq12"`), nil},
		{hex, []byte(`"abc"`), nil},

		{b64, []byte(`"Zm9v"`), []byte("foo")},
		{b64, []byte(`"RCEuM3M="`), []byte("D!.3s")},

		{b64, []byte(`"D4_a--w="`), []byte{0x0f, 0x8f, 0xda, 0xfb, 0xec}},

		{b64, []byte(`"D4/a++1="`), nil},
		{b64, []byte(`0123`), nil},
		{b64, []byte(`"hey!"`), nil},
		{b64, []byte(`"abc"`), nil},

		{rb64, []byte(`"Zm9v"`), []byte("foo")},
		{rb64, []byte(`"RCEuM3M"`), []byte("D!.3s")},

		{rb64, []byte(`"D4_a--w"`), []byte{0x0f, 0x8f, 0xda, 0xfb, 0xec}},

		{rb64, []byte(`"D4/a++1"`), nil},
		{rb64, []byte(`0123`), nil},
		{rb64, []byte(`"hey!"`), nil},
		{rb64, []byte(`"abc="`), nil},
	}

	for _, tc := range cases {
		var output []byte
		err := tc.encoder.Unmarshal(&output, tc.input)
		if tc.expected == nil {
			assert.NotNil(err, tc.input)
		} else if assert.Nil(err, "%s: %+v", tc.input, err) {
			assert.Equal(tc.expected, output, tc.input)
			rev, err := tc.encoder.Marshal(tc.expected)
			if assert.Nil(err, tc.input) {
				assert.Equal(tc.input, rev)
			}
		}
	}
}

func TestString(t *testing.T) {
	assert := assert.New(t)

	hex := data.HexEncoder
	b64 := data.B64Encoder
	rb64 := data.RawB64Encoder
	cases := []struct {
		encoder		data.ByteEncoder
		expected	string
		input		[]byte
	}{

		{hex, "1A2B3C4D", []byte{0x1a, 0x2b, 0x3c, 0x4d}},
		{hex, "DE14", []byte{0xde, 0x14}},
		{b64, "RCEuM3M=", []byte("D!.3s")},
		{rb64, "D4_a--w", []byte{0x0f, 0x8f, 0xda, 0xfb, 0xec}},
	}
	for _, tc := range cases {
		data.Encoder = tc.encoder
		b := data.Bytes(tc.input)
		assert.Equal(tc.expected, b.String())
	}
}

type BData struct {
	Count	int
	Data	data.Bytes
}

type BView struct {
	Count	int
	Data	string
}

func TestBytes(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	cases := []struct {
		encoder		data.ByteEncoder
		data		data.Bytes
		expected	string
	}{
		{data.HexEncoder, []byte{0x1a, 0x2b, 0x3c, 0x4d}, "1A2B3C4D"},
		{data.B64Encoder, []byte("D!.3s"), "RCEuM3M="},
		{data.RawB64Encoder, []byte("D!.3s"), "RCEuM3M"},
	}

	for i, tc := range cases {
		data.Encoder = tc.encoder

		in := BData{Count: 15, Data: tc.data}
		d, err := json.Marshal(in)
		require.Nil(err, "%d: %+v", i, err)

		out := BData{}
		err = json.Unmarshal(d, &out)
		require.Nil(err, "%d: %+v", i, err)
		assert.Equal(in.Count, out.Count, "%d", i)
		assert.Equal(in.Data, out.Data, "%d", i)

		view := BView{}
		err = json.Unmarshal(d, &view)
		require.Nil(err, "%d: %+v", i, err)
		assert.Equal(tc.expected, view.Data)
	}
}

type Dings [5]byte

func (d Dings) MarshalJSON() ([]byte, error) {
	return data.Encoder.Marshal(d[:])
}

func (d *Dings) UnmarshalJSON(enc []byte) error {
	var ref []byte
	err := data.Encoder.Unmarshal(&ref, enc)
	copy(d[:], ref)
	return err
}

func TestByteArray(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	d := Dings{}
	copy(d[:], []byte("D!.3s"))

	cases := []struct {
		encoder		data.ByteEncoder
		data		Dings
		expected	string
	}{
		{data.HexEncoder, Dings{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}, "1A2B3C4D5E"},
		{data.B64Encoder, d, "RCEuM3M="},
		{data.RawB64Encoder, d, "RCEuM3M"},
	}

	for i, tc := range cases {
		data.Encoder = tc.encoder

		d, err := json.Marshal(tc.data)
		require.Nil(err, "%d: %+v", i, err)

		out := Dings{}
		err = json.Unmarshal(d, &out)
		require.Nil(err, "%d: %+v", i, err)
		assert.Equal(tc.data, out, "%d", i)

		view := ""
		err = json.Unmarshal(d, &view)
		require.Nil(err, "%d: %+v", i, err)
		assert.Equal(tc.expected, view)
	}

	invalid := []byte(`"food"`)
	data.Encoder = data.HexEncoder
	ding := Dings{1, 2, 3, 4, 5}
	parsed := ding
	require.Equal(ding, parsed)

	err := json.Unmarshal(invalid, &parsed)
	require.NotNil(err)
	assert.Equal(ding, parsed)
}
