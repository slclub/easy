package log

import (
	"github.com/slclub/go-tips/logf"
	"github.com/slclub/log8q"
)

const (
	EL       = "[ESAY] "
	EL_ERROR = "[ESAY ERROR]"
	EL_FATAL = "[EASY FATAL]"
)

var (
	_log logf.Logger

	/** LEVEL 取值
	0 不打印输出任何日志
	LEVEL_INFO     Level = 1
	LEVEL_DEBUG    Level = 2
	LEVEL_WARNNING Level = 4
	LEVEL_ERROR    Level = 8
	LEVEL_FATAL    Level = 16

	TRACE_INFO     Level = 32
	TRACE_DEBUG    Level = 64
	TRACE_WARNNING Level = 128
	TRACE_ERROR    Level = 256
	TRACE_FATAL    Level = 512
	*/
	LEVEL log8q.Level
)

func init() {
	_log = Log(logf.Log())
	LEVEL = log8q.LEVEL_INFO
}

func Init() {
}

func Info(format string, a ...any) {
	Release(format, a...)
}

func Debug(format string, a ...any) {
	//gLogger.doPrintf(debugLevel, printDebugLevel, format, a...)
	if !LEVEL.Check(log8q.LEVEL_DEBUG) {
		return
	}
	_log.Printf(EL+format, a...)
}

// 与info 同级别
func Release(format string, a ...any) {
	if !LEVEL.Check(log8q.LEVEL_INFO) {
		return
	}
	//gLogger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
	_log.Printf(EL+format, a...)
}

func Error(format string, a ...any) {
	if !LEVEL.Check(log8q.LEVEL_ERROR) {
		return
	}
	//gLogger.doPrintf(errorLevel, printErrorLevel, format, a...)
	_log.Printf(EL_ERROR+format, a...)
}

func Fatal(format string, a ...any) {
	//if !LEVEL.Check(log8q.LEVEL_FATAL) {
	//	return
	//}
	//gLogger.doPrintf(fatalLevel, printFatalLevel, format, a...)
	_log.Printf(EL_FATAL+format, a...)
}

func Log(ls ...logf.Logger) logf.Logger {
	if ls == nil || len(ls) == 0 {
		return _log
	}
	_log = ls[0]
	return _log
}
