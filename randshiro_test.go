package randshiro

import (
	"math"
	"testing"
)

const bound = 1<<16 - 1

func Benchmark512ppNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New512pp()
	}
}

func Benchmark256ppNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New256pp()
	}
}

func Benchmark128ppNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New128pp()
	}
}

func Benchmark512ppUint64(b *testing.B) {
	var rng = New512pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Uint64()
	}
}

func Benchmark256ppUint64(b *testing.B) {
	var rng = New256pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Uint64()
	}
}

func Benchmark128ppUint64(b *testing.B) {
	var rng = New128pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Uint64()
	}
}

func BenchmarkIntn_WorstCase(b *testing.B) {
	var rng = New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Intn(math.MaxInt)
	}
}

func Benchmark512ppIntn(b *testing.B) {
	var rng = New512pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Intn(bound)
	}
}

func Benchmark256ppIntn(b *testing.B) {
	var rng = New256pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Intn(bound)
	}
}

func Benchmark128ppIntn(b *testing.B) {
	var rng = New128pp()
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

func Benchmark256ppFloat64(b *testing.B) {
	var rng = New256pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Float64()
	}
}

func Benchmark128ppFloat64(b *testing.B) {
	var rng = New128pp()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Float64()
	}
}

func BenchmarkNormal(b *testing.B) {
	var rng = New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
