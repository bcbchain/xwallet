package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type level byte

const (
	levelTrace	level	= 1 << iota
	levelDebug
	levelInfo
	levelWarn
	levelError
	levelFatal
)

const (
	DATEFORMAT		= "20060102150405"
	DATETIMEFORMAT		= "[2006-01-02 15:04:05.000]"
	TIMEZONEFORMAT		= "-07"
	DEFAULT_FILE_COUNT	= 10
	DEFAULT_FILE_SIZE	= 20 * 1024 * 1024
)

type filter struct {
	next		Logger
	logfile		string
	loglevel	string
	mtx		sync.Mutex
	allowed		level
	allowedKeyvals	map[keyval]level
}

type keyval struct {
	key	interface{}
	value	interface{}
}

var fl = (*filter)(nil)

func NewFilter(next Logger, options ...Option) Logger {
	l := &filter{
		next:		next,
		allowedKeyvals:	make(map[keyval]level),
	}
	for _, option := range options {
		option(l)
	}
	return l
}

func NewFileFilter(file string, lvl string, next Logger, options ...Option) Logger {
	if len(file) == 0 {
		return nil
	}
	fl = &filter{
		next:		next,
		logfile:	file,
		loglevel:	lvl,
		allowedKeyvals:	make(map[keyval]level),
	}
	for _, option := range options {
		option(fl)
	}
	go rotateRoutine(fl.logfile, fl.loglevel)
	return fl
}

func rotateRoutine(file string, level string) {
	var src, dst *os.File
	var f os.FileInfo
	var err error
	var isRotate bool
	var size int64
	var newfile string

	for {
		f, err = os.Stat(file)
		isRotate = false
		if size = f.Size(); size > DEFAULT_FILE_SIZE {
			isRotate = true
		}

		if isRotate == true {
			checkAndRemoveFile(file)

			fl.mtx.Lock()
			newfile = fmt.Sprintf("%s-%s", file, time.Now().UTC().Format(DATEFORMAT))
			if dst, err = os.OpenFile(newfile, os.O_WRONLY|os.O_CREATE, 0644); err == nil {
				if src, err = os.OpenFile(file, os.O_RDWR, 0644); err == nil {
					io.Copy(dst, src)
					src.Truncate(0)
					src.Close()
				}
				dst.Close()
			}
			fl.mtx.Unlock()
		}
		time.Sleep(10 * time.Second)
	}
	return
}

func checkAndRemoveFile(file string) {
	var count int32
	var odf os.FileInfo

	bn := filepath.Base(file)
	fl, _ := os.Stat(file)
	tm := fl.ModTime()

	dir, _ := filepath.Abs(filepath.Dir(file))
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), bn+"-") {
			count++

			if f.ModTime().Before(tm) {
				tm = f.ModTime()
				odf = f
			}
		}
	}

	if count >= DEFAULT_FILE_COUNT {
		os.Remove(dir + "/" + odf.Name())
	}
}
func rotate(file string, level string) {
	f, err := os.Stat(file)
	isRotate := false
	if err == nil {
		size := f.Size()
		if size > DEFAULT_FILE_SIZE {
			isRotate = true
		}
	}
	if isRotate == true {
		fl.mtx.Lock()
		defer fl.mtx.Unlock()
		newfile := fmt.Sprintf("%s-%s", file, time.Now().UTC().Format(DATEFORMAT))
		dst, err := os.OpenFile(newfile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return
		}
		defer dst.Close()
		src, err := os.OpenFile(file, os.O_RDWR, 0644)
		if err != nil {
			return
		}
		defer src.Close()
		if _, err = io.Copy(dst, src); err != nil {
			fmt.Println("Failed to copy source file to destination file ", err)
			return
		}
		if err = src.Truncate(0); err != nil {
			fmt.Println("Failed to empty source file ")
			return
		}
	}

	return
}

func (l *filter) Trace(msg string, keyvals ...interface{}) {
	levelAllowed := l.allowed&levelTrace != 0
	if !levelAllowed {
		return
	}
	l.next.Trace(msg, keyvals...)
}

func (l *filter) Debug(msg string, keyvals ...interface{}) {
	levelAllowed := l.allowed&levelDebug != 0
	if !levelAllowed {
		return
	}
	l.next.Debug(msg, keyvals...)
}

func (l *filter) Info(msg string, keyvals ...interface{}) {
	levelAllowed := l.allowed&levelInfo != 0
	if !levelAllowed {
		return
	}
	l.next.Info(msg, keyvals...)
}

func (l *filter) Warn(msg string, keyvals ...interface{}) {
	levelAllowed := l.allowed&levelWarn != 0
	if !levelAllowed {
		return
	}
	l.next.Warn(msg, keyvals...)
}

func (l *filter) Error(msg string, keyvals ...interface{}) {
	levelAllowed := l.allowed&levelError != 0
	if !levelAllowed {
		return
	}
	l.next.Error(msg, keyvals...)
}

func (l *filter) Fatal(msg string, keyvals ...interface{}) {
	levelAllowed := l.allowed&levelFatal != 0
	if !levelAllowed {
		return
	}
	l.next.Fatal(msg, keyvals...)
}

func (l *filter) With(keyvals ...interface{}) Logger {
	for i := len(keyvals) - 2; i >= 0; i -= 2 {
		for kv, allowed := range l.allowedKeyvals {
			if keyvals[i] == kv.key && keyvals[i+1] == kv.value {
				return &filter{next: l.next.With(keyvals...), allowed: allowed, allowedKeyvals: l.allowedKeyvals}
			}
		}
	}
	return &filter{next: l.next.With(keyvals...), allowed: l.allowed, allowedKeyvals: l.allowedKeyvals}
}

type Option func(*filter)

func AllowLevel(lvl string) (Option, error) {
	switch lvl {
	case "trace":
		return AllowTrace(), nil
	case "debug":
		return AllowDebug(), nil
	case "info":
		return AllowInfo(), nil
	case "warn":
		return AllowWarn(), nil
	case "error":
		return AllowError(), nil
	case "fatal":
		return AllowFatal(), nil
	case "none":
		return AllowNone(), nil
	default:
		return nil, fmt.Errorf("Expected either \"trace\", \"debug\", \"info\", \"warn\", \"error\", \"fatal\" or \"none\" level, given %s", lvl)
	}
}

func AllowAll() Option {
	return AllowTrace()
}

func AllowTrace() Option {
	return allowed(levelFatal | levelError | levelWarn | levelInfo | levelDebug | levelTrace)
}

func AllowDebug() Option {
	return allowed(levelFatal | levelError | levelWarn | levelInfo | levelDebug)
}

func AllowInfo() Option {
	return allowed(levelFatal | levelError | levelWarn | levelInfo)
}

func AllowWarn() Option {
	return allowed(levelFatal | levelError | levelWarn)
}

func AllowError() Option {
	return allowed(levelFatal | levelError)
}

func AllowFatal() Option {
	return allowed(levelFatal)
}

func AllowNone() Option {
	return allowed(0)
}

func allowed(allowed level) Option {
	return func(l *filter) { l.allowed = allowed }
}

func AllowDebugWith(key interface{}, value interface{}) Option {
	return func(l *filter) { l.allowedKeyvals[keyval{key, value}] = levelError | levelInfo | levelDebug }
}

func AllowInfoWith(key interface{}, value interface{}) Option {
	return func(l *filter) { l.allowedKeyvals[keyval{key, value}] = levelError | levelWarn | levelInfo }
}

func AllowWarnWith(key interface{}, value interface{}) Option {
	return func(l *filter) { l.allowedKeyvals[keyval{key, value}] = levelError | levelWarn }
}

func AllowErrorWith(key interface{}, value interface{}) Option {
	return func(l *filter) { l.allowedKeyvals[keyval{key, value}] = levelError }
}

func AllowNoneWith(key interface{}, value interface{}) Option {
	return func(l *filter) { l.allowedKeyvals[keyval{key, value}] = 0 }
}
