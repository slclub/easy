package option

import (
	"testing"
	"time"
)

func TestOption(t *testing.T) {
	s := &serv1{}
	opt := &Option{}
	opt.Target(s)
	opt.Config(&conf{
		ID: 12,
		F1: "world",
	})

	opt.Apply()

	if s.Dur != 15*time.Second {
		t.Fatal("ERROR option set config error!", *s)
	}
	t.Log("serv1 current value:", *s, s.nameLower)
}

func TestOptionNew(t *testing.T) {
	s := newServ1(OptionWith(&conf{
		ID: 1,
		F1: "World",
	}))

	if s.ID != 1 {
		t.Fatal("set serv1 with newServ1 error")
	}
	t.Log("serv1 current value:", *s, s.nameLower)
}

func TestAnyStruct(t *testing.T) {
	s := newServ1((&Option{}).Config(&struct {
		ID int
		F1 string
	}{22, "sorry"}))

	if s.F1 != "sorry" {
		t.Fatal("set serv1 with any struct error")
	}
	t.Log("serv1 current value:", *s, s.nameLower)
}

func TestNilOption(t *testing.T) {
	s := newServ1(OptionWith(nil).Default(
		OptionFunc(func() (string, any) {
			return "ID", int64(15)
		}),
	))

	if s.ID != 15 {
		t.Fatal("set serv1 with newServ1 error")
	}
	t.Log("serv1 current value:", *s, s.nameLower)
}

// ------------------------------------------------------------------------------------------------
// program for testing.
// ------------------------------------------------------------------------------------------------

type serv1 struct {
	Name      string
	ID        int
	nameLower string
	F1        string
	Dur       time.Duration
}

func newServ1(assignment Assignment) *serv1 {
	s := &serv1{}
	assignment.Target(s)
	assignment.Apply()
	assignment.Default(OptionFunc(func() (string, any) {
		return "Dur", 15 * time.Second
	}))
	return s
}

type conf struct {
	ID        int
	nameLower string
	F1        string
}

func (self *conf) Dur() time.Duration {
	return 15 * time.Second
}

func (self *conf) NameLower() string {
	return "lower"
}
