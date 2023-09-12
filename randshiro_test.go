package randshiro

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

// One in a million odds
const bound = 1000000

func newMathRand() *rand.Rand { return rand.New(rand.NewSource(time.Now().UnixNano())) }

func BenchmarkMathRandNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		newMathRand()
	}
}

func BenchmarkMathRandUint64(b *testing.B) {
	var rng = newMathRand()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Uint64()
	}
}

func BenchmarkMathRandIntn(b *testing.B) {
	var rng = newMathRand()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Intn(bound)
	}
}

func BenchmarkMathRandFloat64(b *testing.B) {
	var rng = newMathRand()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Float64()
	}
}

func BenchmarkMathRandFloat32(b *testing.B) {
	var rng = newMathRand()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Float32()
	}
}

func BenchmarkMathRandNormal(b *testing.B) {
	var rng = newMathRand()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.NormFloat64()
	}
}

func BenchmarkMathRandExponential(b *testing.B) {
	var rng = newMathRand()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.ExpFloat64()
	}
}

func Benchmark512ppNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New512pp()
	}
}

func Benchmark512ppUint64(b *testing.B) {
	var rng = New512pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Uint64()
	}
}

func Benchmark512ppIntn(b *testing.B) {
	var rng = New512pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Intn(bound)
	}
}

func Benchmark512ppFloat64(b *testing.B) {
	var rng = New512pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Float64()
	}
}

func Benchmark512ppFloat32(b *testing.B) {
	var rng = New512pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Float32()
	}
}

func Benchmark512ppFastFloat32(b *testing.B) {
	var rng = New512pp()
	b.ResetTimer()
	for i := 0; i < b.N/2; i++ {
		rng.FastFloat32()
	}
}

func Benchmark256ppNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New256pp()
	}
}

func Benchmark256ppUint64(b *testing.B) {
	var rng = New256pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Uint64()
	}
}

func Benchmark256ppIntn(b *testing.B) {
	var rng = New256pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Intn(bound)
	}
}

func Benchmark256ppFloat64(b *testing.B) {
	var rng = New256pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Float64()
	}
}

func Benchmark256ppFloat32(b *testing.B) {
	var rng = New256pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Float32()
	}
}

func Benchmark256ppFastFloat32(b *testing.B) {
	var rng = New256pp()
	b.ResetTimer()
	for i := 0; i < b.N/2; i++ {
		rng.FastFloat32()
	}
}

func Benchmark128ppNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New128pp()
	}
}

func Benchmark128ppUint64(b *testing.B) {
	var rng = New128pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Uint64()
	}
}

func Benchmark128ppIntn(b *testing.B) {
	var rng = New128pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Intn(bound)
	}
}

func Benchmark128ppFloat64(b *testing.B) {
	var rng = New128pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Float64()
	}
}

func Benchmark128ppFloat32(b *testing.B) {
	var rng = New128pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Float32()
	}
}

func Benchmark128ppFastFloat32(b *testing.B) {
	var rng = New128pp()
	b.ResetTimer()
	for i := 0; i < b.N/2; i++ {
		rng.FastFloat32()
	}
}

func BenchmarkNormal(b *testing.B) {
	var rng = New()
	b.ResetTimer()
	for i := 0; i < b.N/2; i++ {
		rng.Normal()
	}
}

func BenchmarkExponential(b *testing.B) {
	var rng = New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Exponential()
	}
}

func BenchmarkIntnWorstCase(b *testing.B) {
	var rng = New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Intn(math.MaxInt)
	}
}

func BenchmarkPerm3(b *testing.B) {
	var rng = New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Perm(3)
	}
}

func BenchmarkPerm30(b *testing.B) {
	var rng = New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Perm(30)
	}
}

func BenchmarkPerm100(b *testing.B) {
	var rng = New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Perm(100)
	}
}

func BenchmarkShuffleInts3(b *testing.B) {
	var rng = New()
	var slice = rng.Perm(3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Shuffle(rng, slice)
	}
}

func BenchmarkShuffleInts30(b *testing.B) {
	var rng = New()
	var slice = rng.Perm(30)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Shuffle(rng, slice)
	}
}

func BenchmarkShuffleInts100(b *testing.B) {
	var rng = New()
	var slice = rng.Perm(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Shuffle(rng, slice)
	}
}
