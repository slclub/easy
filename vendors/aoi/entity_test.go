package aoi

import (
	"github.com/slclub/easy/nets/agent"
	"github.com/slclub/easy/vendors/option"
)

var global_id = 0

type user struct {
	id        int
	neighbour Neighbour

	// transform
	position_pre []float32
	position     []float32
	position_old []float32

	message AoiMessage
	Object
}

func newUserWithAoi(aoi AOI) *user {
	u1 := newUser(option.OptionWith(nil).Default(
		option.OptionFunc(func() (string, any) {
			return "Neighbour", NewNeighbour(aoi.Option())
		}),
		option.OptionFunc(func() (string, any) {
			return "Position", []float32{0, 0, 0}
		}),
		option.OptionFunc(func() (string, any) {
			return "Position_pre", []float32{0, 0, 0}
		}),
		option.OptionFunc(func() (string, any) {
			return "Position_old", []float32{0, 0, 0}
		}),
	))

	message := &messageUser{master: u1}
	u1.message = message
	u1.neighbour.BindWith(option.OptionWith(nil).Default(
		option.OptionFunc(func() (string, any) {
			return "Master", u1
		}),
	))

	return u1
}

func newUser(assignment option.Assignment) *user {
	global_id++
	u := &user{
		id: global_id,
	}
	if assignment == nil {
		return u
	}
	assignment.Target(u)
	assignment.Apply()
	return u
}

// ---------------------implement aoi.Eneity---------------------

func (self *user) ID() int {
	return self.id
}

func (self *user) PositionPre(args ...float32) []float32 {
	if len(args) > 0 {
		self.position_pre = args
	}
	return self.position_pre
}

func (self *user) Position(args ...float32) []float32 {
	if len(args) > 0 {
		self.position_old = self.position
		self.position = args
	}
	return self.position
}

func (self *user) PositionOld() []float32 {
	return self.position_old
}

func (self *user) AoiMessage() AoiMessage {
	return self.message
}

// ---------------------implement aoi.Eneity---------------------

// ---------------------implement aoi.AoiCaller---------------------
func (self *user) Move() {
	// 这一步很关键
	self.Position(self.PositionPre()...)
}

// ---------------------implement aoi.AoiCaller---------------------

// ---------------------implement aoi.AgentEntity---------------------
// except Entity interface
func (self *user) Agent() agent.Agent {
	return nil
}
func (self *user) Neighbour() Neighbour {
	return self.neighbour
}

// ---------------------implement aoi.AgentEntity---------------------

var _ AgentEntity = &user{}
