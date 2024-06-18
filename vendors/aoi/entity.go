package aoi

type Object struct {
	PID         int
	position    []float32
	positionOld []float32
	MessageUser AoiMessage
}

func newObject() *Object {
	return &Object{
		position:    []float32{0, 0, 0},
		positionOld: []float32{0, 0, 0},
	}
}

func (self *Object) ID() int {
	return self.PID
}

func (self *Object) Enter() {

}
func (self *Object) Move() {

}
func (self *Object) Leave() {

}

func (self *Object) AoiMessage() AoiMessage {
	return self.MessageUser
}

func (self *Object) PositionPre(args ...float32) []float32 {
	return nil
}

func (self *Object) Position(args ...float32) []float32 {
	if len(args) != 3 {
		return self.position
	}

	for i := 0; i < len(args); i++ {
		self.positionOld[i] = self.position[i]
		self.position[i] = args[i]
	}
	return self.position
}
func (self *Object) PositionOld() []float32 {
	return self.positionOld
}
