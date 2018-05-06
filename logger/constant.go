package logger

const (
	LogLevelDebug = iota
	LogLevelTrace
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
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
