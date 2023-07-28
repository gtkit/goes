// @Author xiaozhaofu 2022/11/25 15:54:00
package goes

import (
	"github.com/gtkit/logger"
	"github.com/olivere/elastic/v7"
)

var _ elastic.Logger = (*esLogger)(nil)

type esLogger struct {
}

func (l esLogger) Printf(format string, v ...interface{}) {
	logger.Infof("[ES] "+format, v...)
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

// NewEsLogger 创建 es 日志实例
func SetEsLogger(logger elastic.Logger) elastic.Logger {
	if logger != nil {
		return logger
	}
	initlogger()
	return &esLogger{}
}
