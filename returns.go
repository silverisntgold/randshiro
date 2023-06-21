package randshiro

import (
	"math"
	"math/bits"
	"reflect"
)

// Returns a uint64 in the interval [0, 2^64)
func (rng *Gen) Uint64() uint64 {
	return rng.next()
}

// Returns a uint64 in the interval [0, 2^bitcount)
//
// Makes no range checks on bitcount
func (rng *Gen) Uint64bits(bitcount int) uint64 {
	const bitsInUint64 = 64
	return rng.next() >> (bitsInUint64 - bitcount)
}

// Returns a uint64 in the interval [0, bound)
//
// Makes no range checks on bound
func (rng *Gen) Uint64n(bound uint64) uint64 {
	var high, low = bits.Mul64(rng.next(), bound)
	if low < bound {
		var threshold = -bound % bound
		for low < threshold {
			high, low = bits.Mul64(rng.next(), bound)
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
// use Float32() or FastFloat32() instead
func (rng *Gen) Float64() float64 {
	return float64(rng.Uint64bits(bitsForFloat64)) / float64Denom
}

// Returns a uniformly distributed float32 in the interval [0.0, 1.0)
//
// Don't cast the float32s produced by this function to float64:
// use Float64() instead
func (rng *Gen) Float32() float32 {
	return float32(rng.Uint64bits(bitsForFloat32)) / float32Denom
}

// Returns two independent and uniformly distributed float32s in the interval [0.0, 1.0)
//
// Don't cast the float32s produced by this function to float64:
// use Float64() instead
func (rng *Gen) FastFloat32() (float32, float32) {
	var (
		random48Bits = rng.Uint64bits(bitsForFloat32 * 2)
		bottom24Bits = random48Bits & (1<<bitsForFloat32 - 1)
		upper24Bits  = random48Bits >> bitsForFloat32
		float1       = float32(bottom24Bits) / float32Denom
		float2       = float32(upper24Bits) / float32Denom
	)
	return float1, float2
}

// Returns two independent and normally distributed float64s
// with mean = 0.0 and stddev = 1.0
//
// Use NormalDist() if you need to adjust mean/stddev
func (rng *Gen) Normal() (float64, float64) {
	const bitCount = bitsForFloat64 + 1
	const shiftValues = 1 << bitsForFloat64

	// It's a bit of a mess to have this all manually inlined,
	// but doing so saves ~2ns of runtime
outer_loop:

	// For generating a float64 in the interval (-1.0, 1.0), we roll
	// a random number in the interval [0, 2^54), discard rolls of zero,
	// and subtract 2^53 (we cast to int64 because we need negatives).
	// This gives us an integer in the interval (-2^53, 2^53),
	// which then maps to a float64 in the interval (-1.0, 1.0) when we
	// do our casting and division magic.
inner_loop_1:
	var temp = int64(rng.Uint64bits(bitCount))
	if temp == 0 {
		goto inner_loop_1
	}
	temp -= shiftValues
	var u = float64(temp) / float64Denom

inner_loop_2:
	temp = int64(rng.Uint64bits(bitCount))
	if temp == 0 {
		goto inner_loop_2
	}
	temp -= shiftValues
	var v = float64(temp) / float64Denom

	var s = u*u + v*v
	// We need whatever goes into math.Sqrt() to be positive and non-zero,
	// and since we already have a -2 we need another negative.
	// The s variable itself can't be negative (sum of squares), so the
	// result of math.Log() must be negative.
	// Both of these conditions can only be met when s is between 0 and 1.
	if s >= 1 || s == 0 {
		goto outer_loop
	}
	s = math.Sqrt(-2 * math.Log(s) / s)
	return u * s, v * s
}

// Returns two independent and normally distributed float64s
// with user-defined mean and stddev
//
// Makes no range checks on mean/stddev
func (rng *Gen) NormalDist(mean, stddev float64) (float64, float64) {
	var x, y = rng.Normal()
	return x*stddev + mean, y*stddev + mean
}

// Returns an exponentially distributed float64 with
// a rate constant (lambda) of 1
//
// Lambda can be adjusted with: Exponential() / lambda
func (rng *Gen) Exponential() float64 {
	// Uniformly distributed float64 in the interval (0.0, 1.0]
	var float = float64(rng.Uint64bits(bitsForFloat64)+1) / float64Denom
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
