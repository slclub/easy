package aoi

import "github.com/slclub/easy/log"

type messageUser struct {
	master AgentEntity
	Message
}

func (self *messageUser) Appear(entitys []Entity) {
	if len(entitys) == 0 {
		return
	}
	aes := []AgentEntity{}
	monsters := []Monster{}
	for _, v := range entitys {
		switch val := v.(type) {
		case AgentEntity:
			aes = append(aes, val)
		case Monster:
			monsters = append(monsters, val)
		case Npc:
		}
	}

	appear := &messageAppear{
		ID:      self.master.ID(),
		Players: aes,
	}
	if false {
		log.Info("MSG.APPEAR user.len ID:%v players num :%v", appear.ID, len(appear.Players))
		for _, v := range appear.Players {
			log.Info("=== APPEAR user.player MASTERID:%v ID:=%v position:%v:%v:%v", appear.ID, v.ID(), v.Position()[0], v.Position()[1], v.Position()[2])
		}
	}
}

func (self *messageUser) Disappear(entitys []Entity) {
	if len(entitys) == 0 {
		return
	}
	aes := []AgentEntity{}
	monsters := []Monster{}
	for _, v := range entitys {
		switch val := v.(type) {
		case AgentEntity:
			aes = append(aes, val)
		case Monster:
			monsters = append(monsters, val)
		case Npc:
		}
	}

	msg := &messageDisappear{
		ID:      self.master.ID(),
		Players: aes,
	}
	if false {
		log.Info("MSG.DISAPPEAR user.len ID:%v players num :%v", msg.ID, len(msg.Players))
		for _, v := range msg.Players {
			log.Info("--- DISAPPEAR user MASTERID:%v ID:=%v position:%v:%v:%v", msg.ID, v.ID(), v.Position()[0], v.Position()[1], v.Position()[2])
		}
	}
}

func (self *messageUser) Move(entitys []Entity) {
	msg := &messageMove{
		ID:      self.master.ID(),
		Objects: entitys,
	}
	log.Info("MOVE user.len ID:%v players num :%v", msg.ID, len(msg.Objects))
	//for _, v := range msg.Objects {
	//	log.Info("--- MOVE user MASTERID:%v ID:=%v position:%v:%v:%v", msg.ID, v.ID(), v.Position()[0], v.Position()[1], v.Position()[2])
	//}
}

// ------------------------------------------------------------------------
// message defined
// ------------------------------------------------------------------------

type messageAppear struct {
	ID      int
	Players []AgentEntity
}

type messageDisappear struct {
	ID      int
	Players []AgentEntity
}

type messageMove struct {
	ID      int
	Objects []Entity
}

type messageAppearMonster struct {
	ID       int
	Monsters []Monster
}
