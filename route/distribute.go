package route

import (
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/route/element"
	"github.com/slclub/easy/vendors/encode"
	"github.com/slclub/go-tips"
)

// 路由器 组件 集合器
type DistributeCollection struct {
	binder element.Binder
	coder  encode.Encoder
}

func (self *DistributeCollection) Binder() element.Binder {
	return self.binder
}

func (self *DistributeCollection) Encoder() encode.Encoder {
	return self.coder
}

// func (self *DistributeCollection) Binding(coder typehandle.Encoder, binder element.Binder) {
func (self *DistributeCollection) Binding(plugins ...any) {

	if plugins == nil || len(plugins) == 0 {
		return
	}

	for _, plug := range plugins {
		if tips.IsNil(plug) {
			log.Fatal("EASY.ROUTER.Binding empty plugin")
			continue
		}

		switch val := plug.(type) {
		case encode.Encoder:
			self.coder = val
		case element.Binder:
			self.binder = val
		}
	}
}
