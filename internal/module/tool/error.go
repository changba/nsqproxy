package tool

import (
	"github.com/changba-server/nsqproxy/internal/module/logger"
	"runtime/debug"
)

func PanicHandlerForLog() {
	if err := recover(); err != nil {
		logger.Errorf("recover panic: %v\r\n========\r\n%s", err, string(debug.Stack()))
	}
}
