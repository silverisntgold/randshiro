package randshiro

import (
	"crypto/rand"
	"encoding/binary"
	"time"
	"unsafe"
)

const (
	bitsForFloat64 = 53
	float64Denom   = 1 << bitsForFloat64
	bitsForFloat32 = 24
	float32Denom   = 1 << bitsForFloat32
)

type randomBitGenerator interface {
	next() uint64
	set(uint64)
}

// Instances are not threadsafe
//
// It is recommended that each goroutine needing a source of random numbers
// should create and own a unique Gen instance
type Gen struct {
	randomBitGenerator
}

// Creates and seeds a *Gen with backing Xoshiro256++ instance
//
// Equivalent in all ways to directly calling New256pp()
func New() *Gen {
	return New256pp()
}

// Manually seeds the backing generator of the calling Gen instance
//
// Unless you know for certain that you need to
// manually seed a Gen instance, you don't
func (rng *Gen) ManualSeed(n uint64) {
	rng.set(n)
}

func seed(state []uint64, n uint64, useCustomSeed bool) {
	const bytesInUint64 = 8
	var seed = make([]byte, len(state)*bytesInUint64)
	if useCustomSeed {
		fallbackRead(seed, n, useCustomSeed)
	} else if _, err := rand.Read(seed); err != nil {
		fallbackRead(seed, n, useCustomSeed)
	}
	// Mapping groups of eight bytes from seed to unique indexs of state
	for i := range state {
		var startIndex = i * bytesInUint64
		var endIndex = startIndex + bytesInUint64
		// LittleEndian was chosen arbitrarily
		state[i] = binary.LittleEndian.Uint64(seed[startIndex:endIndex])
	}
}

func fallbackRead(seed []byte, n uint64, useCustomSeed bool) {
	var x = n
	if !useCustomSeed {
		x = uint64(time.Now().UnixMicro()) ^
			uint64(uintptr(unsafe.Pointer(&seed[0])))
	}
	var z uint64
	// SplitMix64
	// https://prng.di.unimi.it/splitmix64.c
	for i := range seed {
		x += 0x9e3779b97f4a7c15
		z = x
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		seed[i] = byte((z ^ (z >> 31)) >> 56)
	}
}
