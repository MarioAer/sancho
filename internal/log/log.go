package log

import "fmt"

var debugEnabled bool

func Init(debug bool) {
	debugEnabled = debug
}

func Debug(msg string, kv ...interface{}) {
	if !debugEnabled {
		return
	}
	fmt.Println(append([]interface{}{"DEBUG: " + msg}, kv...)...)
}

func Info(msg string, kv ...interface{})  {}
func Error(msg string, kv ...interface{}) {}
