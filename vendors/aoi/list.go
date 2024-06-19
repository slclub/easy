package aoi

import (
	"github.com/slclub/go-tips"
	"github.com/slclub/skiplist"
)

const (
	AXIS_X = 0 // x 坐标轴
	AXIS_Y = 1 // y 坐标轴  高度
	AXIS_Z = 2 // z 坐标轴

	CHANGE_FLOAT_TO_INT = 10000 // 做移动4位 就是乘这个值小数变整数
	DEFAULT_MOVE_VALUE  = CHANGE_FLOAT_TO_INT * 1000000

	DEFAULT_SKIP_LINK_LEVEL = 10

	CONST_CONTAINER_WEIGHT_WORK = 1 // 最小查找
	CONST_CONTAINER_WEIGHT_2    = 2
	CONST_CONTAINER_WEIGHT_3    = 3
	CONST_CONTAINER_WEIGHT_4    = 4
	CONST_CONTAINER_WEIGHT_MORE = 5
)

// 将实体坐标转换成 对应链表排序value 值
// 浮点=>整数
// 负数=>整数
// 顺带计算新旧坐标
type HandleIndexFunc func(entity Entity) (uint64, uint64)

/*
*
// 链表
// example ：

	BoxFirst: NewLinkBox(func(entity Entity) (uint64, uint64) {
				return HandleIndexGeneral(entity, AXIS_X) //从x坐标取值
	})
	BoxSecond: NewLinkBox(func(entity Entity) (uint64, uint64) {
		return HandleIndexGeneral(entity, AXIS_Z) //从z坐标取值
	})
*/
type containerList struct {
	list   *skiplist.ConcurrentSkipList
	handle HandleIndexFunc
	state  int   // default: CONST_CONTAINER_WEIGHT_WORK
	rate   []int // 覆盖率
	axis   int
}

func newContainerList(handle HandleIndexFunc, axis int) *containerList {
	skip, _ := skiplist.NewConcurrentSkipList(DEFAULT_SKIP_LINK_LEVEL)
	return &containerList{
		list:   skip,
		handle: handle,
		axis:   axis,
		state:  CONST_CONTAINER_WEIGHT_WORK,
		rate:   []int{0, 0},
	}
}

func (this *containerList) add(entity Entity) {
	index, _ := this.handle(entity)
	this.list.Insert(index, entity)
}

// 存在着内存泄漏的可能性，旧坐标有可能找不到值
func (this *containerList) delete(entity Entity, both bool) {
	index, old_index := this.handle(entity)
	this.list.Delete(old_index)
	if both {
		this.list.Delete(index)
	}
}

func (this *containerList) getSentinelIndex(entity Entity, radius float32) (uint64, uint64) {
	index, index_old := this.handle(entity)
	if index > index_old {
		return index_old - uint64(radius*CHANGE_FLOAT_TO_INT), index + uint64(radius*CHANGE_FLOAT_TO_INT)
	}
	return index - uint64(radius*CHANGE_FLOAT_TO_INT), index_old + uint64(radius*CHANGE_FLOAT_TO_INT)
}

func (this *containerList) getLatestNodes(index uint64) []*skiplist.Node {
	pnode, node := this.list.SearchCloset(index)
	if node != nil {
		return []*skiplist.Node{node}
	}
	if pnode == nil {
		return nil
	}
	if len(pnode) == 0 {
		return nil
	}
	return pnode
}

func (this *containerList) getLatestNode(index uint64) *skiplist.Node {
	pnodes := this.getLatestNodes(index)
	if pnodes == nil {
		return nil
	}
	return pnodes[0]
}
func (this *containerList) Range(fn func(entity Entity) bool) {
	this.list.ForEach(func(node *skiplist.Node) bool {
		entity, ok := node.Value().(Entity)
		if !ok {
			return true
		}
		return fn(entity)
	})
}

func (this *containerList) ResetRate(reset bool) {
	this.rate[0] = 1
	this.rate[1] = 0
}

func (this *containerList) Len() int32 {
	return this.list.Length()
}

func (this *containerList) Index() int {
	return this.axis
}

func (this *containerList) AutoRate(entity, entity2 Entity) {
	this.rate[0]++ // 数量
	long := tips.Int(entity2.Position()[this.Index()]*1000 - entity.Position()[this.Index()]*1000)
	if long < 0 {
		long *= -1
	}
	this.rate[1] += long
}

func (this *containerList) Rate() float32 {
	return float32(this.rate[1]) / float32(this.rate[0])
}
