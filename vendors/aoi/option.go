package aoi

import (
	"github.com/slclub/go-tips/logf"
)

const (
	DEFAULT_RADIUS          = 30 // 20 默认值； 300 测试 大半径
	DEFAULT_NEIGHBOUR_COUNT = 10
)

type Option struct {
	Radius          float32
	NeighbourCount  int
	Axis            []int // 坐标系选择 0, 1, 2 := x,y,z
	NeighbourWeight NeighbourWeight
	Log             logf.Logger
}

func DefaultOption() *Option {
	option := &Option{
		Radius:         DEFAULT_RADIUS,
		NeighbourCount: DEFAULT_NEIGHBOUR_COUNT,
		Log:            &logEmpty{},
	}

	return option
}
