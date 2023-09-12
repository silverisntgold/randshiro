package randshiro

import "math/bits"

type x512pp [8]uint64

// Returns a seeded *Gen with backing Xoshiro512++ instance
func New512pp() *Gen {
	var state x512pp
	seed(state[:])
	return &Gen{&state}
}

//go:noinline
func (state *x512pp) state() []uint64 {
	return state[:]
}

//go:noinline
func (state *x512pp) next() uint64 {
	var result = bits.RotateLeft64(state[0]+state[2], 17) + state[2]
	var temp = state[1] << 11

	state[2] ^= state[0]
	state[5] ^= state[1]
	state[1] ^= state[2]
	state[7] ^= state[3]
	state[3] ^= state[4]
	state[4] ^= state[5]
	state[0] ^= state[6]
	state[6] ^= state[7]

	state[6] ^= temp
	state[7] = bits.RotateLeft64(state[7], 21)

	return result
}
