package MyLog

import (
	"fmt"
	"testing"
)

func TestLog(testing *testing.T) {

	testMessage := "Test Message"

	SetLogLevel(DebugLevel)
	DebugF("This is a debug log")
	InfoF("This is a info log")
	ErrorF("This is a error log")
	NormalF("This is a normal log")
	DebugF("This is a debug log: %s", testMessage)

	fmt.Println("====================================")

	SetLogLevel(InfoLevel)
	DebugF("This is a debug log")
	InfoF("This is a info log")
	ErrorF("This is a error log")
	NormalF("This is a normal log")
	InfoF("This is a info log: %s", testMessage)
	fmt.Println("====================================")

	SetLogLevel(ErrorLevel)
	DebugF("This is a debug log")
	InfoF("This is a info log")
	ErrorF("This is a error log")
	NormalF("This is a normal log")
	ErrorF("This is a error log: %s", testMessage)
	fmt.Println("====================================")

	SetLogLevel(NormalLevel)
	DebugF("This is a debug log")
	InfoF("This is a info log")
	ErrorF("This is a error log")
	NormalF("This is a normal log")
	NormalF("This is a normal log: %s", testMessage)
	fmt.Println("====================================")
}
