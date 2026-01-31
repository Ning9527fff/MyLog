package MyLog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota + 0
	InfoLevel
	NormalLevel
	WarningLevel
	ErrorLevel
)

var (
	currentLogLevel LogLevel
	logEnabled      bool = true // 日志开关，默认开启
	fileLogEnabled  bool = false // 文件日志开关，默认关闭
	logFile         *os.File     // 日志文件句柄
	logFilePath     string       // 日志文件路径
	mu              sync.Mutex
)

func SetLogLevel(level LogLevel) {
	mu.Lock()
	defer mu.Unlock()
	currentLogLevel = level
}

// SetLogEnabled 设置日志开关
func SetLogEnabled(enabled bool) {
	mu.Lock()
	defer mu.Unlock()
	logEnabled = enabled
}

// EnableLog 开启日志，如果文件日志未初始化则自动初始化
func EnableLog() error {
	mu.Lock()
	defer mu.Unlock()

	logEnabled = true

	// 如果文件日志未初始化，自动初始化
	if logFile == nil {
		if err := initFileLogInternal(); err != nil {
			return err
		}
	}

	return nil
}

// DisableLog 关闭日志
func DisableLog() {
	SetLogEnabled(false)
}

// SetLogFile 设置日志文件路径并开启文件日志（自定义路径）
func SetLogFile(filePath string) error {
	mu.Lock()
	defer mu.Unlock()

	// 如果已有打开的文件，先关闭
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}

	// 打开或创建日志文件（追加模式）
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logFile = file
	logFilePath = filePath
	fileLogEnabled = true
	return nil
}

// initFileLogInternal 内部函数：初始化文件日志（不加锁）
func initFileLogInternal() error {
	// 如果已有打开的文件，先关闭
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}

	// 创建 log 目录
	logDir := "log"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 生成时间戳文件名
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("%s.log", timestamp)
	filePath := filepath.Join(logDir, fileName)

	// 打开或创建日志文件（追加模式）
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logFile = file
	logFilePath = filePath
	fileLogEnabled = true
	return nil
}

// InitFileLog 初始化文件日志，以时间戳命名，保存在 log 目录下
func InitFileLog() error {
	mu.Lock()
	defer mu.Unlock()
	return initFileLogInternal()
}

// EnableFileLog 开启文件日志（需要先调用 SetLogFile 设置文件路径）
func EnableFileLog() {
	mu.Lock()
	defer mu.Unlock()
	if logFile != nil {
		fileLogEnabled = true
	}
}

// DisableFileLog 关闭文件日志
func DisableFileLog() {
	mu.Lock()
	defer mu.Unlock()
	fileLogEnabled = false
}

// CloseLogFile 关闭日志文件
func CloseLogFile() error {
	mu.Lock()
	defer mu.Unlock()

	if logFile != nil {
		err := logFile.Close()
		logFile = nil
		fileLogEnabled = false
		logFilePath = ""
		return err
	}
	return nil
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
	case WarningLevel:
		return "Warnning"
	default:
		return "Unknown"
	}
}

func logf(level LogLevel, color string, format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	// 检查日志是否启用
	if !logEnabled {
		return
	}

	if level < currentLogLevel {
		return
	}

	prefix := ""
	suffix := ""
	if color != "" {
		prefix = "\033[" + color + "m"
		suffix = "\033[0m" // 重置颜色
	}

	// 控制台输出（带颜色）
	log.Printf("%s[%s] %s%s\n", prefix, levelName(level), fmt.Sprintf(format, v...), suffix)

	// 文件输出（不带颜色）
	if fileLogEnabled && logFile != nil {
		timestamp := time.Now().Format("2006/01/02 15:04:05")
		logMessage := fmt.Sprintf("%s [%s] %s\n", timestamp, levelName(level), fmt.Sprintf(format, v...))
		logFile.WriteString(logMessage)
	}

}

// DebugF Debug日志 绿色
func DebugF(format string, v ...interface{}) {
	logf(DebugLevel, "32", format, v...) // 32=绿色
}

// InfoF Info日志 蓝色
func InfoF(format string, v ...interface{}) {
	logf(InfoLevel, "34", format, v...) // 34=蓝色
}

// ErrorF Error日志 红色
func ErrorF(format string, v ...interface{}) {
	logf(ErrorLevel, "31", format, v...) // 31=红色
}

// NormalF Normal日志 紫色
func NormalF(format string, v ...interface{}) {
	logf(NormalLevel, "35", format, v...) // 35=紫色
}

// WarnF Warning级别日志，黄色
func WarnF(format string, v ...interface{}) {
	logf(WarningLevel, "33", format, v...) // 33=黄色
}

// PanicF 记录 Panic 日志并触发 panic
func PanicF(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)

	// 记录 panic 日志
	ErrorF("PANIC: %s", msg)

	// 获取堆栈信息并写入文件
	writePanicToFile(msg, debug.Stack())

	// 触发 panic
	panic(msg)
}

// RecoverAndLog 捕获 panic 并记录到日志，用于 defer 中
func RecoverAndLog() {
	if r := recover(); r != nil {
		panicMsg := fmt.Sprintf("%v", r)
		stack := debug.Stack()

		// 输出到控制台
		ErrorF("程序 Panic: %s", panicMsg)

		// 写入文件
		writePanicToFile(panicMsg, stack)
	}
}

// writePanicToFile 将 panic 信息和堆栈写入日志文件
func writePanicToFile(panicMsg string, stack []byte) {
	mu.Lock()
	defer mu.Unlock()

	if fileLogEnabled && logFile != nil {
		timestamp := time.Now().Format("2006/01/02 15:04:05")

		// 写入分隔线
		logFile.WriteString("\n" + strings.Repeat("=", 80) + "\n")

		// 写入 panic 信息
		logFile.WriteString(fmt.Sprintf("%s [PANIC] %s\n", timestamp, panicMsg))

		// 写入堆栈信息
		logFile.WriteString("\n堆栈信息:\n")
		logFile.WriteString(string(stack))

		// 写入结束分隔线
		logFile.WriteString(strings.Repeat("=", 80) + "\n\n")

		// 立即刷新到磁盘
		logFile.Sync()
	}
}
