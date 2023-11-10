package log8q

import (
	"context"
	"github.com/slclub/log8q"
)

var l8 *log8q.Log8

func Log() *log8q.Log8 {
	if l8 == nil {
		l8 = log8q.New(context.Background(), nil)
	}
	return l8
}
