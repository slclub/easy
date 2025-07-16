package aoi

import (
	"github.com/slclub/easy/vendors/option"
	"github.com/slclub/go-tips/logf"
)

//Orthogonal list

const (
	CONST_CROSS_STEP = 100
)

type crossList struct {
	lists      []*containerList
	stepWeight int // 计算权重的步数
	countArr   []int
	radius     float32
	Log        logf.Logger
}

func newCrossList(assignment option.Assignment) *crossList {
	c := &crossList{
		lists:    make([]*containerList, 0, 3),
		countArr: make([]int, 0, 3),
	}
	// 赋值初始化
	assignment.Target(c)
	assignment.Default(
		option.OptionFunc(func() (string, any) {
			return "Radius", DEFAULT_RADIUS
		}),
	)
	assignment.Apply()

	for i := 0; i < len(c.lists); i++ {
		c.countArr = append(c.countArr, 0)
	}
	return c
}

func (self *crossList) step() int {
	self.stepWeight = (self.stepWeight + 1) % CONST_CROSS_STEP
	return self.stepWeight
}

func (self *crossList) rangeByRadius(this *containerList, entity Entity) []Entity {
	index_start, index_end := this.getSentinelIndex(entity, self.radius)
	node_start := this.getLatestNode(index_start)
	//self.Log.Printf("Radius PID:%v index_start:%v index_end:%v node_start:%v", entity.ID(), index_start, index_end, node_start)
	if node_start == nil {
		return nil
	}
	this.ResetRate(true)
	es := []Entity{}
	for node_start.Index() <= index_end {
		other := node_start.Value().(Entity)
		//this.Priority(1)
		if other.ID() == entity.ID() {
			goto NEXT
		}
		//logq.DebugF("AOI.RANGE.BY PID:=%v entity.start:=%v, endtity.end:=%v, other.v:=%v", entity.GetEntityID(), index_start, index_end, node_start.Index())
		es = append(es, other)
		// 计算 玩家平均间隔 rate
		if l := len(es); l >= 2 {
			this.AutoRate(es[l-2], other)
		}
	NEXT:
		node_start = node_start.Next()
		if node_start == nil {
			break
		}
	}
	return es
}

// 差集
func (self *crossList) rangeByRadiusDiff(this *containerList, entity Entity) []Entity {
	// old 哨兵 节点 键
	_, index_old_index := this.handle(entity)
	index_old_start := index_old_index - uint64(self.radius*CHANGE_FLOAT_TO_INT)
	index_old_end := index_old_index + uint64(self.radius*CHANGE_FLOAT_TO_INT)
	// new 哨兵节点键
	index_new_start, index_new_end := this.getSentinelIndex(entity, self.radius)

	// 初始节点建，和 结束节点键
	index_start, index_end := uint64(0), uint64(0)

	if index_new_start == index_old_start {
		return nil
	}

	if index_new_start < index_old_start {
		index_start = index_new_end
		index_end = index_old_end
	} else {
		index_start = index_old_start
		index_end = index_new_start
	}

	node_start := this.getLatestNode(index_start)
	if node_start == nil {
		return nil
	}
	//logq.DebugF("AOI.RANGE.DIFF.BY PID:=%v entity.start:=%v, endtity.end:=%v, other.v:=%v", entity.GetEntityID(), index_start, index_end, node_start.Index())
	es := []Entity{}
	for node_start.Index() < index_end {
		k := node_start.Index()

		other := node_start.Value().(Entity)
		if k < index_start {
			goto NEXT
		}
		//this.Priority(1)
		if other.ID() == entity.ID() {
			goto NEXT
		}
		es = append(es, other)
	NEXT:
		node_start = node_start.Next()
		if node_start == nil {
			break
		}
	}
	return es
}

// 距离检查
func (self *crossList) nearCheck(entity Entity, dest Entity) bool {
	for i, _ := range self.countArr {
		near := compareRadius(entity.Position()[i], dest.Position()[i], self.radius)
		if !near {
			return near
		}
	}
	return true
}

func (self *crossList) nearOldCheck(entity, dest Entity) bool {
	near := true
	for i, _ := range self.countArr {
		near = near && compareRadius(entity.PositionOld()[i], dest.Position()[i], self.radius)
		if !near {
			return near
		}
	}
	return near
}

func (self *crossList) Add(entity Entity) {
	for _, cl := range self.lists {
		cl.add(entity)
	}
}

func (self *crossList) Delete(entity Entity) {
	for _, cl := range self.lists {
		cl.delete(entity, true)
	}
}

func (self *crossList) DeleteCache(entity Entity) {
	for _, cl := range self.lists {
		cl.delete(entity, false)
	}
}

func (self *crossList) RangeByRadius(entity Entity, fn func(other Entity)) {
	self.RangeByRadiusAll(entity, func(other Entity, check int) {
		if check == CONST_COORDINATE_EMPTY || check == CONST_COORDINATE_LEAVE {
			return
		}
		fn(other)
	})
}

func (self *crossList) RangeByRadiusAll(entity Entity, fn func(other Entity, check int)) {
	cl := self.choose(entity)
	entity_arr := self.rangeByRadius(cl, entity)

	for _, one := range entity_arr {
		//// 计算  另外的坐标系
		//if !self.nearCheck(entity, one) {
		//	continue
		//}
		fn(one, self.compareRelation(entity, one))
	}
}

func (self *crossList) choose(entity Entity) *containerList {
	stepnum := self.step()
	var cli *containerList
	var used = 0
	for i, cl := range self.lists {
		if cli == nil {
			cli = cl
			used = i
			continue
		}
		if cli.Rate() < cl.Rate() {
			cli = cl
			used = i
		}
	}
	if stepnum == 0 {
		for i, cl := range self.lists {
			if i == used {
				continue
			}
			self.rangeByRadius(cl, entity)
		}
	}
	return cli
}

func (self *crossList) RangeByAll(entity Entity, fn func(other Entity, check int)) {
	self.lists[0].Range(func(other Entity) bool {
		if entity.ID() == other.ID() {
			return true
		}
		fn(other, CONST_COORDINATE_MOVE)
		return true
	})
}

func (self *crossList) compareRelation(entity, one Entity) int {
	near_new := self.nearCheck(entity, one)
	near_old := self.nearOldCheck(entity, one)

	//log.Debug("compare PID:%v new:%v old:%v  pos:%v opos:%v", entity.ID(), near_new, near_old, entity.Position(), entity.PositionOld())
	if near_old && near_new {
		return CONST_COORDINATE_MOVE
	}
	if near_new && near_old == false {
		return CONST_COORDINATE_INCREASE
	}
	if near_new == false && near_old == true {
		return CONST_COORDINATE_LEAVE
	}
	return CONST_COORDINATE_EMPTY
}

func (self *crossList) RangeByRadiusDiff(entity Entity, fn func(other Entity)) {
	cl := self.choose(entity)

	entity_arr := self.rangeByRadius(cl, entity)

	for _, one := range entity_arr {
		fn(one)
	}
}

func (self *crossList) Range(fn func(entity Entity) bool) {
	self.lists[0].Range(fn)
}
