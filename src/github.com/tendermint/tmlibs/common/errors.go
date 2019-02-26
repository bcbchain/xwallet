package common

import (
	"fmt"
	"runtime"
)

func ErrorWrap(cause interface{}, format string, args ...interface{}) Error {
	msg := Fmt(format, args...)
	if causeCmnError, ok := cause.(*cmnError); ok {
		return causeCmnError.TraceFrom(1, msg)
	}

	return newError(msg, cause, cause).Stacktrace()
}

type Error interface {
	Error() string
	Message() string
	Stacktrace() Error
	Trace(format string, args ...interface{}) Error
	TraceFrom(offset int, format string, args ...interface{}) Error
	Cause() interface{}
	WithT(t interface{}) Error
	T() interface{}
	Format(s fmt.State, verb rune)
}

func NewError(format string, args ...interface{}) Error {
	msg := Fmt(format, args...)
	return newError(msg, nil, format)

}

func NewErrorWithT(t interface{}, format string, args ...interface{}) Error {
	msg := Fmt(format, args...)
	return newError(msg, nil, t)
}

type WithCauser interface {
	WithCause(cause interface{}) Error
}

type cmnError struct {
	msg		string
	cause		interface{}
	t		interface{}
	msgtraces	[]msgtraceItem
	stacktrace	[]uintptr
}

var _ WithCauser = &cmnError{}
var _ Error = &cmnError{}

func newError(msg string, cause interface{}, t interface{}) *cmnError {
	return &cmnError{
		msg:		msg,
		cause:		cause,
		t:		t,
		msgtraces:	nil,
		stacktrace:	nil,
	}
}

func (err *cmnError) Message() string {
	return err.msg
}

func (err *cmnError) Error() string {
	return fmt.Sprintf("%v", err)
}

func (err *cmnError) Stacktrace() Error {
	if err.stacktrace == nil {
		var offset = 3
		var depth = 32
		err.stacktrace = captureStacktrace(offset, depth)
	}
	return err
}

func (err *cmnError) Trace(format string, args ...interface{}) Error {
	msg := Fmt(format, args...)
	return err.doTrace(msg, 0)
}

func (err *cmnError) TraceFrom(offset int, format string, args ...interface{}) Error {
	msg := Fmt(format, args...)
	return err.doTrace(msg, offset)
}

func (err *cmnError) Cause() interface{} {
	return err.cause
}

func (err *cmnError) WithCause(cause interface{}) Error {
	err.cause = cause
	return err
}

func (err *cmnError) WithT(t interface{}) Error {
	err.t = t
	return err
}

func (err *cmnError) T() interface{} {
	return err.t
}

func (err *cmnError) doTrace(msg string, n int) Error {
	pc, _, _, _ := runtime.Caller(n + 2)

	err.msgtraces = append(err.msgtraces, msgtraceItem{
		pc:	pc,
		msg:	msg,
	})
	return err
}

func (err *cmnError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", &err)))
	default:
		if s.Flag('#') {
			s.Write([]byte("--= Error =--\n"))

			s.Write([]byte(fmt.Sprintf("Message: %s\n", err.msg)))

			s.Write([]byte(fmt.Sprintf("Cause: %#v\n", err.cause)))

			s.Write([]byte(fmt.Sprintf("T: %#v\n", err.t)))

			s.Write([]byte(fmt.Sprintf("Msg Traces:\n")))
			for i, msgtrace := range err.msgtraces {
				s.Write([]byte(fmt.Sprintf(" %4d  %s\n", i, msgtrace.String())))
			}

			if err.stacktrace != nil {
				s.Write([]byte(fmt.Sprintf("Stack Trace:\n")))
				for i, pc := range err.stacktrace {
					fnc := runtime.FuncForPC(pc)
					file, line := fnc.FileLine(pc)
					s.Write([]byte(fmt.Sprintf(" %4d  %s:%d\n", i, file, line)))
				}
			}
			s.Write([]byte("--= /Error =--\n"))
		} else {

			if err.cause != nil {
				s.Write([]byte(fmt.Sprintf("Error{`%s` (cause: %v)}", err.msg, err.cause)))
			} else {
				s.Write([]byte(fmt.Sprintf("Error{`%s`}", err.msg)))
			}
		}
	}
}

func captureStacktrace(offset int, depth int) []uintptr {
	var pcs = make([]uintptr, depth)
	n := runtime.Callers(offset, pcs)
	return pcs[0:n]
}

type msgtraceItem struct {
	pc	uintptr
	msg	string
}

func (mti msgtraceItem) String() string {
	fnc := runtime.FuncForPC(mti.pc)
	file, line := fnc.FileLine(mti.pc)
	return fmt.Sprintf("%s:%d - %s",
		file, line,
		mti.msg,
	)
}

func PanicSanity(v interface{}) {
	panic(Fmt("Panicked on a Sanity Check: %v", v))
}

func PanicCrisis(v interface{}) {
	panic(Fmt("Panicked on a Crisis: %v", v))
}

func PanicConsensus(v interface{}) {
	panic(Fmt("Panicked on a Consensus Failure: %v", v))
}

func PanicQ(v interface{}) {
	panic(Fmt("Panicked questionably: %v", v))
}
