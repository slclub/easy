package client

import (
	"errors"
	"github.com/slclub/easy/route"
	"github.com/slclub/easy/route/bind"
	"github.com/slclub/easy/route/encode"
	"github.com/slclub/easy/typehandle"
	"net/url"
	"strings"
)

func RouterWithProtocol(r route.Router, protocol string) {
	switch protocol {
	case typehandle.ENCRIPT_DATA_JSON:
		r.Binding(bind.NewBindJson(r.PathMap()), encode.NewJson(r.PathMap()))
	default:
		r.Binding(bind.NewBindProto(r.PathMap()), encode.NewProtobuf(r.PathMap()))
	}
}

func GinWebSocketSchceme(addr string) (string, error) {
	if len(addr) <= 2 {
		return "", errors.New(("web socket addr is too short"))
	}

	if addr[:1] == ":" {
		addr = "127.0.0.1" + addr
	}
	strings.TrimLeft(addr, "http")
	u := url.URL{Scheme: "ws", Host: addr}

	return u.String(), nil
}
