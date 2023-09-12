package randshiro

import "math/bits"

type x256pp [4]uint64

// Returns a seeded *Gen with backing Xoshiro256++ instance
func New256pp() *Gen {
	var state x256pp
	seed(state[:])
	return &Gen{&state}
}

//go:noinline
func (state *x256pp) state() []uint64 {
	return state[:]
}

//go:noinline
func (state *x256pp) next() uint64 {
	var result = bits.RotateLeft64(state[0]+state[3], 23) + state[0]
	var temp = state[1] << 17

	state[2] ^= state[0]
	state[3] ^= state[1]
	state[1] ^= state[2]
	state[0] ^= state[3]

	state[2] ^= temp
	state[3] = bits.RotateLeft64(state[3], 45)

	return result
}
