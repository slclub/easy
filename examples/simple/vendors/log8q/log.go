package log8q

import (
	"context"
	"github.com/slclub/go-tips/logf"
	"github.com/slclub/log8q"
)

var _log logf.Logger

func Log() logf.Logger {
	if _log == nil {
		_log = log8q.New(context.Background(), &log8q.Config{
			Filename: "log/log8q.log",
		})
	}
	//if _log == nil {
	//	_log = logf.New()
	//}
	return _log
}

func Info(args ...any) {
	if l8, ok := Log().(*log8q.Log8); ok {
		l8.Info(args...)
	}
}

func Debug(args ...any) {
	if l8, ok := Log().(*log8q.Log8); ok {
		l8.Debug(args...)
	}
}
