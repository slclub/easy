package aoi

import (
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/vendors/option"
	"github.com/slclub/go-tips/spinlock"
	"sync"
)

// neighbour
const (
	MESSAGE_EVENT_EMPTY     = 0
	MESSAGE_EVENT_APPEAR    = 1
	MESSAGE_EVENT_DISAPPEAR = 2
	MESSAGE_EVENT_MOVE      = 3

	NEIGHBOUR_BEEN_OBSERVE_RATE = 5 // 被观察集合的人数倍率 按理是无需要限制的，5倍也接近不限制了

	NEIGHBOUR_CLEAN = "clean"
)

// 视野区域集合
type neighbourCollection struct {
	master          Entity        //操作的主对象指针
	observedSet     *neighbourSet // 观察的集合
	beenObservedSet *neighbourSet // 被观察的集合,一般不做限制
	opt             *Option
}

func NewNeighbour(opt *Option) *neighbourCollection {
	if opt == nil {
		opt = DefaultOption()
	}
	return &neighbourCollection{
		//master:              ef(),
		observedSet: newneighbourSet(opt),
		beenObservedSet: newneighbourSet(&Option{
			NeighbourCount: opt.NeighbourCount * NEIGHBOUR_BEEN_OBSERVE_RATE,
			Radius:         opt.Radius,
		}),
		opt: opt,
	}
}

// 基本的视野容器集合
type neighbourSet struct {
	Option        *Option
	list_increase []Entity
	list_move     []Entity
	list_leave    []Entity
	nlock         sync.Locker
}

// ---------neighbourSet-----------------------
// new
func newneighbourSet(opt *Option) *neighbourSet {
	return &neighbourSet{
		Option:        opt,
		list_increase: make([]Entity, 0),
		list_move:     make([]Entity, 0),
		list_leave:    make([]Entity, 0),
		nlock:         spinlock.New(),
	}
}

// 进入/移动 (视野集合 添加操作
func (nc *neighbourSet) add(entity Entity) int {
	limit := nc.Option.NeighbourCount
	nc.nlock.Lock()
	defer nc.nlock.Unlock()
	//leave_len := nc.Len()
	for i, n := 0, len(nc.list_leave); i < n; i++ {
		if entity.ID() != nc.list_leave[i].ID() {
			continue
		}
		if nc.len() >= limit {
			return MESSAGE_EVENT_EMPTY
		}
		nc.list_move = append(nc.list_move, entity)
		// 快速清除 已经存在的元素
		nc.list_leave[i] = nc.list_leave[n-1]
		nc.list_leave = nc.list_leave[:n-1]
		return MESSAGE_EVENT_MOVE
	}
	// 遍历移动
	for _, target := range nc.list_move {
		if entity.ID() == target.ID() {
			return MESSAGE_EVENT_MOVE
		}
	}
	if nc.len() >= limit && limit > 0 {
		return MESSAGE_EVENT_EMPTY
	}
	// 填入到新增
	nc.list_increase = append(nc.list_increase, entity)
	return MESSAGE_EVENT_APPEAR
}

// entity 离开nc.master 的视野
func (nc *neighbourSet) leave(entity Entity) int {
	nc.nlock.Lock()
	defer nc.nlock.Unlock()
	for i, n := 0, len(nc.list_move); i < n; i++ {
		if entity.ID() != nc.list_move[i].ID() {
			continue
		}
		nc.list_leave = append(nc.list_leave, entity)
		nc.list_move[i] = nc.list_move[n-1]
		nc.list_move = nc.list_move[:n-1]
		return MESSAGE_EVENT_DISAPPEAR
	}
	for i, n := 0, len(nc.list_increase); i < n; i++ {
		if entity.ID() != nc.list_increase[i].ID() {
			continue
		}
		nc.list_leave = append(nc.list_leave, entity)
		nc.list_increase[i] = nc.list_increase[n-1]
		nc.list_increase = nc.list_increase[:n-1]
		return MESSAGE_EVENT_DISAPPEAR
	}
	return MESSAGE_EVENT_EMPTY
}

func (nc *neighbourSet) reset() {
	nc.nlock.Lock()
	defer nc.nlock.Unlock()
	// 将过去新增集合与移动集合 合并
	nc.list_move = append(nc.list_move, nc.list_increase[:nc.lenIncrease()]...)
	//
	//nc.list_leave = nc.list_leave[:0]
	// move， 新一轮，先将移动视野内的集合放在 离开的集合里； 后续新增的先查离开的集合 存在如移动的集合，不存在入 新增集合
	nc.list_leave = nc.list_move
	nc.list_move = []Entity{} //nb.list_agent_entity.list_move[0:0]
}

// 清除leave 队列，消息组装完事，将increase 合并到 move
func (nc *neighbourSet) clearLeave() {
	nc.nlock.Lock()
	defer nc.nlock.Unlock()
	nc.list_move = append(nc.list_move, nc.list_increase[:nc.lenIncrease()]...)
	nc.list_leave = nc.list_leave[:0]
	nc.list_increase = nc.list_increase[:0]
}

func (nc *neighbourSet) clear() {
	nc.nlock.Lock()
	defer nc.nlock.Unlock()
	nc.list_leave = []Entity{}    //nc.list_leave[:0]
	nc.list_increase = []Entity{} //nc.list_increase[:0]
	nc.list_move = []Entity{}     //nc.list_move[:0]
}

//func (nc *neighbourSet)

func (nc *neighbourSet) len() int {
	//nc.nlock.Lock()
	//defer nc.nlock.Unlock()
	return len(nc.list_move) //+ len(nc.list_increase)
}

func (nc *neighbourSet) getMove() []Entity {
	nc.nlock.Lock()
	defer nc.nlock.Unlock()
	return nc.list_move
}

func (nc *neighbourSet) rangeIncrease(fn func(entity Entity) bool) {
	nc.nlock.Lock()
	defer nc.nlock.Unlock()
	for _, obj := range nc.list_increase {
		rtn := fn(obj)
		if !rtn {
			break
		}
	}
}

func (nc *neighbourSet) rangeMove(fn func(entity Entity) bool) {
	nc.nlock.Lock()
	defer nc.nlock.Unlock()
	for _, obj := range nc.list_move {
		rtn := fn(obj)
		if !rtn {
			break
		}
	}
}

func (nc *neighbourSet) rangeLeave(fn func(entity Entity) bool) {
	for _, obj := range nc.list_leave {
		rtn := fn(obj)
		if !rtn {
			break
		}
	}
}
func (nc *neighbourSet) lenIncrease() int {
	il := len(nc.list_increase)
	ml := len(nc.list_move)
	if il+ml <= nc.Option.NeighbourCount {
		return il
	}
	cha := nc.Option.NeighbourCount - ml
	if cha < 0 {
		return 0
	}
	return cha
}

// ---------neighbourSet-----------------------

// ---------neighbourSet-----------------------
func (nb *neighbourCollection) BindWith(assignment option.Assignment) {
	//nb.master = e
	if assignment != nil {
		assignment.Target(nb)
		assignment.Apply()
	}
	nb.observedSet.Option = nb.opt
	nb.beenObservedSet.Option = &Option{
		Radius:         nb.opt.Radius,
		NeighbourCount: nb.opt.NeighbourCount * NEIGHBOUR_BEEN_OBSERVE_RATE,
	}
}

func (nb *neighbourCollection) increaseEntitys() []Entity {
	return nb.observedSet.list_increase
}

func (nb *neighbourCollection) RangeIncrease(fn func(entity Entity) bool) {
	nb.observedSet.rangeIncrease(fn)
}

func (nb *neighbourCollection) moveEntitys() []Entity {
	return nb.observedSet.list_move
}

func (nb *neighbourCollection) RangeMove(fn func(entity Entity) bool) {
	nb.observedSet.rangeMove(fn)
}

func (nb *neighbourCollection) leaveEntitys() []Entity {
	return nb.observedSet.list_leave
}
func (nb *neighbourCollection) RangeLeave(fn func(entity Entity) bool) {
	nb.observedSet.rangeLeave(fn)
}

func (nb *neighbourCollection) RangeBeenObservedSet(fn func(entity Entity) bool) {
	nb.beenObservedSet.rangeMove(fn)
	nb.beenObservedSet.rangeIncrease(fn)
	nb.beenObservedSet.rangeLeave(fn)
	//log.Debug("-------RangeBeenObservedSet %v %v %v", len(nb.beenObservedSet.list_increase), len(nb.beenObservedSet.list_move), len(nb.beenObservedSet.list_leave))
}

// --internal

func (nb *neighbourCollection) join(v any) int {
	//return 0 // PPROF.DELETE
	switch val := v.(type) {
	case AgentEntity:
		meCode := nb.observedSet.add(val)
		if meCode == MESSAGE_EVENT_APPEAR {
			val.Neighbour().beenJoin(nb.master)
		}
		return meCode
	}
	return MESSAGE_EVENT_EMPTY
}

// 被观察集合 添加
func (nb *neighbourCollection) beenJoin(v any) int {
	//return 0 // PPROF.DELETE
	switch val := v.(type) {
	case AgentEntity:
		nb.beenObservedSet.add(val)
	}
	log.Debug("-------beenJoin ID:%v %v %v %v", nb.master.ID(),
		len(nb.beenObservedSet.list_increase), len(nb.beenObservedSet.list_move), len(nb.beenObservedSet.list_leave))
	return MESSAGE_EVENT_EMPTY
}

func (nb *neighbourCollection) leave(v any) int {
	switch val := v.(type) {
	case AgentEntity:
		mCode := nb.observedSet.leave(val)
		val.Neighbour().beenLeave(nb.master)
		return mCode
	}
	return MESSAGE_EVENT_EMPTY
}

func (nb *neighbourCollection) beenLeave(v any) int {
	switch val := v.(type) {
	case AgentEntity:
		nb.beenObservedSet.leave(val)
	}
	return MESSAGE_EVENT_EMPTY
}

func (nb *neighbourCollection) reset(v any) {
	if op, ok := v.(string); ok && op == "clean" {
		nb.observedSet.clearLeave()
		nb.beenObservedSet.clearLeave()
		return
	}
	nb.observedSet.reset()
	nb.beenObservedSet.reset()
}

func (nb *neighbourCollection) clear() {
	nb.observedSet.clear()
	nb.beenObservedSet.clear()
}

var _ Neighbour = &neighbourCollection{}

// ---------neighbourSet-----------------------
