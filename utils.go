package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
)

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
