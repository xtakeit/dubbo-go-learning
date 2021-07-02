package component

import (
	"fmt"
	"os"

	"go-server/library/clean"
	"go-server/library/conf"
)

// Conf 全局配置对象
var Conf *conf.Conf

// SetupConf 初始化配置对象
func SetupConf(filename string) (err error) {
	Conf, err = conf.NewConf(filename)
	if err != nil {
		err = fmt.Errorf("conf.NewConf <%s>: %w", filename, err)
		return
	}

	if err = Conf.Load(); err != nil {
		err = fmt.Errorf("Conf.Load: %w", err)
		return
	}

	Conf.SetWatchErrHandleFunc(func(err error) {
		_, _ = fmt.Fprintf(os.Stderr, "Conf.Watch: %v\n", err)
	})

	if err = Conf.Watch(); err != nil {
		err = fmt.Errorf("Conf.Watch: %w", err)
		return
	}

	clean.Push(Conf)

	return
}
