package imager

import (
	"testing"
)

func TestAlignAddress(t *testing.T) {
	tests := []struct {
		address uint64
		alignment uint64
		expected uint64
	}{
		{0, 512, 0},
		{1, 512, 512},
		{512, 512, 512},
		{513, 512, 1024},
	}

	for _, i := range tests {
		result := alignAddress(i.address, i.alignment)
		if result != i.expected {
			t.Fatalf("Address %d, alignement %d: expected %d, got %d", i.address, i.alignment, i.expected, result)
		}
	}
}
