package randshiro

import (
	"math"
	"math/bits"
	"reflect"
)

// Returns an unmodified uint64 directly from the
// backing generator instance
func (rng *Gen) Uint64() uint64 {
	return rng.next()
}

// Returns a uint64 in the interval [0, 2^bitcount)
//
// Makes no range checks on bitcount
func (rng *Gen) Uint64bits(bitcount int) uint64 {
	return rng.Uint64() >> (64 - bitcount)
}

// Returns a uint64 in the interval [0, bound)
//
// Makes no range checks on bound
func (rng *Gen) Uint64n(bound uint64) uint64 {
	var high, low = bits.Mul64(rng.Uint64(), bound)
	if low < bound {
		var threshold = -bound % bound
		for low < threshold {
			high, low = bits.Mul64(rng.Uint64(), bound)
		}
	}
	return high
}

// Returns an int in the interval [0, bound)
//
// Makes no range checks on bound
func (rng *Gen) Intn(bound int) int {
	return int(rng.Uint64n(uint64(bound)))
}

// Returns an int in the interval [lowerBound, upperBound)
//
// Makes no range checks on lowerBound/upperBound
func (rng *Gen) IntRange(lowerBound, upperBound int) int {
	return rng.Intn(upperBound-lowerBound) + lowerBound
}

// Returns a bool in the interval [false, true]
//
// xD
func (rng *Gen) Bool() bool {
	return rng.Uint64bits(1) == 1
}

// Returns a uniformly distributed float64 in the interval [0.0, 1.0)
//
// Don't cast the float64s produced by this function to float32:
// use Float32() instead
func (rng *Gen) Float64() float64 {
	return float64(rng.Uint64bits(53)) / 0x1p53
}

// Returns a uniformly distributed float32 in the interval [0.0, 1.0)
func (rng *Gen) Float32() float32 {
	return float32(rng.Uint64bits(24)) / 0x1p24
}

// Returns two uniformly distributed float64s in the interval (-1.0, 1.0)
//
// Normally float64s are made by shifting the uppermost 53 bits down,
// casting to a float64, and dividing by 0x1p53. Instead, we extract the
// uppermost bit to use as a sign bit, then use the next 53 bits
// to make the floating point. Then bitwise-or that sign bit
// with the bit representation of the generated float.
// The problem with this method is that you then generate
// both 0 and -0, so 0 occurs twice as frequently as other numbers.
// So we just reroll whenever we would get a -0 float.
// This only adds ~1 ns to average execution time and the probability of
// actually rerolling is only (1 / 2^54).
func (rng *Gen) forNormal() (float64, float64) {
	const mask = 1 << 63
	const shift = 64 - 53

loop1:
	var x = rng.Uint64()
	var signbit = x & mask
	x <<= 1
	x >>= shift
	if x == 0 && signbit == mask {
		goto loop1
	}
	var float1 = float64(x) / 0x1p53
	float1 = math.Float64frombits(signbit | math.Float64bits(float1))

loop2:
	x = rng.Uint64()
	signbit = x & mask
	x <<= 1
	x >>= shift
	if x == 0 && signbit == mask {
		goto loop2
	}
	var float2 = float64(x) / 0x1p53
	float2 = math.Float64frombits(signbit | math.Float64bits(float2))

	return float1, float2
}

// Returns two normally distributed float64s
// with mean = 0.0 and stddev = 1.0
//
// Use NormalDist() if you need to adjust mean/stddev
func (rng *Gen) Normal() (float64, float64) {
loop:
	var u, v = rng.forNormal()
	var s = u*u + v*v
	if s >= 1 || s == 0 {
		goto loop
	}
	s = math.Sqrt(-2 * math.Log(s) / s)
	return u * s, v * s
}

// Returns two normally distributed float64s
// with user-defined mean and stddev
//
// Makes no range checks on mean/stddev
func (rng *Gen) NormalDist(mean, stddev float64) (float64, float64) {
	var x, y = rng.Normal()
	return x*stddev + mean, y*stddev + mean
}

// Returns an exponentially distributed float64
// in the interval [0.0, +math.MaxFloat64]
//
// The rate constant (lambda) is 1,
// and can be adjusted with: Exponential() / lambda
func (rng *Gen) Exponential() float64 {
	// Uniformly distributed float64 in the interval (0.0, 1.0]
	var float = float64(rng.Uint64bits(53)+1) / 0x1p53
	return -math.Log(float)
}

// Returns a permutation of ints in the interval [0, n)
//
// Makes no range checks on n
func (rng *Gen) Perm(n int) []int {
	var temp = make([]int, n)
	var j int
	for i := range temp {
		j = rng.Intn(i + 1)
		temp[i] = temp[j]
		temp[j] = i
	}
	return temp
}

// This method only exists to tell you that the real Shuffle()
// is a function belonging to the randshiro package
func (rng *Gen) Shuffle() {}

// Performs a Fisher-Yates shuffle on contents of slice
//
// If len(slice) > 1 and rng == nil then
// Shuffle() will instantiate rng with New512pp()
// before continuing as normal. If Shuffle() is being called as a
// one-off it may be preferable to just pass nil to the rng
// parameter.
func Shuffle[T any](rng *Gen, slice []T) {
	if len(slice) > 1 {
		if rng == nil {
			rng = New512pp()
		}
		var swap = reflect.Swapper(slice)
		for i := len(slice) - 1; i > 0; i-- {
			swap(i, rng.Intn(i+1))
		}
	}
}
