package randshiro

import "math/bits"

type x512pp [8]uint64

// Creates and seeds a *Gen with backing Xoshiro512++ instance
//
// Original C implementation: https://prng.di.unimi.it/xoshiro512plusplus.c
func New512pp() *Gen {
	var temp x512pp
	seed(temp[:])
	return &Gen{&temp}
}

//go:noinline
func (state *x512pp) set(n uint64) {
	for i := range state {
		n += entropy
		state[i] = n
	}
}

//go:noinline
func (state *x512pp) next() uint64 {
	var result = bits.RotateLeft64(state[0]+state[2], 17) + state[2]

	var t = state[1] << 11

	state[2] ^= state[0]
	state[5] ^= state[1]
	state[1] ^= state[2]
	state[7] ^= state[3]
	state[3] ^= state[4]
	state[4] ^= state[5]
	state[0] ^= state[6]
	state[6] ^= state[7]

	state[6] ^= t

	state[7] = bits.RotateLeft64(state[7], 21)

	return result
}
