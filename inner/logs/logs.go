package logs

func Debug(msg string, args ...any) {
	innerLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	innerLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	innerLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	innerLogger.Error(msg, args...)
}
