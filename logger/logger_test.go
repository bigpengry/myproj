package logger

import (
	"fmt"
	"testing"
)

// func TestFileLogger(t *testing.T) {
// 	logger := NewFileLogger(LogLevelDebug, "/Users/david/logs", "test")
// 	logger.Debug("User ID[%d] is come from china.", 1234567)
// 	logger.Warn("test warn log")
// 	logger.Fatal("test fatal log")
// 	logger.Close()
// }

// func TestConsoleLogger(t *testing.T) {
// 	logger := NewConsoleLogger(LogLevelDebug)
// 	logger.Debug("User ID[%d] is come from china.", 1234567)
// 	logger.Warn("test warn log")
// 	logger.Fatal("test fatal log")
// 	logger.Close()
// }

func TestLog(t *testing.T) {
	var log LogInterface
	log = NewFileLogger(LogLevelDebug, "/Users/david/logs", "test")
	fmt.Printf("%+v\n", log)
}
