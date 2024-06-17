package aoi

import (
	"github.com/slclub/easy/vendors/option"
	"github.com/slclub/go-tips/spinlock"
	"sync"
)

// neighbour
const (
	NEIGHBOUR_BEEN_OBSERVE_RATE = 5 // 被观察集合的人数倍率 按理是无需要限制的，5倍也接近不限制了

	NEIGHBOUR_WEIGHT_0 = 0 // 陌生人
	NEIGHBOUR_WEIGHT_1 = 1
	NEIGHBOUR_WEIGHT_2 = 2 // 亲近关系，直接加入到附近兴趣视野
	NEIGHBOUR_WEIGHT_3 = 3 // 同上
	NEIGHBOUR_WEIGHT_4 = 4 // 同上

	NEIGHBOUR_CLEAN = "clean"
)

type NeighbourConfigFunc func(collection *neighbourCollection)

// 视野区域集合
type neighbourCollection struct {
	master          Entity        //操作的主对象指针
	observedSet     *neighbourSet // 观察的集合
	beenObservedSet *neighbourSet // 被观察的集合,一般不做限制
	opt             *Option
}

func NewNeighbour(optfns ...NeighbourConfigFunc) *neighbourCollection {

	opt := DefaultOption()

	nei := &neighbourCollection{
		opt: opt,
	}

	nei.observedSet = newneighbourSet(opt)
	nei.beenObservedSet = newneighbourSet(&Option{
		NeighbourCount: opt.NeighbourCount * NEIGHBOUR_BEEN_OBSERVE_RATE,
		Radius:         opt.Radius,
	})
	nei.observedSet.master = nei.master
	nei.beenObservedSet.master = nei.master
	for _, fn := range optfns {
		fn(nei)
	}
	return nei
}

// 基本的视野容器集合
type neighbourSet struct {
	Option    *Option
	master    Entity
	list_move []Entity
	nlock     sync.Locker
}

// ---------neighbourSet-----------------------
// new
func newneighbourSet(opt *Option) *neighbourSet {
	return &neighbourSet{
		Option:    opt,
		list_move: make([]Entity, 0),
		nlock:     spinlock.New(),
	}
}

func (nc *neighbourSet) add(entity Entity) int {
	if nc.Option.NeighbourWeight == nil {
		return nc.addLimit(entity)
	}
	switch weight := nc.Option.NeighbourWeight.Value(entity, nc.master); weight {
	case NEIGHBOUR_WEIGHT_0, NEIGHBOUR_WEIGHT_1:
		return nc.addLimit(entity)
	case NEIGHBOUR_WEIGHT_2:
		return nc.addDirect(entity)
	case NEIGHBOUR_WEIGHT_3:
		return nc.addDirect(entity)
	case NEIGHBOUR_WEIGHT_4:
		return nc.addDirect(entity)
	default:
		return nc.addLimit(entity)
	}
}

// 进入/移动 (视野集合 添加操作
func (nc *neighbourSet) addLimit(entity Entity) int {
	limit := nc.Option.NeighbourCount
	nc.nlock.Lock()
	defer nc.nlock.Unlock()

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
	nc.list_move = append(nc.list_move, entity)
	return MESSAGE_EVENT_APPEAR
}

func (nc *neighbourSet) addDirect(entity Entity) int {
	nc.nlock.Lock()
	defer nc.nlock.Unlock()

	// 遍历移动
	for _, target := range nc.list_move {
		if entity.ID() == target.ID() {
			return MESSAGE_EVENT_MOVE
		}
	}
	nc.list_move = append(nc.list_move, entity)
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
		nc.list_move[i] = nc.list_move[n-1]
		nc.list_move = nc.list_move[:n-1]
		return MESSAGE_EVENT_DISAPPEAR
	}
	return MESSAGE_EVENT_EMPTY
}

//func (nc *neighbourSet)

func (nc *neighbourSet) len() int {
	//nc.nlock.Lock()
	//defer nc.nlock.Unlock()
	return len(nc.list_move)
}

func (nc *neighbourSet) getMove() []Entity {
	nc.nlock.Lock()
	defer nc.nlock.Unlock()
	return nc.list_move
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

func (nb *neighbourCollection) moveEntitys() []Entity {
	return nb.observedSet.list_move
}

func (nb *neighbourCollection) RangeMove(fn func(entity Entity) bool) {
	nb.observedSet.rangeMove(fn)
}

// --internal

func (nb *neighbourCollection) join(v any) (int, Entity) {
	//return 0 // PPROF.DELETE
	switch val := v.(type) {
	case AgentEntity:
		meCode := nb.observedSet.add(val)
		if meCode == MESSAGE_EVENT_APPEAR {
			val.Neighbour().beenJoin(nb.master)
		}
		return meCode, nil
	}
	return MESSAGE_EVENT_EMPTY, nil
}

// 被观察集合 添加
func (nb *neighbourCollection) beenJoin(v any) int {
	//return 0 // PPROF.DELETE
	switch val := v.(type) {
	case AgentEntity:
		nb.beenObservedSet.add(val)
	}
	//log.Debug("-------beenJoin ID:%v %v %v %v", nb.master.ID(),
	//	len(nb.beenObservedSet.list_increase), len(nb.beenObservedSet.list_move), len(nb.beenObservedSet.list_leave))
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

// @return
// @return 1 关系值
// @return 2 增加的集合
// @return 3 删除的集合
func (nb *neighbourCollection) relation(code int, entity Entity) (int, []Entity, []Entity) {
	if nb.opt.NeighbourCount == 0 {
		return nb.relationZero(code, entity)
	}
	rcode := MESSAGE_EVENT_EMPTY
	adds, leaves := []Entity{}, []Entity{}
	switch code {
	case CONST_COORDINATE_INCREASE:
		rrcode, leave_entity := nb.join(entity)
		rcode = rrcode
		if leave_entity != nil {
			leaves = append(leaves, leave_entity)
		}
		switch rcode {
		case MESSAGE_EVENT_APPEAR:
			adds = append(adds, entity)
		}
	case CONST_COORDINATE_MOVE:
		rrcode, leave_entity := nb.join(entity)
		rcode = rrcode
		if leave_entity != nil {
			leaves = append(leaves, leave_entity)
		}
		switch rcode {
		case MESSAGE_EVENT_APPEAR:
			adds = append(adds, entity)
		case MESSAGE_EVENT_MOVE:
		}
	case CONST_COORDINATE_LEAVE:
		rcode = nb.leave(entity)
		leaves = append(leaves, entity)
	case CONST_COORDINATE_EMPTY:
	}
	return rcode, adds, leaves
}

func (nb *neighbourCollection) relationZero(code int, entity Entity) (int, []Entity, []Entity) {
	adds, leaves := []Entity{}, []Entity{}
	switch code {
	case CONST_COORDINATE_INCREASE:
		adds = append(adds, entity)
	case CONST_COORDINATE_LEAVE:
		leaves = append(leaves, entity)
	case CONST_COORDINATE_MOVE:
	case CONST_COORDINATE_EMPTY:
	}
	return code, adds, leaves
}

func (nb *neighbourCollection) RangeBeenObservedSet(fn func(entity Entity) bool) {

}

var _ Neighbour = &neighbourCollection{}

// ---------neighbourSet-----------------------
