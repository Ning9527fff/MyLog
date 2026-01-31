package MyLog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLog(testing *testing.T) {

	testMessage := "Test Message"

	// 级别顺序: Debug(0) < Info(1) < Normal(2) < Warning(3) < Error(4)

	// 测试 DebugLevel - 应该显示所有级别
	fmt.Println("========== SetLogLevel(DebugLevel) - 显示所有 ==========")
	SetLogLevel(DebugLevel)
	DebugF("This is a debug log: %s", testMessage)
	InfoF("This is a info log")
	NormalF("This is a normal log")
	WarnF("This is a warning log")
	ErrorF("This is a error log")
	fmt.Println()

	// 测试 InfoLevel - 不显示 Debug
	fmt.Println("========== SetLogLevel(InfoLevel) - 不显示 Debug ==========")
	SetLogLevel(InfoLevel)
	DebugF("[不应显示] This is a debug log")
	InfoF("This is a info log: %s", testMessage)
	NormalF("This is a normal log")
	WarnF("This is a warning log")
	ErrorF("This is a error log")
	fmt.Println()

	// 测试 NormalLevel - 只显示 Normal, Warning, Error
	fmt.Println("========== SetLogLevel(NormalLevel) - 只显示 Normal/Warning/Error ==========")
	SetLogLevel(NormalLevel)
	DebugF("[不应显示] This is a debug log")
	InfoF("[不应显示] This is a info log")
	NormalF("This is a normal log: %s", testMessage)
	WarnF("This is a warning log")
	ErrorF("This is a error log")
	fmt.Println()

	// 测试 WarningLevel - 只显示 Warning, Error
	fmt.Println("========== SetLogLevel(WarningLevel) - 只显示 Warning/Error ==========")
	SetLogLevel(WarningLevel)
	DebugF("[不应显示] This is a debug log")
	InfoF("[不应显示] This is a info log")
	NormalF("[不应显示] This is a normal log")
	WarnF("This is a warning log: %s", testMessage)
	ErrorF("This is a error log")
	fmt.Println()

	// 测试 ErrorLevel - 只显示 Error
	fmt.Println("========== SetLogLevel(ErrorLevel) - 只显示 Error ==========")
	SetLogLevel(ErrorLevel)
	DebugF("[不应显示] This is a debug log")
	InfoF("[不应显示] This is a info log")
	NormalF("[不应显示] This is a normal log")
	WarnF("[不应显示] This is a warning log")
	ErrorF("This is a error log: %s", testMessage)
	fmt.Println("====================================")
}

// TestWarnF 测试 WarnF 黄色日志
func TestWarnF(t *testing.T) {
	fmt.Println("\n========== 测试 WarnF ==========")
	SetLogLevel(DebugLevel)
	WarnF("This is a warning log")
	WarnF("Warning: %s", "something might be wrong")
}

// TestLogSwitch 测试日志开关
func TestLogSwitch(t *testing.T) {
	fmt.Println("\n========== 测试日志开关 ==========")

	EnableLog()
	InfoF("日志已开启")

	DisableLog()
	InfoF("这条日志不应该显示")

	EnableLog()
	InfoF("日志重新开启")
}

// TestFileLog 测试文件日志功能
func TestFileLog(t *testing.T) {
	fmt.Println("\n========== 测试文件日志 ==========")

	// 清理测试目录
	defer os.RemoveAll("log")

	// 测试 InitFileLog
	err := InitFileLog()
	if err != nil {
		t.Fatalf("InitFileLog 失败: %v", err)
	}
	defer CloseLogFile()

	// 写入各种级别的日志
	DebugF("文件日志测试 - Debug")
	InfoF("文件日志测试 - Info")
	WarnF("文件日志测试 - Warn")
	ErrorF("文件日志测试 - Error")
	NormalF("文件日志测试 - Normal")

	// 验证文件是否创建
	if logFilePath == "" {
		t.Fatal("日志文件路径为空")
	}

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		t.Fatalf("日志文件不存在: %s", logFilePath)
	}

	fmt.Printf("日志文件已创建: %s\n", logFilePath)
}

// TestFileLogSwitch 测试文件日志开关
func TestFileLogSwitch(t *testing.T) {
	fmt.Println("\n========== 测试文件日志开关 ==========")

	defer os.RemoveAll("log")

	err := InitFileLog()
	if err != nil {
		t.Fatalf("InitFileLog 失败: %v", err)
	}
	defer CloseLogFile()

	InfoF("文件日志已开启")

	// 关闭文件日志
	DisableFileLog()
	InfoF("这条不应该写入文件")

	// 重新开启文件日志
	EnableFileLog()
	InfoF("文件日志重新开启")

	fmt.Println("文件日志开关测试完成")
}

// TestCustomLogFile 测试自定义日志文件路径
func TestCustomLogFile(t *testing.T) {
	fmt.Println("\n========== 测试自定义日志文件路径 ==========")

	customPath := "test_custom.log"
	defer os.Remove(customPath)

	err := SetLogFile(customPath)
	if err != nil {
		t.Fatalf("SetLogFile 失败: %v", err)
	}
	defer CloseLogFile()

	InfoF("自定义路径日志测试")
	WarnF("警告信息: %s", "测试")

	if _, err := os.Stat(customPath); os.IsNotExist(err) {
		t.Fatalf("自定义日志文件不存在: %s", customPath)
	}

	fmt.Printf("自定义日志文件已创建: %s\n", customPath)
}

// TestRecoverAndLog 测试 panic 捕获和日志记录
func TestRecoverAndLog(t *testing.T) {
	fmt.Println("\n========== 测试 RecoverAndLog ==========")

	defer os.RemoveAll("log")

	err := InitFileLog()
	if err != nil {
		t.Fatalf("InitFileLog 失败: %v", err)
	}
	defer CloseLogFile()

	// 测试捕获 panic
	func() {
		defer RecoverAndLog()
		InfoF("准备触发 panic")
		panic("这是一个测试 panic")
	}()

	// 等待日志写入
	time.Sleep(100 * time.Millisecond)

	// 验证文件内容
	content, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "PANIC") {
		t.Error("日志文件中未找到 PANIC 信息")
	}

	if !strings.Contains(contentStr, "堆栈信息") {
		t.Error("日志文件中未找到堆栈信息")
	}

	fmt.Println("RecoverAndLog 测试通过，panic 已成功捕获并记录")
}

// TestPanicF 测试 PanicF 函数
func TestPanicF(t *testing.T) {
	fmt.Println("\n========== 测试 PanicF ==========")

	defer os.RemoveAll("log")

	err := InitFileLog()
	if err != nil {
		t.Fatalf("InitFileLog 失败: %v", err)
	}
	defer CloseLogFile()

	// 使用 recover 捕获 PanicF 触发的 panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("成功捕获 PanicF: %v\n", r)
			}
		}()

		PanicF("测试 PanicF: %s", "错误信息")
	}()

	// 等待日志写入
	time.Sleep(100 * time.Millisecond)

	// 验证文件内容
	content, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "PANIC") {
		t.Error("日志文件中未找到 PANIC 信息")
	}

	if !strings.Contains(contentStr, "堆栈信息") {
		t.Error("日志文件中未找到堆栈信息")
	}

	if !strings.Contains(contentStr, "测试 PanicF: 错误信息") {
		t.Error("日志文件中未找到 panic 消息")
	}

	fmt.Println("PanicF 测试通过")
}

// TestLogFileContent 测试日志文件内容格式
func TestLogFileContent(t *testing.T) {
	fmt.Println("\n========== 测试日志文件内容格式 ==========")

	defer os.RemoveAll("log")

	err := InitFileLog()
	if err != nil {
		t.Fatalf("InitFileLog 失败: %v", err)
	}
	defer CloseLogFile()

	testMsg := "测试消息123"
	InfoF("%s", testMsg)

	// 等待写入
	time.Sleep(100 * time.Millisecond)

	content, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}

	contentStr := string(content)

	// 验证包含时间戳
	if !strings.Contains(contentStr, "/") {
		t.Error("日志文件中未找到日期分隔符")
	}

	// 验证包含日志级别
	if !strings.Contains(contentStr, "[Info]") {
		t.Error("日志文件中未找到 [Info] 标记")
	}

	// 验证包含消息内容
	if !strings.Contains(contentStr, testMsg) {
		t.Error("日志文件中未找到测试消息")
	}

	fmt.Println("日志文件内容格式正确")
	fmt.Printf("文件内容:\n%s\n", contentStr)
}

// TestConcurrentFileLog 测试并发写入文件日志
func TestConcurrentFileLog(t *testing.T) {
	fmt.Println("\n========== 测试并发文件日志 ==========")

	defer os.RemoveAll("log")

	err := InitFileLog()
	if err != nil {
		t.Fatalf("InitFileLog 失败: %v", err)
	}
	defer CloseLogFile()

	// 启动多个 goroutine 并发写入日志
	done := make(chan bool)
	for i := 0; i < 5; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				InfoF("Goroutine %d - Message %d", id, j)
			}
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 5; i++ {
		<-done
	}

	time.Sleep(100 * time.Millisecond)

	// 验证文件存在
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		t.Fatal("并发写入后日志文件不存在")
	}

	fmt.Println("并发文件日志测试完成")
}

// TestLogFileDirectory 测试日志目录创建
func TestLogFileDirectory(t *testing.T) {
	fmt.Println("\n========== 测试日志目录创建 ==========")

	defer os.RemoveAll("log")

	// 确保目录不存在
	os.RemoveAll("log")

	err := InitFileLog()
	if err != nil {
		t.Fatalf("InitFileLog 失败: %v", err)
	}
	defer CloseLogFile()

	// 验证 log 目录存在
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		t.Fatal("log 目录未创建")
	}

	// 验证文件在 log 目录下
	if !strings.HasPrefix(logFilePath, filepath.Join("log", "")) {
		t.Errorf("日志文件路径不正确: %s", logFilePath)
	}

	fmt.Printf("日志目录创建成功: %s\n", logFilePath)
}
