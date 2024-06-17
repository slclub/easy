// Author: slclub@github.com
// Date: 2023-12
// AOI interface definition.
// It is more suitable for mobile games.
// Design:
// 		skip list loop.
//  	cross link.
//      we optimized the cross link using skip list.

package aoi

import (
	"github.com/slclub/easy/nets/agent"
)

// The game player object should implements the Entity interface.
// It is a basic operating unit in aoi.
// Entity interface is mainly compatible with three parts.
// 		Object displacement;
// 		Callback of aoi event handle;
// 		Aoi messages handle;

type Entity interface {
	ID() int
	Position(args ...float32) []float32    //位置
	PositionPre(args ...float32) []float32 //位置
	PositionOld() []float32
	AoiCaller
	AoiMessage() AoiMessage
}

// Entity object with link.
// Your game player should implement AgentEntity interface.
// Neighbour() method return the neighbour object. So, your AgentEntiy object should include a field for neighbour object.
// and use aoi.NewNeighbour(aoi.Option()) to initialized it.

type AgentEntity interface {
	Entity
	Neighbour() Neighbour
	Agent() agent.Agent
}

type Monster interface {
	Entity
	Monster()
}

type Npc interface {
	Entity
	Npc()
}

// The number of people displayed on the same screen in mobile games is limited.
// So, we need neighbour moudule to controll the players of interessting area.

type Neighbour interface {
	neighbourInternel
	RangeMove(fn func(entity Entity) bool)
	RangeBeenObservedSet(fn func(entity Entity) bool)
}

// This module was born solely for the convenience of interal code invoked.
type neighbourInternel interface {
	join(v any) (int, Entity)
	beenJoin(v any) int
	leave(v any) int
	beenLeave(v any) int
	relation(int, Entity) (int, []Entity, []Entity)
	// v:nil将 move和increase 集合的entity 移动到 leave集合中
	// v:clean 将increase 合并到move集合，清空 leave 和 increase 集合
	moveEntitys() []Entity
}

// It is the mainly interface of AOI
// Your scene object should prepare a field for AOI interface.
// It implement the AOIHandle interface.
// The AoiHandle interface is related to your scene logic.

type AOI interface {
	Count(string) int
	Range(func(entity Entity) bool)
	Clear()
	Option() *Option
	AOIHandler
}

type AOIHandler interface {
	Enter(entity Entity)
	Leave(entity Entity)
	Move(entity Entity)
	BroadcastInterstingAll(entity Entity, fn func(mine, other Entity))
	BroadcastAgentAll(fn func(mine, other Entity))
}

// Entity Callback
// Your entity object should implement the AoiCaller interface.
// Please, do not use it as message handler.
// It just your entity should do  when aoi was event happened.
// it is suitable for the changing logic of your entity object when aoi event is happening.

type AoiCaller interface {
	Enter()
	Move()
	Leave()
}

// 消息接收器
// Aoi message handle interface.
// It should be Hang on your entity.
// It is reasonable for different type of entity generate the different messages.

type AoiMessage interface {
	Appear([]Entity)
	Move([]Entity)
	Disappear([]Entity)
}

type NeighbourWeight interface {
	Value(entiry, master Entity) int
}

const (
	MESSAGE_EVENT_EMPTY     = 0
	MESSAGE_EVENT_APPEAR    = 1
	MESSAGE_EVENT_DISAPPEAR = 2
	MESSAGE_EVENT_MOVE      = 3

	// 十字连表 坐标计算状态
	CONST_COORDINATE_EMPTY    = 0
	CONST_COORDINATE_INCREASE = 1
	CONST_COORDINATE_LEAVE    = 2
	CONST_COORDINATE_MOVE     = 3
)
