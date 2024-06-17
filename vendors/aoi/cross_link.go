package aoi

import (
	"github.com/slclub/easy/vendors/option"
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

// 计算权重
func (self *crossList) caculateWeightOne(obj *containerList, count, min int) {
	switch count / min {
	case 0:
		obj.state = CONST_CONTAINER_WEIGHT_WORK
	case 1:
		obj.state = CONST_CONTAINER_WEIGHT_WORK
	case 2:
		obj.state = CONST_CONTAINER_WEIGHT_2
	case 3:
		obj.state = CONST_CONTAINER_WEIGHT_3
	case 4:
		obj.state = CONST_CONTAINER_WEIGHT_4
	default:
		obj.state = CONST_CONTAINER_WEIGHT_MORE
	}
}

func (self *crossList) caculateWeight(min int) {
	if min <= 0 {
		min = 1
	}
	for i, cl := range self.lists {
		self.caculateWeightOne(cl, self.countArr[i], min)
	}
}

func (self *crossList) step() int {
	self.stepWeight = (self.stepWeight + 1) % CONST_CROSS_STEP
	return self.stepWeight
}

func (self *crossList) rangeByRadius(this *containerList, entity Entity) []Entity {
	index_start, index_end := this.getSentinelIndex(entity, self.radius)
	node_start := this.getLatestNode(index_start)

	if node_start == nil {
		return nil
	}
	es := []Entity{}
	for node_start.Index() <= index_end {
		other := node_start.Value().(Entity)
		//this.Priority(1)
		if other.ID() == entity.ID() {
			goto NEXT
		}
		//logq.DebugF("AOI.RANGE.BY PID:=%v entity.start:=%v, endtity.end:=%v, other.v:=%v", entity.GetEntityID(), index_start, index_end, node_start.Index())
		es = append(es, other)
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
		near = near && compareRadius(entity.PositionOld()[i], dest.PositionOld()[i], self.radius)
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
	minCount := 10000
	step := self.step()
	entitys := []Entity{}
	for i, cl := range self.lists {
		if cl.state != CONST_CONTAINER_WEIGHT_WORK && step != 0 {
			continue
		}
		es := self.rangeByRadius(cl, entity)
		self.countArr[i] = len(es)
		if minCount > len(es) {
			minCount = len(es)
			entitys = es
		}
	}

	for _, one := range entitys {
		//// 计算  另外的坐标系
		//if !self.nearCheck(entity, one) {
		//	continue
		//}
		fn(one, self.compareRelation(entity, one))
	}

	if step == 0 {
		self.caculateWeight(minCount)
	}
}

func (self *crossList) compareRelation(entity, one Entity) int {
	near_new := self.nearCheck(entity, one)
	near_old := self.nearOldCheck(entity, one)

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
	entitys := []Entity{}
	for i, cl := range self.lists {
		if cl.state != CONST_CONTAINER_WEIGHT_WORK {
			continue
		}
		es := self.rangeByRadiusDiff(cl, entity)
		self.countArr[i] = len(es)
		if len(entitys) > len(es) {
			entitys = es
		}
	}
	for _, one := range entitys {
		fn(one)
	}
}

func (self *crossList) Range(fn func(entity Entity) bool) {
	self.lists[0].Range(fn)
}
