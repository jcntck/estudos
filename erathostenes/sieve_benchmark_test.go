package sieve

import "testing"

var benchmarkPrimes []int

func BenchmarkPrimesUpTo1e6(b *testing.B) {
	b.Run("regular", func(b *testing.B) {
		runPrimesBenchmark(b, NewRegularSieveOfEratosthenes())
	})

	b.Run("optimized", func(b *testing.B) {
		runPrimesBenchmark(b, NewOptimizedSieveOfEratosthenes())
	})
}

func runPrimesBenchmark(b *testing.B, sieve SieveOfEratosthenes) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkPrimes = sieve.PrimesUpTo(1_000_000)
	}
}
