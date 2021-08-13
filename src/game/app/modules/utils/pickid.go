package utils

import (
	"math/rand"
	"time"
)

// ============================================================================

var rand_b = rand.New(rand.NewSource(time.Now().Unix()))

// ============================================================================

// rand pick id in map ([id]weight)
func PickWeightedMapId(slt map[int32]int32) (id int32) {
	if slt == nil {
		return
	}

	sum := int32(0)
	for _, w := range slt { // total weights
		sum += w
	}

	if sum <= 0 {
		return
	}

	r := rand_b.Int31n(sum)
	for i, v := range slt {
		if r -= v; r < 0 {
			id = i
			return
		}
	}

	return
}
