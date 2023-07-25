// @Author xiaozhaofu 2022/11/25 15:54:00
package goes

import (
	"github.com/gtkit/logger"
	"go.uber.org/zap"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type EsLogger struct {
	ZapLogger *zap.Logger
}

func (l EsLogger) Printf(format string, v ...interface{}) {
	l.ZapLogger.Sugar().Infof("[ES] "+format, v...)
}

func newLogger() EsLogger {
	initlogger()
	return EsLogger{ZapLogger: logger.Zlog()}
}

func initlogger() {
	if logger.Zlog() == nil {
		opt := &logger.Option{
			FileStdout: true,
			Division:   "size",
		}
		logger.NewZap(opt)
	}
}
