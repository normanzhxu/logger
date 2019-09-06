package main

import (
	"fmt"
	"runtime"
	"strings"
	"unicode"
)

// # go-audit - Error Handling, Logging and Instrumentation for Golang

func NotNilErrorAssert(funcname string, err error) {
	if IsNotNil(err) {
		trace := fmt.Sprintf("%s_ERROR", funcname)
		Error(trace, "%s", err.Error())
		Assert(err)
	}
}

func IsNotNil(err error) bool {
	if err != nil {
		return true
	}
	return false
}

// func Throw(msg string, args ...interface{}) {
// 	panic(Trace(msg, args...))
// }
// func Catch(err *error, handler ...func()) {
// 	if e := recover(); e != nil {
// 		*err = e.(error)
// 	}
// 	for _, h := range handler {
// 		h()
// 	}
// }

// func log(traceID string, msg string, args ...interface{}) {
// 	if len(args) > 0 {
// 		msg = fmt.Sprintf(msg, args...)
// 	}
// 	if len(traceID) > 0 {
// 		msg += "\nTRACE_ID:" + traceID
// 	}
// 	lg.Println(msg)
// }

//Assert checks if err is nil or not, if not, it panics with that err.
func Assert(err error) {
	if err != nil {
		panic(err)
	}
}

//Exception error with stack trace
type exception []string

func (e exception) Error() string {
	return strings.Join(e, "\n")
}

//Trace returns the current stack trace
func Trace(msg string, args ...interface{}) (logs exception) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	msg = strings.TrimRightFunc(msg, unicode.IsSpace)
	if len(msg) > 0 {
		logs = exception{msg}
	}
	n := 1
	for {
		n++
		pc, file, line, ok := runtime.Caller(n)
		if !ok {
			break
		}
		f := runtime.FuncForPC(pc)
		name := f.Name()
		if strings.HasPrefix(name, "runtime.") {
			continue
		}
		fn := file[strings.Index(file, "/src/")+5:]
		logs = append(logs, fmt.Sprintf("\t(%s:%d) %s", fn, line, name))
	}
	return
}
