package loggingUtils

import (
	"fmt"
	"strings"

	. "github.com/francoishill/goangi2/utils/debugUtils"
)

type fmtLogger struct{}

func (this *fmtLogger) Emergency(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf("[EMER] "+format, v...))
}

func (this *fmtLogger) Alert(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf("[ALER] "+format, v...))
}

func (this *fmtLogger) Critical(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf("[CRIT] "+format, v...))
}

func (this *fmtLogger) Error(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf("[ERRO] "+format, v...))
}

func (this *fmtLogger) Warning(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf("[WARN] "+format, v...))
}

func (this *fmtLogger) Notice(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf("[NOTI] "+format, v...))
}

func (this *fmtLogger) Informational(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf("[INFO] "+format, v...))
}

func (this *fmtLogger) Debug(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf("[DEBU] "+format, v...))
}

func CreateNewFmtLogger() *fmtLogger { return &fmtLogger{} }

func RecoverAndLogStackTrace_Error(messagePrefix string, loggerToUse ILogger) {
	originalRecoveredObj := recover()
	if originalRecoveredObj == nil {
		return
	}

	hadLogger := loggerToUse != nil
	l := loggerToUse
	if l == nil {
		l = CreateNewFmtLogger()
	}

	func() {
		finalPrefix := strings.TrimRight(messagePrefix, ":")
		if hadLogger {
			defer func() {
				if r2 := recover(); r2 != nil {
					fmt.Println(fmt.Sprintf("ERROR while trying to log error: %+v", r2))
					fmt.Println(fmt.Sprintf("Original error. %s: %+v. Stack trace:\n%s", finalPrefix, originalRecoveredObj, GetFullStackTrace_Pretty()))
				}
			}()
		}
		l.Warning(finalPrefix+": %+v. Stack trace:\n%s", originalRecoveredObj, GetFullStackTrace_Pretty())
	}()
}
