package loggingUtils

type multiLogger struct {
	loggers []ILogger
}

func (this *multiLogger) Emergency(format string, v ...interface{}) {
	for _, l := range this.loggers {
		l.Emergency(format, v...)
	}
}

func (this *multiLogger) Alert(format string, v ...interface{}) {
	for _, l := range this.loggers {
		l.Alert(format, v...)
	}
}

func (this *multiLogger) Critical(format string, v ...interface{}) {
	for _, l := range this.loggers {
		l.Critical(format, v...)
	}
}

func (this *multiLogger) Error(format string, v ...interface{}) {
	for _, l := range this.loggers {
		l.Error(format, v...)
	}
}

func (this *multiLogger) Warning(format string, v ...interface{}) {
	for _, l := range this.loggers {
		l.Warning(format, v...)
	}
}

func (this *multiLogger) Notice(format string, v ...interface{}) {
	for _, l := range this.loggers {
		l.Notice(format, v...)
	}
}

func (this *multiLogger) Informational(format string, v ...interface{}) {
	for _, l := range this.loggers {
		l.Informational(format, v...)
	}
}

func (this *multiLogger) Debug(format string, v ...interface{}) {
	for _, l := range this.loggers {
		l.Debug(format, v...)
	}
}

func CreateNewMultiLogger(loggers ...ILogger) *multiLogger {
	return &multiLogger{
		loggers: loggers,
	}
}
