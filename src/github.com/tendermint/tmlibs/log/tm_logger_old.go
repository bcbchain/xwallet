package log

import (
	"fmt"
	"io"

	kitlog "github.com/go-kit/kit/log"
	kitlevel "github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
)

const (
	msgKey		= "_msg"
	moduleKey	= "module"
)

type tmLogger struct {
	srcLogger kitlog.Logger
}

var _ Logger = (*tmLogger)(nil)

func NewOldTMLogger(w io.Writer) Logger {

	colorFn := func(keyvals ...interface{}) term.FgBgColor {
		if keyvals[0] != kitlevel.Key() {
			panic(fmt.Sprintf("expected level key to be first, got %v", keyvals[0]))
		}
		switch keyvals[1].(kitlevel.Value).String() {
		case "debug":
			return term.FgBgColor{Fg: term.DarkGray}
		case "error":
			return term.FgBgColor{Fg: term.Red}
		default:
			return term.FgBgColor{}
		}
	}

	return &tmLogger{term.NewLogger(w, NewTMFmtLogger, colorFn)}
}

func NewOldTMLoggerWithColorFn(w io.Writer, colorFn func(keyvals ...interface{}) term.FgBgColor) Logger {
	return &tmLogger{term.NewLogger(w, NewTMFmtLogger, colorFn)}
}

func (l *tmLogger) Trace(msg string, keyvals ...interface{}) {
	lWithLevel := kitlevel.Debug(l.srcLogger)
	if err := kitlog.With(lWithLevel, msgKey, msg).Log(keyvals...); err != nil {
		errLogger := kitlevel.Error(l.srcLogger)
		kitlog.With(errLogger, msgKey, msg).Log("err", err)
	}
}

func (l *tmLogger) Debug(msg string, keyvals ...interface{}) {
	lWithLevel := kitlevel.Debug(l.srcLogger)
	if err := kitlog.With(lWithLevel, msgKey, msg).Log(keyvals...); err != nil {
		errLogger := kitlevel.Error(l.srcLogger)
		kitlog.With(errLogger, msgKey, msg).Log("err", err)
	}
}

func (l *tmLogger) Info(msg string, keyvals ...interface{}) {
	lWithLevel := kitlevel.Info(l.srcLogger)
	if err := kitlog.With(lWithLevel, msgKey, msg).Log(keyvals...); err != nil {
		errLogger := kitlevel.Error(l.srcLogger)
		kitlog.With(errLogger, msgKey, msg).Log("err", err)
	}
}

func (l *tmLogger) Warn(msg string, keyvals ...interface{}) {
	lWithLevel := kitlevel.Warn(l.srcLogger)
	if err := kitlog.With(lWithLevel, msgKey, msg).Log(keyvals...); err != nil {
		errLogger := kitlevel.Error(l.srcLogger)
		kitlog.With(errLogger, msgKey, msg).Log("err", err)
	}
}

func (l *tmLogger) Error(msg string, keyvals ...interface{}) {
	lWithLevel := kitlevel.Error(l.srcLogger)
	lWithMsg := kitlog.With(lWithLevel, msgKey, msg)
	if err := lWithMsg.Log(keyvals...); err != nil {
		lWithMsg.Log("err", err)
	}
}

func (l *tmLogger) Fatal(msg string, keyvals ...interface{}) {
	lWithLevel := kitlevel.Error(l.srcLogger)
	lWithMsg := kitlog.With(lWithLevel, msgKey, msg)
	if err := lWithMsg.Log(keyvals...); err != nil {
		lWithMsg.Log("err", err)
	}
}

func (l *tmLogger) With(keyvals ...interface{}) Logger {
	return &tmLogger{kitlog.With(l.srcLogger, keyvals...)}
}
