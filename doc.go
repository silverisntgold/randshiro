/*
Package randshiro is a Go implementation of the Xoroshiro128++, Xoshiro256++, and Xoshiro512++
pseudo random number generators, a subset of PRNGs from the Xoroshiro/Xoshiro family,
created by David Blackman and Sebastiano Vigna.

# Useful links

Their "PRNG shootout" can be found here:
https://prng.di.unimi.it/

Their paper on this family of PRNGs can be found here:
https://vigna.di.unimi.it/papers.php#BlVSLPNG

# Explanation

This implementation provides an API somewhat similar to that of the math/rand package,
but it is not a drop-in replacement. This is primarily a byproduct of how seeding/creation
is handled: randshiro automatically seeds each *Gen on creation using crypto/rand,
with the ability to fall back to a SplitMix64 implementation if that fails.
The internally-used generators themselves are also not accesible by end-users.
The user instead calls a factory function that returns a *Gen,
which internally manages one of the backing generators.
That being said, if you currently use math/rand and your code only calls Intn(), Float32(),
and/or Float64(), the only change you'll need to make is replacing your math/rand instantiations
with one of the randshiro factory functions.

Local initialization of a math/rand instance may look like:

	var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

But using randshiro:

	var rngDefault = randshiro.New()
	var rng128pp   = randshiro.New128pp()
	var rng256pp   = randshiro.New256pp()
	var rng512pp   = randshiro.New512pp()

New() is a convenience wrapper for New256pp(); they are equal in all ways.

Unlike math/rand, randshiro does not provide a global generator instance,
and it lacks a threadsafe variant.
I believe that a large part of the reason math/rand includes these is the high cost
of creating and maintaining instances. Since randshiro generators are cheap to
create and maintain it is entirely reasonable to create a unique generator
for each function or goroutine that needs one. Of course if your program is serial
and you want to create a single randshiro *Gen and pass it to functions that need it
(or, if you're feeling adventurous, put it in global scope),
there's nothing stopping you from doing so.

If a *Gen needs to be manually seeded, there is a ManualSeed()

	var rng = randshiro.New()
	var newSeed uint64 = 69420
	rng.ManualSeed(newSeed)

Unpredictability can only be restored by re-initializing the generator with another call to a factory function.
A dedicated Reseed() method is not provided, since I believe that would make manually seeding too accessible,
and you shouldn't be manually seeding these generators unless you are absolutely certain that you need to.
Since backing generator instances live behind an interface, it is not required that the factory function used
to re-seed is the same that was originally used for that *Gen.

Methods belonging to Gen do no verification of the variables passed into them.

# Performance

The main driving factor behind this implementation was the relatively poor performance/memory
utilization of math/rand.

At time of writing math/rand uses an additive lagged fibonacci generator
(https://en.wikipedia.org/wiki/Lagged_Fibonacci_generator) with a length of 607 and a tap of 273.
This means that for each math/rand instance an array of 607 int64 values (4856 bytes)
needs to be seeded and maintained.
Xoroshiro128++, Xoshiro256++, and Xoshiro512++ only use 16, 32, and 64 bytes, respectively.
Not only are they cheap to keep around in memory, they are [relatively] computationally
cheap to properly seed.

The math/rand package also uses an old method for bounding integers/floats to maintain
backwards compatable value streams. Replacing those methods is where most of the speed improvements come from.
The method this package uses for bounding integers is Daniel Lemire's nearly divisionless method:
https://arxiv.org/abs/1805.10941.

On my machine (AMD R7 2700 at 4.05 Ghz with 3200MHz memory on Go 1.20.5):

	Xoroshiro128++
		New: ~225 ns
		Uint64: ~2.7 ns
		Intn/Uint64n: ~3.2 ns
		Float64/Float32: ~2.8 ns
		FastFloat32: ~3.1ns (~1.55 ns per float32)
	Xoshiro256++ (randshiro default)
		New: ~250 ns
		Uint64: ~3.1 ns
		Intn/Uint64n: ~3.8 ns
		Float64/Float32: ~3.1 ns
		FastFloat32: ~3.8ns (~1.9 ns per float32)
		IntnWorstCase: ~8.2 ns
		Normal: ~32 ns (~16 ns per float64)
		Exponential: ~15 ns
	Xoshiro512++
		New: ~300 ns
		Uint64: ~4.2 ns
		Intn/Uint64n: ~5.3 ns
		Float64/Float32: ~4.2 ns
		FastFloat32: ~4.9ns (~2.45 ns per float32)
	math/rand
		New: ~8100 ns
		Uint64: ~4.2 ns
		Uint64Parallel: ~90ns (with 2 goroutines)
		Intn: ~9.6 ns
		Float64: ~4.3 ns
		Float32: ~4.9 ns
		Normal: ~6.6 ns
		Exponential: ~6 ns

All benchmarks for math/rand (with the exception of Uint64Parallel) are for local instances.
Using the global instance in parallel goroutines will absolutely decimate performance. If you're
currently doing so, moving to randshiro will likely result in orders of magnitude more performance.

Float64() and Float32() are both extremely fast due to being just a shift and multiplication,
and both are capable of generating all unique real numbers their type can accurately represent
in the interval [0.0, 1.0). That is, multiplying their output by 2^53 (float64) or 2^24 (float32)
will restore the exact ouput of the generator before it was converted to a float.
This is in contrast to the floating-point methods from math/rand,
whose outputs values are known to be far denser towards 0.
If you need to batch-generate float32s, FastFloat32() provides two float32s for a little more
than the cost of one Float32() call. Since the backing generators all output 64 bits and we only need
24 bits to create a float32 value, one call to the backing generator is used to create two float32s.
Don't retrieve float32s by calling Float64() and casting the result to float32,
performance will be exactly the same but some values will be unreachable due to how
float64s are rounded when casting to float32.
If you need float64s use Float64(); if you need float32s use Float32()/FastFloat32().
Explanantion of method used can be found at:
https://lemire.me/blog/2017/02/28/how-many-floating-point-numbers-are-in-the-interval-01/

IntnWorstCase calls it's Intn() with math.IntMax as the bound (and Intn() just wraps Uint64n()).
Intn()/Uint64n() execution time should never exceed this time for *Gen instances backed by Xoshiro256++,
but should be much closer to that of the 256ppIntn benchmark (bound of 1,000,000) for reasonable bounds.
Generally speaking, Intn()/Uint64n() excecution time increases the closer the bound is to math.IntMax.
If you need a bound that happens to be a power of two prefer using Uint64bits().

Normal() is provided for generating normally distributed float64s.
It uses the Marsaglia polar method (https://en.wikipedia.org/wiki/Marsaglia_polar_method)
and returns pairs of independent float64s. NormalDist() can be used to adjust the mean and stddev
of the returned float64s.

Exponential() is provided for generating exponentially distributed float64s.
It uses inverse transform sampling to generate its output:
https://en.wikipedia.org/wiki/Exponential_distribution#Random_variate_generation.

Other members of the Xoroshiro/Xoshiro PRNG family (+ and ** variants) were not included
in this package because testing showed near zero performance benefit for doing so.
In C/C++ they make more of a difference since a 0.14 ns jump is relatively big when your
initial time is already low (0.75 ns -> 0.61 ns for 256++ -> 256+, from the PRNG shootout).
But in Go much of the execution time is runtime overhead that can't be avoided,
so the jump is less meaningful. The jump is also smaller:
when testing the 256+ variant I was only seeing around 0.08 ns faster test results vs
the 256++ counterpart (~3.8 ns -> ~3.72 ns for Intn()).
No Xoroshiro1024 variant was included because there is no use for them in this package.

# Extra

A Fisher-Yates shuffle is also provided as Shuffle(), but it is a function belonging to the randshiro package
instead of a method belonging to *Gen. This is done to work around the inability to use generics in methods.
*/
package randshiro
