package logger

const (
	LogLevelDebug = iota
	LogLevelTrace
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

const (
	LogSplitTypeHour = iota
	LogSplitTypeSize
)

func getLevelText(level int) string {
	switch level {
	case LogLevelDebug:
		return "Debug"
	case LogLevelTrace:
		return "Trace"
	case LogLevelInfo:
		return "Info"
	case LogLevelWarn:
		return "Warn"
	case LogLevelError:
		return "Error"
	case LogLevelFatal:
		return "Fatal"
	}
	return "Unknow"
}

func getLogLevel(level string) int {
	switch level {
	case "Debug":
		return LogLevelDebug
	case "Trace":
		return LogLevelTrace
	case "Info":
		return LogLevelInfo
	case "Warn":
		return LogLevelWarn
	case "Error":
		return LogLevelError
	case "Fatal":
		return LogLevelFatal
	}
	return LogLevelDebug
}
