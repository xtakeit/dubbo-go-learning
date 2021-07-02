package component

import (
	"fmt"

	"go-server/library/clean"
	"go-server/library/log"
)

// InfLogger INFO级别日志容器
var InfLogger *log.LoggerContainer

type InfLoggerConfig struct {
	Level  string `env:"INF_LOG_LEVEL"`
	Output string `env:"INF_LOG_OUTPUT"`
}

// SetupInfLogger 配置INFO级别日志
func SetupInfLogger() (err error) {
	InfLogger, err = log.NewLoggerContainer(getInfLoggerConf)
	if err != nil {
		err = fmt.Errorf("log.NewLoggerContainer: %w", err)
		return
	}

	clean.Push(InfLogger)
	Conf.PushUpdater(InfLogger)

	return
}

func getInfLoggerConf() (cf *log.LoggerConf, err error) {
	cfg := &InfLoggerConfig{}

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
