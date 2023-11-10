package route

import (
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/route/bind"
	"github.com/slclub/easy/route/encode"
	"reflect"
	"testing"
	"time"
)

func TestRouterRegister(t *testing.T) {
	r := NewRouter()

	// 修改程JSON
	r.Binding(bind.NewBindJson(r.PathMap()), encode.NewJson(r.PathMap()))

	r.Register(MSG_LOGIN_REQ, &JsonLoginReq{}, loginRequest)
	r.Register(MSG_LOGIN_RES, &JsonLoginReS{}, nil)

	e := r.PathMap().GetNewByMID(MSG_LOGIN_REQ)

	if e == nil {
		t.Fatal("EASY.ROUTE.REGISTER register binding error")
		return
	}

	if e.MID != MSG_LOGIN_REQ {
		t.Fatal("EASY.ROUTE.BIND message bind with mid error")
		return
	}
	loginReqType := reflect.TypeOf(&JsonLoginReq{})
	if e.Type != loginReqType {
		t.Fatal("EASY.ROUTE.BIND message type binded error")
	}
}

func TestBinderAndEncodeForJsonType(t *testing.T) {
	r := NewRouter()
	forTestRouterBinderClassName(t, r, "BindProto")
	forTestRouterEncoderClassName(t, r, "Protobuf")
	// 修改程JSON
	r.Binding(bind.NewBindJson(r.PathMap()), encode.NewJson(r.PathMap()))

	forTestRouterBinderClassName(t, r, "BindJson")
	forTestRouterEncoderClassName(t, r, "Json")

	r.Register(MSG_LOGIN_REQ, &JsonLoginReq{}, loginRequest)
	//r.Binder().Register(MSG_LOGIN_REQ, &JsonLoginReq{})
	r.Register(MSG_LOGIN_RES, &JsonLoginReS{}, loginResponse)

	req := &JsonLoginReq{
		MID:  MSG_LOGIN_REQ,
		Name: "Are you",
		Sex:  1,
	}

	data, err := r.Encoder().Marshal(req)
	if err != nil {
		t.Fatal("Encoder.Marshal error")
		return
	}
	res, err := r.Encoder().Unmarshal(data)
	if err != nil {
		t.Fatal("Encoder.Unmarshal error")
	}
	resv, ok := res.(*JsonLoginReq)
	if !ok {
		t.Fatal("Encoder.Testing Unmarshal data con not conver to JsonLoginRes struct")
	}
	if resv.MID != MSG_LOGIN_REQ {
		t.Fatal("Encoder.Testing data lose")
	}

	nilAgent := agent.NewAgent(nil)
	r.Route(req, nilAgent)

	response := &JsonLoginReS{}
	r.Route(response, nilAgent)
	time.Sleep(1 * time.Millisecond)
}

// for testing functions

func forTestRouterBinderClassName(t *testing.T, r Router, name string) {
	typeValue := reflect.TypeOf(r.Binder())
	if name != typeValue.Elem().Name() {
		t.Fatal("Router Binder verify error", name)
	}
}

func forTestRouterEncoderClassName(t *testing.T, r Router, name string) {
	typeValue := reflect.TypeOf(r.Encoder())
	if name != typeValue.Elem().Name() {
		t.Fatal("Router Encoder verify error", name)
	}
}
