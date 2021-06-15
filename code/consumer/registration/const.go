package registration

import (
	"github.com/0chain/gosdk/core/logger"
)

// intLevelFromStr converts string log level to gosdk logger level int value.
func intLevelFromStr(level string) int {
	switch level {
	case "none":
		return logger.NONE
	case "fatal":
		return logger.FATAL
	case "error":
		return logger.ERROR
	case "info":
		return logger.INFO
	case "debug":
		return logger.DEBUG

	default:
		return -1
	}
}
