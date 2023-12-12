package aoi

import (
	"github.com/slclub/easy/vendors/option"
)

// 取出某个坐标轴的值，且转换成 链表的键值
func handleIndexForList(index int) HandleIndexFunc {
	return func(entity Entity) (uint64, uint64) {
		//value := entity.Position()[index]
		//value_int := int(value * 1000) + DEFAULT_MOVE_VALUE

		// 坐标冲突时，用此修正数值避免冲突，获取玩家id 的后两位数
		fix := uint64(entity.ID() % 100)
		return naomalIndexGeneral(entity.Position()[index]) + fix, naomalIndexGeneral(entity.PositionOld()[index]) + fix
	}
}

func handleNewListWithAxis(axisArr []int) option.OptionFunc {
	return func() (string, any) {
		lists := []*containerList{}
		for _, axis := range axisArr {
			lists = append(lists, newContainerList(handleIndexForList(axis)))
		}
		return "Lists", lists
	}
}
