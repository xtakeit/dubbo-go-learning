package log

import (
	"errors"
	"fmt"
	"go-server/library/conf"
)

// LoggerContainer 实现了热更新的日志工具容器
type LoggerContainer struct {
	*conf.Container
}

// GetLoggerConfFunc 定义了获取日志工具配置的函数类型
type GetLoggerConfFunc func() (*LoggerConf, error)

// ErrGetLoggerConfFuncIsNil 初始化容器时配置函数未指定时返回
var ErrGetLoggerConfFuncIsNil = errors.New(
	"get logger conf func is nil",
)

// NewLoggerContainer 初始化日志容器
func NewLoggerContainer(getLoggerConf GetLoggerConfFunc) (ct *LoggerContainer, err error) {
	if getLoggerConf == nil {
		err = ErrGetLoggerConfFuncIsNil
		return
	}

	getObjConf := func() (icf conf.IConf, err error) {
		icf, err = getLoggerConf()
		if err != nil {
			err = fmt.Errorf("get log conf: %w", err)
			return
		}
		return
	}

	ict, err := conf.NewContainer(
		getObjConf, compareLoggerConf, newLoggerObj,
		resetLoggerObj,
	)
	if err != nil {
		err = fmt.Errorf("new conf container: %w", err)
		return
	}

	ct = &LoggerContainer{
		Container: ict,
	}

	return
}

// newLoggerObj 初始日志工具函数
func newLoggerObj(icf conf.IConf) (iobj conf.IObject, err error) {
	cf, ok := icf.(*LoggerConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	iobj, err = NewLogger(cf)
	if err != nil {
		err = fmt.Errorf("new log: %w", err)
		return
	}

	return
}

// compareLoggerConf 日志工具配置比较函数
func compareLoggerConf(iocf, incf conf.IConf) (rst conf.CompareObjConfRst, err error) {
	ocf, ncf := iocf.(*LoggerConf), incf.(*LoggerConf)
	if *ocf != *ncf {
		rst = conf.CompareObjConfRstNeedReset
		return
	}
	rst = conf.CompareObjConfRstNoNeed
	return
}

// resetLoggerObj 重设日志工具函数
func resetLoggerObj(iobj conf.IObject, iocf, incf conf.IConf) (err error) {
	ocf, ncf, logger := iocf.(*LoggerConf), incf.(*LoggerConf), iobj.(*Logger)
	if ncf.Level != ocf.Level {
		if err = logger.SetLevel(ncf); err != nil {
			err = fmt.Errorf("set level: %w", err)
			return
		}
		ocf.Level = ncf.Level
	}
	if ncf.Output != ocf.Output {
		if err = logger.SetOutput(ncf); err != nil {
			err = fmt.Errorf("set output: %w", err)
			return
		}
		ocf.Output = ncf.Output
	}
	return
}

// MustGetLogger 获取容器包装的日志工具, 类型断言失败将导致panic
func (ct *LoggerContainer) MustGetLogger() (logger *Logger) {
	logger = ct.MustGetObj().(*Logger)
	return
}

// PutLogger 回收日志工具
func (ct *LoggerContainer) PutLogger(logger *Logger) {
	ct.PutObj(logger)
}
