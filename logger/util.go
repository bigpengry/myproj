package logger

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

type LogData struct {
	Message      string
	TimeStr      string
	LevelStr     string
	FileName     string
	FuncName     string
	LineNum      int
	WarnAndFatal bool
}

func GetLineInfo() (fileName string, funcName string, lineNum int) {
	pc, file, line, ok := runtime.Caller(4)
	if ok {
		fileName = file
		funcName = runtime.FuncForPC(pc).Name()
		lineNum = line
	}
	return
}

/*
1.当业务调用打印日志时，我们把日志相关的数据写入到chan
2.然后我们有一个后台的线程，不断的从chan读取数据
*/

func writeLog(level int, format string, args ...interface{}) *LogData {
	now := time.Now()
	nowStr := now.Format("2006/01/02 15:04:05.999 ")
	levelStr := getLevelText(level)
	fileName, funcName, lineNum := GetLineInfo()
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	msg := fmt.Sprintf(format, args...)

	logData := &LogData{
		Message:      msg,
		TimeStr:      nowStr,
		LevelStr:     levelStr,
		FileName:     fileName,
		FuncName:     funcName,
		LineNum:      lineNum,
		WarnAndFatal: false,
	}

	if level >= LogLevelWarn {
		logData.WarnAndFatal = true
	}
	return logData
	//fmt.Fprintf(file, "%s %s [%s:%s:%d] %s\n", nowStr, levelStr, fileName, funcName, lineNum, msg)
}
