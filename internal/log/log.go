package log

var debugEnabled bool

func Init(debug bool) {
	debugEnabled = debug
}

func Debug(msg string, kv ...interface{}) {}
func Info(msg string, kv ...interface{})  {}
func Error(msg string, kv ...interface{}) {}
