package aoi

const ()

func naomalIndexGeneral(f float32) uint64 {
	value_uint := int(f*CHANGE_FLOAT_TO_INT) + DEFAULT_MOVE_VALUE
	return uint64(value_uint)
}

// f1 和 f2 的距离小于半径内的
func compareRadius(f1, f2 float32, radius float32) bool {
	f := f1 - f2
	if f < 0 {
		f = f * -1
	}
	return f < radius
}
