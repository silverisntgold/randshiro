package randshiro

import (
	"crypto/rand"
	"encoding/binary"
	"time"
	"unsafe"
)

const (
	float64Bits  = 53
	float64Denom = 1 << float64Bits
	float32Bits  = 24
	float32Denom = 1 << float32Bits
)

type randomBitGenerator interface {
	getState() []uint64
	next() uint64
}

// Instances are not threadsafe
//
// It is recommended that each goroutine needing a source of random values
// should create and own a unique Gen instance
type Gen struct {
	randomBitGenerator
}

// Creates and cryptographically seeds a *Gen with backing Xoshiro256++ instance
//
// Equivalent in all ways to directly calling New256pp()
func New() *Gen {
	return New256pp()
}

// Manually seeds the backing generator of the calling Gen instance
//
// Unless you are absolutely certain that you need to
// manually seed a Gen instance, you don't
func (rng *Gen) ManualSeed(seed uint64) {
	alternateSeed(rng.getState(), seed)
}

func seed(state []uint64) {
	const bytesInUint64 = 8
	var randBytes = make([]byte, len(state)*bytesInUint64)
	if _, err := rand.Read(randBytes); err == nil {
		// Mapping sequences of eight bytes from randBytes to unique indexs of state
		for i := range state {
			var start = i * bytesInUint64
			var end = start + bytesInUint64
			// LittleEndian was chosen arbitrarily
			state[i] = binary.LittleEndian.Uint64(randBytes[start:end])
		}
	} else {
		var randSeed = uint64(time.Now().UnixMicro()) ^
			uint64(uintptr(unsafe.Pointer(&randBytes[0])))
		alternateSeed(state, randSeed)
	}
}

// SplitMix64
//
// https://prng.di.unimi.it/splitmix64.c
func alternateSeed(state []uint64, x uint64) {
	for i := range state {
		x += 0x9e3779b97f4a7c15
		var z = x
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		state[i] = z ^ (z >> 31)
	}
}
