package logger

//以下功能依赖systemconfig中的配置
var logger *Logger

func Init(filename, level string) {
	logger = NewLogger(filename, "", level)
}

func Debugf(format string, v ...interface{}) {
	logger.WithLevelf(LOG_DEBUG, format, v...)
}

func Infof(format string, v ...interface{}) {
	logger.WithLevelf(LOG_INFO, format, v...)
}

func Warningf(format string, v ...interface{}) {
	logger.WithLevelf(LOG_WARNING, format, v...)
}

func Errorf(format string, v ...interface{}) {
	logger.WithLevelf(LOG_ERROR, format, v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.WithLevelf(LOG_FATAL, format, v...)
}

func Close() {
	logger.Close()
}
