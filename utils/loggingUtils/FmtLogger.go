package loggingUtils

import (
	"fmt"
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
