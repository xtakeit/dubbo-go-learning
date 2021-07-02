// Author: Steve Zhang
// Date: 2020/9/16 4:55 下午

package log

import (
	"github.com/sirupsen/logrus"
)

func (ct *LoggerContainer) Info(fields F, args ...interface{}) {
	logger := ct.MustGetLogger()
	defer ct.PutLogger(logger)
	logger.WithFields(logrus.Fields(fields)).Info(args...)
}

func (ct *LoggerContainer) Warn(fields F, args ...interface{}) {
	logger := ct.MustGetLogger()
	defer ct.PutLogger(logger)
	logger.WithFields(logrus.Fields(fields)).Warn(args...)
}

func (ct *LoggerContainer) Error(fields F, args ...interface{}) {
	logger := ct.MustGetLogger()
	defer ct.PutLogger(logger)
	logger.WithFields(logrus.Fields(fields)).Error(args...)
}

func (ct *LoggerContainer) Debug(fields F, args ...interface{}) {
	logger := ct.MustGetLogger()
	defer ct.PutLogger(logger)
	logger.WithFields(logrus.Fields(fields)).Debug(args...)
}

func (ct *LoggerContainer) Fatal(fields F, args ...interface{}) {
	logger := ct.MustGetLogger()
	defer ct.PutLogger(logger)
	logger.WithFields(logrus.Fields(fields)).Fatal(args...)
}

func (ct *LoggerContainer) Panic(fields F, args ...interface{}) {
	logger := ct.MustGetLogger()
	defer ct.PutLogger(logger)
	logger.WithFields(logrus.Fields(fields)).Panic(args...)
}
