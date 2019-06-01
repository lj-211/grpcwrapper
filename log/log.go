package log

import (
	"fmt"
	"github.com/lj-211/common/cmodel"
	"github.com/lj-211/common/ecode"
)

var logger cmodel.Logger = nil

func Info(args ...interface{}) {
	if logger != nil {
		logger.Info(args...)
	}
	fmt.Println(ecode.String("1"))
}

func Infof(format string, args ...interface{}) {
	if logger != nil {
		logger.Infof(format, args...)
	}
}

func Debug(args ...interface{}) {
	if logger != nil {
		logger.Debug(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if logger != nil {
		logger.Debugf(format, args...)
	}
}

func Warning(args ...interface{}) {
	if logger != nil {
		logger.Warning(args...)
	}
}

func Warningf(format string, args ...interface{}) {
	if logger != nil {
		logger.Warningf(format, args...)
	}
}

func Error(args ...interface{}) {
	if logger != nil {
		logger.Error(args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if logger != nil {
		logger.Errorf(format, args...)
	}
}

func SetLogger(log cmodel.Logger) {
	logger = log
}
