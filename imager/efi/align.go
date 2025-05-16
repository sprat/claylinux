package efi

func align(value, alignment uint64) uint64 {
	offset := alignment - 1
	return (value + offset) &^ offset
}
