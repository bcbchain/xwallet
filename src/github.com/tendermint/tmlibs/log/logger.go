package log

import (
	"io"

	kitlog "github.com/go-kit/kit/log"
)

type Logger interface {
	Trace(msg string, keyvals ...interface{})
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	Fatal(msg string, keyvals ...interface{})

	With(keyvals ...interface{}) Logger
}

type Loggerf interface {
	Logger

	Tracef(fmtStr string, vals ...interface{})
	Debugf(fmtStr string, vals ...interface{})
	Infof(fmtStr string, vals ...interface{})
	Warnf(fmtStr string, vals ...interface{})
	Errorf(fmtStr string, vals ...interface{})
	Fatalf(fmtStr string, vals ...interface{})

	AllowLevel(lvl string)
	SetOutputToFile(isToFile bool)
	SetOutputToScreen(isToScreen bool)
	SetOutputAsync(isAsync bool)
	SetOutputFileSize(maxFileSize int)
	SetWithThreadID(with bool)

	Flush()
}

func NewSyncWriter(w io.Writer) io.Writer {
	return kitlog.NewSyncWriter(w)
}
