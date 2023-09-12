package randshiro

import "math/bits"

type x128pp [2]uint64

// Returns a seeded *Gen with backing Xoroshiro128++ instance
func New128pp() *Gen {
	var state x128pp
	seed(state[:])
	return &Gen{&state}
}

//go:noinline
func (state *x128pp) state() []uint64 {
	return state[:]
}

//go:noinline
func (state *x128pp) next() uint64 {
	var s0 = state[0]
	var s1 = state[1]
	var result = bits.RotateLeft64(s0+s1, 17) + s0

	s1 ^= s0
	state[0] = bits.RotateLeft64(s0, 49) ^ s1 ^ (s1 << 21)
	state[1] = bits.RotateLeft64(s1, 28)

	return result
}
