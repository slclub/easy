package element

import (
	"github.com/slclub/easy/vendors/encode"
)

// 为route 绑定插件
type Distributer interface {
	DistributePlug
	// 绑定 解码器=typehandle.Encoder ; 绑定器 = element.Binder
	// Binding(encoder typehandle.Encoder, binder Binder)
	Binding(plugins ...any)
}

type DistributePlug interface {
	Binder() Binder
	Encoder() encode.Encoder
}
