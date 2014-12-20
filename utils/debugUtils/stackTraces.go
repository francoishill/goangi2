package debugUtils

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
)

func GetFullStackTrace_Normal(getAllGoRoutines bool) string {
	stackBuf := make([]byte, 1<<16)
	runtime.Stack(stackBuf, getAllGoRoutines) //False to not print all goroutines
	stackBuf = bytes.Trim(stackBuf, "\x00")
	return string(stackBuf)
}

func GetFullStackTrace_Pretty() string {
	var buf bytes.Buffer
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		buf.WriteString(fmt.Sprintln(fmt.Sprintf("%s:%d", file, line)))
	}
	return buf.String()
}

func CatchPrintAndExitIfStartupError(errorCode int) {
	if r := recover(); r != nil {
		fmt.Println(fmt.Sprintf("Unable to start application, ERROR:\n%+v.\nSTACK TRACE:\n%s", r, GetFullStackTrace_Pretty()))
		os.Exit(errorCode)
	}
}
