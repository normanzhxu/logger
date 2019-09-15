package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
)

type logEntry struct {
	t   time.Time
	msg []string
}

//LogFilter is called when a message is being sent to the
//buffer.  If it returns nil, the message will be dropped.
//This callback function is useful for augmenting messages
//or to implement alternative log destinations.
type LogFilter func([]string) []string

type logConfig struct {
	Path string `json:"path"`
	Size int    `json:"size"`
	Keep int    `json:"keep"`
	QLen int    `json:"qlen"`
	Errs string `json:"errs"`
	hook LogFilter
	ch   chan *logEntry
	sync.RWMutex
}

var (
	debugging bool
	logCfg    logConfig
	rv        *regexp.Regexp
	mx        sync.Mutex
	termSig   chan int
)

func init() {
	logCfg.Size = 10 * 1024 * 1024 //single log file size 10M
	logCfg.Keep = 10               //keep 10 log files
	logCfg.QLen = 8192             //buffer length for log channel
	logCfg.ch = make(chan *logEntry, logCfg.QLen)
	rv = regexp.MustCompile(`.func\d+(.\d+)?\s*$`)
	termSig = make(chan int)
}

func emit(prefix string, msgs ...string) {
	logCfg.RLock()
	defer logCfg.RUnlock()
	if logCfg.hook != nil {
		fmt.Println("logCfg.hook:", logCfg.hook)
		msgs = logCfg.hook(msgs)
		if len(msgs) == 0 {
			return
		}
	}
	if logCfg.ch == nil || len(logCfg.ch) == logCfg.QLen {
		fmt.Println("logCfg.ch:", logCfg.ch)
		fmt.Println("len(logCfg.ch):", len(logCfg.ch))
		logCfg.Errs = fmt.Sprintf("%s channel blocked, %d messages dropped",
			time.Now().Format(time.RFC3339), len(msgs))
		return
	}
	le := &logEntry{t: time.Now()}
	for _, msg := range msgs {
		for _, m := range strings.Split(msg, "\n") {
			if prefix != "" {
				m = prefix + " " + m
			}
			le.msg = append(le.msg, m)
		}
		fmt.Println("le.msg:", le.msg)
	}
	logCfg.ch <- le
}

func doLog(prefix, msg string, args ...interface{}) {
	if len(args) > 0 {
		emit(prefix, fmt.Sprintf(msg, args...))
	} else {
		emit(prefix, msg)
	}
}

//Error output message with stack trace. prefix is added to every single
//line of log output for tracing purpose.
func ErrorEmit(prefix, msg string, args ...interface{}) {
	emit(prefix, Trace(msg, args...)...)
}

//Log output message. prefix is added to every single line of log
//output for tracing purpose.
func LogEmit(prefix, msg string, args ...interface{}) {
	mx.Lock()
	defer mx.Unlock()
	doLog(prefix, msg, args...)
}

//Dbg output message if current function is targeted for debugging. prefix
//is added to every single line of log output for tracing purpose.
func DbgEmit(prefix, msg string, args ...interface{}) {
	if !debugging {
		return
	}
	var caller string
	mx.Lock()
	defer mx.Unlock()
	log := Trace("")
	for _, l := range log {
		if l != "" {
			caller = l
			break
		}
	}
	caller = rv.ReplaceAllString(caller, "")
	doLog(prefix+" "+strings.TrimSpace(caller)+">", msg, args...)
}

//SetDebugging turn debugging on or off.
func SetDebugging(onoff bool) {
	logCfg.Lock()
	debugging = onoff
	logCfg.Unlock()
}

//GetDebugging get debugging switch status
func GetDebugging() bool {
	logCfg.RLock()
	defer logCfg.RUnlock()
	return debugging
}

//SetLogFile specifies desination for logs. If filePath is set to empty,
//logging will be effectively disabled.
func SetLogFile(filePath string) error {
	logCfg.Lock()
	defer logCfg.Unlock()
	logCfg.Path = strings.TrimSpace(filePath)
	if logCfg.Path == "" {
		return nil
	}
	err := os.MkdirAll(path.Dir(logCfg.Path), 0755)
	if err != nil {
		logCfg.Path = ""
		return err
	}
	f, err := os.OpenFile(logCfg.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		f.Close()
		return nil
	}
	logCfg.Path = ""
	return err
}

//SetLogRotate sets size of single log file, and the number of log files to keep.
func SetLogRotate(size, keep int) error {
	logCfg.Lock()
	defer logCfg.Unlock()
	if size <= 0 {
		return errors.New("size of log file must be positive")
	}
	logCfg.Size = size
	if keep <= 0 {
		return errors.New("number of log files to keep must be positive")
	}
	logCfg.Keep = keep
	return nil
}

//SetLogBuffer sets the size of log buffer. If logging is so frequent as to fill
//the buffer before flushing happens every 1 or 2 seconds, new messages will be
//dropped.  Dropping of log messages can be monitored by /debug/vars if exposed.
func SetLogBuffer(size int) error {
	logCfg.Lock()
	defer logCfg.Unlock()
	if size < 0 {
		return errors.New("length of log channel cannot be negative")
	}
	logCfg.QLen = size
	logCfg.ch = make(chan *logEntry, logCfg.QLen)
	return nil
}

//SetLogFilter sets log filter.
func SetLogFilter(filter LogFilter) {
	logCfg.Lock()
	logCfg.hook = filter
	logCfg.Unlock()
}

//FlushLogs ensures log messages are flushed on program termination
func FlushLogs(timeout int) {
	termSig <- 1
	select {
	case <-termSig:
	case <-time.After(time.Duration(timeout) * time.Second):
	}
}
