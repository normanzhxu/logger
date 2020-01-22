package main

import (
	"fmt"
	lg "log"
	"strings"
	"time"
)

func log(traceID string, msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	if len(traceID) > 0 {
		msg += "\nTRACE_ID:" + traceID
	}
	lg.Println(msg)
}

// Error output message with stack trace. prefix is added to every single
// line of log output for tracing purpose.
func Error(traceID, msg string, args ...interface{}) {
	log(traceID, Trace("[ERROR]"+msg, args...).Error())
}

// Log output message. prefix is added to every single line of log
// output for tracing purpose.
func Log(traceID, msg string, args ...interface{}) {
	if traceID == "Dbg" || strings.Contains(traceID, ">") {
		log(traceID, "[Debug]"+msg, args...)
	} else {
		log(traceID, "[INFO]"+msg, args...)
	}
}

// LogTaketime take time
func LogTaketime(startTime int64) int {
	endTimd := time.Now().UnixNano() / 1e6
	logtime := int(endTimd - startTime)
	return logtime
}

// Perf 计算work func 使用时间
func Perf(tag string, work func()) {
	start := time.Now()
	DbgPro("prefix", "[EXEC]%s", tag)
	work()
	elapsed := time.Since(start).Seconds()
	DbgPro("prefix", "[DONE]%s (elapsed: %f)", tag, elapsed)
}

//SetDebugPro turn debugging on or off.
func SetDebugPro(onoff bool) {
	debugging = onoff
}

// DbgPro pro
func DbgPro(prefix, msg string, args ...interface{}) {
	if !debugging {
		return
	}
	var caller string
	mx.Lock()
	defer mx.Unlock()
	log := Trace("")
	fmt.Println("log:", log)
	for _, l := range log {
		if l != "" {
			caller = l
			break
		}
	}
	fmt.Println("call:", caller)

	caller = rv.ReplaceAllString(caller, "")
	fmt.Println("caller:", caller)
	fmt.Println("strings.TrimSpace(caller):", strings.TrimSpace(caller))

	Log(prefix+" "+strings.TrimSpace(caller)+" > ", msg, args...)
}

// DEBUG_TARGETS backup
var DEBUG_TARGETS []string

// SetDebugTargets backup
func SetDebugTargets(targets string) {
	DEBUG_TARGETS = []string{}
	for _, t := range strings.Split(targets, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			DEBUG_TARGETS = append(DEBUG_TARGETS, t)
		}
	}
}

// Dbg backup
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
