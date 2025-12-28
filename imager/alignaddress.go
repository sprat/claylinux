package imager

func alignAddress(value, alignment uint64) uint64 {
	offset := alignment - 1
	return (value + offset) &^ offset
}
