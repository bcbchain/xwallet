package wire

import (
	"io"
	"time"

	cmn "github.com/tendermint/tmlibs/common"
)

func WriteTime(t time.Time, w io.Writer, n *int, err *error) {
	nanosecs := t.UnixNano()
	millisecs := nanosecs / 1000000
	if nanosecs < 0 {
		cmn.PanicSanity("can't encode times below 1970")
	} else {
		WriteInt64(millisecs*1000000, w, n, err)
	}
}

func ReadTime(r io.Reader, n *int, err *error) time.Time {
	t := ReadInt64(r, n, err)
	if t < 0 {
		*err = ErrBinaryReadInvalidTimeNegative
		return time.Time{}
	}
	if t%1000000 != 0 {
		*err = ErrBinaryReadInvalidTimeSubMillisecond
		return time.Time{}
	}
	return time.Unix(0, t)
}
