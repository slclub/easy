package aoi

import (
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/vendors/option"
	"math"
	"testing"
	"time"
)

func TestCrossList(t *testing.T) {
	cross := initCrosslink()
	obj1 := newObjectTest()
	obj2 := newObjectTest()
	obj3 := newObjectTest()
	obj1.MessageUser = &messageUser{}
	obj2.MessageUser = &messageUser{}
	obj3.MessageUser = &messageUser{}

	obj1.Position(2.6, 1, -4)
	obj1.Position(2.6, 1, -4)
	obj2.Position(0.3, 1, -6)
	obj2.Position(0.3, 1, -6)

	obj3.Position(26, 1, -10)

	log.Info("cross init %v", cross.radius)
	testMockCrossMove(cross, obj3)
	for i := 0; i < 100; i++ {
		ob := obj1
		if i%2 == 0 {
			ob = obj2
			recalue(obj2, i)

		}
		time.Sleep(time.Millisecond * 1)
		testMockCrossMove(cross, ob)
	}

}

func recalue(obj *Object, i int) {
	z := obj.Position()
	d := math.Pow(-1, float64(i/10))
	k := z[2] + float32(float64(i%10)*d)
	obj.Position(z[0], z[1], k)
	//log.Debug("recaule PID:%v new:%v old:%v", obj.ID(), obj.Position(), obj.PositionOld())
}

func testMockCrossMove(cross *crossList, entity Entity) {
	//log.Debug("-----Handle------- Move cross PID:%v POS:%v", entity.ID(), entity.Position())
	cross.DeleteCache(entity)
	cross.Add(entity)
	cross.RangeByRadiusAll(entity, func(other Entity, check int) {
		//log.Debug("RANGE PID:%v POS:%v O-PID:%v POS:%v check:%v", entity.ID(), entity.Position(), other.ID(), other.Position(), check)
	})
}

var _pid = 1

func newObjectTest() *Object {
	obj := newObject()
	obj.PID = _pid
	_pid++
	return obj
}

func initCrosslink() *crossList {
	ax := []int{0, 1, 2}
	co := newCrossList(option.OptionWith(struct{ Radius float32 }{
		Radius: 15,
	}).Default(
		option.OptionFunc(func() (string, any) {
			return "Axis", ax
		}),
		option.OptionFunc(func() (string, any) {
			return "Log", log.Log()
		}),
		handleNewListWithAxis(ax),
		option.DEFAULT_IGNORE_ZERO,
	))
	return co
}
