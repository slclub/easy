package route

import (
	"github.com/slclub/easy/route/element"
	"github.com/slclub/easy/typehandle"
)

type Router interface {
	element.Distributer
	Register(ID element.MID, msg any, handle typehandle.HandleMessage)
	Route(msg any, ag any) error
	PathMap() *element.PathMap
}
