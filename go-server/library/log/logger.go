package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	timeLayout = "2006-01-02 15:04:05"
)

// F 日志字段快捷类型
type F logrus.Fields

// LoggerConf 日志配置
type LoggerConf struct {
	Level  string // 级别
	Output string // 输出
}

// Logger 日志工具
type Logger struct {
	*logrus.Logger

	outputCloser io.WriteCloser
}

// NewLogger 返回日志工具实例
func NewLogger(cf *LoggerConf) (logger *Logger, err error) {
	logger = &Logger{
		Logger: logrus.New(),
	}

	// 设置日志输出
	if err = logger.SetOutput(cf); err != nil {
		err = fmt.Errorf("set output: %w", err)
		return
	}

	// 设置日志格式
	logger.SetFormatter()

	// 设置日志级别
	if err = logger.SetLevel(cf); err != nil {
		err = fmt.Errorf("set level: %w", err)
		return
	}

	return
}

// SetOutput 设置日志输出
func (logger *Logger) SetOutput(cf *LoggerConf) (err error) {
	var output io.WriteCloser
	switch strings.ToLower(cf.Output) {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		output, err = os.OpenFile(cf.Output, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			err = fmt.Errorf("open log file: %w", err)
			return
		}
	}
	logger.Logger.SetOutput(output)
	if logger.outputCloser != nil && logger.outputCloser != os.Stdout && logger.outputCloser != os.Stderr {
		_ = logger.outputCloser.Close()
	}
	logger.outputCloser = output
	return
}

// SetFormatter 设置格式
func (logger *Logger) SetFormatter() {
	logger.Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: timeLayout,
	})
}

// SetLevel 设置级别
func (logger *Logger) SetLevel(cf *LoggerConf) (err error) {
	level, err := logrus.ParseLevel(cf.Level)
	if err != nil {
		err = fmt.Errorf("parse level: %w", err)
		return
	}
	logger.Logger.SetLevel(level)
	return
}

// Close 如果日志是写到文件的, 关闭该文件
func (logger *Logger) Close() (err error) {
	if logger.outputCloser != nil && logger.outputCloser != os.Stdout && logger.outputCloser != os.Stderr {
		if err = logger.outputCloser.Close(); err != nil {
			err = fmt.Errorf("close output: %w", err)
			return
		}
	}
	return
}
