package servers

import (
	"github.com/slclub/easy/nets/agent"
	"testing"
)

func TestConnSet(t *testing.T) {
	cs := ConnBox{
		conns:  make(ConnSet),
		server: &NewWSServer().Server,
	}
	cs.server.Init(&agent.Gate{})
	a, b, c := 1, 2, 3
	adderr := cs.Add(&a)
	cs.Add(&b)
	cs.Add(&c)
	if len(cs.conns) != 3 {
		t.Fatal("ESAY.SERVER.WS ConnSet.Add", len(cs.conns), adderr)
	}

	cs.Del(&b)
	if len(cs.conns) != 2 {
		t.Fatal("ESAY.SERVER.WS ConnSet.Del", len(cs.conns), adderr)
	}
}

func TestConnSetWithTCP(t *testing.T) {
	cs := ConnBox{
		conns:  make(ConnSet),
		server: &NewTCPServer().Server,
	}
	cs.server.Init(&agent.Gate{})
	a, b, c := 1, 2, 3
	adderr := cs.Add(&a)
	cs.Add(&b)
	cs.Add(&c)
	if len(cs.conns) != 3 {
		t.Fatal("ESAY.SERVER.TCP ConnSet.Add", len(cs.conns), adderr)
	}

	cs.Del(&b)
	if len(cs.conns) != 2 {
		t.Fatal("ESAY.SERVER.TCP ConnSet.Del", len(cs.conns), adderr)
	}
}
