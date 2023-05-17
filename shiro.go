package randshiro

import (
	"crypto/rand"
	"encoding/binary"
	"time"
	"unsafe"
)

const entropy = 0x9e3779b97f4a7c15

type nextable interface {
	next() uint64
	set(uint64)
}

// Methods belonging to Gen do no verification of the variables passed into them;
// passing negative numbers or calling methods with nil *Gen instances is undefined behavior
//
// INSTANCES ARE NOT THREADSAFE
type Gen struct {
	nextable
}

// Creates and seeds a *Gen with backing Xoshiro256++ instance
//
// Equivalent in all ways to directly calling New256pp()
func New() *Gen {
	return New256pp()
}

// Manually seeds *Gen
//
// Restoring the unpredictability of the instance must be done
// by reinitializing with one of the factory functions
// (doesn't have to be the same as was originally used)
//
// Unless you know for absolute certain that you need to
// use this, you don't
func (rng *Gen) ManualSeed(n uint64) {
	const buildEntropy = 1 << 4
	for i := 0; i < buildEntropy; i++ {
		n += entropy
	}
	rng.set(n)
	for i := 0; i < buildEntropy; i++ {
		rng.next()
	}
}

func seed(slice []uint64) {
	// bytes in a uint64
	const n = 8
	var seed = make([]byte, len(slice)*n)
	if _, err := rand.Read(seed); err != nil {
		fallbackRead(seed)
	}
	for i := range slice {
		var startIndex = i * n
		var endIndex = startIndex + n
		slice[i] = binary.BigEndian.Uint64(seed[startIndex:endIndex])
	}
}

func fallbackRead(slice []byte) {
	var x = uint64(time.Now().UnixMicro()) ^
		uint64(uintptr(unsafe.Pointer(&slice[0])))
	var z uint64
	// SplitMix64
	// https://prng.di.unimi.it/splitmix64.c
	for i := range slice {
		x += entropy
		z = x
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		slice[i] = byte((z ^ (z >> 31)) >> 56)
	}
}
