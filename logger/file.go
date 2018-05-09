/*日志输出到文件模块*/
package logger

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type FileLogger struct {
	level        int
	logPath      string
	logName      string
	file         *os.File
	warnFile     *os.File
	LogDataChan  chan *LogData
	logSplitType int
	logSplitSize int64
	logSplitHour int
}

func NewFileLogger(config map[string]string) (log LogInterface, err error) {
	logPath, ok := config["log_path"]
	if !ok {
		err = fmt.Errorf("not found log_path")
		return
	}
	logName, ok := config["log_name"]
	if !ok {
		err = fmt.Errorf("not found log_name")
		return
	}
	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("not found log_path")
		return
	}
	logChanSize, ok := config["log_chan_size"]
	if !ok {
		logChanSize = "50000"
	}

	var logSplitType int = LogSplitTypeHour
	var logSplitSize int64
	logSplitStr, ok := config["log_split_Type"]
	if !ok {
		logSplitStr = "Hour "
	} else {
		if logSplitStr == "Size" {
			logSplitSizeStr, ok := config["log_split_size"]
			if !ok {
				logSplitSizeStr = "104857600"
			}

			logSplitSize, err = strconv.ParseInt(logSplitSizeStr, 10, 64)
			if err != nil {
				logSplitSize = 104857600
			}

			logSplitType = LogSplitTypeSize
		} else {
			logSplitType = LogSplitTypeHour
		}
	}

	chanSize, err := strconv.Atoi(logChanSize)
	if err != nil {
		chanSize = 50000
	}
	level := getLogLevel(logLevel)
	log = &FileLogger{
		level:        level,
		logPath:      logPath,
		logName:      logName,
		LogDataChan:  make(chan *LogData, chanSize),
		logSplitType: logSplitType,
		logSplitSize: logSplitSize,
		logSplitHour: time.Now().Hour(),
	}
	log.Init()
	return
}

func (f *FileLogger) Init() {
	fileName := fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed,err:%v", fileName, err))
	}

	f.file = file

	//写错误日志和Fatal日志的文件
	fileName = fmt.Sprintf("%s/%s.wf", f.logPath, f.logName)
	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed,err:%v", fileName, err))
	}

	f.warnFile = file
	go f.WriteLogBackground()
}

//按小时切分日志
func (f *FileLogger) splitFileHour(warnFile bool) {
	now := time.Now()
	hour := now.Hour()
	if hour == f.logSplitHour {
		return
	}

	f.logSplitHour = hour
	var backupFileName string
	var fileName string
	if warnFile {
		backupFileName = fmt.Sprintf("%s/%s.wf_%04d%02d%02d%02d", f.logPath, f.logName,
			now.Year(), now.Month, now.Day(), f.logSplitHour)
		fileName = fmt.Sprintf("%s/%s.wf", f.logPath, f.logName)
	} else {
		backupFileName = fmt.Sprintf("%s/%s.log_%04d%02d%02d%02d", f.logPath, f.logName,
			now.Year(), now.Month, now.Day(), f.logSplitHour)
		fileName = fmt.Sprintf("%s/%s.wf", f.logPath, f.logName)
	}

	file := f.file
	if warnFile {
		file = f.warnFile
	}
	file.Close()
	os.Rename(fileName, backupFileName)

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}

	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}
}

//按大小切分日志
func (f *FileLogger) splitFileSize(warnFile bool) {
	now := time.Now()
	file := f.file
	if warnFile {
		f.warnFile = file
	}

	statInfo, err := file.Stat()
	if err != nil {
		return
	}

	fileSize := statInfo.Size()
	if fileSize <= f.logSplitSize {
		return
	}

	var backupFileName string
	var fileName string
	if warnFile {
		backupFileName = fmt.Sprintf("%s/%s.wf_%04d%02d%02d%02d", f.logPath, f.logName,
			now.Year(), now.Month, now.Day(), now.Hour(), now.Minute)
		fileName = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	} else {
		backupFileName = fmt.Sprintf("%s/%s.log_%04d%02d%02d%02d", f.logPath, f.logName,
			now.Year(), now.Month, now.Day(), now.Hour(), now.Minute)
		fileName = fmt.Sprintf("%s/%s.wf", f.logPath, f.logName)
	}
	file.Close()
	os.Rename(fileName, backupFileName)

	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}
	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}
}

func (f *FileLogger) checkSplitFile(warnFile bool) {
	if f.logSplitType == LogSplitTypeHour {
		f.splitFileHour(warnFile)
		return
	}

	f.splitFileSize(warnFile)
}

//后台写入
func (f *FileLogger) WriteLogBackground() {
	for logData := range f.LogDataChan {
		var file *os.File = f.file
		if logData.WarnAndFatal {
			f.file = f.warnFile
		}

		f.checkSplitFile(logData.WarnAndFatal)

		fmt.Fprintf(file, "%s %s [%s:%s:%d] %s\n", logData.TimeStr, logData.LevelStr,
			logData.FileName, logData.FuncName, logData.LineNum, logData.Message)
	}
}

func (f *FileLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		level = LogLevelDebug
	}
	f.level = level
}

func (f *FileLogger) Debug(format string, args ...interface{}) {
	if f.level > LogLevelDebug {
		return
	}
	logData := writeLog(LogLevelDebug, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	if f.level > LogLevelTrace {
		return
	}
	logData := writeLog(LogLevelTrace, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}
	logData := writeLog(LogLevelInfo, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	if f.level > LogLevelWarn {
		return
	}
	logData := writeLog(LogLevelWarn, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	if f.level > LogLevelError {
		return
	}
	logData := writeLog(LogLevelError, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	if f.level > LogLevelFatal {
		return
	}
	logData := writeLog(LogLevelFatal, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Close() {
	f.file.Close()
	f.warnFile.Close()
}
