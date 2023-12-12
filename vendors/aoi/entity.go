package aoi

type Object struct {
}

func (self *Object) ID() int {
	return 0
}

func (self *Object) Enter() {

}
func (self *Object) Move() {

}
func (self *Object) Leave() {

}

func (self *Object) Position(args ...float32) []float32 {
	return nil
}
func (self *Object) PositionOld() []float32 {
	return nil
}