package aoi

import (
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/vendors/option"
	"github.com/slclub/go-tips/logf"
	"github.com/slclub/log8q"
	"math/rand"
	"testing"
	"time"
)

func TestJoinAoi1(t *testing.T) {
	aoiObject, users := tInitAoiAndUsers(0)

	for i, _ := range users {
		aoiObject.Enter(users[i])
	}

	aoiObject.Clear()
}

func TestLeave1(t *testing.T) {
	aoiObject, users := tInitAoiAndUsers(900)

	for i, _ := range users {
		aoiObject.Enter(users[i])
	}
	aoiArea, _ := aoiObject.(*AoiArea)

	//time.Sleep(time.Millisecond * 5)
	for i, _ := range users {
		for j := 0; j < 10; j++ {
			aoiArea.cross.choose(users[i])
		}

		t.Logf("TestLevea1 crose.choose.index:%v u.x:%v u.y:%v u.z:%v  rate0=%v rate1=%v rate2=%v \n", aoiArea.cross.choose(users[i]).Index(),
			users[i].Position()[0],
			users[i].Position()[1],
			users[i].Position()[2],
			aoiArea.cross.lists[0].Rate(),
			aoiArea.cross.lists[1].Rate(),
			aoiArea.cross.lists[2].Rate(),
		)
		if i >= 55 {
			continue
		}

		aoiObject.Leave(users[i])
	}

	aoiObject.Clear()

	t.Logf("TestLevea1 aoi.Len:%v \n", aoiArea.Count(COUNT_AGENT))

	for i, cl := range aoiArea.cross.lists {
		t.Logf("TestLevea1 cross[%v].Len:%v \n", i, cl.Len())
	}
}

func TestMove1(t *testing.T) {
	log.LEVEL = log8q.ALL_LEVEL
	aoiObject, users := tInitAoiAndUsers(4)

	for i, _ := range users {
		aoiObject.Enter(users[i])
	}
	time.Sleep(time.Millisecond * 5)
	positions := [][]float32{
		[]float32{18, 0, 20},
		[]float32{19, 0, 17},
		[]float32{19.5, 0, 20.2},
	}

	// 预设 移动坐标
	users[0].PositionPre(positions[0]...)
	users[1].PositionPre(positions[1]...)
	users[3].PositionPre(positions[2]...)

	// 具体移动交给 entity.Move 设置到 Position才算真正的移动
	//aoiObject.Move(users[0])

	log.Debug(" --------- user2.Move Master.ID:%v", users[1].ID())
	aoiObject.Move(users[1])
	//aoiObject.Move(users[3])
	time.Sleep(time.Millisecond * 5)
	aoiObject.Clear()
}

func TestMove2(t *testing.T) {
	//log.LEVEL = log8q.ALL_LEVEL
	//aoiObject, users := tInitAoiAndUsers()
	//
	//for i, _ := range users {
	//	aoiObject.Enter(users[i])
	//}
	//time.Sleep(time.Millisecond * 5)
	//positions := [][]float32{
	//	[]float32{1, 0, 1},
	//	[]float32{5, 0, 50},
	//	[]float32{19.5, 0, 20.2},
	//}
	//for i := 0; i < 20; i++ {
	//
	//}
}

func tInitAoiAndUsers(n int) (AOI, []*user) {

	if n == 0 {
		n = 5
	}

	aoiObject := New(option.OptionWith(&struct {
		Radius float32
		Log    logf.Logger
		Axis   []int
	}{
		Radius: 10,
		Log:    log.Log(),
		Axis:   []int{0, 1, 2},
	}))

	users := []*user{}

	for i := 0; i < n; i++ {
		users = append(users, newUserWithAoi(aoiObject))
	}

	rand.Seed(time.Now().UnixNano())
	posotions := [][]float32{}
	for i := 0; i < n; i++ {
		randomFloatX := rand.Float32()
		randomFloatZ := rand.Float32()
		pos := []float32{float32(i) + randomFloatX, 0, float32(n-i) + randomFloatZ}
		posotions = append(posotions, pos)
	}

	for i, _ := range users {
		log.Info("AOI.Enter master.ID:%v", users[i].ID())
		users[i].PositionPre(posotions[i]...)
		users[i].Move()
	}
	return aoiObject, users
}
