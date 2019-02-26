package log

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	kitlog "github.com/go-kit/kit/log"
	kitlevel "github.com/go-kit/kit/log/level"
	"github.com/go-logfmt/logfmt"
)

type tmfmtEncoder struct {
	*logfmt.Encoder
	buf	bytes.Buffer
}

func (l *tmfmtEncoder) Reset() {
	l.Encoder.Reset()
	l.buf.Reset()
}

var tmfmtEncoderPool = sync.Pool{
	New: func() interface{} {
		var enc tmfmtEncoder
		enc.Encoder = logfmt.NewEncoder(&enc.buf)
		return &enc
	},
}

type tmfmtLogger struct {
	w io.Writer
}

func NewTMFmtLogger(w io.Writer) kitlog.Logger {
	return &tmfmtLogger{w}
}

func (l tmfmtLogger) Log(keyvals ...interface{}) error {
	enc := tmfmtEncoderPool.Get().(*tmfmtEncoder)
	enc.Reset()
	defer tmfmtEncoderPool.Put(enc)

	const unknown = "unknown"
	lvl := "none"
	msg := unknown
	module := unknown

	excludeIndexes := make([]int, 0)

	for i := 0; i < len(keyvals)-1; i += 2 {

		if keyvals[i] == kitlevel.Key() {
			excludeIndexes = append(excludeIndexes, i)
			switch keyvals[i+1].(type) {
			case string:
				lvl = keyvals[i+1].(string)
			case kitlevel.Value:
				lvl = keyvals[i+1].(kitlevel.Value).String()
			default:
				panic(fmt.Sprintf("level value of unknown type %T", keyvals[i+1]))
			}

		} else if keyvals[i] == msgKey {
			excludeIndexes = append(excludeIndexes, i)
			msg = keyvals[i+1].(string)

		} else if keyvals[i] == moduleKey {
			excludeIndexes = append(excludeIndexes, i)
			module = keyvals[i+1].(string)
		}
	}

	enc.buf.WriteString(fmt.Sprintf("[%s][%s] [%s] %-44s ", time.Now().UTC().Format("2006-01-02|15:04:05.000"), time.Now().Format("-07:00"), lvl, msg))

	if module != unknown {
		enc.buf.WriteString("module=" + module + " ")
	}

	if module != unknown {
		enc.buf.WriteString("module=" + module + " ")
	}

KeyvalueLoop:
	for i := 0; i < len(keyvals)-1; i += 2 {
		for _, j := range excludeIndexes {
			if i == j {
				continue KeyvalueLoop
			}
		}

		err := enc.EncodeKeyval(keyvals[i], keyvals[i+1])
		if err == logfmt.ErrUnsupportedValueType {
			enc.EncodeKeyval(keyvals[i], fmt.Sprintf("%+v", keyvals[i+1]))
		} else if err != nil {
			return err
		}
	}

	if err := enc.EndRecord(); err != nil {
		return err
	}

	if _, err := l.w.Write(enc.buf.Bytes()); err != nil {
		return err
	}
	return nil
}
