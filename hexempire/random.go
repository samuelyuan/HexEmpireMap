package hexempire

import "math"

// SeededRandom is a deterministic linear congruential generator.
// It deliberately avoids math/rand so that a given seed always reproduces
// the exact same generated map (relied on by TestGeneratedBoard).
type SeededRandom struct {
	seed int
}

func NewSeededRandom(seed int) *SeededRandom {
	return &SeededRandom{seed: seed}
}

// Intn returns a deterministic pseudo-random int in [0, n).
func (r *SeededRandom) Intn(n int) int {
	r.seed = (r.seed*9301 + 49297) % 233280
	return int(math.Floor(float64(r.seed) / 233280 * float64(n)))
}
