package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

func GetLineInfo() (fileName string, funcName string, lineNum int) {
	pc, file, line, ok := runtime.Caller(3)
	if ok {
		fileName = file
		funcName = runtime.FuncForPC(pc).Name()
		lineNum = line
	}
	return
}

func writeLog(file *os.File, level int, format string, args ...interface{}) {
	now := time.Now()
	nowStr := now.Format("2006/01/02 15:04:05.999 ")
	levelStr := getLevelText(level)
	fileName, funcName, lineNum := GetLineInfo()
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	msg := fmt.Sprintf(format, args...)

	fmt.Fprintf(file, "%s %s [%s:%s:%d] %s\n", nowStr, levelStr, fileName, funcName, lineNum, msg)
}
