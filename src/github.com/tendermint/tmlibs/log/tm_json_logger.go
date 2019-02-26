package log

import (
	"io"

	kitlog "github.com/go-kit/kit/log"
)

func NewTMJSONLogger(w io.Writer) Logger {
	return &tmLogger{kitlog.NewJSONLogger(w)}
}
