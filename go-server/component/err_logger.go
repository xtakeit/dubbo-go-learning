package component

import (
	"fmt"

	"go-server/library/clean"
	"go-server/library/log"
)

// ErrLogger ERR级别日志容器
var ErrLogger *log.LoggerContainer

type ErrLoggerConfig struct {
	Level  string `env:"ERR_LOG_LEVEL"`
	Output string `env:"ERR_LOG_OUTPUT"`
}

// SetupErrLogger 配置INFO级别日志
func SetupErrLogger() (err error) {
	ErrLogger, err = log.NewLoggerContainer(getErrLoggerConf)
	if err != nil {
		err = fmt.Errorf("log.NewLoggerContainer: %w", err)
		return
	}

	clean.Push(ErrLogger)
	Conf.PushUpdater(ErrLogger)

	return
}

func getErrLoggerConf() (cf *log.LoggerConf, err error) {
	cfg := &ErrLoggerConfig{}

	if err = Conf.Scan(cfg, "env"); err != nil {
		err = fmt.Errorf("Conf.Scan: %w", err)
		return
	}

	cf = &log.LoggerConf{
		Level:  cfg.Level,
		Output: cfg.Output,
	}

	return
}
