package randshiro

import "math/bits"

type x256pp [4]uint64

// Creates and seeds a *Gen with backing Xoshiro256++ instance
//
// Original C implementation: https://prng.di.unimi.it/xoshiro256plusplus.c
func New256pp() *Gen {
	var temp x256pp
	seed(temp[:])
	return &Gen{&temp}
}

//go:noinline
func (state *x256pp) set(n uint64) {
	for i := range state {
		n += entropy
		state[i] = n
	}
}

//go:noinline
func (state *x256pp) next() uint64 {
	var result = bits.RotateLeft64(state[0]+state[3], 23) + state[0]

	var t = state[1] << 17

	state[2] ^= state[0]
	state[3] ^= state[1]
	state[1] ^= state[2]
	state[0] ^= state[3]

	state[2] ^= t

	state[3] = bits.RotateLeft64(state[3], 45)

	return result
}
