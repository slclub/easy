/*
 * Aoi module.
 * you should hang it on your scene object
 */
package aoi

import (
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/vendors/option"
)

const (
	COUNT_ALL     = "all"
	COUNT_AGENT   = "agent"
	COUNT_NPC     = "npc"
	COUNT_MONSTER = "monster"
)

type AoiArea struct {
	option   *Option
	chanMove chan *OperEvent
	cross    *crossList

	countAgent   int // 带链接的对象，一般指人的数量
	countNpc     int // NPC 数量
	countMonster int // 怪物数量
}

func New(assignment option.Assignment) AOI {
	a := &AoiArea{
		option:       DefaultOption(),
		chanMove:     make(chan *OperEvent, 20), // 接收的事件队列
		countAgent:   0,
		countMonster: 0,
		countNpc:     0,
	}

	assignment.Target(a.option)
	assignment.Default(
		option.OptionFunc(func() (string, any) {
			return "Axis", []int{0, 1, 2}
		}),
	)
	assignment.Apply()
	// option 0值不包含指针
	if a.option.Log == nil {
		a.option.Log = &logEmpty{}
	}

	a.cross = newCrossList(option.OptionWith(a.option).Default(
		option.OptionFunc(func() (string, any) {
			return "Axis", a.option.Axis
		}),
		option.OptionFunc(func() (string, any) {
			return "Radius", DEFAULT_RADIUS
		}),
		option.OptionFunc(func() (string, any) {
			return "Log", a.option.Log
		}),
		handleNewListWithAxis(a.option.Axis),
		option.DEFAULT_IGNORE_ZERO,
	))

	a.init()
	return a
}

func (self *AoiArea) init() {
	go self.tickUpdate()
}

func (self *AoiArea) tickUpdate() {
	defer self.clear()
	defer close(self.chanMove)
	for {
		select {
		case op := <-self.chanMove:
			switch op.Op {
			case OPER_AGENT_JOIN:
				entity, _ := op.Entity.(AgentEntity)
				self.enter(entity)
			case OPER_AGENT_LEAVE:
				entity, _ := op.Entity.(AgentEntity)
				self.leave(entity)
			case OPER_AGENT_MOVE:
				entity, _ := op.Entity.(AgentEntity)

				self.move(entity)
			case OPER_NOAGENT_JOIN:

			case OPER_ACTION_AOI:
			// 例如广播
			//self.actionWithOper(op)
			case OPER_ACTION_INTERSTING_ALL:
				self.broadcastInterstingAll(op)
			case OPER_ACTION_AGENT_ALL:
				self.broadcastAgentAll(op)
			case OPER_QUIT:
				op = nil
				log.Info("[AOI][TickUpdate][QUIT]")
				//self.Clear()
				return
			}
			op = nil
		}
	}
}

func (self *AoiArea) Count(count string) int {
	switch count {
	case COUNT_AGENT:
		return self.countAgent
	case COUNT_NPC:
		return self.countNpc
	case COUNT_MONSTER:
		return self.countMonster
	case COUNT_ALL:
		return self.countAgent + self.countNpc + self.countMonster
	}
	return self.countAgent + self.countNpc + self.countMonster
}

func (self *AoiArea) Range(fn func(entity Entity) bool) {
	self.cross.Range(fn)
}

func (self *AoiArea) Clear() {
	self.chanMove <- &OperEvent{Op: OPER_QUIT}
}

func (self *AoiArea) Option() *Option {
	return self.option
}

func (self *AoiArea) clear() {
}

func (self *AoiArea) enter(entity Entity) {
	self.cross.Delete(entity)
	self.cross.Add(entity)

	arr := []Entity{}
	// 查询半径范围内Entity
	switch me := entity.(type) {
	case AgentEntity:
		self.cross.RangeByRadius(entity, func(other Entity) {
			if self.Option().NeighbourCount == 0 {
				arr = append(arr, other)
			} else {
				// 处理玩家邻居
				ecode := me.Neighbour().join(other)
				if ecode == MESSAGE_EVENT_APPEAR {
					arr = append(arr, other)
				}
			}

			switch him := other.(type) {
			case AgentEntity:
				if self.Option().NeighbourCount == 0 {
					him.AoiMessage().Appear([]Entity{me})
				} else {
					eecode := him.Neighbour().join(me)
					self.handleMessageEvent(eecode, him, me)
				}
			case Monster:
				him.AoiMessage().Appear([]Entity{me})
			case Npc:
				him.AoiMessage().Appear([]Entity{me})
			}
		})
		//log.Info("--------- appear after cross.RangeByRadius")
		// 处理怪物和NPC
		me.AoiMessage().Appear(arr)
		self.countAgent++
	case Monster:
		all := []Entity{}
		self.cross.RangeByRadius(entity, func(other Entity) {
			all = append(all, other)
			other.AoiMessage().Appear([]Entity{me})
		})
		me.AoiMessage().Appear(all)
		self.countMonster++
	case Npc:
		all := []Entity{}
		self.cross.RangeByRadius(entity, func(other Entity) {
			all = append(all, other)
			other.AoiMessage().Appear([]Entity{me})
		})
		me.AoiMessage().Appear(all)
	}

	entity.Enter()
}

func (self *AoiArea) leave(entity Entity) {

	defer self.cross.Delete(entity)
	arr := []Entity{}
	switch me := entity.(type) {
	case AgentEntity:
		self.cross.RangeByRadius(entity, func(other Entity) {
			if self.Option().NeighbourCount == 0 {
				arr = append(arr, other)
			} else {
				mcode := me.Neighbour().leave(other)
				if mcode == MESSAGE_EVENT_DISAPPEAR {
					arr = append(arr, other)
				}
			}

			switch him := other.(type) {
			case AgentEntity:
				if self.option.NeighbourCount == 0 {
					him.AoiMessage().Disappear([]Entity{me})
				} else {
					ecode := him.Neighbour().leave(me)
					self.handleMessageEvent(ecode, him, me)
				}
			case Monster:
				arr = append(arr, him)
				him.AoiMessage().Disappear([]Entity{me})
			case Npc:
				him.AoiMessage().Disappear([]Entity{me})
			}
		})
		me.AoiMessage().Disappear(arr)
		self.countAgent--
		// 处理entity 消息
		//this.msgEvent.Dest(at).From(nil).AoiDisappearWithEntity(at.GetNeighbour().GetLeaveEntity(true))

	case Npc, Monster:
		self.cross.RangeByRadius(entity, func(other Entity) {
			arr = append(arr, other)
			other.AoiMessage().Disappear([]Entity{me})
		})
		me.AoiMessage().Disappear(arr)
		self.countMonster--
	}
	entity.Leave()
}

func (self *AoiArea) move(entity Entity) {
	entity.Move()
	// 从旧节点中删除
	self.cross.DeleteCache(entity)
	// 从新坐标中，再去
	self.cross.Add(entity)

	switch me := entity.(type) {
	case AgentEntity:
		increases_agents, increases := []Entity{}, []Entity{}
		decrease_agents, decrease := []Entity{}, []Entity{}
		//self.cross.RangeByRadiusAll(me, func(other Entity, nearcheck int) {
		self.cross.RangeByRadiusAll(me, func(other Entity, nearcheck int) {
			switch him := other.(type) {
			case AgentEntity:
				if self.Option().NeighbourCount == 0 {
					self.handleWithCoordinate(nearcheck, him, me, func() {
						switch nearcheck {
						case CONST_COORDINATE_INCREASE:
							increases_agents = append(increases_agents, him)
						case CONST_COORDINATE_LEAVE:
							decrease_agents = append(decrease_agents, him)
						}
					})
				} else {
					me.Neighbour().relation(nearcheck, him, func(rcode int) {
						switch rcode {
						case MESSAGE_EVENT_APPEAR:
							increases_agents = append(increases_agents, him)
						case MESSAGE_EVENT_DISAPPEAR:
							decrease_agents = append(decrease_agents, him)
						}
					})
					him.Neighbour().relation(nearcheck, me, func(rcode int) {
						self.handleMessageEvent(rcode, him, me)
					})
				}

				self.Option().Log.Printf("ME.PID:%v ME.Pos:%v him.ID:%v him.Pos:%v near:%v", me.ID(), me.Position(), him.ID(), him.Position(), nearcheck)
			case Monster:
				self.handleWithCoordinate(nearcheck, him, me, func() {
					switch nearcheck {
					case CONST_COORDINATE_INCREASE:
						increases = append(increases, him)
					case CONST_COORDINATE_LEAVE:
						decrease = append(decrease, him)
					}
				})

			case Npc:
				self.handleWithCoordinate(nearcheck, him, me, func() {
					switch nearcheck {
					case CONST_COORDINATE_INCREASE:
						increases = append(increases, him)
					case CONST_COORDINATE_LEAVE:
						decrease = append(decrease, him)
					}
				})
			}
		})

		// 出视野
		me.AoiMessage().Disappear(decrease_agents)
		//me.AoiMessage().Disappear(decrease)
		// 入视野
		me.AoiMessage().Appear(increases_agents)
		//me.AoiMessage().Appear(increases)
	}
}

func (self *AoiArea) handleWithCoordinate(nearcheck int, target, from Entity, fn func()) {
	switch nearcheck {
	case CONST_COORDINATE_MOVE:
		target.AoiMessage().Move([]Entity{from})
	case CONST_COORDINATE_INCREASE:
		target.AoiMessage().Appear([]Entity{from})
	case CONST_COORDINATE_LEAVE:
		target.AoiMessage().Disappear([]Entity{from})
	}
	fn()
}

func (self *AoiArea) handleMessageEvent(ecode int, target, from AgentEntity) {
	switch ecode {
	case MESSAGE_EVENT_APPEAR:
		target.AoiMessage().Appear([]Entity{from})
	case MESSAGE_EVENT_DISAPPEAR:
		target.AoiMessage().Disappear([]Entity{from})
	case MESSAGE_EVENT_MOVE:
		target.AoiMessage().Move([]Entity{from})
	}
}

// 依据半径广播 半径范围内的所有 人
func (self *AoiArea) broadcastInterstingAll(oper *OperEvent) {

	self.cross.RangeByRadius(oper.Entity, func(other Entity) {
		if oper.Entity.ID() == other.ID() {
			return
		}

		//if _, ok := other.(AgentEntity); !ok {
		//	return
		//}
		if oper.Handle != nil {
			oper.Handle(oper.Entity, other)
		}
	})
	oper = nil
}

func (self *AoiArea) broadcastAgentAll(oper *OperEvent) {

	self.cross.Range(func(one Entity) bool {
		if _, ok := one.(AgentEntity); !ok {
			return true
		}
		if oper.Handle != nil {
			oper.Handle(one, nil)
		}
		return true
	})
	oper = nil
}

// 依据 人物的 被观察列表广播
func (self *AoiArea) broadcastIntersting(oper OperEvent) {
	entity, _ := oper.Entity.(Entity)
	switch me := entity.(type) {
	case AgentEntity:
		me.Neighbour().RangeBeenObservedSet(func(other Entity) bool {
			if oper.Handle != nil {
				oper.Handle(me, other)
			}
			return true
		})
	}
}

// 主动事件方法，都要通过它去执行
func (self *AoiArea) actionWithOper(op *OperEvent) {
	self.chanMove <- op
}

func (self *AoiArea) Enter(entity Entity) {
	op := &OperEvent{
		Op:     OPER_AGENT_JOIN,
		Entity: entity,
	}
	self.actionWithOper(op)
}

func (self *AoiArea) Move(entity Entity) {
	op := &OperEvent{
		Op:     OPER_AGENT_MOVE,
		Entity: entity,
	}
	self.actionWithOper(op)
}

func (self *AoiArea) Leave(entity Entity) {
	op := &OperEvent{
		Op:     OPER_AGENT_LEAVE,
		Entity: entity,
	}
	self.actionWithOper(op)
}

// 广播视野区域内的所有
func (self *AoiArea) BroadcastInterstingAll(entity Entity, fn func(mine, other Entity)) {
	op := &OperEvent{
		Op:     OPER_ACTION_INTERSTING_ALL,
		Entity: entity,
		Handle: fn,
	}
	self.actionWithOper(op)
}

// 广播给所有的用户
func (self *AoiArea) BroadcastAgentAll(fn func(mine, other Entity)) {
	op := &OperEvent{
		Op:     OPER_ACTION_AGENT_ALL,
		Handle: fn,
	}
	self.actionWithOper(op)
}

var _ AOI = &AoiArea{}
