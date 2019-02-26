package log

import (
	"os"
	"testing"

	"github.com/go-kit/kit/log/term"
)

var (
	_testingLogger Logger
)

func TestingLogger() Logger {
	if _testingLogger != nil {
		return _testingLogger
	}

	if testing.Verbose() {
		_testingLogger = NewOldTMLogger(NewSyncWriter(os.Stdout))
	} else {
		_testingLogger = NewNopLogger()
	}

	return _testingLogger
}

func TestingLoggerWithColorFn(colorFn func(keyvals ...interface{}) term.FgBgColor) Logger {
	if _testingLogger != nil {
		return _testingLogger
	}

	if testing.Verbose() {
		_testingLogger = NewOldTMLoggerWithColorFn(NewSyncWriter(os.Stdout), colorFn)
	} else {
		_testingLogger = NewNopLogger()
	}

	return _testingLogger
}
