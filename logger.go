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
	debugging     bool
	lc            logConfig
	rv            *regexp.Regexp
	mx            sync.Mutex
	termSig       chan int
	DEBUG_TARGETS []string
)

func init() {
	lc.Size = 10 * 1024 * 1024 //single log file size 10M
	lc.Keep = 10               //keep 10 log files
	lc.QLen = 8192             //buffer length for log channel
	lc.ch = make(chan *logEntry, lc.QLen)
	rv = regexp.MustCompile(`.func\d+(.\d+)?\s*$`)
	termSig = make(chan int)
}

func emit(prefix string, msgs ...string) {
	lc.RLock()
	defer lc.RUnlock()
	if lc.hook != nil {
		fmt.Println("lc.hook:", lc.hook)
		msgs = lc.hook(msgs)
		if len(msgs) == 0 {
			return
		}
	}
	if lc.ch == nil || len(lc.ch) == lc.QLen {
		fmt.Println("lc.ch:", lc.ch)
		fmt.Println("len(lc.ch):", len(lc.ch))
		lc.Errs = fmt.Sprintf("%s channel blocked, %d messages dropped",
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
	lc.ch <- le
}

func doLog(prefix, msg string, args ...interface{}) {
	if len(args) > 0 {
		emit(prefix, fmt.Sprintf(msg, args...))
	} else {
		emit(prefix, msg)
	}
}

// func Error(err error) {
// 	fmt.Fprintln(os.Stderr, Trace(err.Error()))
// }

// func Log(msg string, args ...interface{}) {
// 	msg = strings.TrimRightFunc(fmt.Sprintf(msg, args...), unicode.IsSpace)
// 	fmt.Println(msg)
// }

//Error output message with stack trace. prefix is added to every single
//line of log output for tracing purpose.
func Error(prefix, msg string, args ...interface{}) {
	emit(prefix, Trace(msg, args...)...)
}

//Error output message with stack trace. prefix is added to every single
//line of log output for tracing purpose.
// func Error(traceID, msg string, args ...interface{}) {
// 	log(traceID, Trace("[ERROR]"+msg, args...).Error())
// }

//Log output message. prefix is added to every single line of log
//output for tracing purpose.
func Log(prefix, msg string, args ...interface{}) {
	mx.Lock()
	defer mx.Unlock()
	doLog(prefix, msg, args...)
}

//Log output message. prefix is added to every single line of log
//output for tracing purpose.
// func Log(traceID, msg string, args ...interface{}) {
// 	log(traceID, msg, args...)
// }

//Dbg output message if current function is targeted for debugging. prefix
//is added to every single line of log output for tracing purpose.
func Dbg(prefix, msg string, args ...interface{}) {
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

// func Dbg(msg string, args ...interface{}) {
// 	if len(DEBUG_TARGETS) == 0 {
// 		return
// 	}
// 	var wanted bool
// 	caller := ""
// 	log := Trace("")
// 	for _, l := range log {
// 		if l != "" {
// 			caller = l
// 			break
// 		}
// 	}
// 	caller = rv.ReplaceAllString(caller, "")
// 	if DEBUG_TARGETS[0] == "*" {
// 		wanted = true
// 	} else {
// 		if caller == "" {
// 			wanted = true
// 		} else {
// 			for _, t := range DEBUG_TARGETS {
// 				if strings.HasSuffix(caller, t) {
// 					wanted = true
// 					break
// 				}
// 			}
// 		}
// 	}
// 	if wanted {
// 		Log("Dbg", strings.TrimSpace(caller)+"> "+msg, args...)
// 	}
// }

//SetDebugging turn debugging on or off.
func SetDebugging(onoff bool) {
	lc.Lock()
	debugging = onoff
	lc.Unlock()
}

// func SetDebugTargets(targets string) {
// 	DEBUG_TARGETS = []string{}
// 	for _, t := range strings.Split(targets, ",") {
// 		t = strings.TrimSpace(t)
// 		if t != "" {
// 			DEBUG_TARGETS = append(DEBUG_TARGETS, t)
// 		}
// 	}
// }

//GetDebugging get debugging switch status
func GetDebugging() bool {
	lc.RLock()
	defer lc.RUnlock()
	return debugging
}

//SetLogFile specifies desination for logs. If filePath is set to empty,
//logging will be effectively disabled.
func SetLogFile(filePath string) error {
	lc.Lock()
	defer lc.Unlock()
	lc.Path = strings.TrimSpace(filePath)
	if lc.Path == "" {
		return nil
	}
	err := os.MkdirAll(path.Dir(lc.Path), 0755)
	if err != nil {
		lc.Path = ""
		return err
	}
	f, err := os.OpenFile(lc.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		f.Close()
		return nil
	}
	lc.Path = ""
	return err
}

//SetLogRotate sets size of single log file, and the number of log files to keep.
func SetLogRotate(size, keep int) error {
	lc.Lock()
	defer lc.Unlock()
	if size <= 0 {
		return errors.New("size of log file must be positive")
	}
	lc.Size = size
	if keep <= 0 {
		return errors.New("number of log files to keep must be positive")
	}
	lc.Keep = keep
	return nil
}

//SetLogBuffer sets the size of log buffer. If logging is so frequent as to fill
//the buffer before flushing happens every 1 or 2 seconds, new messages will be
//dropped.  Dropping of log messages can be monitored by /debug/vars if exposed.
func SetLogBuffer(size int) error {
	lc.Lock()
	defer lc.Unlock()
	if size < 0 {
		return errors.New("length of log channel cannot be negative")
	}
	lc.QLen = size
	lc.ch = make(chan *logEntry, lc.QLen)
	return nil
}

//SetLogFilter sets log filter.
func SetLogFilter(filter LogFilter) {
	lc.Lock()
	lc.hook = filter
	lc.Unlock()
}

//FlushLogs ensures log messages are flushed on program termination
func FlushLogs(timeout int) {
	termSig <- 1
	select {
	case <-termSig:
	case <-time.After(time.Duration(timeout) * time.Second):
	}
}

// 计算work func 使用时间
func Perf(tag string, work func()) {
	start := time.Now()
	Dbg("[EXEC]%s", tag)
	work()
	elapsed := time.Since(start).Seconds()
	Dbg("[DONE]%s (elapsed: %f)", tag, elapsed)
}
