package MyLog

import (
	"fmt"
	"log"
	"sync"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota + 0
	InfoLevel
	ErrorLevel
	NormalLevel
)

var (
	currentLogLevel LogLevel
	mu              sync.Mutex
)

func SetLogLevel(level LogLevel) {
	mu.Lock()
	defer mu.Unlock()
	currentLogLevel = level
}

func levelName(level LogLevel) string {
	switch level {
	case DebugLevel:
		return "Debug"
	case InfoLevel:
		return "Info"
	case ErrorLevel:
		return "Error"
	case NormalLevel:
		return "Normal"
	default:
		return "Unknown"
	}
}

func logf(level LogLevel, color string, format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	if level < currentLogLevel {
		return
	}

	// 拼接颜色代码
	prefix := ""
	suffix := ""
	if color != "" {
		prefix = "\033[" + color + "m"
		suffix = "\033[0m" // 重置颜色
	}

	//输出日志
	log.Printf("%s[%s] %s%s\n", prefix, levelName(level), fmt.Sprintf(format, v...), suffix)

}

// 各级别日志输出函数
func DebugF(format string, v ...interface{}) {
	logf(DebugLevel, "32", format, v...) // 32=绿色
}

func InfoF(format string, v ...interface{}) {
	logf(InfoLevel, "34", format, v...) // 34=蓝色
}

func ErrorF(format string, v ...interface{}) {
	logf(ErrorLevel, "31", format, v...) // 31=红色
}

func NormalF(format string, v ...interface{}) {
	logf(NormalLevel, "", format, v...) // 无颜色
}
