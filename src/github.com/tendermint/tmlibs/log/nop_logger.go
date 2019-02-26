package log

type nopLogger struct{}

var _ Logger = (*nopLogger)(nil)

func NewNopLogger() Logger	{ return &nopLogger{} }

func (nopLogger) Trace(string, ...interface{})	{}
func (nopLogger) Debug(string, ...interface{})	{}
func (nopLogger) Info(string, ...interface{})	{}
func (nopLogger) Warn(string, ...interface{})	{}
func (nopLogger) Error(string, ...interface{})	{}
func (nopLogger) Fatal(string, ...interface{})	{}

func (l *nopLogger) With(...interface{}) Logger {
	return l
}
