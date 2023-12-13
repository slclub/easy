# AOI

## Overview

It is service for mobile games. Application scenarios where multiple people on the
same screen are synchronizated in the big world and limited number of people. It will
be usefull for it.

AOI components should be embedded in your scene or room objects or used as their fields.

## Install

```go get github.com/slclub/easy/vendors/aoi```

## Design

- Orthogonal List
- Skip list

Efficient search through cross links and skip list.

## AoiTest

I belive that the testing code will help you and me. The testing code also is a guide.

### New Aoi Object

```aoiObject := New(option.OptionWith(nil))```

You will get an aoi object with default configuration.

### Create testing users

Here we create five testing users and set their positions (x, y, z).

- user.PositionPre(): Prepare user coordinate paramters.
- user.Move(): This method  will be called by aoi object. Here we are only to set user coordinates.

```go
func tInitAoiAndUsers() (AOI, []*user) {
	aoiObject := New(option.OptionWith(nil))

	users := []*user{
		newUserWithAoi(aoiObject),
		newUserWithAoi(aoiObject),
		newUserWithAoi(aoiObject),
		newUserWithAoi(aoiObject),
		newUserWithAoi(aoiObject),
	}
	posotions := [][]float32{
		[]float32{1, 0, 1},
		[]float32{2, 0, 1.2},
		[]float32{5, 0, 7.2},
		[]float32{100, 0, 6.2},
		[]float32{14, 0, 10.2},
	}
	for i, _ := range users {
		log.Info("AOI.Enter master.ID:%v", users[i].ID())
		users[i].PositionPre(posotions[i]...)
		users[i].Move()
	}
	return aoiObject, users
}
```

### Aoi Enter Test

I am used the fmt package to print out log infomation and verify the results.

```go
func TestJoinAoi1(t *testing.T) {
	aoiObject, users := tInitAoiAndUsers()

	for i, _ := range users {
		aoiObject.Enter(users[i])
	}

	aoiObject.Clear()
}
```

### Aoi Leave Test

```go
func TestLeave1(t *testing.T) {
	aoiObject, users := tInitAoiAndUsers()

	for i, _ := range users {
		aoiObject.Enter(users[i])
	}

	//time.Sleep(time.Millisecond * 5)
	for i, _ := range users {
		if i > 3 {
			break
		}
		aoiObject.Leave(users[i])
	}

	aoiObject.Clear()
}
```

### Aoi Move Test

```go
func TestMove1(t *testing.T) {
	log.LEVEL = log8q.ALL_LEVEL
	aoiObject, users := tInitAoiAndUsers()

	for i, _ := range users {
		aoiObject.Enter(users[i])
	}
	time.Sleep(time.Millisecond * 5)
	positions := [][]float32{
		[]float32{18, 0, 20},
		[]float32{19, 0, 100},
		[]float32{19.5, 0, 20.2},
	}

	// 预设 移动坐标
	users[0].PositionPre(positions[0]...)
	users[1].PositionPre(positions[1]...)
	users[3].PositionPre(positions[2]...)

	// 具体移动交给 entity.Move 设置到 Position才算真正的移动
	aoiObject.Move(users[0])
	time.Sleep(time.Millisecond * 5)
	log.Debug(" --------- user2.Move Master.ID:%v", users[1].ID())
	aoiObject.Move(users[1])
	aoiObject.Move(users[3])

	aoiObject.Clear()
}
```


## User Demo

It is just an example. You do not reuse it in an actual project. It implements 
the ```aoi.AgentEntity``` interface.

- FILE ```vendors/aoi/entity_test.go```

### 1 New An Test User

Create a  testing user.

```go

func newUserWithAoi(aoi AOI) *user {
	u1 := newUser(option.OptionWith(nil).Default(
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
    
	// binding the AoiMessage
	message := &messageUser{master: u1}
	u1.message = message
	
	// binding the Neighbour for aoi
	u1.neighbour = NewNeighbour(
		NeighbourWithOption(aoi.Option()),
		NeighbourWithMaster(u1),
	)

	return u1
}
```

```go
// Not recommended to use option package here.
// It has a lower performence.
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
```

> It is a wrong decision to using option package here. But it does not affect the code running and tested results. 


### 2 define
```go

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
```

## MessageUser Demo

It is an example. It impelements the ```aoi.AoiMessage``` interface.

- FILE ```vendors/aoi/message_test.go```

```go

import "github.com/slclub/easy/log"

type messageUser struct {
	master AgentEntity
	Message
}

func (self *messageUser) Appear(entitys []Entity) {
	aes := []AgentEntity{}
	monsters := []Monster{}
	for _, v := range entitys {
		switch val := v.(type) {
		case AgentEntity:
			aes = append(aes, val)
		case Monster:
			monsters = append(monsters, val)
		case Npc:
		}
	}

	appear := &messageAppear{
		ID:      self.master.ID(),
		Players: aes,
	}
	log.Info("APPEAR user.len ID:%v players num :%v", appear.ID, len(appear.Players))
	for _, v := range appear.Players {
		log.Info("=== APPEAR user.player MASTERID:%v ID:=%v position:%v:%v:%v", appear.ID, v.ID(), v.Position()[0], v.Position()[1], v.Position()[2])
	}
}

func (self *messageUser) Disappear(entitys []Entity) {
	aes := []AgentEntity{}
	monsters := []Monster{}
	for _, v := range entitys {
		switch val := v.(type) {
		case AgentEntity:
			aes = append(aes, val)
		case Monster:
			monsters = append(monsters, val)
		case Npc:
		}
	}

	msg := &messageDisappear{
		ID:      self.master.ID(),
		Players: aes,
	}
	log.Info("DISAPPEAR user.len ID:%v players num :%v", msg.ID, len(msg.Players))
	for _, v := range msg.Players {
		log.Info("--- DISAPPEAR user MASTERID:%v ID:=%v position:%v:%v:%v", msg.ID, v.ID(), v.Position()[0], v.Position()[1], v.Position()[2])
	}
}

func (self *messageUser) Move(entitys []Entity) {
	msg := &messageMove{
		ID:      self.master.ID(),
		Objects: entitys,
	}
	log.Info("MOVE user.len ID:%v players num :%v", msg.ID, len(msg.Objects))
	for _, v := range msg.Objects {
		log.Info("--- MOVE user MASTERID:%v ID:=%v position:%v:%v:%v", msg.ID, v.ID(), v.Position()[0], v.Position()[1], v.Position()[2])
	}
}

// ------------------------------------------------------------------------
// message defined
// ------------------------------------------------------------------------

type messageAppear struct {
	ID      int
	Players []AgentEntity
}

type messageDisappear struct {
	ID      int
	Players []AgentEntity
}

type messageMove struct {
	ID      int
	Objects []Entity
}

type messageAppearMonster struct {
	ID       int
	Monsters []Monster
}

```
