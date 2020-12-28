package logger

//log库

import (
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	LOG_DEBUG LogLevel =  iota
	LOG_INFO
	LOG_WARNING
	LOG_ERROR
	LOG_FATAL
)

func (ll LogLevel) String() string {
	switch ll {
	case LOG_DEBUG:
		return "debug"
	case LOG_INFO:
		return "info"
	case LOG_WARNING:
		return "warning"
	case LOG_ERROR:
		return "error"
	case LOG_FATAL:
		return "fatal"
	default:
		return "info"
	}
}

func getLevelByString(levelString string) LogLevel {
	switch strings.ToLower(levelString) {
	case "debug":
		return LOG_DEBUG
	case "info":
		return LOG_INFO
	case "warning":
		return LOG_WARNING
	case "error":
		return LOG_ERROR
	case "fatal":
		return LOG_FATAL
	default:
		return LOG_INFO
	}
}

type Logger struct {
	*log.Logger
	Level   LogLevel
	logFile *os.File
}

func NewLogger(fileName string, prefix string, level string) *Logger {
	//目录是否存在
	pathSliceList := strings.Split(fileName, "/")
	//不是当前目录
	if len(pathSliceList) > 1 {
		pathSliceList = pathSliceList[:len(pathSliceList)-1]
		path := strings.Join(pathSliceList, "/")
		_, err := os.Stat(path)
		//创建目录
		if err != nil && os.IsNotExist(err) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				panic("create dir failed: " + err.Error())
			}
		}
	}
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("open log file " + fileName + " error: " + err.Error())
	}
	logger := log.New(logFile, prefix, log.LstdFlags)
	logger.SetFlags(log.Ldate | log.Ltime)
	return &Logger{
		Logger:  logger,
		Level:   getLevelByString(level),
		logFile: logFile,
	}
}

func (l *Logger) WithLevelf(lev LogLevel, format string, v ...interface{}){
	if l == nil || lev < l.Level{
		return
	}
	format = "["+lev.String()+"]" + " " + format
	l.logf(format, v...)
}

func (l *Logger) logf(format string, v ...interface{}){
	l.Printf(format, v...)
}

func (l *Logger) Close() bool {
	_ = l.logFile.Close()
	return true
}