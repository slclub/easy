package aoi

import (
	"github.com/slclub/easy/log"
	"github.com/slclub/easy/vendors/option"
	"github.com/slclub/log8q"
	"testing"
	"time"
)

func TestJoinAoi1(t *testing.T) {
	aoiObject, users := tInitAoiAndUsers()

	for i, _ := range users {
		aoiObject.Enter(users[i])
	}

	aoiObject.Clear()
}

func TestLeave1(t *testing.T) {
	aoiObject, users := tInitAoiAndUsers()

	for i, _ := range users {
		aoiObject.Enter(users[i])
	}

	//time.Sleep(time.Millisecond * 5)
	for i, _ := range users {
		if i > 3 {
			break
		}
		aoiObject.Leave(users[i])
	}

	aoiObject.Clear()
}

func TestMove1(t *testing.T) {
	log.LEVEL = log8q.ALL_LEVEL
	aoiObject, users := tInitAoiAndUsers()

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

func tInitAoiAndUsers() (AOI, []*user) {
	aoiObject := New(option.OptionWith(&struct{ Radius float32 }{
		Radius: 10,
	}))

	users := []*user{
		newUserWithAoi(aoiObject),
		newUserWithAoi(aoiObject),
		newUserWithAoi(aoiObject),
		newUserWithAoi(aoiObject),
		newUserWithAoi(aoiObject),
	}
	posotions := [][]float32{
		[]float32{1, 0, 1},
		[]float32{2, 0, 1.2},
		[]float32{5, 0, 7.2},
		[]float32{100, 0, 6.2},
		[]float32{14, 0, 10.2},
	}
	for i, _ := range users {
		log.Info("AOI.Enter master.ID:%v", users[i].ID())
		users[i].PositionPre(posotions[i]...)
		users[i].Move()
	}
	return aoiObject, users
}
