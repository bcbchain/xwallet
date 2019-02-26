package algorithm

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T)	{ TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestBytesToInt(c *C) {
	var testcases = []struct {
		index		int
		data		[]byte
		expectedErrCode	int
	}{
		{0, []byte{1}, 0x01},
		{0, []byte{1, 2}, 0x0102},
		{0, []byte{1, 2, 3}, 0x010203},
		{0, []byte{1, 2, 3, 4}, 0x01020304},
		{0, []byte{1, 2, 3, 4, 5}, 0x0102030405},
		{0, []byte{1, 2, 3, 4, 5, 6}, 0x010203040506},
		{0, []byte{1, 2, 3, 4, 5, 6, 7}, 0x01020304050607},
		{0, []byte{1, 2, 3, 4, 5, 6, 7, 8}, 0x0102030405060708},
		{0, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, 0x0102030405060708},
	}

	for _, test := range testcases {
		c.Check(test.expectedErrCode, Equals, BytesToInt(test.data))
	}
}
