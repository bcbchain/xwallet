package log

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type logInfo struct {
	level		string
	msg		string
	keyvals		[]interface{}
	fmt		bool
	fmtStr		string
	vals		[]interface{}
	goroutineID	string
}

type giLogger struct {
	path			string
	module			string
	allowed			level
	isOutputAsync		bool
	isOutputToScreen	bool
	isOutputToFile		bool
	isOutput		bool
	withThreadID		bool
	maxFileSize		int
	curFileSize		int
	timeZone		string
	fileName		string
	logFile			*os.File
	caStop			chan bool
	caStopped		chan bool
	cf			chan *logInfo
	mutex			*sync.Mutex
}

var _ Loggerf = (*giLogger)(nil)

func NewTMLogger(path, module string) Loggerf {
	return &giLogger{
		path:			path,
		module:			module,
		allowed:		levelFatal | levelError | levelWarn | levelInfo,
		isOutputAsync:		false,
		isOutputToScreen:	false,
		isOutputToFile:		true,
		isOutput:		true,
		maxFileSize:		DEFAULT_FILE_SIZE,
		curFileSize:		0,
		timeZone:		fmt.Sprintf("%v", time.Now().Format(TIMEZONEFORMAT)),
		fileName:		"",
		logFile:		nil,
		caStop:			make(chan bool, 1),
		caStopped:		make(chan bool, 1),
		cf:			make(chan *logInfo, 10000),
		mutex:			new(sync.Mutex),
		withThreadID:		true,
	}
}

func (log *giLogger) SetOutputToFile(isToFile bool) {
	log.isOutputToFile = isToFile
	log.isOutput = log.isOutputToFile || log.isOutputToScreen
}

func (log *giLogger) SetOutputToScreen(isToScreen bool) {
	log.isOutputToScreen = isToScreen
	log.isOutput = log.isOutputToFile || log.isOutputToScreen
}

func (log *giLogger) SetOutputFileSize(maxFileSize int) {
	if maxFileSize > 0 {
		log.maxFileSize = maxFileSize
	}
}

func (log *giLogger) SetWithThreadID(with bool) {
	log.withThreadID = with
}

func (log *giLogger) AllowLevel(lvl string) {
	switch strings.ToLower(lvl) {
	case "trace":
		log.allowed = levelFatal | levelError | levelWarn | levelInfo | levelDebug | levelTrace
	case "debug":
		log.allowed = levelFatal | levelError | levelWarn | levelInfo | levelDebug
	case "info":
		log.allowed = levelFatal | levelError | levelWarn | levelInfo
	case "warn":
		log.allowed = levelFatal | levelError | levelWarn
	case "error":
		log.allowed = levelFatal | levelError
	case "fatal":
		log.allowed = levelFatal
	case "none":
		log.allowed = 0
	default:
		fmt.Printf("Expected either \"trace\", \"debug\", \"info\", \"warn\", \"error\", \"fatal\" or \"none\" level, given %s", lvl)
	}
}

func (log *giLogger) SetOutputAsync(isAsync bool) {
	if log.isOutputAsync == isAsync {
		return
	}
	if log.isOutputAsync {
		log.isOutputAsync = false
		log.caStop <- true
		<-log.caStopped
	} else {
		log.isOutputAsync = true
		go asyncRun(log)
	}
}

func (log *giLogger) Trace(msg string, keyvals ...interface{}) {
	levelAllowed := log.allowed&levelTrace != 0
	if !levelAllowed || !log.isOutput {
		return
	}
	log.Log("Trace", msg, keyvals...)
}

func (log *giLogger) Debug(msg string, keyvals ...interface{}) {
	levelAllowed := log.allowed&levelDebug != 0
	if !levelAllowed || !log.isOutput {
		return
	}
	log.Log("Debug", msg, keyvals...)
}

func (log *giLogger) Info(msg string, keyvals ...interface{}) {
	levelAllowed := log.allowed&levelInfo != 0
	if !levelAllowed || !log.isOutput {
		return
	}
	log.Log("Info", msg, keyvals...)
}

func (log *giLogger) Warn(msg string, keyvals ...interface{}) {
	levelAllowed := log.allowed&levelWarn != 0
	if !levelAllowed || !log.isOutput {
		return
	}
	log.Log("Warn", msg, keyvals...)
}

func (log *giLogger) Error(msg string, keyvals ...interface{}) {
	levelAllowed := log.allowed&levelError != 0
	if !levelAllowed || !log.isOutput {
		return
	}
	log.Log("Error", msg, keyvals...)
}

func (log *giLogger) Fatal(msg string, keyvals ...interface{}) {
	levelAllowed := log.allowed&levelFatal != 0
	if !levelAllowed || !log.isOutput {
		return
	}
	log.Log("Fatal", msg, keyvals...)
}

func (log *giLogger) Tracef(fmtStr string, vals ...interface{}) {
	levelAllowed := log.allowed&levelTrace != 0
	if !levelAllowed || !log.isOutput {
		return
	}
	log.LogEx("Trace", fmtStr, vals...)
}

func (log *giLogger) Debugf(fmtStr string, vals ...interface{}) {
	levelAllowed := log.allowed&levelDebug != 0
	if !levelAllowed || !log.isOutput {
		return
	}
	log.LogEx("Debug", fmtStr, vals...)
}

func (log *giLogger) Infof(fmtStr string, vals ...interface{}) {
	levelAllowed := log.allowed&levelInfo != 0
	if !levelAllowed || !log.isOutput {
		return
	}
	log.LogEx("Info", fmtStr, vals...)
}

func (log *giLogger) Warnf(fmtStr string, vals ...interface{}) {
	levelAllowed := log.allowed&levelWarn != 0
	if !levelAllowed || !log.isOutput {
		return
	}
	log.LogEx("Warn", fmtStr, vals...)
}

func (log *giLogger) Errorf(fmtStr string, vals ...interface{}) {
	levelAllowed := log.allowed&levelError != 0
	if !levelAllowed || !log.isOutput {
		return
	}
	log.LogEx("Error", fmtStr, vals...)
}

func (log *giLogger) Fatalf(fmtStr string, vals ...interface{}) {
	levelAllowed := log.allowed&levelFatal != 0
	if !levelAllowed || !log.isOutput {
		return
	}
	log.LogEx("Fatal", fmtStr, vals...)
}

func GetGID() string {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	return string(b)
}

func (log *giLogger) genLogText(goroutineID, level, msg string, keyvals []interface{}) (logText string) {
	var kvs []interface{}
	var tag string

	if goroutineID != "" {
		tag = fmt.Sprintf(
			"%v[%v][%-5v][%5s] ",
			time.Now().UTC().Format(DATETIMEFORMAT),
			log.timeZone,
			level,
			goroutineID)
	} else {
		tag = fmt.Sprintf(
			"%v[%v][%-5v] ",
			time.Now().UTC().Format(DATETIMEFORMAT),
			log.timeZone,
			level)
	}

	kvs = append(kvs, tag, msg)
	kvs = append(kvs, keyvals...)
	if len(kvs)%2 != 0 {
		kvs = append(kvs, "<Missing something!>")
	}

	logText = ""
	for i, v := range kvs {
		if i <= 1 {
			logText += fmt.Sprintf("%v", v)
		} else if i%2 == 0 {
			logText += fmt.Sprintf("[%v=", v)
		} else {
			logText += fmt.Sprintf("%v]", v)
		}
		if i%2 != 0 && i != len(kvs)-1 {
			logText += fmt.Sprintf(", ")
		}
	}
	logText += "\n"
	return
}

func (log *giLogger) genLogTextEx(goroutineID, level, fmtStr string, vals []interface{}) (logText string) {
	if goroutineID != "" {
		logText = fmt.Sprintf(
			"%v[%v][%-5v][%5s] ",
			time.Now().UTC().Format(DATETIMEFORMAT),
			log.timeZone,
			level,
			goroutineID)
	} else {
		logText = fmt.Sprintf(
			"%v[%v][%-5v] ",
			time.Now().UTC().Format(DATETIMEFORMAT),
			log.timeZone,
			level)
	}
	logText += fmt.Sprintf(fmtStr, vals...)
	logText += "\n"
	return
}

func (log *giLogger) newLogFile() {

	_, err := os.Stat(log.path)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(log.path, 0775)
		}
	}

	log.fileName = log.path + "/" + log.module + time.Now().UTC().Format(DATEFORMAT) + ".log"
	log.logFile, err = os.OpenFile(log.fileName, os.O_WRONLY|os.O_CREATE, 0644)
	log.curFileSize = 0

	if err == nil {
		dup2(int(log.logFile.Fd()), int(os.Stderr.Fd()))
	}
}

func (log *giLogger) flush(logBuf *bytes.Buffer) {
	if log.isOutputToScreen {
		fmt.Printf(logBuf.String())
	}
	if log.isOutputToFile {
		if log.logFile == nil {

			log.newLogFile()
		}
		_, err := os.Stat(log.fileName)
		if err != nil {

			log.newLogFile()
		}
		if log.logFile == nil {

			if !log.isOutputToScreen {

				fmt.Printf(logBuf.String())
			}
			return
		}
		log.logFile.Write(logBuf.Bytes())
		log.logFile.Sync()

		if log.curFileSize += logBuf.Len(); log.curFileSize >= log.maxFileSize {
			log.logFile.Close()
			log.newLogFile()
		}
	} else {
		if log.logFile != nil {
			log.logFile.Close()
			log.logFile = nil
			log.curFileSize = 0
		}
	}
}

func (log *giLogger) flushMutex(logBuf *bytes.Buffer) {
	log.mutex.Lock()
	defer log.mutex.Unlock()
	log.flush(logBuf)
}

func asyncRun(log *giLogger) {
	for {
		if len(log.cf) > 0 {
			logBuf := bytes.NewBuffer(nil)
			logText := ""
			if li := <-log.cf; li.fmt == false {
				logText = log.genLogText(li.goroutineID, li.level, li.msg, li.keyvals)
			} else {
				logText = log.genLogTextEx(li.goroutineID, li.level, li.fmtStr, li.vals)
			}
			logBuf.WriteString(logText)
			log.flush(logBuf)
		} else {
			if len(log.caStop) > 0 {
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
	log.caStopped <- true
}

func (log *giLogger) Log(level, msg string, keyvals ...interface{}) {
	if log.isOutputAsync {
		li := &logInfo{
			level:	level,
			fmt:	false,
			msg:	msg,
		}
		if log.withThreadID {
			li.goroutineID = GetGID()
		}
		li.keyvals = append(li.keyvals, keyvals...)
		log.cf <- li
	} else {
		rid := ""
		if log.withThreadID {
			rid = GetGID()
		}
		logText := log.genLogText(rid, level, msg, keyvals)
		log.flushMutex(bytes.NewBuffer([]byte(logText)))
	}
}

func (log *giLogger) LogEx(level, fmtStr string, vals ...interface{}) {
	if log.isOutputAsync {
		li := &logInfo{
			level:	level,
			fmt:	true,
			fmtStr:	fmtStr,
			vals:	vals,
		}
		if log.withThreadID {
			li.goroutineID = GetGID()
		}
		log.cf <- li
	} else {
		rid := ""
		if log.withThreadID {
			rid = GetGID()
		}
		logText := log.genLogTextEx(rid, level, fmtStr, vals)
		log.flushMutex(bytes.NewBuffer([]byte(logText)))
	}
}

func (log *giLogger) Flush() {
	if log.isOutputAsync == true {
		log.SetOutputAsync(false)
	}
	if len(log.cf) > 0 {
		logBuf := bytes.NewBuffer(nil)

		for {
			if len(log.cf) > 0 {
				logText := ""
				if li := <-log.cf; li.fmt == false {
					logText = log.genLogText(li.goroutineID, li.level, li.msg, li.keyvals)
				} else {
					logText = log.genLogTextEx(li.goroutineID, li.level, li.fmtStr, li.vals)
				}
				logBuf.WriteString(logText)
			} else {
				log.flushMutex(logBuf)
				break
			}
		}
	}
}

func (log *giLogger) With(keyvals ...interface{}) Logger {
	return log
}
