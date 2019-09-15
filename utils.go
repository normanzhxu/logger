package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func ErrorFmt(err error) {
	fmt.Fprintln(os.Stderr, Trace(err.Error()))
}

func LogFmt(msg string, args ...interface{}) {
	msg = strings.TrimRightFunc(fmt.Sprintf(msg, args...), unicode.IsSpace)
	fmt.Println(msg)
}
func Throw(msg string, args ...interface{}) {
	panic(Trace(msg, args...))
}

func Catch(err *error, handler ...func()) {
	if e := recover(); e != nil {
		*err = e.(error)
	}
	for _, h := range handler {
		h()
	}
}

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

func Jsonify(value interface{}, indent ...string) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	padding := "    "
	if len(indent) > 0 {
		padding = indent[0]
	}
	enc.SetIndent("", padding)
	err := enc.Encode(value)
	if err != nil {
		return err.Error()
	}
	return buf.String()
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
