package logger

import (
	"crypto/rand"
	"regexp"
	"strings"
	"time"
)

var DEBUG_TARGETS []string
var rv *regexp.Regexp

func init() {
	rv = regexp.MustCompile(`.func\d+(.\d+)?\s*$`)
}

func UUID(n int) string {
	const charMap = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	buf := make([]byte, n)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	for i := 0; i < n; i++ {
		ch := buf[i]
		buf[i] = charMap[int(ch)%62]
	}
	return string(buf)

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
func Error(traceID, msg string, args ...interface{}) {
	log(traceID, Trace("[ERROR]"+msg, args...).Error())
}

//Log output message. prefix is added to every single line of log
//output for tracing purpose.
func Log(traceID, msg string, args ...interface{}) {
	log(traceID, msg, args...)
}

func Dbg(msg string, args ...interface{}) {
	if len(DEBUG_TARGETS) == 0 {
		return
	}
	var wanted bool
	caller := ""
	log := Trace("")
	for _, l := range log {
		if l != "" {
			caller = l
			break
		}
	}
	caller = rv.ReplaceAllString(caller, "")
	if DEBUG_TARGETS[0] == "*" {
		wanted = true
	} else {
		if caller == "" {
			wanted = true
		} else {
			for _, t := range DEBUG_TARGETS {
				if strings.HasSuffix(caller, t) {
					wanted = true
					break
				}
			}
		}
	}
	if wanted {
		Log("Dbg", strings.TrimSpace(caller)+"> "+msg, args...)
	}
}

func SetDebugTargets(targets string) {
	DEBUG_TARGETS = []string{}
	for _, t := range strings.Split(targets, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			DEBUG_TARGETS = append(DEBUG_TARGETS, t)
		}
	}
}

func Perf(tag string, work func()) {
	start := time.Now()
	Dbg("[EXEC]%s", tag)
	work()
	elapsed := time.Since(start).Seconds()
	Dbg("[DONE]%s (elapsed: %f)", tag, elapsed)
}
