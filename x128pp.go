package randshiro

import "math/bits"

type x128pp [2]uint64

// Creates and seeds a *Gen with backing Xoroshiro128++ instance
//
// Original C implementation: https://prng.di.unimi.it/xoroshiro128plusplus.c
func New128pp() *Gen {
	var temp x128pp
	seed(temp[:], 0 /* seed */, false /* seed used? */)
	return &Gen{&temp}
}

//go:noinline
func (state *x128pp) set(n uint64) {
	seed(state[:], n /* seed */, true /* seed used? */)
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
