package logger

import (
	"log"
	"os"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info 信息日志
func Info(v ...interface{}) {
	InfoLogger.Println(v...)
}

// Warning 警告日志
func Warning(v ...interface{}) {
	WarningLogger.Println(v...)
}

// Error 错误日志
func Error(v ...interface{}) {
	ErrorLogger.Println(v...)
}

// Infof 格式化信息日志
func Infof(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

// Warningf 格式化警告日志
func Warningf(format string, v ...interface{}) {
	WarningLogger.Printf(format, v...)
}

// Errorf 格式化错误日志
func Errorf(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}
